package raw_capture

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/admin-else/strom/mc/api"
	"github.com/admin-else/strom/mc/client"
	"github.com/admin-else/strom/mc/data"
	"github.com/admin-else/strom/mc/event"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto/replay"
	"github.com/admin-else/strom/mc/proto_base"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_11"
	"github.com/admin-else/strom/mc/server"
	"github.com/admin-else/strom/mc/text"
)

var Cmd = flag.NewFlagSet("raw-capture", flag.ContinueOnError)

var OfflineNameFlag = Cmd.String("name", "raw_capture", "offline player name")
var Token = Cmd.String("token", "", "token to use")

var TargetAddr = Cmd.String("addr", "127.0.0.1:25565", "address to connect to")
var ListenAddr = Cmd.String("listen", "127.0.0.1:25566", "address to listen on")
var ClientboundFile = Cmd.String("clientbound-file", "clientbound.tmcpr", "file to write clientbound packets to")
var ServerboundFile = Cmd.String("serverbound-file", "serverbound.tmcpr", "file to write serverbound packets to")
var Version = Cmd.String("version", "", "client protocol version (e.g. 26.1); if empty, detect from server status")

var StatusResponse = server.StatusResponse{
	Version: server.StatusResponseVersion{
		Name:     text.Pretty("STROM raw-capture"),
		Protocol: 0,
	},
	Players:            server.StatusResponsePlayers{},
	Description:        text.Pretty("STROM raw-capture proxy"),
	Favicon:            "",
	EnforcesSecureChat: false,
}

type Proxy struct {
	Client, Servee *proto.Conn

	clientboundFile *os.File
	serverboundFile *os.File

	clientboundWriter *replay.Writer
	serverboundWriter *replay.Writer

	startTime time.Time
}

func (p *Proxy) Start() (err error) {
	p.startTime = time.Now()

	p.Client.RegisterCritical(p.OnAnything)
	p.Servee.RegisterCritical(p.OnAnything)

	p.Client.RegisterCriticalUntilLatest(p.OnCompress)
	p.Servee.RegisterCriticalUntilLatest(p.OnHandshake)
	p.Servee.RegisterCriticalUntilLatest(p.OnLoginAcknowledged)
	p.Servee.RegisterCriticalUntilLatest(p.OnFinishConfiguration)

	errChan := make(chan error)
	go func() {
		errChan <- p.Client.StartConn()
	}()
	go func() {
		errChan <- p.Servee.StartConn()
	}()
	err = <-errChan
	return
}

func (p *Proxy) timestamp() uint32 {
	return uint32(time.Since(p.startTime).Milliseconds())
}

func (p *Proxy) OnAnything(e event.Anything) (err error) {
	packet, ok := e.Val.(proto_base.EncodeDecodeAble)
	if !ok {
		return
	}

	packetInfo, ok := proto.LookupPacketInfoByType(packet)
	if !ok {
		return
	}

	if packetInfo.Name == "encryption_begin" {
		return fmt.Errorf("encryption requested by remote server; raw capture cannot continue")
	}

	switch packetInfo.Direction {
	case proto_base.ToServer:
		err = p.serverboundWriter.WritePacket(packet, p.timestamp())
		if err != nil {
			return
		}
		err = p.Client.Send(packet)
	case proto_base.ToClient:
		err = p.clientboundWriter.WritePacket(packet, p.timestamp())
		if err != nil {
			return
		}
		err = p.Servee.Send(packet)
	}
	return
}

func (p *Proxy) OnCompress(packet *v1_21_11.LoginToClientPacketCompress) (err error) {
	err = p.clientboundWriter.WritePacket(packet, p.timestamp())
	if err != nil {
		return
	}
	err = p.Servee.Send(packet)
	if err != nil {
		return
	}
	p.Servee.SetCompressionThreshold(packet.Threshold)
	p.Client.SetCompressionThreshold(packet.Threshold)
	return event.DontForwardErr
}

func (p *Proxy) OnHandshake(packet *v1_21_11.HandshakingToServerPacketSetProtocol) (err error) {
	state := proto_base.State(packet.NextState)
	p.Servee.SetState(state)
	return event.DontForwardErr
}

