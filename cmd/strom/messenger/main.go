package messenger

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/admin-else/strom/cmd/strom/cmd_util"
	"github.com/admin-else/strom/mc/bot/chat"
	"github.com/admin-else/strom/mc/bot/control"
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
	chat    *chat.Module
	world   *world.Module
	player  *player.Module
	control *control.Controller
	legacy  bool
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
	case "lookat", "lookhere":
		return m.handleLookAt(args)
	case "move", "press":
		return m.handleMove(args)
	case "jump":
		return m.handleJump(args)
	case "sneak":
		return m.handleSneak(args)
	case "sprint":
		return m.handleSprint(args)
	case "stop":
		return m.handleStop()
	case "tp":
		return m.handleTp(args)
	case "fly":
		return m.handleFly(args)
	case "viewdistance":
		return m.handleViewDistance(args)
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
	logArgs := []any{
		"x", x, "y", y, "z", z,
		"yaw", yaw, "pitch", pitch,
		"gamemode", m.player.Gamemode(),
		"flying", m.player.Flying(),
		"center", m.world.World().Center(),
	}
	if m.control != nil {
		lx, ly, lz := m.control.Position()
		cyaw, cpitch := m.control.Rotation()
		logArgs = append(logArgs, "localX", lx, "localY", ly, "localZ", lz, "localYaw", cyaw, "localPitch", cpitch, "onGround", m.control.Entity().OnGround())
	}
	m.Log.Info("pos", logArgs...)
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
	if m.control != nil {
		m.control.Look(float32(values[0]), float32(values[1]))
	}
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

// handleViewDistance queues a new render distance, which the player module sends
// in its next ServerboundClientInformationPacket.
func (m *Messenger) handleViewDistance(args string) (err error) {
	if m.player == nil {
		m.Log.Info("viewdistance", "error", "player unavailable for this version")
		return nil
	}
	fields := strings.Fields(args)
	if len(fields) != 1 {
		m.Log.Info("viewdistance", "error", "expected one integer")
		return nil
	}
	distance, err := strconv.Atoi(fields[0])
	if err != nil {
		m.Log.Info("viewdistance", "error", "bad number "+fields[0])
		return nil
	}
	m.player.SetViewDistance(distance)
	m.Log.Info("viewdistance", "distance", m.player.ViewDistance())
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
	return
}

// parseSeconds parses a seconds argument into a time.Duration.
func parseSeconds(field string) (d time.Duration, err error) {
	value, err := strconv.ParseFloat(field, 64)
	if err != nil {
		return 0, fmt.Errorf("bad seconds %q", field)
	}
	if value < 0 {
		value = 0
	}
	return time.Duration(value * float64(time.Second)), nil
}

func (m *Messenger) controlUnavailable(command string) (err error) {
	m.Log.Info(command, "error", "physics controller unavailable for this version")
	return nil
}

// handleLookAt implements "look here in X time": turn toward (x, y, z) after an
// optional delay in seconds.
func (m *Messenger) handleLookAt(args string) (err error) {
	if m.control == nil {
		return m.controlUnavailable("lookat")
	}
	fields := strings.Fields(args)
	if len(fields) != 3 && len(fields) != 4 {
		m.Log.Info("lookat", "error", "expected: lookat <x> <y> <z> [delaySeconds]")
		return nil
	}
	coords, err := parseFloats(fields[:3], 3)
	if err != nil {
		m.Log.Info("lookat", "error", err.Error())
		return nil
	}
	var delay time.Duration
	if len(fields) == 4 {
		delay, err = parseSeconds(fields[3])
		if err != nil {
			m.Log.Info("lookat", "error", err.Error())
			return nil
		}
	}
	m.control.LookAt(coords[0], coords[1], coords[2], delay)
	m.Log.Info("lookat", "x", coords[0], "y", coords[1], "z", coords[2], "delay", delay.Seconds())
	return nil
}

// handleMove implements "press wasd for X time": hold a set of movement keys for
// the given duration (keys is a subset of wasd, plus 'j' to jump).
func (m *Messenger) handleMove(args string) (err error) {
	if m.control == nil {
		return m.controlUnavailable("move")
	}
	fields := strings.Fields(args)
	if len(fields) != 2 {
		m.Log.Info("move", "error", "expected: move <wasdj> <seconds>")
		return nil
	}
	duration, err := parseSeconds(fields[1])
	if err != nil {
		m.Log.Info("move", "error", err.Error())
		return nil
	}
	unknown := m.control.Move(strings.ToLower(fields[0]), duration)
	m.Log.Info("move", "keys", strings.ToLower(fields[0]), "seconds", duration.Seconds(), "unknown", string(unknown))
	return nil
}

// handleJump implements "jump X seconds": hold the jump key for the duration
// (default 0.25s).
func (m *Messenger) handleJump(args string) (err error) {
	if m.control == nil {
		return m.controlUnavailable("jump")
	}
	fields := strings.Fields(args)
	duration := 250 * time.Millisecond
	if len(fields) == 1 {
		duration, err = parseSeconds(fields[0])
		if err != nil {
			m.Log.Info("jump", "error", err.Error())
			return nil
		}
	} else if len(fields) > 1 {
		m.Log.Info("jump", "error", "expected: jump [seconds]")
		return nil
	}
	m.control.Jump(duration)
	m.Log.Info("jump", "seconds", duration.Seconds())
	return nil
}

func parseToggle(args string) (value bool, ok bool) {
	switch strings.ToLower(strings.TrimSpace(args)) {
	case "on", "true", "1":
		return true, true
	case "off", "false", "0":
		return false, true
	default:
		return false, false
	}
}

func (m *Messenger) handleSneak(args string) (err error) {
	if m.control == nil {
		return m.controlUnavailable("sneak")
	}
	value, ok := parseToggle(args)
	if !ok {
		m.Log.Info("sneak", "error", "expected: sneak on|off")
		return nil
	}
	m.control.SetSneak(value)
	m.Log.Info("sneak", "value", value)
	return nil
}

func (m *Messenger) handleSprint(args string) (err error) {
	if m.control == nil {
		return m.controlUnavailable("sprint")
	}
	value, ok := parseToggle(args)
	if !ok {
		m.Log.Info("sprint", "error", "expected: sprint on|off")
		return nil
	}
	m.control.SetSprint(value)
	m.Log.Info("sprint", "value", value)
	return nil
}

func (m *Messenger) handleStop() (err error) {
	if m.control == nil {
		return m.controlUnavailable("stop")
	}
	m.control.Stop()
	m.Log.Info("stop", "ok", true)
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
		m.control = control.New(control.NewLevel(w, c.Version), m.player)
		ticker := &event.Timer{}
		ticker.Every(50*time.Millisecond, m.control.Tick)
		ticker.Start(m.Conn.Loop)
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
