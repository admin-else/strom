// Command weather connects to a server, drives /weather and teleports, and
// prints the world weather, biome, fluid and heightmap state the clouds and
// weather renderers need. It asserts each query against the server's own data.
package weather

import (
	"flag"
	"fmt"
	"math"
	"time"

	"github.com/admin-else/strom/cmd/strom/cmd_util"
	"github.com/admin-else/strom/mc/bot/keepalive"
	"github.com/admin-else/strom/mc/bot/player"
	"github.com/admin-else/strom/mc/bot/world"
	"github.com/admin-else/strom/mc/client"
	"github.com/admin-else/strom/mc/data"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_generated/v26_4_snapshot_2"
	"github.com/admin-else/strom/mc/registry"
	"github.com/admin-else/strom/mc/text"
)

var (
	cmd         = flag.NewFlagSet("weather", flag.ContinueOnError)
	addrFlag    = cmd.String("addr", "127.0.0.1:25565", "address to connect to")
	accFlag     = cmd.String("acc", "dev", "account name (offline mode, must be op)")
	versionFlag = cmd.String("version", "26.4-snapshot-2", "version to use")
	noAssert    = cmd.Bool("no-assert", false, "skip the built-in assertion checks")
)

// biomeSearchCoords are tried in order until the biome at the column differs
// from the spawn biome, so the probe proves BiomeAt tracks a real change.
var biomeSearchCoords = [][2]float64{
	{2000, 2000},
	{-3000, 4000},
	{6000, -6000},
	{-8000, -8000},
	{10000, 10000},
}

type probe struct {
	*proto.Conn
	world  *world.Module
	player *player.Module
	fail   int
}

func (p *probe) failed(format string, args ...any) {
	p.fail++
	fmt.Printf("FAIL: "+format+"\n", args...)
}

func (p *probe) check(ok bool, format string, args ...any) {
	if ok {
		fmt.Printf("ok: "+format+"\n", args...)
		return
	}
	p.failed(format, args...)
}

// OnSystemChat logs server chat and command feedback.
func (p *probe) OnSystemChat(m *v26_4_snapshot_2.PlayToClientPacketSystemChat) (err error) {
	if m.IsActionBar {
		return nil
	}
	var raw text.RawComponent
	if err = raw.FromNBT(m.Content.Value); err != nil {
		return nil
	}
	fmt.Printf("[chat] %s\n", text.ComponentFromRaw(&raw).GetString())
	return nil
}

func (p *probe) sendCommand(command string) {
	fmt.Printf("[cmd] /%s\n", command)
	if err := p.Send(&v26_4_snapshot_2.PlayToServerPacketChatCommand{Command: command}); err != nil {
		p.failed("send command %q: %v", command, err)
	}
}

