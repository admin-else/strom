package keepalive

import (
	"git.anygate.cloud/anygatecloud/strom/mc/proto"
	"git.anygate.cloud/anygatecloud/strom/mc/proto_generated/v1_21_8"
)

type mod struct {
	*proto.Conn
}

func (m *mod) onKeepAlive(p *v1_21_8.PlayToClientPacketKeepAlive) (err error) {
	return m.Send(&v1_21_8.PlayToServerPacketKeepAlive{KeepAliveId: p.KeepAliveId})
}

// Start begins responding to keep-alive packets on the connection.
func Start(c *proto.Conn) {
	m := &mod{Conn: c}
	m.RegisterUntilLatest(m.onKeepAlive)
}
