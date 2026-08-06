package replay

import (
	"bytes"
	"encoding/binary"
	"io"

	proto2 "git.anygate.cloud/anygatecloud/strom/mc/proto"
	"git.anygate.cloud/anygatecloud/strom/mc/proto_base"
)

type Writer struct {
	w io.Writer
}

// NewWriter creates a new replay Writer from an io.Writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

func (w *Writer) WriteRaw(p []byte, timestamp uint32) (err error) {
	err = binary.Write(w.w, binary.BigEndian, timestamp)
	if err != nil {
		return
	}
	err = binary.Write(w.w, binary.BigEndian, uint32(len(p)))
	if err != nil {
		return
	}
	_, err = w.w.Write(p)
	return
}

func (w *Writer) WritePacket(packet proto_base.EncodeDecodeAble, timestamp uint32) (err error) {
	b := bytes.NewBuffer(nil)

	info, ok := proto2.LookupPacketInfoByType(packet)
	if !ok {
		return proto2.BadPacketIdErr
	}

	err = proto_base.EncodeVarInt(b, info.PacketId)
	if err != nil {
		return
	}

	err = packet.Encode(b)
	if err != nil {
		return
	}
	return w.WriteRaw(b.Bytes(), timestamp)
}
