package client

import (
	"net"

	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/proto_base"
)

// Connect connects to the given address note that it does not resolve the SRV records its just a raw net.Dial
func Connect(addr string) (ret *proto.Conn, err error) {
	ret = &proto.Conn{}
	err = ret.SetVersion("1.21.8")
	if err != nil {
		return
	}
	ret.SetState(proto_base.Handshaking)
	ret.SetCompressionThreshold(-1)
	ret.Actor = proto_base.Client
	ret.Conn, err = net.Dial("tcp", addr)
	if err != nil {
		return
	}
	ret.R = ret.Conn
	ret.W = ret.Conn
	return
}
