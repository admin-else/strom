package world

import (
	"bytes"
	"errors"
	"sync"

	"github.com/admin-else/strom/mc/level"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_11"
	"github.com/admin-else/strom/mc/proto_generated/v26_4_snapshot_2"
	"github.com/admin-else/strom/mc/registry"
)

var (
	ChunkNotLoadedErr  = errors.New("chunk not loaded")
	OutOfBoundsErr     = errors.New("coordinates out of bounds")
	HeightmapAbsentErr = errors.New("heightmap not present")
)

// ChunkPos identifies a chunk column.
type ChunkPos struct {
	X, Z int32
}

// World stores a shareable in-memory view of the Minecraft world.
// Multiple bots can subscribe to the same World and update it from their
// respective connections.
type World struct {
	mu         sync.RWMutex
	chunks     map[ChunkPos]*level.Chunk
	revisions  map[ChunkPos]int64
	version    string
	minY       int
	height     int
	center     ChunkPos
	registries *registry.Store

	// Weather state mirrors net.minecraft.world.level.Level. The current levels
	// are lerped against the previous tick's levels with the partial tick.
	oRainLevel, rainLevel       float32
	oThunderLevel, thunderLevel float32
}

// NewWorld creates a new World for the given protocol version and vertical
// bounds. For a 1.21 overworld use minY=-64 and height=384.
func NewWorld(version string, minY, height int) *World {
	return &World{
		chunks:    make(map[ChunkPos]*level.Chunk),
		revisions: make(map[ChunkPos]int64),
		version:   version,
		minY:      minY,
		height:    height,
	}
}

// MinY returns the world's minimum block Y.
func (w *World) MinY() int { return w.minY }

// Height returns the world's total block height.
func (w *World) Height() int { return w.height }

// Version returns the world's protocol version string.
func (w *World) Version() string { return w.version }

// SetRegistries attaches the dynamic registries captured during configuration,
// letting consumers resolve wire ids to resource ids and back.
func (w *World) SetRegistries(store *registry.Store) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.registries = store
}

// Registries returns the attached registry store, or nil when none was set.
func (w *World) Registries() (store *registry.Store) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.registries
}

// RegistryID returns the wire id of resourceID inside registryName, resolving
// through the attached registry store.
func (w *World) RegistryID(registryName, resourceID string) (id int, ok bool) {
	store := w.Registries()
	if store == nil {
		return
	}
	return store.RegistryID(registryName, resourceID)
}

// ResourceID returns the resource id of wire id inside registryName, resolving
// through the attached registry store.
func (w *World) ResourceID(registryName string, id int) (resourceID string, ok bool) {
	store := w.Registries()
	if store == nil {
		return
	}
	return store.ResourceID(registryName, id)
}

// Center returns the last chunk-cache center received via UpdateViewPosition.
func (w *World) Center() ChunkPos {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.center
}

// Chunks returns a snapshot of the loaded chunk positions.
func (w *World) Chunks() (ret []ChunkPos) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	for pos := range w.chunks {
		ret = append(ret, pos)
	}
	return
}

// ChunkCount returns the number of loaded chunks.
func (w *World) ChunkCount() (n int) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return len(w.chunks)
}

// ChunkAt returns the loaded chunk at the given chunk position.
func (w *World) ChunkAt(cx, cz int32) (chunk *level.Chunk, err error) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	chunk, ok := w.chunks[ChunkPos{cx, cz}]
	if !ok {
		err = ChunkNotLoadedErr
		return
	}
	return
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

// BlockStateIdAt returns the global block state ID and whether the chunk is loaded.
func (w *World) BlockStateIdAt(x, y, z int32) (stateId int32, ok bool) {
	id, err := w.GetBlock(x, y, z)
	if err != nil {
		return 0, false
	}
	return id, true
}

// LightAt returns the block and sky light at the given world coordinates.
func (w *World) LightAt(x, y, z int32) (block, sky uint8, err error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

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
	block, sky = chunk.LightAt(sectionIndex, int(lx), int(ly), int(lz))
	return
}

// HeightAt returns the height for the given heightmap kind at world column (x, z).
func (w *World) HeightAt(kind string, x, z int32) (h int32, err error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	chunkX := floorDiv(x, level.ChunkWidth)
	chunkZ := floorDiv(z, level.ChunkWidth)
	chunk, ok := w.chunks[ChunkPos{chunkX, chunkZ}]
	if !ok {
		err = ChunkNotLoadedErr
		return
	}
	lx, _, lz := blockToLocal(x, 0, z)
	hv, ok := chunk.HeightAt(kind, int(lx), int(lz))
	if !ok {
		err = HeightmapAbsentErr
		return
	}
	// Heightmap values are stored relative to the world's minimum Y
	// (vanilla Heightmap.getFirstAvailable adds chunk.getMinY()).
	return hv + int32(w.minY), nil
}

