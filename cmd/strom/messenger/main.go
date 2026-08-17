package messenger

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"git.anygate.cloud/anygatecloud/strom/cmd/strom/cmd_util"
	"git.anygate.cloud/anygatecloud/strom/mc/bot/chat"
	"git.anygate.cloud/anygatecloud/strom/mc/bot/keepalive"
	"git.anygate.cloud/anygatecloud/strom/mc/bot/world"
	"git.anygate.cloud/anygatecloud/strom/mc/client"
	"git.anygate.cloud/anygatecloud/strom/mc/data"
	"git.anygate.cloud/anygatecloud/strom/mc/event"
	"git.anygate.cloud/anygatecloud/strom/mc/proto"
	"git.anygate.cloud/anygatecloud/strom/mc/proto_generated/v1_21_11"
	"git.anygate.cloud/anygatecloud/strom/mc/proto_generated/v1_8"
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

func (m *Messenger) OnStdin(e event.Stdin) (err error) {
	input := strings.TrimSpace(e.Val)
	if strings.HasPrefix(input, ">say ") {
		return m.sendChat(strings.TrimSpace(strings.TrimPrefix(input, ">say")))
	}
	if !m.legacy && strings.HasPrefix(input, ">getblock ") {
		return m.handleGetBlock(strings.TrimSpace(strings.TrimPrefix(input, ">getblock")))
	}
	return m.sendChat(input)
}

func (m *Messenger) sendChat(message string) (err error) {
	if m.legacy {
		return m.Send(&v1_8.PlayToServerPacketChat{Message: message})
	}
	return m.chat.SendMessage(message)
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

	if !isLegacy {
		w := world.NewWorld(c.Version, -64, 384)
		m.world = world.Start(c, w)
	}

	event.StartListingStdin(m.Loop)
	m.Register(m.OnStdin)

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
