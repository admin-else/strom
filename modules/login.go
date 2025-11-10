package modules

import (
	"fmt"
	"strconv"
	"strings"
	"strom"

	"github.com/admin-else/queser"
	"github.com/admin-else/queser/data"
	"github.com/admin-else/queser/generated/v1_21_8"
	"github.com/google/uuid"
)

type Account struct {
	Username string
	UUID     uuid.UUID
	Token    string
}
type LoginClient struct {
	*strom.Connection
	Account Account
}

func (s *LoginClient) Default(event any) (err error) {
	err = fmt.Errorf("unexpected event: %T%v", event, event)
	return
}

func (s *LoginClient) OnStart(_ strom.OnStart) (err error) {
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
		NextState:       queser.VarInt(queser.Login),
	})
	if err != nil {
		return
	}
	s.State = queser.Login
	err = s.Send(v1_21_8.LoginToServerPacketLoginStart{Username: s.Account.Username, PlayerUUID: s.Account.UUID})
	return
}

func (s *LoginClient) OnCompress(compress v1_21_8.LoginToClientPacketCompress) (err error) {
	s.CompressionThreshold = int32(compress.Threshold)
	return
}

func (s *LoginClient) OnSuccess(success v1_21_8.LoginToClientPacketSuccess) (err error) {
	err = s.Send(v1_21_8.LoginToServerPacketLoginAcknowledged{})
	if err != nil {
		return
	}
	s.State = queser.Configuration
	return
}
