package replay

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"

	proto2 "github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_base"
)

type Reader struct {
	r       io.Reader
	State   proto_base.State
	Version int32
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r: r}
}

func (r *Reader) ReadRaw() (packetData []byte, timestamp, length uint32, err error) {
	err = binary.Read(r.r, binary.BigEndian, &timestamp)
	if err != nil {
		return
	}
	err = binary.Read(r.r, binary.BigEndian, &length)
	if err != nil {
		return
	}

	// Read the packet data into a byte slice
	packetData = make([]byte, length)
	_, err = io.ReadFull(r.r, packetData)
	return
}

func (r *Reader) ReadPacket() (packet proto_base.EncodeDecodeAble, err error) {
	packetData, _, _, err := r.ReadRaw()
	if err != nil {
		return
	}
	b := bytes.NewReader(packetData)

	packetId, err := proto_base.DecodeVarInt(b)
	if err != nil {
		return
	}

	// Calculate position after reading the packet ID
	// indexafterid represents how many bytes we've read so far
	indexafterid := int(b.Size()) - b.Len()

	i, ok := proto2.LookUpTypeByPacketInfo(proto_base.ToClient, r.State, packetId, r.Version)
	if !ok {
		packet = &proto2.UnCodablePacket{Err: proto2.BadPacketIdErr, Data: packetData, Info: i}
		return
	}
	packet = reflect.New(reflect.TypeOf(i.Type).Elem()).Interface().(proto_base.EncodeDecodeAble)
	err = packet.Decode(b)
	if err != nil {
		packet = &proto2.UnCodablePacket{Err: err, Data: packetData[indexafterid:], Info: i}
		return
	}
	return
}
