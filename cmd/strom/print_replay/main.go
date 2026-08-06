package print_replay

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io"
	"log/slog"
	"os"
	"reflect"
	"strings"

	"git.anygate.cloud/anygatecloud/strom/mc/proto"
	"git.anygate.cloud/anygatecloud/strom/mc/proto_base"
)

var (
	cmd             = flag.NewFlagSet("print-replay", flag.ContinueOnError)
	FileFlag        = cmd.String("file", "", "The tmcpr file to print")
	ProtocolVersion = cmd.Int("protocol", -1, "The protocol version")
	CreateTestsPath = cmd.String("create-tests", "", "")
	StateFlag       = cmd.String("state", "login", "Initial state: handshake, login, config, status, play")
	DirectionFlag   = cmd.String("direction", "clientbound", "Packet direction: clientbound (ToClient) or serverbound (ToServer)")
)

// TODO: maybe use the proto replay api

//go:embed failed_packet.go.tmpl
var TestSrcF string

func SaveUnCodeAbleAsTest(packet *proto.UnCodablePacket) (err error) {
	if packet.Info.Type == nil {
		slog.Warn("skipping test for unknown packet id", "packet_id", packet.Info.PacketId, "state", packet.Info.State, "error", packet.Err)
		return
	}

	hUntrimmed := sha256.Sum256(packet.Data)
	h := hUntrimmed[:8]
	if err = os.MkdirAll(*CreateTestsPath, 0755); err != nil {
		return
	}

	typeStr := reflect.TypeOf(packet.Info.Type).String()
	goPackage := typeStr[1:strings.LastIndex(typeStr, ".")]
	typeName := reflect.TypeOf(packet.Info.Type).Elem().Name()

	src := fmt.Sprintf(TestSrcF, goPackage, packet.Err, h, typeName, packet.Data)
	formatted, formatErr := format.Source([]byte(src))
	if formatErr != nil {
		formatted = []byte(src)
	}

	filePath := fmt.Sprintf(*CreateTestsPath+"/%v_%10x_test.go", len(packet.Data), h)
	if err = os.WriteFile(filePath, formatted, 0644); err != nil {
		return
	}
	slog.Info("Packet failed to parse", "error", packet.Err, "saved", filePath, "data", fmt.Sprintf("%q", packet.Data))
	return
}

var state = proto_base.Login
var direction = proto_base.ToClient

func UnpackPacket(r io.Reader) (packet proto_base.EncodeDecodeAble, i proto_base.PacketInfo, err error) {
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

	var ok bool
	i, ok = proto.LookUpTypeByPacketInfo(direction, state, packetId, int32(*ProtocolVersion))
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
	if b.Len() != 0 {
		packet = &proto.UnCodablePacket{Err: proto.PacketNotFullyDecodedErr, Data: packetData[indexafterid:], Info: i}
		return
	}
	trackState(packet, i)
	return
}

func trackState(packet proto_base.EncodeDecodeAble, i proto_base.PacketInfo) {
	switch i.Name {
	case "set_protocol":
		v := reflect.ValueOf(packet).Elem()
		nsField := v.FieldByName("NextState")
		if nsField.IsValid() && nsField.Kind() == reflect.Int32 {
			switch nsField.Int() {
			case 1:
				state = proto_base.Status
			case 2:
				state = proto_base.Login
			}
		}
	case "login_start":
		if int32(*ProtocolVersion) <= 47 {
			state = proto_base.Play
		}
	case "success":
		if int32(*ProtocolVersion) >= 764 {
			state = proto_base.Configuration
		} else {
			state = proto_base.Play
		}
	case "finish_configuration":
		state = proto_base.Play
	case "login_acknowledged":
		state = proto_base.Configuration
	}
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
	case "handshake":
		state = proto_base.Handshaking
	case "login":
		state = proto_base.Login
	case "config":
		state = proto_base.Configuration
	case "play":
		state = proto_base.Play
	case "status":
		state = proto_base.Status
	default:
		return fmt.Errorf("unknown state: %s", *StateFlag)
	}
	switch *DirectionFlag {
	case "clientbound", "to-client":
		direction = proto_base.ToClient
	case "serverbound", "to-server":
		direction = proto_base.ToServer
	default:
		return fmt.Errorf("unknown direction: %s", *DirectionFlag)
	}
	f, err := os.Open(*FileFlag)
	if err != nil {
		return
	}
	defer f.Close()

	var totalPackets, failedPackets int
	for {
		var packet proto_base.EncodeDecodeAble
		var info proto_base.PacketInfo
		packet, info, err = UnpackPacket(f)
		unCodablePacket, ok := packet.(*proto.UnCodablePacket)
		if ok {
			slog.Error("failed to decode packet", "error", err, "packet_name", info.Name)
			if *CreateTestsPath != "" {
				err = SaveUnCodeAbleAsTest(unCodablePacket)
				if err != nil {
					slog.Error("failed to save uncodable packet as test", "error", err)
				}
			}
			failedPackets++
		}

		fmt.Printf("Packet: %#v\n", packet)
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			break
		}
		totalPackets++
	}
	if totalPackets > 0 {
		fmt.Printf("Total packets: %d, failed packets: %d, %%failed: %.2f%%\n", totalPackets, failedPackets, float64(failedPackets)/float64(totalPackets)*100)
	} else {
		fmt.Printf("Total packets: %d, failed packets: %d\n", totalPackets, failedPackets)
	}
	return
}
