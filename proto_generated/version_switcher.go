package proto_generated

import (
	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/proto_generated/v1_21_8"
	"io"
)

func PacketIdentifierToType(v string, d proto_base.Direction, s proto_base.State, i string) (t any) {
	switch v {
	case "1.21.8":
		t = v1_21_8.PacketIdentifierToType(d, s, i)
	default:
	}
	return
}
func TypeToPacketIdentifier(v string, d proto_base.Direction, s proto_base.State, t any) (i string) {
	switch v {
	case "1.21.8":
		i = v1_21_8.TypeToPacketIdentifier(d, s, t)
	}
	return
}
func DecodePacket(v string, d proto_base.Direction, s proto_base.State, r io.Reader) (packet any, err error) {
	switch v {
	case "1.21.8":
		packet, err = v1_21_8.DecodePacket(d, s, r)
	}
	return
}
func EncodePacket(v string, d proto_base.Direction, s proto_base.State, i string, p any, w io.Writer) (err error) {
	switch v {
	case "1.21.8":
		err = v1_21_8.EncodePacket(d, s, i, p, w)
	}
	return
}
