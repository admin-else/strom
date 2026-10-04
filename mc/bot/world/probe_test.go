package world

import (
	"bytes"
	"net"
	"testing"
	"time"

	"github.com/admin-else/strom/mc/api"
	"github.com/admin-else/strom/mc/client"
	"github.com/admin-else/strom/mc/level"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_generated/v26_4_snapshot_2"
)

const probeAddr = "127.0.0.1:25565"

type probe struct {
	*proto.Conn
	t *testing.T
	n int
}

func (p *probe) OnMap(m *v26_4_snapshot_2.PlayToClientPacketMapChunk) (err error) {
	if p.n == 0 {
		chunk, err := level.ReadChunkFromChunkPacketData(bytes.NewReader(m.ChunkData.Val), "26.4-snapshot-2", 384)
		p.t.Logf("err=%v sections=%d", err, len(chunk.Sections))
	}
	p.n++
	return
}

// TestProbe decodes a live 26.4 MapChunk against a local test server. It skips
// when no server is listening so the suite stays green without one.
func TestProbe(t *testing.T) {
	conn, err := net.DialTimeout("tcp", probeAddr, time.Second)
	if err != nil {
		t.Skipf("no test server on %s: %v", probeAddr, err)
	}
	conn.Close()

	acc := &api.Account{Name: "probe26_4"}
	c, err := client.Login(probeAddr, acc, client.WithVersion("26.4-snapshot-2"))
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	p := &probe{Conn: c, t: t}
	p.RegisterUntil("26.4-snapshot-2", p.OnMap)
	go func() { time.Sleep(9 * time.Second); c.Close() }()
	c.StartConn()
}
