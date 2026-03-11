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
}

func (s *StatusClient) OnStart() (err error) {
	s.RegisterUntilLatest(s.OnStatus)
	s.RegisterUntilLatest(s.OnPong)

	p, err := MakeHandshakePacket(s.Conn, proto_base.Status)
	if err != nil {
		return
	}
	err = s.Send(p)
	if err != nil {
		return
	}
	s.SetState(proto_base.Status)
	err = s.Send(&v1_21_8.StatusToServerPacketPingStart{})
	return
}

func (s *StatusClient) OnStatus(p *v1_21_8.StatusToClientPacketServerInfo) (err error) {
	s.Status = p.Response
	s.PingSendTime = time.Now()
	err = s.Send(&v1_21_8.StatusToServerPacketPing{Time: s.PingSendTime.UnixMilli()})
	return
}

func (s *StatusClient) OnPong(p *v1_21_8.StatusToClientPacketPing) (err error) {
	s.PingReceiveTime = time.Now()
	err = event.HandlerDoneErr{}
	return
}

func StatusRaw(addr string) (s *StatusClient, err error) {
	c, err := ConnectVersionLess(addr)
	if err != nil {
		return
	}
	err = c.SetVersion("1.21.8")
	if err != nil {
		return
	}
	s = &StatusClient{
		Conn: c,
	}
	err = s.Start(s)
	defer c.Close()
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
