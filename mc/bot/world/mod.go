package world

import (
	"bytes"
	"errors"
	"sync"

	"github.com/admin-else/strom/mc/level"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_11"
)

var (
	ChunkNotLoadedErr = errors.New("chunk not loaded")
	OutOfBoundsErr    = errors.New("coordinates out of bounds")
)

// ChunkPos identifies a chunk column.
type ChunkPos struct {
	X, Z int32
}

// World stores a shareable in-memory view of the Minecraft world.
// Multiple bots can subscribe to the same World and update it from their
// respective connections.
type World struct {
	mu      sync.RWMutex
	chunks  map[ChunkPos]*level.Chunk
	version string
	minY    int
	height  int
}

// NewWorld creates a new World for the given protocol version and vertical
// bounds. For a 1.21 overworld use minY=-64 and height=384.
func NewWorld(version string, minY, height int) *World {
	return &World{
		chunks:  make(map[ChunkPos]*level.Chunk),
		version: version,
		minY:    minY,
		height:  height,
	}
}

// GetBlock returns the global block state ID at the given world coordinates.
func (w *World) GetBlock(x, y, z int32) (stateId int32, err error) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.getBlockLocked(x, y, z)
}

func (w *World) getBlockLocked(x, y, z int32) (stateId int32, err error) {
	chunkX := floorDiv(x, level.ChunkWidth)
	chunkZ := floorDiv(z, level.ChunkWidth)
	chunk, ok := w.chunks[ChunkPos{chunkX, chunkZ}]
	if !ok {
		err = ChunkNotLoadedErr
		return
	}

	sectionIndex := (int(y) - w.minY) / level.ChunkWidth
	if sectionIndex < 0 || sectionIndex >= len(chunk.Sections) {
		err = OutOfBoundsErr
		return
	}

	lx, ly, lz := blockToLocal(x, y, z)
	index := ly*level.ChunkWidth*level.ChunkWidth + lz*level.ChunkWidth + lx
	return chunk.Sections[sectionIndex].Blocks.Get(index)
}

// SetBlock updates the global block state ID at the given world coordinates.
func (w *World) SetBlock(x, y, z int32, stateId int32) (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	chunkX := floorDiv(x, level.ChunkWidth)
	chunkZ := floorDiv(z, level.ChunkWidth)
	chunk, ok := w.chunks[ChunkPos{chunkX, chunkZ}]
	if !ok {
		return ChunkNotLoadedErr
	}

	sectionIndex := (int(y) - w.minY) / level.ChunkWidth
	if sectionIndex < 0 || sectionIndex >= len(chunk.Sections) {
		return OutOfBoundsErr
	}

	lx, ly, lz := blockToLocal(x, y, z)
	index := ly*level.ChunkWidth*level.ChunkWidth + lz*level.ChunkWidth + lx
	return chunk.Sections[sectionIndex].Blocks.Set(index, stateId)
}

// Module subscribes to chunk and block-change packets and updates a shared
// World. Multiple Modules can point to the same World.
type Module struct {
	*proto.Conn
	world *World
}

// GetBlock returns the block state ID at the given world coordinates via the
// shared World.
func (m *Module) GetBlock(x, y, z int32) (int32, error) {
	return m.world.GetBlock(x, y, z)
}

// Start registers packet handlers on c that update the shared World w.
func Start(c *proto.Conn, w *World) *Module {
	m := &Module{Conn: c, world: w}
	m.RegisterUntilLatest(m.onMapChunk)
	m.RegisterUntilLatest(m.onBlockChange)
	m.RegisterUntilLatest(m.onMultiBlockChange)
	return m
}

func (m *Module) onMapChunk(p *v1_21_11.PlayToClientPacketMapChunk) (err error) {
	chunk, err := level.ReadChunkFromChunkPacketData(
		bytes.NewReader(p.ChunkData.Val),
		m.world.version,
		m.world.height,
	)
	if err != nil {
		return
	}

	m.world.mu.Lock()
	m.world.chunks[ChunkPos{p.X, p.Z}] = chunk
	m.world.mu.Unlock()
	return nil
}

func (m *Module) onBlockChange(p *v1_21_11.PlayToClientPacketBlockChange) (err error) {
	return m.world.SetBlock(p.Location.X, int32(p.Location.Y), p.Location.Z, p.Type)
}

func (m *Module) onMultiBlockChange(p *v1_21_11.PlayToClientPacketMultiBlockChange) (err error) {
	chunkX := p.ChunkCoordinates.X
	chunkZ := p.ChunkCoordinates.Z
	sectionY := p.ChunkCoordinates.Y

	for _, record := range p.Records {
		packedPos := record & 0xFFF
		stateId := record >> 12

		localX := (packedPos >> 8) & 0xF
		localY := (packedPos >> 4) & 0xF
		localZ := packedPos & 0xF

		x := chunkX*level.ChunkWidth + int32(localX)
		y := sectionY*level.ChunkWidth + int32(localY)
		z := chunkZ*level.ChunkWidth + int32(localZ)

		if err = m.world.SetBlock(x, y, z, stateId); err != nil {
			return
		}
	}
	return nil
}

func floorDiv(a, b int32) int32 {
	return (a - ((a%b + b) % b)) / b
}

func blockToLocal(x, y, z int32) (lx, ly, lz int32) {
	lx = x % level.ChunkWidth
	if lx < 0 {
		lx += level.ChunkWidth
	}
	ly = y % level.ChunkWidth
	if ly < 0 {
		ly += level.ChunkWidth
	}
	lz = z % level.ChunkWidth
	if lz < 0 {
		lz += level.ChunkWidth
	}
	return
}
