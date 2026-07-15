package status

import (
	"flag"
	"fmt"
	"log/slog"
	"time"

	"github.com/admin-else/strom/mc/client"
	"github.com/admin-else/strom/mc/event"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_base"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_8"
)

var cmd = flag.NewFlagSet("status", flag.ContinueOnError)
var addrFlag = cmd.String("addr", "localhost:25565", "address to connect to")
var srvFlag = cmd.Bool("srv", false, "use SRV records to resolve the address")
var versionFlag = cmd.String("version", "1.21.8", "version to use")

type StatusClient struct {
	*proto.Conn
	Status                        string
	PingSendTime, PingReceiveTime time.Time
	DoPingRoundTripTime           bool
}

func (s *StatusClient) OnStatus(p *v1_21_8.StatusToClientPacketServerInfo) (err error) {
	s.Status = p.Response
	s.PingSendTime = time.Now()
	if s.DoPingRoundTripTime {
		err = s.Send(&v1_21_8.StatusToServerPacketPing{Time: s.PingSendTime.UnixMilli()})
	} else {
		err = event.HandlerDoneErr{}
	}
	return
}

func (s *StatusClient) OnPong(p *v1_21_8.StatusToClientPacketPing) (err error) {
	s.PingReceiveTime = time.Now()
	err = event.HandlerDoneErr{}
	return
}

// StatusRaw returns a StatusClient that is not connected to a server.
// So ignore the resource leak warning. And maybe attach a warn ignore comment.
// It does not resolve SRV records.
func StatusRaw(addr string) (s *StatusClient, err error) {
	c, err := client.ConnectVersionLess(nil, addr)
	if err != nil {
		return
	}
	defer c.Close()
	err = c.SetVersion(*versionFlag)
	if err != nil {
		return
	}
	s = &StatusClient{
		Conn: c,
	}
	s.RegisterCriticalUntilLatest(s.OnStatus)
	s.RegisterCriticalUntilLatest(s.OnPong)

	p, err := client.MakeHandshakePacketAddr(s.Conn, proto_base.Status, addr)
	if err != nil {
		return
	}
	err = s.Send(p)
	if err != nil {
		return
	}
	s.SetState(proto_base.Status)
	err = s.Send(&v1_21_8.StatusToServerPacketPingStart{})

	err = s.StartConn()
	return
}

func Run(args []string) (err error) {
	err = cmd.Parse(args)
	if err != nil {
		return
	}
	if *srvFlag {
		*addrFlag, _, _, err = client.DoDNSChecked(nil, *addrFlag, nil)
		if err != nil {
			return
		}
	}
	slog.Debug("resolved address to ", "addr", *addrFlag)

	status, err := StatusRaw(*addrFlag)
	if err != nil {
		return
	}
	fmt.Println(status.Status)
	return
}