func (p *probe) waitForChunk(cx, cz int32) bool {
	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := p.world.World().ChunkAt(cx, cz); err == nil {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}

func (p *probe) waitForChunks() bool {
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		if p.world.World().ChunkCount() > 0 && p.player.Joined() {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}

// logWeather prints the current weather state.
func (p *probe) logWeather(label string) {
	w := p.world.World()
	fmt.Printf("weather %-14s rain=%.3f (raining=%v) thunder=%.3f (thundering=%v)\n",
		label, w.RainLevel(), w.IsRaining(), w.ThunderLevel(), w.IsThundering())
}

// runWeather drives /weather rain then /weather clear and watches the level
// ramp via the GameStateChange level-change events.
func (p *probe) runWeather() {
	w := p.world.World()
	p.logWeather("before")
	before := w.RainLevel()

	p.sendCommand("weather rain")
	maxRain := before
	last := before
	deadline := time.Now().Add(8 * time.Second)
	for time.Now().Before(deadline) {
		if now := w.RainLevel(); math.Abs(float64(now-last)) > 0.05 {
			p.logWeather("raining")
			last = now
		}
		if now := w.RainLevel(); now > maxRain {
			maxRain = now
		}
		time.Sleep(100 * time.Millisecond)
	}
	p.logWeather("during")
	p.check(maxRain > 0.5, "rain level ramped up (max %.3f)", maxRain)
	p.check(w.IsRaining(), "IsRaining after /weather rain")

	p.sendCommand("weather thunder")
	maxThunder := w.ThunderLevel()
	last = w.ThunderLevel()
	deadline = time.Now().Add(8 * time.Second)
	for time.Now().Before(deadline) {
		if now := w.ThunderLevel(); math.Abs(float64(now-last)) > 0.05 {
			p.logWeather("thundering")
			last = now
		}
		if now := w.ThunderLevel(); now > maxThunder {
			maxThunder = now
		}
		time.Sleep(100 * time.Millisecond)
	}
	p.check(maxThunder > 0.5, "thunder level ramped up (max %.3f)", maxThunder)
	p.check(w.IsThundering(), "IsThundering after /weather thunder")

	p.sendCommand("weather clear")
	minRain := w.RainLevel()
	last = w.RainLevel()
	deadline = time.Now().Add(8 * time.Second)
	for time.Now().Before(deadline) {
		if now := w.RainLevel(); math.Abs(float64(now-last)) > 0.05 {
			p.logWeather("clearing")
			last = now
		}
		if now := w.RainLevel(); now < minRain {
			minRain = now
		}
		time.Sleep(100 * time.Millisecond)
	}
	p.logWeather("after")
	p.check(minRain < 0.5, "rain level ramped down (min %.3f)", minRain)
	p.check(!w.IsRaining(), "not raining after /weather clear")
	p.check(!w.IsThundering(), "not thundering after /weather clear")
}

// biomeAtPlayer queries the biome at the player's feet.
func (p *probe) biomeAtPlayer() (world.BiomeInfo, error) {
	x, y, z := p.player.Position()
	return p.world.World().BiomeAt(int32(math.Floor(x)), int32(math.Floor(y)), int32(math.Floor(z)))
}

// runBiome teleports across the overworld until the biome changes and asserts
// every biome resolves through the registry and minecraft-data.
func (p *probe) runBiome() {
	start, err := p.biomeAtPlayer()
	if err != nil {
		p.failed("spawn biome: %v", err)
		return
	}
	fmt.Printf("biome spawn: %s (id=%d precip=%v temp=%.2f)\n", start.ID, start.NumericID, start.HasPrecipitation, start.Temperature)
	p.checkBiomeData("spawn", start)

	var changed *world.BiomeInfo
	for _, c := range biomeSearchCoords {
		p.sendCommand(fmt.Sprintf("execute in minecraft:overworld run tp @s %.1f 100 %.1f", c[0], c[1]))
		p.waitForChunk(int32(math.Floor(c[0]))>>4, int32(math.Floor(c[1]))>>4)
		time.Sleep(500 * time.Millisecond)
		got, err := p.biomeAtPlayer()
		if err != nil {
			p.failed("biome at %.0f,%.0f: %v", c[0], c[1], err)
			continue
		}
		fmt.Printf("biome @%.0f,%.0f: %s (id=%d precip=%v temp=%.2f)\n", c[0], c[1], got.ID, got.NumericID, got.HasPrecipitation, got.Temperature)
		p.checkBiomeData("tp", got)
		if got.ID != start.ID {
			changed = &got
			break
		}
	}
	p.check(changed != nil, "biome changed after teleport (spawn %s)", start.ID)
}

// checkBiomeData asserts the biome resource id is in the registry and the
// precipitation/temperature equal the minecraft-data definition.
func (p *probe) checkBiomeData(label string, info world.BiomeInfo) {
	p.check(info.ID != "", "%s biome id resolved", label)
	definition, ok := data.LookupBiomeByName(p.Version, trimNamespace(info.ID))
	p.check(ok, "%s biome %s present in minecraft-data", label, info.ID)
	if !ok {
		return
	}
	p.check(info.HasPrecipitation == definition.HasPrecipitation,
		"%s %s has_precipitation=%v", label, info.ID, info.HasPrecipitation)
	p.check(info.Temperature == definition.Temperature,
		"%s %s temperature=%.2f", label, info.ID, info.Temperature)
}

func trimNamespace(id string) (ret string) {
	if len(id) > 10 && id[:10] == "minecraft:" {
		return id[10:]
	}
	return id
}

// motionBlockingScan recomputes the MOTION_BLOCKING height from the loaded
// blocks: the first Y above the topmost motion-blocking (or fluid) block.
func (p *probe) motionBlockingScan(x, z int32) (h int32, err error) {
	w := p.world.World()
	for y := int32(w.MinY() + w.Height() - 1); y >= int32(w.MinY()); y-- {
		stateID, stateErr := w.GetBlock(x, y, z)
		if stateErr != nil {
			return 0, stateErr
		}
		blocksMotion := false
		if boxes, ok := data.CollisionShapeBoxesAt(p.Version, stateID); ok && len(boxes) > 0 {
			blocksMotion = true
		}
		if !blocksMotion {
			if fluid, fluidErr := w.FluidAt(x, y, z); fluidErr == nil && fluid.Kind != world.FluidEmpty {
				blocksMotion = true
			}
		}
		if blocksMotion {
			return y + 1, nil
		}
	}
	return int32(w.MinY()), nil
}

// runHeightmap asserts MotionBlockingHeightAt matches a scan of the loaded
// blocks.
func (p *probe) runHeightmap() {
	w := p.world.World()
	x, _, z := p.player.Position()
	bx, bz := int32(math.Floor(x)), int32(math.Floor(z))
	got, err := w.MotionBlockingHeightAt(bx, bz)
	if err != nil {
		p.failed("MotionBlockingHeightAt(%d,%d): %v", bx, bz, err)
		return
	}
	want, err := p.motionBlockingScan(bx, bz)
	if err != nil {
		p.failed("heightmap scan: %v", err)
		return
	}
	fmt.Printf("heightmap motion_blocking @%d,%d = %d (scan %d)\n", bx, bz, got, want)
	p.check(got == want, "MotionBlockingHeightAt matches block scan (%d == %d)", got, want)
}

// runFluid places a water column at the player and teleports inside it, then
// asserts FluidAt reports water at the eye and empty in the air above.
func (p *probe) runFluid() {
	x, y, z := p.player.Position()
	bx, by, bz := int32(math.Floor(x)), int32(math.Floor(y)), int32(math.Floor(z))
	if !p.waitForChunk(bx>>4, bz>>4) {
		p.failed("fluid: chunk %d,%d not loaded", bx>>4, bz>>4)
		return
	}
	p.sendCommand(fmt.Sprintf("fill %d %d %d %d %d %d water", bx, by, bz, bx, by+3, bz))
	p.sendCommand(fmt.Sprintf("tp @s %d.5 %d %d.5", bx, by, bz))

	w := p.world.World()
	var eye world.FluidState
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		var err error
		eye, err = w.FluidAt(bx, by+1, bz)
		if err == nil && eye.Kind == world.FluidWater {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	air, err := w.FluidAt(bx, by+6, bz)
	if err != nil {
		p.failed("FluidAt air: %v", err)
		return
	}
	fmt.Printf("fluid @%d,%d,%d eye=%s (%s) air=%s (%s)\n", bx, by, bz, eye.Kind, eye.ID, air.Kind, air.ID)
	p.check(eye.Kind == world.FluidWater && eye.ID == "minecraft:water", "eye fluid is water")
	p.check(air.Kind == world.FluidEmpty, "open air fluid is empty")
}

func Run(args []string) (err error) {
	if err = cmd.Parse(args); err != nil {
		return
	}
	acc, err := cmd_util.Account(*accFlag)
	if err != nil {
		return
	}

	store := registry.NewStore()
	c, err := client.Login(*addrFlag, acc,
		client.WithVersion(*versionFlag),
		client.WithRegistries(store),
	)
	if err != nil {
		return
	}
	defer c.Close()

	w := world.NewWorld(c.Version, -64, 384)
	w.SetRegistries(store)
	m := world.Start(c, w)
	p := &probe{
		Conn:   c,
		world:  m,
		player: player.Start(c),
	}
	keepalive.Start(c)
	p.RegisterUntilLatest(p.OnSystemChat)

	done := make(chan error, 1)
	go func() { done <- c.StartConn() }()

	if !p.waitForChunks() {
		return fmt.Errorf("timed out waiting for chunks (count=%d joined=%v)", w.ChunkCount(), p.player.Joined())
	}
	time.Sleep(time.Second)

	p.runWeather()
	p.runHeightmap()
	p.runBiome()
	p.runFluid()

	if *noAssert {
		fmt.Println("\nassertions skipped (--no-assert)")
	} else if p.fail > 0 {
		return fmt.Errorf("%d weather probe assertions failed", p.fail)
	} else {
		fmt.Println("\nall weather probe assertions passed")
	}
	return nil
}
