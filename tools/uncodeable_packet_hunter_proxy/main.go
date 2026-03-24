package main

import (
	"crypto/sha256"
	_ "embed"
	"fmt"
	"log/slog"
	"os"

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

//go:embed failed_packet.go.tmpl
var TestSrcF string

type Proxy struct {
	Client, Servee *proto.Conn
	Data           nbt.Tag
}

func (p *Proxy) Start() (err error) {
	p.Client.RegisterCritical(p.OnAnything)
	p.Servee.RegisterCritical(p.OnAnything)

	p.Servee.Register(p.OnUnCodeAble)
	p.Client.Register(p.OnUnCodeAble)

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

func (p *Proxy) OnUnCodeAble(packet *proto.UnCodablePacket) (err error) {
	SaveUnCodeAbleAsTest(packet.Direction, packet)
	return
}

func SaveUnCodeAbleAsTest(d proto_base.Direction, packet *proto.UnCodablePacket) {
	hUntrimmed := sha256.Sum256(packet.Data)
	h := hUntrimmed[:8]
	f, err := os.Create(fmt.Sprintf(".failed_packets/%v_%10x_test.go", len(packet.Data), h))
	if err != nil {
		panic(err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer f.Close()
	_, err = fmt.Fprintf(f, TestSrcF, packet.Err, h, d.Opposite(), packet.Data)
	if err != nil {
		panic(err)
	}
	slog.Info("Packet failed to parse", "error", packet.Err, "saved", f.Name(), "data", fmt.Sprintf("%q", packet.Data))
}

var StatusResponse = server.StatusResponse{
	Version: server.StatusResponseVersion{
		Name:     text.Pretty("STROM"),
		Protocol: 772,
	},
	Players:            server.StatusResponsePlayers{},
	Description:        text.Pretty("Hunt the un-code-able packets"),
	Favicon:            "",
	EnforcesSecureChat: false,
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	err := server.StartServerWithOnConn(":25566", func(serveeConn *proto.Conn) (err error) {
		p := &Proxy{Servee: serveeConn, Client: nil}
		acc := api.NewOfflineAccount("sigma")
		_, err = server.ServeLogin(serveeConn, server.WithOtherAccount(acc), server.WithStatus(StatusResponse))
		if err != nil {
			return
		}
		c, err := client.ConnectAndLogin("127.0.0.1:25565", acc)
		if err != nil {
			return
		}
		defer c.Close()
		p.Client = c
		err = p.Start()
		return
	})
	if err != nil {
		panic(err)
	}
}
