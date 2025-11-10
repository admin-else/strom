package main

import (
	"fmt"
	"strom"

	"github.com/admin-else/queser"
	"github.com/admin-else/queser/generated/v1_21_8"
)

type StatusClient struct {
	*strom.Connection
}

func (s *StatusClient) Default(event any) (err error) {
	err = fmt.Errorf("unexpected event: %v", event)
	return
}

func (s *StatusClient) OnStart(_ strom.OnStart) (err error) {
	err = s.Send(v1_21_8.HandshakingToServerPacketSetProtocol{
		ProtocolVersion: 772,
		ServerHost:      "127.0.0.1",
		ServerPort:      25566,
		NextState:       queser.VarInt(queser.Status),
	})
	if err != nil {
		return
	}
	s.State = queser.Status
	err = s.Send(v1_21_8.StatusToServerPacketPingStart{})
	return
}

func (s *StatusClient) OnStatus(p v1_21_8.StatusToClientPacketServerInfo) (err error) {
	fmt.Println(p)
	err = strom.HandlerDone
	return
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	c, err := strom.ClientConnect("127.0.0.1:25566")
	must(err)
	err = c.Start(&StatusClient{c})
	must(err)
}
