package messenger

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/admin-else/strom/cmd/strom/cmd_util"
	"github.com/admin-else/strom/mc/bot/chat"
	"github.com/admin-else/strom/mc/bot/keepalive"
	"github.com/admin-else/strom/mc/bot/player"
	"github.com/admin-else/strom/mc/bot/world"
	"github.com/admin-else/strom/mc/client"
	"github.com/admin-else/strom/mc/data"
	"github.com/admin-else/strom/mc/event"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_11"
	"github.com/admin-else/strom/mc/proto_generated/v1_8"
)

func stror(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

var (
	cmd           = flag.NewFlagSet("status", flag.ContinueOnError)
	connectToFlag = cmd.String("addr", "", "server address to connect to")
	accFlag       = cmd.String("acc", stror(os.Getenv("STROM_ACC"), "strom_messenger"), "the name to use if longer will be treated as ygg token")
	versionFlag   = cmd.String("version", "26.2", "minecraft version to use (e.g. 1.8, 1.21.11)")
)

type Messenger struct {
	*proto.Conn
	chat   *chat.Module
	world  *world.Module
	player *player.Module
	legacy bool
}

func (m *Messenger) OnChat(e *v1_21_11.PlayToClientPacketPlayerChat) (err error) {
	m.Log.Info("chat", "author", e.SenderUuid, "m", e.PlainMessage)
	return
}

func (m *Messenger) OnChatUnsigned(e *v1_21_11.PlayToClientPacketProfilelessChat) (err error) {
	m.Log.Info("chatP", "author", e.Name, "m", e.Message)
	return
}

func (m *Messenger) OnChatSystem(e *v1_21_11.PlayToClientPacketSystemChat) (err error) {
	m.Log.Info("chatS", "m", e.Content)
	return
}

func (m *Messenger) OnChatLegacy(e *v1_8.PlayToClientPacketChat) (err error) {
	m.Log.Info("chat", "m", e.Message)
	return
}

// OnUnhandled surfaces packets that reached a registered decoder but could not
// be decoded. Packets with no handler are silently dropped by the connection.
func (m *Messenger) OnUnhandled(e event.Unhandled) (err error) {
	u, ok := e.Val.(*proto.UnCodablePacket)
	if ok && !errors.Is(u.Err, proto.NoHandlerRegisteredErr) {
		m.Log.Debug("undecodable packet", "name", u.Info.Name, "err", u.Err.Error())
	}
	return
}

func (m *Messenger) OnStdin(e event.Stdin) (err error) {
	input := strings.TrimSpace(e.Val)
	if input == "" {
		return nil
	}
	if !strings.HasPrefix(input, ">") {
		return m.sendChat(input)
	}
	command, args, _ := strings.Cut(strings.TrimPrefix(input, ">"), " ")
	args = strings.TrimSpace(args)
	switch command {
	case "say":
		return m.sendChat(args)
	case "pos":
		return m.handlePos()
	case "look":
		return m.handleLook(args)
	case "tp":
		return m.handleTp(args)
	case "fly":
		return m.handleFly(args)
	case "getblock":
		return m.handleGetBlock(args)
	case "light":
		return m.handleLight(args)
	case "height":
		return m.handleHeight(args)
	case "chunk":
		return m.handleChunk(args)
	case "chunks":
		return m.handleChunks()
	default:
		m.Log.Info("unknown command", "command", command)
		return nil
	}
}

func (m *Messenger) sendChat(message string) (err error) {
	if m.legacy {
		return m.Send(&v1_8.PlayToServerPacketChat{Message: message})
	}
	return m.chat.SendMessage(message)
}

func (m *Messenger) handlePos() (err error) {
	if m.player == nil {
		m.Log.Info("pos", "known", false)
		return
	}
	x, y, z := m.player.Position()
	yaw, pitch := m.player.Rotation()
	m.Log.Info("pos",
		"x", x, "y", y, "z", z,
		"yaw", yaw, "pitch", pitch,
		"gamemode", m.player.Gamemode(),
		"flying", m.player.Flying(),
		"center", m.world.World().Center(),
	)
	return
}

func parseFloats(fields []string, n int) (values []float64, err error) {
	if len(fields) != n {
		return nil, fmt.Errorf("expected %d numbers", n)
	}
	values = make([]float64, n)
	for i, f := range fields {
		values[i], err = strconv.ParseFloat(f, 64)
		if err != nil {
			return nil, fmt.Errorf("bad number %q", f)
		}
	}
	return values, nil
}

func (m *Messenger) handleLook(args string) (err error) {
	if m.player == nil {
		m.Log.Info("look", "error", "player unavailable for this version")
		return nil
	}
	values, err := parseFloats(strings.Fields(args), 2)
	if err != nil {
		m.Log.Info("look", "error", err.Error())
		return nil
	}
	m.player.Look(float32(values[0]), float32(values[1]))
	m.Log.Info("look", "yaw", values[0], "pitch", values[1])
	return nil
}

func (m *Messenger) handleTp(args string) (err error) {
	if m.player == nil {
		m.Log.Info("tp", "error", "player unavailable for this version")
		return nil
	}
	values, err := parseFloats(strings.Fields(args), 3)
	if err != nil {
		m.Log.Info("tp", "error", err.Error())
		return nil
	}
	m.player.SetPosition(values[0], values[1], values[2])
	m.Log.Info("tp", "x", values[0], "y", values[1], "z", values[2])
	return nil
}

func (m *Messenger) handleFly(args string) (err error) {
	if m.player == nil {
		m.Log.Info("fly", "error", "player unavailable for this version")
		return nil
	}
	values, err := parseFloats(strings.Fields(args), 3)
	if err != nil {
		m.Log.Info("fly", "error", err.Error())
		return nil
	}
	m.player.Fly(values[0], values[1], values[2])
	x, y, z := m.player.Position()
	m.Log.Info("fly", "dx", values[0], "dy", values[1], "dz", values[2], "x", x, "y", y, "z", z)
	return nil
}

func parseCoords(fields []string, n int) (coords []int32, err error) {
	if len(fields) != n {
		return nil, fmt.Errorf("expected %d coordinates", n)
	}
	coords = make([]int32, n)
	for i, f := range fields {
		var v int64
		v, err = strconv.ParseInt(f, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("bad coordinate %q", f)
		}
		coords[i] = int32(v)
	}
	return coords, nil
}

func (m *Messenger) handleGetBlock(args string) (err error) {
	coords, err := parseCoords(strings.Fields(args), 3)
	if err != nil {
		m.Log.Info("getblock", "error", err.Error())
		return nil
	}
	stateId, err := m.world.GetBlock(coords[0], coords[1], coords[2])
	if err != nil {
		m.Log.Info("getblock", "error", err.Error())
		return nil
	}
	block, props, err := data.FromBlockState(m.Version, stateId)
	if err != nil {
		m.Log.Info("getblock", "x", coords[0], "y", coords[1], "z", coords[2], "stateId", stateId, "error", err.Error())
		return nil
	}
	m.Log.Info("getblock", "x", coords[0], "y", coords[1], "z", coords[2], "name", block.Name, "stateId", stateId, "properties", props)
	return nil
}

func (m *Messenger) handleLight(args string) (err error) {
	coords, err := parseCoords(strings.Fields(args), 3)
	if err != nil {
		m.Log.Info("light", "error", err.Error())
		return nil
	}
	block, sky, err := m.world.World().LightAt(coords[0], coords[1], coords[2])
	if err != nil {
		m.Log.Info("light", "error", err.Error())
		return nil
	}
	m.Log.Info("light", "x", coords[0], "y", coords[1], "z", coords[2], "block", block, "sky", sky)
	return nil
}

func (m *Messenger) handleHeight(args string) (err error) {
	coords, err := parseCoords(strings.Fields(args), 2)
	if err != nil {
		m.Log.Info("height", "error", err.Error())
		return nil
	}
	kind := "world_surface"
	h, err := m.world.World().HeightAt(kind, coords[0], coords[1])
	if err != nil {
		m.Log.Info("height", "error", err.Error())
		return nil
	}
	m.Log.Info("height", "kind", kind, "x", coords[0], "z", coords[1], "value", h)
	return nil
}

func (m *Messenger) handleChunk(args string) (err error) {
	coords, err := parseCoords(strings.Fields(args), 2)
	if err != nil {
		m.Log.Info("chunk", "error", err.Error())
		return nil
	}
	chunk, err := m.world.World().ChunkAt(coords[0], coords[1])
	if err != nil {
		m.Log.Info("chunk", "error", err.Error())
		return nil
	}
	sectionsWithLight := 0
	for _, s := range chunk.Light {
		if s.Sky.Data != nil || s.Block.Data != nil {
			sectionsWithLight++
		}
	}
	m.Log.Info("chunk", "cx", coords[0], "cz", coords[1], "sections", len(chunk.Sections), "lightSections", sectionsWithLight, "heightmaps", len(chunk.Heightmaps), "blockEntities", len(chunk.BlockEntities))
	return nil
}

func (m *Messenger) handleChunks() (err error) {
	w := m.world.World()
	m.Log.Info("chunks", "count", w.ChunkCount(), "center", w.Center(), "minY", w.MinY(), "height", w.Height())
	return nil
}

func Run(args []string) (err error) {
	err = cmd.Parse(args)
	if err != nil {
		return
	}

	acc, err := cmd_util.Account(*accFlag)
	if err != nil {
		return
	}
	c, err := client.Login(*connectToFlag, acc, client.WithVersion(*versionFlag))
	if err != nil {
		return
	}
	defer c.Close()

	isLegacy := c.ProtocolVersion < 764

	var chatMod *chat.Module
	if !isLegacy {
		chatMod, err = chat.Start(c, acc)
		if err != nil {
			return
		}
	}

	m := &Messenger{
		Conn:   c,
		chat:   chatMod,
		legacy: isLegacy,
	}

	if debugPackets := os.Getenv("STROM_DEBUG_PACKETS"); debugPackets != "" {
		logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
		slog.SetDefault(logger)
		c.Log = logger
		m.DebugPrintPackets = strings.Split(debugPackets, ",")
	}

	if !isLegacy {
		w := world.NewWorld(c.Version, -64, 384)
		m.world = world.Start(c, w)
		m.player = player.Start(c)
	}

	event.StartListingStdin(m.Loop)
	m.Register(m.OnStdin)
	m.Register(m.OnUnhandled)

	if isLegacy {
		m.Register(m.OnChatLegacy)
	} else {
		m.RegisterUntilLatest(m.OnChat, m.OnChatUnsigned, m.OnChatSystem)
	}

	if !isLegacy {
		keepalive.Start(m.Conn)
	}

	return m.StartConn()
}