// ChunkRevision returns a per-chunk counter that is bumped whenever the stored
// chunk changes (a new chunk packet or a block update). Consumers compare it
// against the revision they last read to detect stale cached data.
func (w *World) ChunkRevision(cx, cz int32) (revision int64) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.revisions[ChunkPos{cx, cz}]
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
	if err = chunk.Sections[sectionIndex].Blocks.Set(index, stateId); err != nil {
		return
	}
	w.revisions[ChunkPos{chunkX, chunkZ}]++
	return nil
}

func (w *World) storeChunk(pos ChunkPos, chunk *level.Chunk) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.chunks[pos] = chunk
	w.revisions[pos]++
}

// updateLight applies a standalone light-update packet to a loaded chunk. Unlike
// a full chunk packet, a light update only carries the sections whose light
// changed, so sections without a bit set in either mask keep their current
// layer.
func (w *World) updateLight(pos ChunkPos, skyMask, blockMask, emptySkyMask, emptyBlockMask []byte, skyLight, blockLight [][]byte) {
	w.mu.Lock()
	defer w.mu.Unlock()
	chunk, ok := w.chunks[pos]
	if !ok {
		return
	}
	lightSections := len(chunk.Sections) + 2
	if len(chunk.Light) != lightSections {
		chunk.Light = make([]level.SectionLight, lightSections)
		for i := range chunk.Light {
			chunk.Light[i].Sky = level.DataLayer{DefaultValue: 15}
		}
	}
	skyIndex, blockIndex := 0, 0
	for i := range lightSections {
		if bitSet(skyMask, i) && skyIndex < len(skyLight) {
			chunk.Light[i].Sky = level.DataLayer{Data: skyLight[skyIndex]}
			skyIndex++
		} else if bitSet(emptySkyMask, i) {
			chunk.Light[i].Sky = level.DataLayer{}
		}
		if bitSet(blockMask, i) && blockIndex < len(blockLight) {
			chunk.Light[i].Block = level.DataLayer{Data: blockLight[blockIndex]}
			blockIndex++
		} else if bitSet(emptyBlockMask, i) {
			chunk.Light[i].Block = level.DataLayer{}
		}
	}
}

func (w *World) dropChunk(pos ChunkPos) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.chunks, pos)
	delete(w.revisions, pos)
}

func (w *World) setCenter(pos ChunkPos) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.center = pos
}

// Module subscribes to chunk and block-change packets and updates a shared
// World. Multiple Modules can point to the same World.
type Module struct {
	*proto.Conn
	world        *World
	desiredBatch float32
}

// NewModule creates a Module around an existing World without registering handlers.
func NewModule(c *proto.Conn, w *World) (m *Module) {
	return &Module{Conn: c, world: w, desiredBatch: 8}
}

// World returns the shared World backing this module.
func (m *Module) World() *World { return m.world }

// SetDesiredChunksPerTick configures the batch size sent in ChunkBatchReceived.
func (m *Module) SetDesiredChunksPerTick(n float32) { m.desiredBatch = n }

// GetBlock returns the block state ID at the given world coordinates via the
// shared World.
func (m *Module) GetBlock(x, y, z int32) (int32, error) {
	return m.world.GetBlock(x, y, z)
}

// Start registers packet handlers on c that update the shared World w.
func Start(c *proto.Conn, w *World) *Module {
	m := NewModule(c, w)
	m.RegisterUntil("26.2", m.onMapChunk)
	m.RegisterUntil("26.4-snapshot-2", m.onMapChunk26_4)
	m.RegisterUntil("26.2", m.onUpdateLight)
	m.RegisterUntil("26.4-snapshot-2", m.onUpdateLight26_4)
	m.RegisterUntilLatest(m.onBlockChange)
	m.RegisterUntilLatest(m.onMultiBlockChange)
	m.RegisterUntilLatest(m.onUnloadChunk)
	m.RegisterUntilLatest(m.onUpdateViewPosition)
	m.RegisterUntilLatest(m.onChunkBatchFinished)
	m.RegisterUntilLatest(m.onChunkBatchStart)
	m.RegisterUntil("26.4-snapshot-2", m.onGameStateChange)
	return m
}

func (m *Module) onMapChunk(p *v1_21_11.PlayToClientPacketMapChunk) (err error) {
	return m.applyChunkPacket(chunkPacketFromV1_21_11(p))
}

func (m *Module) onMapChunk26_4(p *v26_4_snapshot_2.PlayToClientPacketMapChunk) (err error) {
	return m.applyChunkPacket(chunkPacketFrom26_4(p))
}

func (m *Module) onUpdateLight(p *v1_21_11.PlayToClientPacketUpdateLight) (err error) {
	return m.applyLightUpdate(
		p.ChunkX, p.ChunkZ,
		longArrayToBitSet(p.SkyLightMask), longArrayToBitSet(p.BlockLightMask),
		longArrayToBitSet(p.EmptySkyLightMask), longArrayToBitSet(p.EmptyBlockLightMask),
		p.SkyLight, p.BlockLight,
	)
}

