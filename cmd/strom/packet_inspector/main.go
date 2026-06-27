package packet_inspector

import (
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"time"

	api2 "github.com/admin-else/strom/mc/api"
	client2 "github.com/admin-else/strom/mc/client"
	"github.com/admin-else/strom/mc/event"
	"github.com/admin-else/strom/mc/nbt"
	proto2 "github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_base"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_11"
	server2 "github.com/admin-else/strom/mc/server"
	"github.com/admin-else/strom/mc/text"
)

var Cmd = flag.NewFlagSet("packet-inspector", flag.ContinueOnError)

var OfflineNameFlag = Cmd.String("name", "packet_inspector", "offline player name")
var Token = Cmd.String("token", "", "token to use")

var TargetAddr = Cmd.String("addr", "127.0.0.1:25565", "address to connect to")
var ListenAddr = Cmd.String("listen", "127.0.0.1:25566", "address to listen on")
var TrueForward = Cmd.Bool("true-forward", false, "forward packets without modification -name and custom server status wont work")

type Proxy struct {
	Client, Servee *proto2.Conn
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

type PacketLogEvent struct {
	Time       time.Time
	Packet     proto_base.EncodeDecodeAble
	Info       PacketInfoNoTypeJson
	ClientAddr string
}

type PacketInfoNoTypeJson struct {
	Type            proto_base.EncodeDecodeAble `json:"-"`
	Name            string
	Direction       proto_base.Direction
	State           proto_base.State
	PacketId        int32
	ProtocolVersion int32
}

func (p *Proxy) OnAnything(e event.Anything) (err error) {
	packet, ok := e.Val.(proto_base.EncodeDecodeAble)
	if !ok {
		return
	}
	packetInfo, ok := proto2.LookupPacketInfoByType(packet)
	if !ok {
		return
	}
	jsonBytes, err := json.Marshal(PacketLogEvent{
		Time:       time.Now(),
		Packet:     packet,
		ClientAddr: p.Client.RemoteAddr().String(),
		Info:       PacketInfoNoTypeJson(packetInfo),
	})
	if err != nil {
		return
	}
	// :pray:
	_, _ = os.Stdout.Write(jsonBytes)
	_, _ = os.Stdout.Write([]byte("\n"))

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

var StatusResponse = server2.StatusResponse{
	Version: server2.StatusResponseVersion{
		Name:     text.Pretty("STROM"),
		Protocol: 772,
	},
	Players:            server2.StatusResponsePlayers{},
	Description:        text.Pretty("Packet spy forwards to " + *TargetAddr),
	Favicon:            "",
	EnforcesSecureChat: false,
}

func Run(args []string) (err error) {
	err = Cmd.Parse(args)
	if err != nil {
		return
	}
	slog.Info("Starting packet inspector", "listen", *ListenAddr)
	err = server2.StartServerWithOnConn(*ListenAddr, func(serveeConn *proto2.Conn) (err error) {
		p := &Proxy{Servee: serveeConn, Client: nil}
		var acc = api2.NewOfflineAccount(*OfflineNameFlag)
		if !*TrueForward {
			if *Token != "" {
				acc, err = api2.NewAccountFromYGG(*Token)
				if err != nil {
					return
				}
			}
			_, err = server2.ServeLogin(serveeConn, server2.WithOtherAccount(acc), server2.WithStatus(StatusResponse), server2.WithoutOnlineMode(), server2.WithoutEncryption())
			if err != nil {
				return
			}
		}

		*TargetAddr, err = client2.DoDnsSimple(*TargetAddr)
		if err != nil {
			return
		}
		c, err := client2.Connect(*TargetAddr)
		if err != nil {
			return
		}
		defer c.Close()
		if !*TrueForward {
			err = client2.LoginRaw(c, acc)
			if err != nil {
				return
			}
		}
		p.Client = c
		p.Client.AlwaysDecode = true
		p.Servee.AlwaysDecode = true
		err = p.Start()
		return
	})
	return
}
