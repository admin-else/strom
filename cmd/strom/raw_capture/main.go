package raw_capture

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/admin-else/strom/mc/client"
)

var Cmd = flag.NewFlagSet("raw-capture", flag.ContinueOnError)

var TargetAddr = Cmd.String("addr", "127.0.0.1:25565", "address to connect to")
var ListenAddr = Cmd.String("listen", "127.0.0.1:25566", "address to listen on")
var ClientboundFile = Cmd.String("clientbound-file", "clientbound.tmcpr", "file to write clientbound packets to")
var ServerboundFile = Cmd.String("serverbound-file", "serverbound.tmcpr", "file to write serverbound packets to")

func Handle(c net.Conn, addr string) {
	sockPuppet, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	io.MultiWriter()
	go io.Copy(c, sockPuppet)
	io.Copy(sockPuppet, c)
}
func Run(args []string) (err error) {
	err = Cmd.Parse(args)
	if err != nil {
		return
	}

	clientboundFile, err := os.Create(*ClientboundFile)
	if err != nil {
		return fmt.Errorf("create clientbound file: %w", err)
	}
	defer clientboundFile.Close()

	serverboundFile, err := os.Create(*ServerboundFile)
	if err != nil {
		return fmt.Errorf("create serverbound file: %w", err)
	}
	defer serverboundFile.Close()

	resolvedAddr, _, _, err := client.DoDNSChecked(nil, *TargetAddr, nil)
	if err != nil {
		return
	}
	l, err := net.Listen("tcp", *ListenAddr)
	if err != nil {
		return
	}
	for {
		c, err := l.Accept()
		if err != nil {
			continue
		}
		go Handle(c, resolvedAddr)
	}
}
