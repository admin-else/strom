package messenger

import (
	"flag"

	"github.com/admin-else/strom/mc/api"
	"github.com/admin-else/strom/mc/client"
	"github.com/admin-else/strom/mc/event"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_11"
)

var (
	cmd           = flag.NewFlagSet("status", flag.ContinueOnError)
	connectToFlag = cmd.String("addr", "", "server address to connect to")
	nameFlag      = cmd.String("name", "strom_messenger", "the name to use")
)

type Messenger struct {
	*proto.Conn
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
	return m.Send(&v1_21_11.PlayToServerPacketChatMessage{Message: e.Val})
}

func Run(args []string) (err error) {
	err = cmd.Parse(args)

	c, err := client.Login(*connectToFlag, api.NewOfflineAccount(*nameFlag))
	if err != nil {
		return
	}
	defer c.Close()

	m := Messenger{Conn: c}
	event.StartListingStdin(m.Loop)
	m.Register(m.OnStdin)
	m.RegisterUntilLatest(m.OnChat, m.OnChatUnsigned, m.OnChatSystem)

	return m.StartConn()
}
