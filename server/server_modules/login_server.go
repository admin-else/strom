package server_modules

import (
	"fmt"

	"github.com/admin-else/strom/api"
	"github.com/admin-else/strom/event"
	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/proto_generated/v1_21_8"
	"github.com/google/uuid"
)

var UnexpectedStatusRequest = UnexpectedNextStateError{proto_base.Status}

type NameAndUUID struct {
	Name string
	UUID uuid.UUID
}

type UnexpectedNextStateError struct {
	NextState proto_base.State
}

func (u UnexpectedNextStateError) Error() string {
	return fmt.Sprintf("unexpected next state: %v", u.NextState)
}

type LoginServer struct {
	*proto.Conn
	//OnlineMode           bool TODO: implement this
	Requested            NameAndUUID
	Given                *NameAndUUID
	CompressionThreshold int32
}

func (l *LoginServer) OnHandshake(packet v1_21_8.HandshakingToServerPacketSetProtocol) (err error) {
	l.State = proto_base.State(packet.NextState)
	if packet.NextState != int32(proto_base.Login) {
		err = UnexpectedNextStateError{proto_base.State(packet.NextState)}
		return
	}
	return
}

func (l *LoginServer) SetCompressionThreshold(threshold int32) (err error) {
	err = l.Send(v1_21_8.LoginToClientPacketCompress{Threshold: threshold})
	if err != nil {
		return
	}
	l.CompressionThreshold = threshold
	return
}

func (l *LoginServer) OnLoginStart(packet v1_21_8.LoginToServerPacketLoginStart) (err error) {
	l.Requested = NameAndUUID{packet.Username, packet.PlayerUUID}
	//err = l.SetCompressionThreshold(l.CompressionThreshold)
	//if err != nil {
	//	return
	//}
	if l.Given == nil {
		l.Given = &l.Requested
	}
	err = l.Send(v1_21_8.LoginToClientPacketSuccess{
		Uuid:       l.Given.UUID,
		Username:   l.Given.Name,
		Properties: nil,
	})
	return
}

func (l *LoginServer) OnLoginAcknowledged(_ v1_21_8.LoginToServerPacketLoginAcknowledged) (err error) {
	l.State = proto_base.Configuration
	err = event.HandlerDone
	return
}

func (l *LoginServer) Default(event any) (err error) {
	err = fmt.Errorf("unexpected event during login: %#v", event)
	fmt.Println(err)
	return
}

func (l *LoginServer) OnStart(_ event.OnStart) (err error) {
	return
}

func (l *LoginServer) OnCycle(_ event.OnLoopCycle) (err error) {
	return
}

func ServeLogin(c *proto.Conn) (ret *LoginServer, err error) {
	ret = &LoginServer{Conn: c}
	ret.CompressionThreshold = 256
	err = ret.Start(ret)
	return
}

func ServeLoginWithOtherAccount(c *proto.Conn, a *api.Account) (ret *LoginServer, err error) {
	ret = &LoginServer{Conn: c, Given: &NameAndUUID{
		Name: a.Name,
		UUID: a.Uuid,
	}}
	ret.CompressionThreshold = 256
	err = ret.Start(ret)
	return
}
