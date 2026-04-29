package client

import (
	"context"
	"net"

	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/proto_base"
)

// Connect connects to the given address note that it does not resolve the SRV records it's just a raw net.Dial
// You must close the connection when you're done with it
func Connect(addr string) (ret *proto.Conn, err error) {
	s, err := Status(addr)
	if err != nil {
		return
	}
	ret, err = ConnectVersionLess(addr)
	if err != nil {
		return
	}
	err = ret.SetProtocolVersion(s.Version.Protocol)
	return
}

func ConnectVersionLess(addr string) (ret *proto.Conn, err error) {
	return ConnectVersionLessCtx(context.Background(), addr)
}

// ConnectVersionLess connects to the given address without setting the version you should set it
// when using this function.
// You must close the connection when you're done with it
func ConnectVersionLessCtx(ctx context.Context, addr string) (ret *proto.Conn, err error) {
	ret = proto.NewConnCtx(ctx)
	ret.Actor = proto_base.Client
	ret.Conn, err = (&net.Dialer{}).DialContext(ctx, "tcp", addr)
	if err != nil {
		return
	}
	ret.R = ret.Conn
	ret.W = ret.Conn
	return
}