func (m *Module) onUpdateLight26_4(p *v26_4_snapshot_2.PlayToClientPacketUpdateLight) (err error) {
	return m.applyLightUpdate(
		p.ChunkX, p.ChunkZ,
		p.SkyLightMask.Val, p.BlockLightMask.Val,
		p.EmptySkyLightMask.Val, p.EmptyBlockLightMask.Val,
		byteArrays(p.SkyLight), byteArrays(p.BlockLight),
	)
}

// byteArrays unwraps a list of generated ByteArray into raw byte slices.
func byteArrays(in []v26_4_snapshot_2.ByteArray) (ret [][]byte) {
	for _, b := range in {
		ret = append(ret, b.Val)
	}
	return
}

// chunkPacket is the version-independent subset of MapChunk needed by the world store.
type chunkPacket struct {
	X                   int32
	Z                   int32
	ChunkData           []byte
	Heightmaps          []level.Heightmap
	BlockEntities       []level.BlockEntity
	SkyLightMask        []byte
	BlockLightMask      []byte
	EmptySkyLightMask   []byte
	EmptyBlockLightMask []byte
	SkyLight            [][]byte
	BlockLight          [][]byte
}

func (m *Module) applyChunkPacket(p *chunkPacket) (err error) {
	chunk, err := level.ReadChunkFromChunkPacketData(
		bytes.NewReader(p.ChunkData),
		m.world.version,
		m.world.height,
	)
	if err != nil {
		return
	}
	chunk.Heightmaps = p.Heightmaps
	chunk.BlockEntities = p.BlockEntities
	chunk.Light = buildSectionLight(
		p.SkyLightMask, p.BlockLightMask,
		p.EmptySkyLightMask, p.EmptyBlockLightMask,
		p.SkyLight, p.BlockLight, len(chunk.Sections),
	)
	m.world.storeChunk(ChunkPos{p.X, p.Z}, chunk)
	return nil
}

// applyLightUpdate applies a standalone light-update packet to a loaded chunk.
func (m *Module) applyLightUpdate(x, z int32, skyMask, blockMask, emptySkyMask, emptyBlockMask []byte, skyLight, blockLight [][]byte) (err error) {
	m.world.updateLight(ChunkPos{x, z}, skyMask, blockMask, emptySkyMask, emptyBlockMask, skyLight, blockLight)
	return nil
}

// buildSectionLight expands the light bitmasks and arrays into per-light-section
// layers, faithfully following ClientboundLightUpdatePacketData semantics.
//
// The masks are bitsets over the sectionsCount+2 light sections that vanilla's
// LevelLightEngine spans (minSectionY-1 through maxSectionY+1). Arrays appear in
// increasing light-section order for the set bits of the corresponding mask. A
// section whose bit is only set in an Empty*LightMask is a homogeneous all-zero
// layer. A sky section with neither bit set is above the column's data (open
// sky), which the vanilla sky light engine reports as 15; a block section with
// neither bit set is unlit (0).
func buildSectionLight(
	skyMask, blockMask []byte,
	emptySkyMask, emptyBlockMask []byte,
	skyLight, blockLight [][]byte,
	sections int,
) (ret []level.SectionLight) {
	lightSections := sections + 2
	ret = make([]level.SectionLight, lightSections)
	skyIndex, blockIndex := 0, 0
	for i := range lightSections {
		switch {
		case bitSet(skyMask, i) && skyIndex < len(skyLight):
			ret[i].Sky = level.DataLayer{Data: skyLight[skyIndex]}
			skyIndex++
		case bitSet(emptySkyMask, i):
			ret[i].Sky = level.DataLayer{}
		default:
			ret[i].Sky = level.DataLayer{DefaultValue: 15}
		}
		switch {
		case bitSet(blockMask, i) && blockIndex < len(blockLight):
			ret[i].Block = level.DataLayer{Data: blockLight[blockIndex]}
			blockIndex++
		default:
			ret[i].Block = level.DataLayer{}
		}
	}
	return
}

// bitSet reports whether bit i is set in a little-endian packed bit set.
func bitSet(mask []byte, i int) bool {
	if i < 0 || i/8 >= len(mask) {
		return false
	}
	return mask[i/8]>>(uint(i)%8)&1 == 1
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

func (m *Module) onUnloadChunk(p *v1_21_11.PlayToClientPacketUnloadChunk) (err error) {
	m.world.dropChunk(ChunkPos{p.ChunkX, p.ChunkZ})
	return nil
}

func (m *Module) onUpdateViewPosition(p *v1_21_11.PlayToClientPacketUpdateViewPosition) (err error) {
	m.world.setCenter(ChunkPos{p.ChunkX, p.ChunkZ})
	return nil
}

func (m *Module) onChunkBatchFinished(p *v1_21_11.PlayToClientPacketChunkBatchFinished) (err error) {
	return m.Send(&v26_4_snapshot_2.PlayToServerPacketChunkBatchReceived{ChunksPerTick: m.desiredBatch})
}

func (m *Module) onChunkBatchStart(_ *v1_21_11.PlayToClientPacketChunkBatchStart) (err error) {
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
