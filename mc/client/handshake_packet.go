package client

import (
	"net"
	"strconv"

	"git.anygate.cloud/anygatecloud/strom/mc/proto"
	"git.anygate.cloud/anygatecloud/strom/mc/proto_base"
	"git.anygate.cloud/anygatecloud/strom/mc/proto_generated/v1_21_8"
)

// MakeHandshakePacket creates a handshake packet using the connection's remote address as the server host.
func MakeHandshakePacket(c *proto.Conn, nextState proto_base.State) (p *v1_21_8.HandshakingToServerPacketSetProtocol, err error) {
	return MakeHandshakePacketAddr(c, nextState, c.RemoteAddr().String())
}

// MakeHandshakePacketAddr gets a custom addr
func MakeHandshakePacketAddr(s *proto.Conn, nextState proto_base.State, addr string) (p *v1_21_8.HandshakingToServerPacketSetProtocol, err error) {
	p = &v1_21_8.HandshakingToServerPacketSetProtocol{NextState: int32(nextState), ProtocolVersion: s.ProtocolVersion}
	var portStr string
	p.ServerHost, portStr, err = net.SplitHostPort(addr)
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
