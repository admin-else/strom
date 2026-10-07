// Command entities connects to a server, summons and equips an armor stand and
// a sheep, and prints the entity metadata and equipment the strom entity
// tracker exposes. It asserts the tracked equipment item ids match the items
// minecraft-data assigns, and the tracked metadata matches the summon NBT.
package entities

import (
	"flag"
	"fmt"
	"math"
	"time"

	"github.com/admin-else/strom/cmd/strom/cmd_util"
	"github.com/admin-else/strom/mc/bot/entity"
	"github.com/admin-else/strom/mc/bot/keepalive"
	"github.com/admin-else/strom/mc/bot/player"
	"github.com/admin-else/strom/mc/client"
	"github.com/admin-else/strom/mc/data"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_generated/v26_4_snapshot_2"
	"github.com/admin-else/strom/mc/text"
)

var (
	cmd         = flag.NewFlagSet("entities", flag.ContinueOnError)
	addrFlag    = cmd.String("addr", "127.0.0.1:25565", "address to connect to")
	accFlag     = cmd.String("acc", "dev", "account name (offline mode, must be op)")
	versionFlag = cmd.String("version", "26.4-snapshot-2", "version to use")
	noAssert    = cmd.Bool("no-assert", false, "skip the built-in assertion checks")
)

type probe struct {
	*proto.Conn
	entities *entity.Module
	player   *player.Module
	fail     int
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

// waitForEntity polls until a tracked entity is within tolerance of the
// summoned position, which identifies it without hardcoding entity type ids.
func (p *probe) waitForEntity(x, y, z, tolerance float64, timeout time.Duration) (found *entity.Entity) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		for _, e := range p.entities.Entities() {
			if math.Abs(e.Position.X-x) <= tolerance &&
				math.Abs(e.Position.Y-y) <= tolerance &&
				math.Abs(e.Position.Z-z) <= tolerance {
				return e
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	return nil
}

// runEquipment summons an armor stand, equips a diamond helmet and asserts the
// tracked head slot matches.
func (p *probe) runEquipment() {
	p.sendCommand("kill @e[type=!player]")
	time.Sleep(300 * time.Millisecond)

	x, y, z := p.player.Position()
	sx, sy, sz := x+2, y, z
	p.sendCommand(fmt.Sprintf("summon minecraft:armor_stand %.2f %.2f %.2f", sx, sy, sz))

	stand := p.waitForEntity(sx, sy, sz, 1.5, 5*time.Second)
	if stand == nil {
		p.failed("armor stand at %.1f,%.1f,%.1f was not tracked", sx, sy, sz)
		return
	}
	fmt.Printf("armor_stand tracked id=%d type=%d pos=(%.2f,%.2f,%.2f)\n",
		stand.Id, stand.Type, stand.Position.X, stand.Position.Y, stand.Position.Z)

	helmet, ok := data.LookupItemByName(p.Version, "diamond_helmet")
	if !ok {
		p.failed("diamond_helmet missing from minecraft-data")
		return
	}
	p.sendCommand("item replace entity @e[type=armor_stand,limit=1] armor.head with minecraft:diamond_helmet")

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		stack, tracked := stand.Equipment(entity.EquipmentHead)
		if tracked && stack.Count > 0 {
			fmt.Printf("head slot item id=%d count=%d components=%d\n", stack.ItemId, stack.Count, len(stack.Components))
			p.check(stack.ItemId == helmet.Id,
				"head item id %d matches diamond_helmet id %d", stack.ItemId, helmet.Id)
			p.check(stack.Count == 1, "head count is 1")
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	p.failed("head slot was never tracked")
}

// runMetadata summons a baby lime sheep and asserts the tracked metadata.
func (p *probe) runMetadata() {
	x, y, z := p.player.Position()
	sx, sy, sz := x-2, y, z
	p.sendCommand(fmt.Sprintf("summon minecraft:sheep %.2f %.2f %.2f {Age:-24000,Color:5b}", sx, sy, sz))

	sheep := p.waitForEntity(sx, sy, sz, 1.5, 5*time.Second)
	if sheep == nil {
		p.failed("sheep at %.1f,%.1f,%.1f was not tracked", sx, sy, sz)
		return
	}
	fmt.Printf("sheep tracked id=%d type=%d pos=(%.2f,%.2f,%.2f)\n",
		sheep.Id, sheep.Type, sheep.Position.X, sheep.Position.Y, sheep.Position.Z)

	var flags int8
	hasFlags := false
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if got, ok := sheep.SharedFlags(); ok {
			flags, hasFlags = got, true
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	p.check(hasFlags, "shared flags metadata tracked: %d", flags)

	baby, ok := sheep.Baby()
	fmt.Printf("baby=%v (ok=%v) wool=%v\n", baby, ok, sheepWoolString(sheep))
	p.check(ok && baby, "sheep baby flag tracked")

	color, sheared, ok := sheep.SheepWool()
	p.check(ok, "sheep wool metadata tracked")
	p.check(color == 5, "sheep wool colour = lime (5), got %d", color)
	p.check(!sheared, "sheep is not sheared")
}

func sheepWoolString(sheep *entity.Entity) (ret string) {
	color, sheared, ok := sheep.SheepWool()
	if !ok {
		return "untracked"
	}
	return fmt.Sprintf("color=%d sheared=%v", color, sheared)
}

func Run(args []string) (err error) {
	if err = cmd.Parse(args); err != nil {
		return
	}
	acc, err := cmd_util.Account(*accFlag)
	if err != nil {
		return
	}

	c, err := client.Login(*addrFlag, acc, client.WithVersion(*versionFlag))
	if err != nil {
		return
	}
	defer c.Close()

	p := &probe{
		Conn:     c,
		entities: entity.Start(c),
		player:   player.Start(c),
	}
	keepalive.Start(c)
	p.RegisterUntilLatest(p.OnSystemChat)

	done := make(chan error, 1)
	go func() { done <- c.StartConn() }()

	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		if p.player.Joined() {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if !p.player.Joined() {
		return fmt.Errorf("timed out waiting for the player to join")
	}
	time.Sleep(time.Second)

	p.runEquipment()
	p.runMetadata()

	if *noAssert {
		fmt.Println("\nassertions skipped (--no-assert)")
	} else if p.fail > 0 {
		return fmt.Errorf("%d entities probe assertions failed", p.fail)
	} else {
		fmt.Println("\nall entities probe assertions passed")
	}
	return nil
}
