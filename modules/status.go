package modules

import (
	"fmt"
	"strconv"
	"strings"
	"strom"

	"github.com/admin-else/queser"
	"github.com/admin-else/queser/data"
	"github.com/admin-else/queser/generated/v1_21_8"
)

type StatusClient struct {
	*strom.Conn
	Status string
}

func (s *StatusClient) Default(event any) (err error) {
	err = fmt.Errorf("unexpected event: %v", event)
	return
}

func (s *StatusClient) OnStart(_ strom.OnStart) (err error) {
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
		ProtocolVersion: queser.VarInt(versionData.Version),
		ServerHost:      parts[0],
		ServerPort:      uint16(port),
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
	s.Status = p.Response
	err = strom.HandlerDone
	return
}
