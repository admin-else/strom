package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/admin-else/strom/api"
	"github.com/admin-else/strom/client"
	"github.com/admin-else/strom/client/modules"
	"github.com/admin-else/strom/data"
	"github.com/admin-else/strom/event"
	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/proto_generated/v1_21_8"
	"github.com/admin-else/strom/server"
	"github.com/admin-else/strom/text"
)

func wontFail[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

var (
	CompatibleProtocolVersion = wontFail(data.LookUpVersionByName("1.21.8")).Version
	FinishedConfigErr         = errors.New("finished configuration")
)

var Status = wontFail(json.Marshal(server.StatusResponse{
	Version: server.StatusResponseVersion{
		Name:     nil,
		Protocol: CompatibleProtocolVersion,
	},
	Description: text.Pretty("Registry data hunter"),
}))

type Proxy struct {
	Client, Servee *proto.Conn
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
	slog.Debug("packet", "type", packetInfo.Name, "packet", fmt.Sprintf("%#v", packet))
	switch packetInfo.Direction {
	case proto_base.ToServer:
		err = p.Client.Send(packet)
	case proto_base.ToClient:
		err = p.Servee.Send(packet)
	}
	return
}

func (p *Proxy) OnFinishConfiguration(_ *v1_21_8.ConfigurationToServerPacketFinishConfiguration) (err error) {
	err = FinishedConfigErr
	return
}

func (p *Proxy) OnData(data *v1_21_8.ConfigurationToClientPacketRegistryData) (err error) {
	return
}

func handleClient(c *proto.Conn) (err error) {
	acc := api.NewOfflineAccount("RegistryHunter")

	// ignore resource leak warning we own c conn so we should not close it here
	_, err = server.ServeLogin(c, server.WithOtherAccount(acc), server.WithCompatibleVersions(CompatibleProtocolVersion), server.WithRawStatus(Status))
	if err != nil {
		return
	}
	sockPuppet, err := client.Connect(":25566")
	if err != nil {
		return
	}
	defer sockPuppet.Close()
	err = sockPuppet.Start(&modules.LoginClient{
		Conn:    sockPuppet,
		Account: acc,
	})
	if err != nil {
		return
	}
	errChan := make(chan error)
	p := &Proxy{Client: sockPuppet, Servee: c}
	go func() {
		errChan <- p.Client.Start(p)
	}()
	go func() {
		errChan <- p.Servee.Start(p)
	}()
	err = <-errChan
	return
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	err := server.StartServerWithOnConn(":25565", handleClient)
	if err != nil {
		panic(err)
	}
}
