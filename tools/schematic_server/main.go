package main

import (
	"log/slog"

	"github.com/admin-else/strom/event"
	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/server"
	"github.com/admin-else/strom/text"
)

func wontFail[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

var StatusResponse = server.StatusResponse{
	Version: server.StatusResponseVersion{
		Name:     text.Pretty("STROM"),
		Protocol: 772,
	},
	Players:            server.StatusResponsePlayers{},
	Description:        text.Pretty("GO TO LIMBO"),
	Favicon:            "",
	EnforcesSecureChat: false,
}.MustMarshal()

type Server struct {
	*proto.Conn
}

func (s *Server) OnStart(_ event.OnStart) (err error) {
	//goland:noinspection GoResourceLeak captures the conn object which should NOT be closed
	_, err = server.ServeLogin(s.Conn, server.WithRawStatus(StatusResponse))
	if err != nil {
		return
	}
	return
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	err := server.StartServerWithFactory(":25565", func(c *proto.Conn) (h any, err error) {
		h = &Server{Conn: c}
		return
	})
	if err != nil {
		panic(err)
	}
}
