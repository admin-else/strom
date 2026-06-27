package serve_tmcpr

import (
	"errors"
	"flag"
	"io"
	"log/slog"
	"math/rand/v2"
	"os"
	"time"

	"github.com/admin-else/strom/mc/event"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto/replay"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_8"
	server2 "github.com/admin-else/strom/mc/server"
)

type Server struct {
	*proto.Conn
	*event.Timer
}

var (
	cmd       = flag.NewFlagSet("serve-tmcpr", flag.ContinueOnError)
	tmcprFile = cmd.String("file", "/tmp/world.tmcpr", "Path to .tmcpr replay file")
	bindAddr  = cmd.String("bind", ":25566", "Address to listen on")
)

func Limbo(c *proto.Conn, tmcprPath string) (err error) {
	s := &Server{Conn: c}

	_, err = server2.ServeLogin(s.Conn, server2.WithoutOnlineMode())
	if err != nil {
		return
	}

	err = server2.ServeConfig(s.Conn)
	if err != nil {
		return
	}

	f, err := os.Open(tmcprPath)
	if err != nil {
		return
	}
	defer f.Close()

	r := replay.NewReader(f)
	r.Version = c.ProtocolVersion

	start := time.Now()

	for {
		var packetData []byte
		var timestamp uint32
		packetData, timestamp, _, err = r.ReadRaw()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		time.Sleep(time.Duration(timestamp)*time.Millisecond - time.Since(start))
		err = s.SendRaw(packetData)
		if err != nil {
			return err
		}
	}

	s.Timer = &event.Timer{}
	s.Timer.Every(time.Second*15, s.SendKeepAlive)
	s.Timer.Start(c.Loop)

	return s.StartConn()
}

func (s *Server) SendKeepAlive() error {
	return s.Send(&v1_21_8.PlayToClientPacketKeepAlive{KeepAliveId: rand.Int64()})
}

func Run(args []string) error {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	err := cmd.Parse(args)
	if err != nil {
		return err
	}
	if *tmcprFile == "" {
		slog.Error("no file specified")
		return nil
	}

	return server2.StartServerWithOnConn(*bindAddr, func(c *proto.Conn) error {
		return Limbo(c, *tmcprFile)
	})
}
