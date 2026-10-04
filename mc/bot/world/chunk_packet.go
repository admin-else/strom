package world

import (
	"github.com/admin-else/strom/mc/level"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_11"
	"github.com/admin-else/strom/mc/proto_generated/v26_4_snapshot_2"
)

// chunkPacketFromV1_21_11 adapts the pre-26.4 MapChunk wire form, where light
// masks are long arrays and light payloads are raw nibble slices.
func chunkPacketFromV1_21_11(p *v1_21_11.PlayToClientPacketMapChunk) (ret *chunkPacket) {
	ret = &chunkPacket{
		X:         p.X,
		Z:         p.Z,
		ChunkData: p.ChunkData.Val,
	}
	for _, hm := range p.Heightmaps {
		ret.Heightmaps = append(ret.Heightmaps, level.Heightmap{Type: hm.Type, Data: hm.Data})
	}
	for _, be := range p.BlockEntities {
		ret.BlockEntities = append(ret.BlockEntities, level.BlockEntity{
			LocalX: be.Anon.X, LocalZ: be.Anon.Z, Y: be.Y, Type: be.Type, NbtData: be.NbtData.Value,
		})
	}
	ret.SkyLightMask = longArrayToBitSet(p.SkyLightMask)
	ret.BlockLightMask = longArrayToBitSet(p.BlockLightMask)
	ret.EmptySkyLightMask = longArrayToBitSet(p.EmptySkyLightMask)
	ret.EmptyBlockLightMask = longArrayToBitSet(p.EmptyBlockLightMask)
	for _, l := range p.SkyLight {
		ret.SkyLight = append(ret.SkyLight, l)
	}
	for _, l := range p.BlockLight {
		ret.BlockLight = append(ret.BlockLight, l)
	}
	return
}

// chunkPacketFrom26_4 adapts the 26.4 MapChunk wire form, where light masks and
// payloads are encoded as byte arrays (BitSet.toByteArray / DataLayer data).
func chunkPacketFrom26_4(p *v26_4_snapshot_2.PlayToClientPacketMapChunk) (ret *chunkPacket) {
	ret = &chunkPacket{
		X:                   p.X,
		Z:                   p.Z,
		ChunkData:           p.ChunkData.Val,
		SkyLightMask:        p.SkyLightMask.Val,
		BlockLightMask:      p.BlockLightMask.Val,
		EmptySkyLightMask:   p.EmptySkyLightMask.Val,
		EmptyBlockLightMask: p.EmptyBlockLightMask.Val,
	}
	for _, hm := range p.Heightmaps {
		ret.Heightmaps = append(ret.Heightmaps, level.Heightmap{Type: hm.Type, Data: hm.Data})
	}
	for _, be := range p.BlockEntities {
		ret.BlockEntities = append(ret.BlockEntities, level.BlockEntity{
			LocalX: be.Anon.X, LocalZ: be.Anon.Z, Y: be.Y, Type: be.Type, NbtData: be.NbtData.Value,
		})
	}
	for _, l := range p.SkyLight {
		ret.SkyLight = append(ret.SkyLight, l.Val)
	}
	for _, l := range p.BlockLight {
		ret.BlockLight = append(ret.BlockLight, l.Val)
	}
	return
}

// longArrayToBitSet converts a long array into the little-endian packed byte
// representation used by Java's BitSet.toByteArray.
func longArrayToBitSet(longs []int64) (ret []byte) {
	ret = make([]byte, len(longs)*8)
	for i, v := range longs {
		for b := range 8 {
			ret[i*8+b] = byte(v >> (uint(b) * 8))
		}
	}
	return
}
