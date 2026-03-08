package server

import (
	"errors"
	"log/slog"
	"net"

	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/text"
)

func Servee(c net.Conn) (ret *proto.Conn) {
	ret = &proto.Conn{}
	_ = ret.SetVersion("1.21.8") // cant fail
	ret.SetState(proto_base.Handshaking)
	ret.SetCompressionThreshold(-1)
	ret.Actor = proto_base.Servee
	ret.Conn = c
	ret.R = c
	ret.W = c
	return
}

func ServeClient(cNet net.Conn, factory func(c *proto.Conn) (h any, err error)) {
	c := Servee(cNet)
	defer c.Close()
	h, err := factory(c)
	if err != nil {
		return
	}
	err = c.Start(h)
	if errors.Is(err, StatusServedErr) {
		return
	}
	if err != nil {
		slog.Error("Error while handling client", "error", err, "client", c.Conn.RemoteAddr())
		err = Kick(c, text.Pretty(err.Error()))
		if err != nil {
			slog.Error("Error while kicking client", "error", err, "client", c.Conn.RemoteAddr())
		}
	}
	return
}

func StartServerWithFactory(listenAddr string, factory func(c *proto.Conn) (h any, err error)) (err error) {
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return
	}
	var cNet net.Conn
	for {
		cNet, err = l.Accept()
		if err != nil {
			return
		}
		go ServeClient(cNet, factory)
	}
}

func StartServerWithOnConn(listenAddr string, onConn func(c *proto.Conn) (err error)) (err error) {
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return
	}
	for {
		var cNet net.Conn
		cNet, err = l.Accept()
		if err != nil {
			return
		}
		go func() {
			c := Servee(cNet)
			defer c.Close()
			connErr := onConn(c)
			if connErr != nil {
				if c.State() == proto_base.Status || c.State() == proto_base.Handshaking {
					return // we dont care about status packets
				}
				slog.Error("Error while handling client", "error", connErr, "client", c.Conn.RemoteAddr())
				connErr = Kick(c, text.Pretty(connErr.Error()))
				if connErr != nil {
					slog.Error("Error while kicking client", "error", connErr, "client", c.Conn.RemoteAddr())
				}
			}
		}()
	}
}
