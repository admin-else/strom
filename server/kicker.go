package server

import (
	"errors"

	"github.com/admin-else/strom/nbt"
	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/proto_generated/v1_21_8"
	"github.com/admin-else/strom/text"
)

var (
	NotClientErr       = errors.New("not a servee")
	CantKickInStateErr = errors.New("cant kick in this state")
)

func Kick(c *proto.Conn, reason *text.Component) (err error) {
	if c.Actor != proto_base.Servee {
		err = NotClientErr
		return
	}
	switch c.State() {
	case proto_base.Handshaking, proto_base.Status:
		err = CantKickInStateErr
	case proto_base.Login:
		var b []byte
		b, err = reason.MarshalJSON()
		if err != nil {
			return
		}
		err = c.Send(&v1_21_8.LoginToClientPacketDisconnect{Reason: string(b)})
	case proto_base.Configuration:
		err = c.Send(&v1_21_8.ConfigurationToClientPacketDisconnect{Reason: nbt.Anon{Value: reason.ToNBT()}})
	case proto_base.Play:
		err = c.Send(&v1_21_8.PlayToClientPacketKickDisconnect{Reason: nbt.Anon{Value: reason.ToNBT()}})
	}
	return
}
