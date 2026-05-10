package print_replay

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/proto_base"
)

var (
	cmd             = flag.NewFlagSet("print-replay", flag.ContinueOnError)
	FileFlag        = cmd.String("file", "", "The tmcpr file to print")
	ProtocolVersion = cmd.Int("protocol", -1, "The protocol version")
)

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
		err = nil
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
		if err != nil {
			return
		}
		fmt.Printf("Packet: %#v\n", packet)
	}
}
