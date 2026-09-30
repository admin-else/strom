package server

import (
	"errors"

	"fmt"

	"github.com/admin-else/strom/mc/nbt"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_base"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_8"
	"github.com/admin-else/strom/mc/text"
)

var (
	NotClientErr       = errors.New("not a servee")
	CantKickInStateErr = errors.New("cant kick in this state")
)

// Kick sends a disconnect/kick packet appropriate for the connection's current state.
func Kick(c *proto.Conn, reason *text.Component, status string) (err error) {
	if c.Actor != proto_base.Servee {
		err = NotClientErr
		return
	}
	switch c.State() {
	case proto_base.Handshaking:
		err = CantKickInStateErr
	case proto_base.Status:
		if status == "" {
			var b []byte
			b, err = reason.MarshalJSON()
			if err != nil {
				return
			}
			status = fmt.Sprintf("{\"description\":%s}", string(b))
		}
		err = c.Send(&v1_21_8.StatusToClientPacketServerInfo{Response: status})
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
