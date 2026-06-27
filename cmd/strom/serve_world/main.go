package serve_world

import (
	"bytes"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"time"

	"github.com/admin-else/strom/mc/event"
	"github.com/admin-else/strom/mc/level"
	anvil2 "github.com/admin-else/strom/mc/level/anvil"
	nbt2 "github.com/admin-else/strom/mc/nbt"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_8"
	server2 "github.com/admin-else/strom/mc/server"
	"github.com/admin-else/strom/mc/util"
)

type Server struct {
	*proto.Conn
	*event.Timer

	reconnectCalled bool
}

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
	xmca := util.FloorDiv(x, anvil2.ChunksWidthRegionFile)
	zmca := util.FloorDiv(z, anvil2.ChunksWidthRegionFile)
	x = x % anvil2.ChunksWidthRegionFile
	if x < 0 {
		x += anvil2.ChunksWidthRegionFile
	}
	z = z % anvil2.ChunksWidthRegionFile
	if z < 0 {
		z += anvil2.ChunksWidthRegionFile
	}

	f, err := os.Open(fmt.Sprintf("./level/anvil/testdata/securechattestworld/region/r.%v.%v.mca", xmca, zmca))
	if err != nil {
		return
	}
	defer f.Close()
	mca, err := anvil2.ReadChunkHeader(f)
	if err != nil {
		return
	}

	var n *nbt2.Tag
	n, err = mca.ChunkAt(f, x, z)
	if err != nil {
		return
	}
	if n == nil {
		return
	}

	var chunk anvil2.Chunk
	err = nbt2.Format.Decode(n.Value, &chunk)
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

func Limbo(c *proto.Conn) (err error) {
	s := &Server{
		Conn: c,
	}
	_, err = server2.ServeLogin(s.Conn, server2.WithoutOnlineMode())
	if err != nil {
		return
	}

	err = server2.ServeConfig(s.Conn)
	if err != nil {
		return
	}
	login := &v1_21_8.PlayToClientPacketLogin{EntityId: 73, IsHardcore: false, WorldNames: []string{"minecraft:overworld", "minecraft:the_end", "minecraft:the_nether"}, MaxPlayers: 100, ViewDistance: 10, SimulationDistance: 10, ReducedDebugInfo: false, EnableRespawnScreen: true, DoLimitedCrafting: false, WorldState: v1_21_8.PlayToClientSpawnInfo{Dimension: 0, Name: "minecraft:overworld", HashedSeed: -8365436745629239539, Gamemode: "survival", PreviousGamemode: 0xff, IsDebug: false, IsFlat: false, Death: nil, PortalCooldown: 0, SeaLevel: 63}, EnforcesSecureChat: false}
	err = s.Send(login)
	if err != nil {
		return
	}
	err = s.Send(&v1_21_8.PlayToClientPacketPosition{
		TeleportId: 0,
		X:          0.5,
		Y:          100,
		Z:          0.5,
		Dx:         0,
		Dy:         0,
		Dz:         0,
		Yaw:        0,
		Pitch:      0,
		Flags:      v1_21_8.PlayToClientPositionUpdateRelatives{},
	})
	if err != nil {
		return
	}
	err = s.Send(&v1_21_8.PlayToClientPacketUpdateViewPosition{
		ChunkX: 0,
		ChunkZ: 0,
	})
	if err != nil {
		return
	}
	err = s.Send(&v1_21_8.PlayToClientPacketGameStateChange{
		Reason: "level_chunks_load_start",
	})

	radius := int32(10)

	for x := -radius; x < radius; x++ {
		for z := -radius; z < radius; z++ {
			var chunkData []byte
			chunkData, err = GetChunkFromMca(x, z, c.Version)
			if err != nil {
				return
			}
			if chunkData == nil {
				continue
			}
			OLDpacket := &v1_21_8.PlayToClientPacketMapChunk{}
			OLDpacket.X = x
			OLDpacket.Z = z
			OLDpacket.ChunkData = v1_21_8.ByteArray{Val: chunkData}
			OLDpacket.BlockLight = FullLight
			OLDpacket.BlockLightMask = []int64{(1 << 24) - 1}
			err = s.Send(OLDpacket)
			if err != nil {
				return
			}
		}
	}

	s.Timer = &event.Timer{}
	s.Timer.Every(time.Second*15, s.SendKeepAlive)
	s.Timer.Start(c.Loop)

	return s.StartConn()
}
func (s *Server) SendKeepAlive() (err error) {
	err = s.Send(&v1_21_8.PlayToClientPacketKeepAlive{KeepAliveId: rand.Int64()})
	if err != nil {
		return
	}
	return
}

func Run(args []string) error {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	return server2.StartServerWithOnConn(":25566", Limbo)
}
