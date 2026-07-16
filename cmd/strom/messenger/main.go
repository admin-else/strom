package messenger

import (
	"flag"

	"github.com/admin-else/strom/cmd/strom/cmd_util"
	"github.com/admin-else/strom/mc/bot/chat"
	"github.com/admin-else/strom/mc/bot/keepalive"
	"github.com/admin-else/strom/mc/client"
	"github.com/admin-else/strom/mc/event"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_11"
)

var (
	cmd           = flag.NewFlagSet("status", flag.ContinueOnError)
	connectToFlag = cmd.String("addr", "", "server address to connect to")
	accFlag       = cmd.String("acc", "strom_messenger", "the name to use if longer will be treated as ygg token")
)

type Messenger struct {
	*proto.Conn
	chat *chat.Module
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
	return m.chat.SendMessage(e.Val)
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
	m := &Messenger{
		Conn: c,
		chat: chatMod,
	}

	event.StartListingStdin(m.Loop)
	m.Register(m.OnStdin)
	m.RegisterUntilLatest(m.OnChat, m.OnChatUnsigned, m.OnChatSystem)
	keepalive.Start(m.Conn)

	return m.StartConn()
}
