package server

import (
	"log"
	"net"

	"github.com/admin-else/strom/event"
	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/proto_base"
)

func Servee(c net.Conn) (ret *proto.Conn) {
	ret = &proto.Conn{}
	ret.Version = "1.21.8"
	ret.State = proto_base.Handshaking
	ret.CompressionThreshold = -1
	ret.Actor = proto_base.Server
	ret.Conn = c
	ret.R = c
	ret.W = c
	return
}

func ServeClient(cNet net.Conn, factory func(c *proto.Conn) (h event.HandlerInst, err error)) (err error) {
	c := Servee(cNet)
	h, err := factory(c)
	if err != nil {
		return
	}
	err = c.Start(h)
	return
}

func StartServerWithFactory(listenAddr string, factory func(c *proto.Conn) (h event.HandlerInst, err error)) (err error) {
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
			connErr := onConn(Servee(cNet))
			if connErr != nil {
				log.Println(connErr)
			}
		}()
	}
}