func (p *Proxy) OnLoginAcknowledged(packet *v1_21_11.LoginToServerPacketLoginAcknowledged) (err error) {
	err = p.serverboundWriter.WritePacket(packet, p.timestamp())
	if err != nil {
		return
	}
	err = p.Client.Send(packet)
	if err != nil {
		return
	}
	p.Servee.SetState(proto_base.Configuration)
	p.Client.SetState(proto_base.Configuration)
	return event.DontForwardErr
}

func (p *Proxy) OnFinishConfiguration(packet *v1_21_11.ConfigurationToServerPacketFinishConfiguration) (err error) {
	err = p.serverboundWriter.WritePacket(packet, p.timestamp())
	if err != nil {
		return
	}
	err = p.Client.Send(packet)
	if err != nil {
		return
	}
	p.Servee.SetState(proto_base.Play)
	p.Client.SetState(proto_base.Play)
	return event.DontForwardErr
}

func Run(args []string) (err error) {
	err = Cmd.Parse(args)
	if err != nil {
		return
	}

	clientboundFile, err := os.Create(*ClientboundFile)
	if err != nil {
		return fmt.Errorf("create clientbound file: %w", err)
	}
	defer clientboundFile.Close()

	serverboundFile, err := os.Create(*ServerboundFile)
	if err != nil {
		return fmt.Errorf("create serverbound file: %w", err)
	}
	defer serverboundFile.Close()

	resolvedAddr, _, _, err := client.DoDNSChecked(nil, *TargetAddr, nil)
	if err != nil {
		return
	}

	var targetVersion string
	var targetProtocol int32
	if *Version == "" {
		status, err := client.Status(nil, resolvedAddr)
		if err != nil {
			return fmt.Errorf("failed to detect server version: %w", err)
		}
		targetProtocol = status.Version.Protocol
		versionInfo, err := data.LookUpVersionByProtocolVersion(targetProtocol)
		if err != nil {
			return fmt.Errorf("unknown protocol version %d: %w", targetProtocol, err)
		}
		targetVersion = versionInfo.MinecraftVersion
	} else {
		versionInfo, err := data.LookUpVersionByName(*Version)
		if err != nil {
			return fmt.Errorf("unknown version %s: %w", *Version, err)
		}
		targetVersion = versionInfo.MinecraftVersion
		targetProtocol = versionInfo.Version
	}

	StatusResponse.Version.Name = text.Pretty("STROM raw-capture " + targetVersion)
	StatusResponse.Version.Protocol = targetProtocol
	StatusResponse.Description = text.Pretty(fmt.Sprintf("Raw capture proxy forwards to %s (%s)", resolvedAddr, targetVersion))

	slog.Info("Starting raw capture", "listen", *ListenAddr, "target", resolvedAddr, "version", targetVersion, "protocol", targetProtocol, "clientbound", *ClientboundFile, "serverbound", *ServerboundFile)

	err = server.StartServerWithOnConn(*ListenAddr, func(serveeConn *proto.Conn) (err error) {
		proxy := &Proxy{
			Servee:            serveeConn,
			clientboundFile:   clientboundFile,
			serverboundFile:   serverboundFile,
			clientboundWriter: replay.NewWriter(clientboundFile),
			serverboundWriter: replay.NewWriter(serverboundFile),
		}

		acc := api.NewOfflineAccount(*OfflineNameFlag)
		if *Token != "" {
			acc, err = api.NewAccountFromYGG(*Token)
			if err != nil {
				return
			}
		}

		_, err = server.ServeLogin(serveeConn, server.WithOtherAccount(acc), server.WithStatus(StatusResponse), server.WithoutOnlineMode(), server.WithoutEncryption())
		if err != nil {
			return
		}

		c, err := client.Connect(nil, resolvedAddr, targetVersion)
		if err != nil {
			return
		}
		defer c.Close()

		err = client.LoginRaw(c, acc)
		if err != nil {
			return
		}

		proxy.Client = c
		proxy.Client.AlwaysDecode = true
		proxy.Servee.AlwaysDecode = true

		return proxy.Start()
	})
	return
}
