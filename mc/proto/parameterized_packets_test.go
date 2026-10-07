package proto_test

import (
	"bytes"
	"testing"

	"github.com/admin-else/strom/mc/data"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_base"
	"github.com/admin-else/strom/mc/proto_generated/v1_20_4"
	"github.com/admin-else/strom/mc/proto_generated/v26_4_snapshot_2"
)

// entityMetadataPacketBytes crafts a clientbound entity_metadata body: key 0,
// type "byte", value 42, end marker. Both 1.20.4 (parameterized
// entityMetadataItem) and 26.4 (inlined entityMetadataEntry) share this shape.
func entityMetadataPacketBytes(t *testing.T, version string) (b []byte, info proto_base.PacketInfo) {
	t.Helper()
	ver := data.MustLookupProtocolVersion(version)
	info, ok := proto.LookupPacketInfoByNameProtocolVersionStateAndDirection(
		"entity_metadata", ver, proto_base.Play, proto_base.ToClient,
	)
	if !ok {
		t.Fatalf("entity_metadata packet info not found for %s", version)
	}
	var buf bytes.Buffer
	if err := proto_base.EncodeVarInt(&buf, info.PacketId); err != nil {
		t.Fatalf("encode packet id: %v", err)
	}
	// entityId=7, key=0, type "byte", value 42, end marker.
	buf.Write([]byte{0x07, 0x00, 0x00, 0x2A, 0xFF})
	return buf.Bytes(), info
}

// TestEntityMetadataDecode120 proves the parameterized entityMetadataItem switch
// on 1.20.4 resolves to a real decoder instead of UnCodablePacket.
func TestEntityMetadataDecode120(t *testing.T) {
	raw, info := entityMetadataPacketBytes(t, "1.20.4")

	decoded, err := proto.SimpleBytesToPacket(raw, info.ProtocolVersion, info.Direction, info.State)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if unc, isUnc := decoded.(*proto.UnCodablePacket); isUnc {
		t.Fatalf("entity_metadata wrapped as UnCodablePacket: %v", unc.Err)
	}
	md := decoded.(*v1_20_4.PlayToClientPacketEntityMetadata)
	entry, found := md.Metadata.Val[0]
	if !found {
		t.Fatalf("metadata key 0 missing: %#v", md.Metadata.Val)
	}
	if entry.Type != "byte" {
		t.Errorf("metadata type = %q, want byte", entry.Type)
	}
	if v, isInt := entry.Value.(int8); !isInt || v != 42 {
		t.Errorf("metadata value = %#v, want int8(42)", entry.Value)
	}
}

// TestEntityMetadataDecode264 executes the same probe against 26.4-snapshot-2.
func TestEntityMetadataDecode264(t *testing.T) {
	raw, info := entityMetadataPacketBytes(t, "26.4-snapshot-2")

	decoded, err := proto.SimpleBytesToPacket(raw, info.ProtocolVersion, info.Direction, info.State)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if unc, isUnc := decoded.(*proto.UnCodablePacket); isUnc {
		t.Fatalf("entity_metadata wrapped as UnCodablePacket: %v", unc.Err)
	}
	md := decoded.(*v26_4_snapshot_2.PlayToClientPacketEntityMetadata)
	entry, found := md.Metadata.Val[0]
	if !found {
		t.Fatalf("metadata key 0 missing: %#v", md.Metadata.Val)
	}
	if entry.Type != "byte" {
		t.Errorf("metadata type = %q, want byte", entry.Type)
	}
	if v, isInt := entry.Value.(int8); !isInt || v != 42 {
		t.Errorf("metadata value = %#v, want int8(42)", entry.Value)
	}
}

// TestWorldParticlesParameterizedDecode round-trips a world_particles packet
// whose particleData switch is parameterized by the particle id.
func TestWorldParticlesParameterizedDecode(t *testing.T) {
	ver := data.MustLookupProtocolVersion("1.20.4")
	info, ok := proto.LookupPacketInfoByNameProtocolVersionStateAndDirection(
		"world_particles", ver, proto_base.Play, proto_base.ToClient,
	)
	if !ok {
		t.Fatal("world_particles packet info not found for 1.20.4")
	}

	dust := struct {
		Red   float32
		Green float32
		Blue  float32
		Scale float32
	}{0.25, 0.5, 0.75, 4}
	pkt := &v1_20_4.PlayToClientPacketWorldParticles{
		ParticleId: 14,
		X:          1,
		Y:          2,
		Z:          3,
		Particles:  1,
		Data:       dust,
	}

	raw, err := proto.SimplePacketToBytes(pkt)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	decoded, err := proto.SimpleBytesToPacket(raw, info.ProtocolVersion, info.Direction, info.State)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if unc, isUnc := decoded.(*proto.UnCodablePacket); isUnc {
		t.Fatalf("world_particles wrapped as UnCodablePacket: %v", unc.Err)
	}
	got := decoded.(*v1_20_4.PlayToClientPacketWorldParticles)
	if got.Data != dust {
		t.Errorf("particle data = %#v, want %#v", got.Data, dust)
	}
}
