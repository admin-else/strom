package modules

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/admin-else/strom/bot"
	"github.com/admin-else/strom/proto_base"

	"github.com/admin-else/strom/data"
	"github.com/admin-else/strom/proto_generated/v1_21_8"
)

type StatusClient struct {
	*bot.Conn
	Status string
}

func (s *StatusClient) Default(event any) (err error) {
	err = fmt.Errorf("unexpected event: %v", event)
	return
}

func (s *StatusClient) OnStart(_ bot.OnStart) (err error) {
	parts := strings.Split(s.RemoteAddr().String(), ":")
	versionData, err := data.LookUpProtocolVersionByName(s.Version)
	if err != nil {
		return
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return
	}
	err = s.Send(v1_21_8.HandshakingToServerPacketSetProtocol{
		ProtocolVersion: proto_base.VarInt(versionData.Version),
		ServerHost:      parts[0],
		ServerPort:      uint16(port),
		NextState:       proto_base.VarInt(proto_base.Status),
	})
	if err != nil {
		return
	}
	s.State = proto_base.Status
	err = s.Send(v1_21_8.StatusToServerPacketPingStart{})
	return
}

func (s *StatusClient) OnStatus(p v1_21_8.StatusToClientPacketServerInfo) (err error) {
	s.Status = p.Response
	err = bot.HandlerDone
	return
}
