package client

import (
	"net"
	"strconv"

	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/proto_generated/v1_21_8"
)

func MakeHandshakePacket(s *proto.Conn, nextState proto_base.State) (p *v1_21_8.HandshakingToServerPacketSetProtocol, err error) {
	p = &v1_21_8.HandshakingToServerPacketSetProtocol{NextState: int32(nextState), ProtocolVersion: s.ProtocolVersion}
	var portStr string
	p.ServerHost, portStr, err = net.SplitHostPort(s.RemoteAddr().String())
	if err != nil {
		return
	}
	var portUint uint64
	portUint, err = strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return
	}
	p.ServerPort = uint16(portUint)
	return
}
