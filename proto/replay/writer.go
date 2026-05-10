package replay

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/admin-else/strom/proto_base"
)

type Writer struct {
	w io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w}
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
	err = packet.Encode(b)
	if err != nil {
		return
	}
	return w.WriteRaw(b.Bytes(), timestamp)
}
