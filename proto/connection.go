package proto

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"
	"net"

	"github.com/admin-else/strom/event"
	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/proto_generated"
)

var BadPacketTypeError = errors.New("bad packet type")

// UnCodablePacket represents a packet that could not be decoded due to proto_generated not supporting all packets.
type UnCodablePacket struct {
	Err  error
	Data []byte
}

type Conn struct {
	net.Conn
	R                                io.Reader
	W                                io.Writer
	CompressionThreshold             int32
	State                            proto_base.State
	Actor                            proto_base.Actor
	Version                          string
	DontDecodePacketsWithoutHandlers bool
}

func (c *Conn) Send(packet any) (err error) {
	var rawPacketBytes []byte
	switch packet := packet.(type) {
	case UnCodablePacket:
		rawPacketBytes = packet.Data
	default:
		packetIdentifier := proto_generated.TypeToPacketIdentifier(c.Version, c.Actor.SendDirection(), c.State, packet)
		if packetIdentifier == "" {
			err = BadPacketTypeError
			return
		}
		rawPacketBuffer := bytes.NewBuffer(nil)
		err = proto_generated.EncodePacket(c.Version, c.Actor.SendDirection(), c.State, packetIdentifier, packet, rawPacketBuffer)
		if err != nil {
			return
		}
		rawPacketBytes = rawPacketBuffer.Bytes()
	}
	var packetBytes []byte
	if c.CompressionThreshold > 0 {
		packetBuffer := bytes.NewBuffer(nil)
		if int32(len(packetBytes)) >= c.CompressionThreshold {
			err = proto_base.EncodeVarInt(packetBuffer, int32(len(packetBytes)))
			if err != nil {
				return
			}
			_, err = zlib.NewWriter(packetBuffer).Write(rawPacketBytes)
			if err != nil {
				return
			}
		} else {
			err = proto_base.EncodeVarInt(packetBuffer, 0)
			if err != nil {
				return
			}
			_, err = packetBuffer.Write(rawPacketBytes)
			if err != nil {
				return
			}
		}
		packetBytes = packetBuffer.Bytes()
	} else {
		packetBytes = rawPacketBytes
	}
	err = proto_base.EncodeVarInt(c.W, int32(len(packetBytes)))
	if err != nil {
		return
	}
	_, err = c.W.Write(packetBytes)
	return
}

func (c *Conn) Receive() (packet any, err error) {
	rawPacketLen, err := proto_base.DecodeVarInt(c.R)
	if err != nil {
		return
	}
	rawPacketBytes, err := io.ReadAll(io.LimitReader(c.R, int64(rawPacketLen)))
	if err != nil {
		return
	}
	rawPacketBuffer := bytes.NewBuffer(rawPacketBytes)
	var packetBytes []byte
	if c.CompressionThreshold > 0 {
		var packetLen int32
		packetLen, err = proto_base.DecodeVarInt(rawPacketBuffer)
		if err != nil {
			return
		}
		if packetLen == 0 {
			packetBytes, err = io.ReadAll(rawPacketBuffer)
			if err != nil {
				return
			}
		} else {
			var zReader io.ReadCloser
			zReader, err = zlib.NewReader(rawPacketBuffer)
			if err != nil {
				return
			}
			defer zReader.Close()
			packetBytes, err = io.ReadAll(zReader)
			if err != nil {
				return
			}
		}
	} else {
		packetBytes = rawPacketBytes
	}
	packetBuff := bytes.NewBuffer(packetBytes)
	packet, err = proto_generated.DecodePacket(c.Version, c.Actor.ReceiveDirection(), c.State, packetBuff)
	if err != nil {
		packet = UnCodablePacket{Err: err, Data: packetBytes}
		err = nil
	}
	return
}

func (c *Conn) Start(inst event.HandlerInst) (err error) {
	_ = *c // exit early on nil connection
	handlers := event.FindHandlers(inst)
	err = event.Fire(inst, event.OnStart{}, handlers)
	for err == nil {
		var packet any
		packet, err = c.Receive()
		if err != nil {
			return
		}
		err = event.Fire(inst, packet, handlers)
		if err != nil {
			break
		}
		err = event.Fire(inst, event.OnLoopCycle{}, handlers)
	}
	if errors.Is(err, event.HandlerDone) {
		err = nil
	}
	return
}
