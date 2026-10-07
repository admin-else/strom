package world

import (
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/admin-else/strom/mc/api"
	"github.com/admin-else/strom/mc/client"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_generated/v26_4_snapshot_2"
)

type particleProbe struct {
	*proto.Conn
	t         *testing.T
	metadata  atomic.Int32
	particles atomic.Int32
	flame     atomic.Int32
}

func (p *particleProbe) OnEntityMetadata(m *v26_4_snapshot_2.PlayToClientPacketEntityMetadata) (err error) {
	if p.metadata.Add(1) == 1 {
		p.t.Logf("entity_metadata decoded: entityId=%d entries=%d", m.EntityId, len(m.Metadata.Val))
	}
	return
}

func (p *particleProbe) OnWorldParticles(pkt *v26_4_snapshot_2.PlayToClientPacketWorldParticles) (err error) {
	p.particles.Add(1)
	if pkt.Particle.Type == "flame" {
		p.flame.Add(1)
	}
	return
}

// TestProbeEntityMetadataAndParticles connects to the local 26.4 server, summons
// an armor stand and spawns a flame particle, and asserts both clientbound
// packets decode into their typed form instead of UnCodablePacket. It skips
// when no server is listening.
func TestProbeEntityMetadataAndParticles(t *testing.T) {
	conn, err := net.DialTimeout("tcp", probeAddr, time.Second)
	if err != nil {
		t.Skipf("no test server on %s: %v", probeAddr, err)
	}
	conn.Close()

	acc := api.NewOfflineAccount("dev")
	c, err := client.Login(probeAddr, acc, client.WithVersion("26.4-snapshot-2"))
	if err != nil {
		t.Fatal(err)
	}
	p := &particleProbe{Conn: c, t: t}
	p.RegisterUntil("26.4-snapshot-2", p.OnEntityMetadata)
	p.RegisterUntil("26.4-snapshot-2", p.OnWorldParticles)

	go func() {
		deadline := time.Now().Add(20 * time.Second)
		for time.Now().Before(deadline) {
			if p.metadata.Load() > 0 && p.flame.Load() > 0 {
				break
			}
			time.Sleep(1500 * time.Millisecond)
			_ = c.Send(&v26_4_snapshot_2.PlayToServerPacketChatCommand{
				Command: "summon minecraft:armor_stand ~ ~ ~",
			})
			_ = c.Send(&v26_4_snapshot_2.PlayToServerPacketChatCommand{
				Command: "particle minecraft:flame ~ ~2 ~ 1 1 1 0.1 30 force",
			})
		}
		c.Close()
	}()

	c.StartConn()

	if p.metadata.Load() == 0 {
		t.Error("no entity_metadata packet decoded (still UnCodablePacket?)")
	}
	if p.particles.Load() == 0 {
		t.Error("no world_particles packet decoded (still UnCodablePacket?)")
	} else if p.flame.Load() == 0 {
		t.Errorf("world_particles decoded but no flame particle (mapper ids wrong?)")
	}
}
