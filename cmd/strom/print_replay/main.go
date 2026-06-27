package print_replay

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"reflect"
	"strings"

	"github.com/admin-else/strom/mc/data"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_base"
)

var (
	cmd             = flag.NewFlagSet("print-replay", flag.ContinueOnError)
	FileFlag        = cmd.String("file", "", "The tmcpr file to print")
	ProtocolVersion = cmd.Int("protocol", -1, "The protocol version")
	CreateTestsPath = cmd.String("create-tests", "", "")
	StateFlag       = cmd.String("state", "login", "Initial state: login, config, play")
)

// TODO: maybe use the proto replay api

//go:embed failed_packet.go.tmpl
var TestSrcF string

func SaveUnCodeAbleAsTest(packet *proto.UnCodablePacket) {
	hUntrimmed := sha256.Sum256(packet.Data)
	h := hUntrimmed[:8]
	f, err := os.Create(fmt.Sprintf(*CreateTestsPath+"/%v_%10x_test.go", len(packet.Data), h))
	if err != nil {
		panic(err)
	}

	ret, err := data.LookUpVersionByProtocolVersion(packet.Info.ProtocolVersion)
	if err != nil {
		panic(err)
	}
	goPackage := "v" + strings.ReplaceAll(ret.MinecraftVersion, ".", "_")

	//goland:noinspection GoUnhandledErrorResult
	defer f.Close()
	_, err = fmt.Fprintf(f, TestSrcF, goPackage, packet.Err, h, reflect.TypeOf(packet.Info.Type).Elem().Name(), packet.Data)
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

	// Read the packet data into a byte slice
	packetData := make([]byte, length)
	_, err = io.ReadFull(r, packetData)
	if err != nil {
		return
	}

	// Use bytes.Reader instead of bytes.Buffer
	b := bytes.NewReader(packetData)

	packetId, err := proto_base.DecodeVarInt(b)
	if err != nil {
		return
	}

	// Calculate position after reading the packet ID
	// indexafterid represents how many bytes we've read so far
	indexafterid := int(b.Size()) - b.Len()

	i, ok := proto.LookUpTypeByPacketInfo(proto_base.ToClient, state, packetId, int32(*ProtocolVersion))
	if !ok {
		packet = &proto.UnCodablePacket{Err: proto.BadPacketIdErr, Data: packetData, Info: i}
		return
	}
	packet = reflect.New(reflect.TypeOf(i.Type).Elem()).Interface().(proto_base.EncodeDecodeAble)
	err = packet.Decode(b)
	if err != nil {
		packet = &proto.UnCodablePacket{Err: err, Data: packetData[indexafterid:], Info: i}
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
	switch *StateFlag {
	case "login":
		state = proto_base.Login
	case "config":
		state = proto_base.Configuration
	case "play":
		state = proto_base.Play
	default:
		return fmt.Errorf("unknown state: %s", *StateFlag)
	}
	f, err := os.Open(*FileFlag)
	if err != nil {
		return
	}
	defer f.Close()

	var totalPackets, failedPackets int
	for {
		var packet proto_base.EncodeDecodeAble
		packet, err = UnpackPacket(f)
		unCodablePacket, ok := packet.(*proto.UnCodablePacket)
		if ok {
			if *CreateTestsPath != "" {
				SaveUnCodeAbleAsTest(unCodablePacket)
			}
			failedPackets++
		}

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Packet: %#v\n", packet)
		}
		if errors.Is(err, io.EOF) {
			break
		}
		totalPackets++
	}
	fmt.Printf("Total packets: %d, failed packets: %d, %%failed: %.2f%%\n", totalPackets, failedPackets, float64(failedPackets)/float64(totalPackets)*100)
	return
}
