package packet_inspector

import (
	"flag"
	"log/slog"

	"github.com/admin-else/strom/api"
	"github.com/admin-else/strom/client"
	"github.com/admin-else/strom/event"
	"github.com/admin-else/strom/nbt"
	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/proto_generated/v1_21_11"
	"github.com/admin-else/strom/server"
	"github.com/admin-else/strom/text"
)

var Cmd = flag.NewFlagSet("packet-inspector", flag.ContinueOnError)

var OfflineNameFlag = Cmd.String("name", "packet_inspector", "offline player name")
var Token = Cmd.String("token", "", "token to use")

var TargetAddr = Cmd.String("addr", "127.0.0.1:25565", "address to connect to")
var ListenAddr = Cmd.String("listen", "127.0.0.1:25566", "address to listen on")
var TrueForward = Cmd.Bool("true-forward", false, "forward packets without modification -name and custom server status wont work")

type Proxy struct {
	Client, Servee *proto.Conn
	Data           nbt.Tag
}

func (p *Proxy) Start() (err error) {
	p.Client.RegisterCritical(p.OnAnything)
	p.Servee.RegisterCritical(p.OnAnything)

	p.Servee.RegisterCriticalUntilLatest(p.OnFinishConfiguration)
	p.Servee.RegisterCriticalUntilLatest(p.OnFinishLogin)
	p.Servee.RegisterCriticalUntilLatest(p.OnHandshake)
	//p.Servee.RegisterCriticalUntilLatest(p.OnCompress)
	p.Client.RegisterCriticalUntilLatest(p.OnCompress)

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

func (p *Proxy) OnAnything(e event.Anything) (err error) {
	packet, ok := e.Val.(proto_base.EncodeDecodeAble)
	if !ok {
		return
	}
	packetInfo, ok := proto.LookupPacketInfoByType(packet)
	if !ok {
		return
	}
	switch packetInfo.Direction {
	case proto_base.ToServer:
		err = p.Client.Send(packet)
	case proto_base.ToClient:
		err = p.Servee.Send(packet)
	}
	return
}

func (p *Proxy) OnFinishConfiguration(packet *v1_21_11.ConfigurationToServerPacketFinishConfiguration) (err error) {
	err = p.Client.Send(packet)
	if err != nil {
		return
	}
	p.Servee.SetState(proto_base.Play)
	p.Client.SetState(proto_base.Play)
	err = event.DontForwardErr
	return
}

func (p *Proxy) OnFinishLogin(packet *v1_21_11.LoginToServerPacketLoginAcknowledged) (err error) {
	err = p.Client.Send(packet)
	if err != nil {
		return
	}
	p.Servee.SetState(proto_base.Configuration)
	p.Client.SetState(proto_base.Configuration)
	err = event.DontForwardErr
	return
}

func (p *Proxy) OnHandshake(packet *v1_21_11.HandshakingToServerPacketSetProtocol) (err error) {
	err = p.Client.Send(packet)
	state := proto_base.State(packet.NextState)
	if err != nil {
		return
	}
	p.Servee.SetState(state)
	p.Client.SetState(state)
	err = event.DontForwardErr
	return
}

func (p *Proxy) OnCompress(packet *v1_21_11.LoginToClientPacketCompress) (err error) {
	err = p.Servee.Send(packet)
	if err != nil {
		return
	}
	p.Servee.SetCompressionThreshold(packet.Threshold)
	p.Client.SetCompressionThreshold(packet.Threshold)
	err = event.DontForwardErr
	return
}

var StatusResponse = server.StatusResponse{
	Version: server.StatusResponseVersion{
		Name:     text.Pretty("STROM"),
		Protocol: 772,
	},
	Players:            server.StatusResponsePlayers{},
	Description:        text.Pretty("Packet spy forwards to " + *TargetAddr),
	Favicon:            "",
	EnforcesSecureChat: false,
}

func Run(args []string) (err error) {
	err = Cmd.Parse(args)
	if err != nil {
		return
	}
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Info("Starting packet inspector", "listen", *ListenAddr)
	err = server.StartServerWithOnConn(*ListenAddr, func(serveeConn *proto.Conn) (err error) {
		serveeConn.DebugPrintPackets = []string{""}
		p := &Proxy{Servee: serveeConn, Client: nil}
		var acc = api.NewOfflineAccount(*OfflineNameFlag)
		if !*TrueForward {
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
		}

		*TargetAddr, err = client.DoDnsSimple(*TargetAddr)
		if err != nil {
			return
		}
		c, err := client.Connect(*TargetAddr)
		if err != nil {
			return
		}
		defer c.Close()
		c.DebugPrintPackets = []string{""}
		if !*TrueForward {
			err = client.LoginRaw(c, acc)
			if err != nil {
				return
			}
		}
		p.Client = c
		err = p.Start()
		return
	})
	return
}
