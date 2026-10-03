package client

import (
	"context"
	"encoding/json"
	"time"

	"github.com/admin-else/strom/mc/event"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_base"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_8"
	"github.com/admin-else/strom/mc/proto_generated/v1_8"
	"github.com/admin-else/strom/mc/server"
)

type StatusClient struct {
	*proto.Conn
	Status                        string
	PingSendTime, PingReceiveTime time.Time
	DoPingRoundTripTime           bool
}

func (s *StatusClient) OnStatus(p *v1_8.StatusToClientPacketServerInfo) (err error) {
	s.Status = p.Response
	s.PingSendTime = time.Now()
	if s.DoPingRoundTripTime {
		err = s.Send(&v1_21_8.StatusToServerPacketPing{Time: s.PingSendTime.UnixMilli()})
	} else {
		err = event.HandlerDoneErr{}
	}
	return
}

func (s *StatusClient) OnPong(p *v1_8.StatusToClientPacketPing) (err error) {
	s.PingReceiveTime = time.Now()
	err = event.HandlerDoneErr{}
	return
}

// StatusRaw returns a StatusClient that is not connected to a server.
// So ignore the resource leak warning. And maybe attach a warn ignore comment.
// It does not resolve SRV records.
func StatusRaw(ctx context.Context, addr string) (s *StatusClient, err error) {
	return StatusRawWithVersion(ctx, addr, "1.21.8")
}

// StatusRawWithVersion is like StatusRaw but uses the given Minecraft version for the handshake.
func StatusRawWithVersion(ctx context.Context, addr, version string) (s *StatusClient, err error) {
	return statusRawWithVersion(ctx, addr, version, false)
}

// StatusRawWithVersionAndPing is like StatusRawWithVersion but also sends the
// ping packet and records PingSendTime/PingReceiveTime, so PingMillis returns a
// real round-trip. It is opt-in: StatusRawWithVersion keeps the
// status-response-only behaviour and avoids the extra round-trip.
func StatusRawWithVersionAndPing(ctx context.Context, addr, version string) (s *StatusClient, err error) {
	return statusRawWithVersion(ctx, addr, version, true)
}

func statusRawWithVersion(ctx context.Context, addr, version string, ping bool) (s *StatusClient, err error) {
	c, err := ConnectVersionLess(ctx, addr)
	if err != nil {
		return
	}
	defer c.Close()
	err = c.SetVersion(version)
	if err != nil {
		return
	}
	s = &StatusClient{
		Conn:                c,
		DoPingRoundTripTime: ping,
	}
	s.RegisterUntil("26.2", s.OnStatus, s.OnPong)

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

// PingMillis returns the pong round-trip in milliseconds, or 0 when the ping was
// not requested or no pong was received.
func (s *StatusClient) PingMillis() int64 {
	if s.PingSendTime.IsZero() || s.PingReceiveTime.IsZero() {
		return 0
	}
	return s.PingReceiveTime.Sub(s.PingSendTime).Milliseconds()
}

// StatusNoDns is like Status but does not resolve SRV records.
func StatusNoDns(ctx context.Context, addr string) (status server.StatusResponse, err error) {
	s, err := StatusRaw(ctx, addr)
	if err != nil {
		return
	}
	defer s.Close()
	err = json.Unmarshal([]byte(s.Status), &status)
	return
}

// Status resolves the status of the minecraft server at the given address.
func Status(ctx context.Context, addr string) (status server.StatusResponse, err error) {
	addr, _, _, err = DoDNSChecked(ctx, addr, nil)
	if err != nil {
		return
	}
	return StatusNoDns(ctx, addr)
}
