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
	ret = proto.NewConn()
	ret.Actor = proto_base.Servee
	ret.Conn = c
	ret.R = c
	ret.W = c

	ret.Log = slog.Default().With("conn", c.RemoteAddr())
	return
}

type ConnAcceptor func(c net.Conn) (err error)

type Factory func(c *proto.Conn) (h ConnAcceptor, err error)

func ServeClient(cNet net.Conn, factory Factory) {
	c := Servee(cNet)
	defer c.Close()
	h, err := factory(c)
	if err != nil {
		return
	}
	err = h(c)
	if errors.Is(err, StatusServedErr) {
		return
	}
	if err != nil {
		err = Kick(c, text.Pretty(err.Error()), "")
	}
}

func StartServerWithFactory(listenAddr string, factory Factory) (err error) {
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

var ShutDownServerErr = errors.New("server shutting down")

func StartServerWithOnConn(listenAddr string, onConn func(c *proto.Conn) (err error)) (err error) {
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return
	}
	defer l.Close()
	running := true
	for running {
		var cNet net.Conn
		cNet, err = l.Accept()
		if err != nil {
			return
		}
		go func() {
			c := Servee(cNet)
			defer c.Close()
			connErr := onConn(c)
			if errors.Is(connErr, ShutDownServerErr) {
				running = false
			}
			c.Log.Error("closing connection", "err", connErr)
			if connErr != nil {
				connErr = Kick(c, text.Pretty(connErr.Error()), "")
			}
		}()
	}
	err = ShutDownServerErr
	return
}
