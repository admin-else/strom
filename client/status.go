package client

import (
	"context"
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
func StatusRaw(ctx context.Context, addr string) (s *StatusClient, err error) {
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
	s.RegisterCriticalUntilLatest(s.OnStatus)
	s.RegisterCriticalUntilLatest(s.OnPong)

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

// StatusNoDns is like Status but does not resolve SRV records.
func StatusNoDns(addr string) (status server.StatusResponse, err error) {
	s, err := StatusRaw(context.Background(), addr)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(s.Status), &status)
	return
}

// Status resolves the status of the minecraft server at the given address.
func Status(addr string) (status server.StatusResponse, err error) {
	addr, err = DoDnsSimple(addr)
	if err != nil {
		return
	}
	return StatusNoDns(addr)
}
