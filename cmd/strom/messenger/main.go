package messenger

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/admin-else/strom/cmd/strom/cmd_util"
	"github.com/admin-else/strom/mc/bot/chat"
	"github.com/admin-else/strom/mc/bot/keepalive"
	"github.com/admin-else/strom/mc/bot/world"
	"github.com/admin-else/strom/mc/client"
	"github.com/admin-else/strom/mc/data"
	"github.com/admin-else/strom/mc/event"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_11"
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
)

type Messenger struct {
	*proto.Conn
	chat  *chat.Module
	world *world.Module
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

func (m *Messenger) OnStdin(e event.Stdin) (err error) {
	input := strings.TrimSpace(e.Val)
	if strings.HasPrefix(input, ">say ") {
		return m.chat.SendMessage(strings.TrimSpace(strings.TrimPrefix(input, ">say")))
	}
	if strings.HasPrefix(input, ">getblock ") {
		return m.handleGetBlock(strings.TrimSpace(strings.TrimPrefix(input, ">getblock")))
	}
	return m.chat.SendMessage(input)
}

func (m *Messenger) handleGetBlock(args string) (err error) {
	fields := strings.Fields(args)
	if len(fields) != 3 {
		m.Log.Info("getblock", "error", "expected >getblock <x> <y> <z>")
		return nil
	}

	coords := make([]int64, 3)
	for i, f := range fields {
		coords[i], err = strconv.ParseInt(f, 10, 32)
		if err != nil {
			m.Log.Info("getblock", "error", fmt.Sprintf("bad coordinate %q", f))
			return nil
		}
	}

	stateId, err := m.world.GetBlock(int32(coords[0]), int32(coords[1]), int32(coords[2]))
	if err != nil {
		m.Log.Info("getblock", "error", err.Error())
		return nil
	}

	block, ok := data.LookupBlockByStateId(m.Version, stateId)
	if !ok {
		m.Log.Info("getblock", "stateId", stateId, "name", "unknown")
		return nil
	}

	m.Log.Info("getblock", "x", coords[0], "y", coords[1], "z", coords[2], "name", block.Name, "stateId", stateId)
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
	c, err := client.Login(*connectToFlag, acc)
	if err != nil {
		return
	}
	defer c.Close()

	chatMod, err := chat.Start(c, acc)
	if err != nil {
		return
	}

	// Shareable world storage. Other bots can subscribe to the same World.
	w := world.NewWorld(c.Version, -64, 384)

	m := &Messenger{
		Conn:  c,
		chat:  chatMod,
		world: world.Start(c, w),
	}

	event.StartListingStdin(m.Loop)
	m.Register(m.OnStdin)
	m.RegisterUntilLatest(m.OnChat, m.OnChatUnsigned, m.OnChatSystem)
	keepalive.Start(m.Conn)

	return m.StartConn()
}
