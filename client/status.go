package client

import (
	"encoding/json"
	"time"

	"github.com/admin-else/strom/event"
	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/server"

	"github.com/admin-else/strom/proto_generated/v1_21_8"
)

type StatusClient struct {
	*proto.Conn
	Status                        string
	PingSendTime, PingReceiveTime time.Time
	DoPingRoundTripTime           bool
}

func (s *StatusClient) OnStart() (err error) {
	return
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
// So ignore the resource leak warning. And maybe attach a warn ignore comment
func StatusRaw(addr string) (s *StatusClient, err error) {
	c, err := ConnectVersionLess(addr)
	if err != nil {
		return
	}
	defer c.Close()
	err = c.SetVersion("1.21.8")
	if err != nil {
		return
	}
	s = &StatusClient{
		Conn: c,
	}
	s.Loop = event.NewLoop()
	s.RegisterUntilLatest(s.OnStatus)
	s.RegisterUntilLatest(s.OnPong)

	p, err := MakeHandshakePacketAddr(s.Conn, proto_base.Status, addr)
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

func Status(addr string) (status server.StatusResponse, err error) {
	s, err := StatusRaw(addr)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(s.Status), &status)
	return
}
