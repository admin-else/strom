package print_replay

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"reflect"

	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/proto_base"
)

var (
	cmd             = flag.NewFlagSet("print-replay", flag.ContinueOnError)
	FileFlag        = cmd.String("file", "", "The tmcpr file to print")
	ProtocolVersion = cmd.Int("protocol", -1, "The protocol version")
	CreateTestsPath = cmd.String("create-tests", "", "The path to create tests from failed packets, if empty no tests will be created")
)

//go:embed failed_packet.go.tmpl
var TestSrcF string

func SaveUnCodeAbleAsTest(d proto_base.Direction, packet *proto.UnCodablePacket) {
	hUntrimmed := sha256.Sum256(packet.Data)
	h := hUntrimmed[:8]
	f, err := os.Create(fmt.Sprintf(*CreateTestsPath+"/%v_%10x_test.go", len(packet.Data), h))
	if err != nil {
		panic(err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer f.Close()
	_, err = fmt.Fprintf(f, TestSrcF, packet.Err, h, d.Opposite(), packet.Data)
	if err != nil {
		panic(err)
	}
	slog.Info("Packet failed to parse", "error", packet.Err, "saved", f.Name(), "data", fmt.Sprintf("%q", packet.Data))
}

var state = proto_base.Login

func UnpackPacket(r io.Reader) (packet proto_base.EncodeDecodeAble, err error) {
	var timestamp, length uint32
	err = binary.Read(r, binary.BigEndian, &timestamp)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return
	}
	b := bytes.NewBuffer(nil)
	_, err = io.CopyN(b, r, int64(length))
	if err != nil {
		return
	}

	packetId, err := proto_base.DecodeVarInt(b)
	if err != nil {
		return
	}

	i, ok := proto.LookUpTypeByPacketInfo(proto_base.ToClient, state, packetId, int32(*ProtocolVersion))
	if !ok {
		packet = &proto.UnCodablePacket{Err: proto.BadPacketIdErr, Data: b.Bytes(), Direction: proto_base.ToClient}
		return
	}
	packet = reflect.New(reflect.TypeOf(i.Type).Elem()).Interface().(proto_base.EncodeDecodeAble)
	err = packet.Decode(b)
	if err != nil {
		packet = &proto.UnCodablePacket{Err: err, Data: b.Bytes(), Direction: proto_base.ToClient}
		return
	}
	switch i.Name {
	case "success":
		state = proto_base.Configuration
	case "finish_configuration":
		state = proto_base.Play
	}
	return
}

func Run(args []string) (err error) {
	err = cmd.Parse(args)
	if err != nil {
		return
	}
	if *ProtocolVersion == -1 {
		return fmt.Errorf("no protocol version specified")
	}
	if *FileFlag == "" {
		return fmt.Errorf("no file specified")
	}
	f, err := os.Open(*FileFlag)
	if err != nil {
		return
	}
	defer f.Close()

	for {
		var packet proto_base.EncodeDecodeAble
		packet, err = UnpackPacket(f)
		unCodablePacket, ok := packet.(*proto.UnCodablePacket)
		if ok && *CreateTestsPath != "" {
			SaveUnCodeAbleAsTest(proto_base.ToClient, unCodablePacket)
		}

		if err != nil {
			return
		}
		fmt.Printf("Packet: %#v\n", packet)
	}
}
