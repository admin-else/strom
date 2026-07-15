package client

import (
	"context"
	"net"

	"github.com/admin-else/strom/mc/data"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_base"
	"github.com/admin-else/strom/mc/server"
)

// Connect connects to the given address note that it does not resolve the SRV records it's just a raw net.Dial
// You must close the connection when you're done with it
// version must be a vaild minecraft version eg 1.21.11 if left empty it will be detected with a status-request.
func Connect(ctx context.Context, addr string, version string) (ret *proto.Conn, err error) {
	var protocolVersion int32
	if version == "" {
		var s server.StatusResponse
		s, err = Status(ctx, addr)
		if err != nil {
			return
		}
		protocolVersion = s.Version.Protocol
	} else {
		var v data.VersionInfo
		v, err = data.LookUpVersionByName(version)
		if err != nil {
			return
		}
		protocolVersion = v.Version
	}
	ret, err = ConnectVersionLess(ctx, addr)
	if err != nil {
		return
	}
	err = ret.SetProtocolVersion(protocolVersion)
	return
}

// ConnectVersionLess connects to the given address without setting the version you should set it
// when using this function.
// You must close the connection when you're done with it
func ConnectVersionLess(ctx context.Context, addr string) (ret *proto.Conn, err error) {
	if ctx == nil {
		ctx = context.Background()
	}
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
