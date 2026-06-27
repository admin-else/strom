package mca_to_tmcpr

import (
	"bytes"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/admin-else/strom/level"
	"github.com/admin-else/strom/level/anvil"
	"github.com/admin-else/strom/nbt"
	"github.com/admin-else/strom/proto/replay"
	"github.com/admin-else/strom/proto_generated/v1_21_11"
	"github.com/admin-else/strom/util"
)

var (
	cmd       = flag.NewFlagSet("mca-to-tmcpr", flag.ContinueOnError)
	worldPath = cmd.String("world", "", "Path to world directory")
	output    = cmd.String("output", "", "Output .tmcpr file path")
	radius    = cmd.Int("radius", 14, "Chunk radius around spawn")
	spawnX    = cmd.Float64("spawn-x", 0.5, "Spawn X coordinate")
	spawnY    = cmd.Float64("spawn-y", 100, "Spawn Y coordinate")
	spawnZ    = cmd.Float64("spawn-z", 0.5, "Spawn Z coordinate")
)

var FullLight [][]uint8

func init() {
	light := make([]uint8, 2048)
	for i := range light {
		light[i] = 0xFF
	}
	sectionsPerChunk := 384 / 16
	FullLight = make([][]uint8, sectionsPerChunk)
	for i := range sectionsPerChunk {
		FullLight[i] = light
	}
}

func GetChunkFromMca(x, z int32, version string) (b []byte, err error) {
	xmca := util.FloorDiv(x, anvil.ChunksWidthRegionFile)
	zmca := util.FloorDiv(z, anvil.ChunksWidthRegionFile)
	x = x % anvil.ChunksWidthRegionFile
	if x < 0 {
		x += anvil.ChunksWidthRegionFile
	}
	z = z % anvil.ChunksWidthRegionFile
	if z < 0 {
		z += anvil.ChunksWidthRegionFile
	}

	f, err := os.Open(fmt.Sprintf("./level/anvil/testdata/securechattestworld/region/r.%v.%v.mca", xmca, zmca))
	if err != nil {
		return
	}
	defer f.Close()
	mca, err := anvil.ReadChunkHeader(f)
	if err != nil {
		return
	}

	var n *nbt.Tag
	n, err = mca.ChunkAt(f, x, z)
	if err != nil {
		return
	}
	if n == nil {
		return
	}

	var chunk anvil.Chunk
	err = nbt.Format.Decode(n.Value, &chunk)
	if err != nil {
		return
	}

	var packetChunk *level.Chunk
	packetChunk, err = chunk.ToStorage(version)
	if err != nil {
		return
	}
	buffer := bytes.NewBuffer(nil)
	err = packetChunk.WriteChunkData(buffer)
	if err != nil {
		return
	}
	b = buffer.Bytes()
	return
}

func Run(args []string) error {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	err := cmd.Parse(args)
	if err != nil {
		return err
	}
	if *worldPath == "" {
		slog.Error("no world path specified")
		return nil
	}
	if *output == "" {
		slog.Error("no output path specified")
		return nil
	}

	outFile, err := os.Create(*output)
	if err != nil {
		slog.Error("failed to create output file", "error", err)
		return err
	}
	defer outFile.Close()

	w := replay.NewWriter(outFile)

	err = w.WritePacket(&v1_21_11.PlayToClientPacketLogin{
		EntityId:            73,
		IsHardcore:          false,
		WorldNames:          []string{"minecraft:overworld", "minecraft:the_end", "minecraft:the_nether"},
		MaxPlayers:          100,
		ViewDistance:        10,
		SimulationDistance:  10,
		ReducedDebugInfo:    false,
		EnableRespawnScreen: true,
		DoLimitedCrafting:   false,
		WorldState: v1_21_11.PlayToClientSpawnInfo{
			Dimension:        0,
			Name:             "minecraft:overworld",
			HashedSeed:       -8365436745629239539,
			Gamemode:         "survival",
			PreviousGamemode: 0xff,
			IsDebug:          false,
			IsFlat:           false,
			Death:            nil,
			PortalCooldown:   0,
			SeaLevel:         63,
		},
		EnforcesSecureChat: false,
	}, 0)
	if err != nil {
		slog.Error("failed to write login packet", "error", err)
		return err
	}

	err = w.WritePacket(&v1_21_11.PlayToClientPacketPosition{
		TeleportId: 0,
		X:          *spawnX,
		Y:          *spawnY,
		Z:          *spawnZ,
		Dx:         0,
		Dy:         0,
		Dz:         0,
		Yaw:        0,
		Pitch:      0,
		Flags:      v1_21_11.PlayToClientPositionUpdateRelatives{},
	}, 0)
	if err != nil {
		slog.Error("failed to write position packet", "error", err)
		return err
	}

	chunkX := int32(*spawnX) >> 4
	chunkZ := int32(*spawnZ) >> 4

	err = w.WritePacket(&v1_21_11.PlayToClientPacketUpdateViewPosition{
		ChunkX: chunkX,
		ChunkZ: chunkZ,
	}, 0)
	if err != nil {
		slog.Error("failed to write view position packet", "error", err)
		return err
	}

	err = w.WritePacket(&v1_21_11.PlayToClientPacketGameStateChange{
		Reason: "level_chunks_load_start",
	}, 0)
	if err != nil {
		slog.Error("failed to write game state packet", "error", err)
		return err
	}

	radius32 := int32(*radius)

	for x := -radius32; x < radius32; x++ {
		for z := -radius32; z < radius32; z++ {
			var chunkData []byte
			chunkData, err = GetChunkFromMca(x, z, "1.21.11")
			if err != nil {
				return err
			}
			if chunkData == nil {
				continue
			}
			OLDpacket := &v1_21_11.PlayToClientPacketMapChunk{}
			OLDpacket.X = x
			OLDpacket.Z = z
			OLDpacket.ChunkData = v1_21_11.ByteArray{Val: chunkData}
			OLDpacket.BlockLight = FullLight
			OLDpacket.BlockLightMask = []int64{(1 << 24) - 1}
			err = w.WritePacket(OLDpacket, 0)
			if err != nil {
				return err
			}
		}
	}

	slog.Info("conversion complete", "output", *output)
	return nil
}
