package proto

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"
	"net"

	"github.com/admin-else/strom/data"
	"github.com/admin-else/strom/event"
	"github.com/admin-else/strom/proto_base"
)

var BadPacketTypeError = errors.New("bad packet type")
var BadPacketIdError = errors.New("bad packet id")
var PacketNotFullyDecodedError = errors.New("packet not fully decoded")

// UnCodablePacket represents a packet that could not be decoded due to proto_generated not supporting all packets.
type UnCodablePacket struct {
	Err       error
	Data      []byte
	Direction proto_base.Direction
}

func (u *UnCodablePacket) Encode(w io.Writer) (err error) {
	_, err = w.Write(u.Data)
	return
}

func (u *UnCodablePacket) Decode(r io.Reader) (err error) {
	err = errors.New("how did you even get here")
	return
}

type Conn struct {
	net.Conn
	R                    io.Reader
	W                    io.Writer
	CompressionThreshold int32
	State                proto_base.State
	Actor                proto_base.Actor
	Version              string
	ProtocolVersion      int32
}

func (c *Conn) SetVersion(version string) (err error) {
	c.Version = version
	versionData, err := data.LookUpProtocolVersionByName(c.Version)
	if err != nil {
		return
	}
	c.ProtocolVersion = int32(versionData.Version)
	return
}

func (c *Conn) Send(packet proto_base.EncodeDecodeAble) (err error) {
	switch packet := packet.(type) {
	case *UnCodablePacket:
		err = c.SendRaw(packet.Data)
	default:
		i, ok := LookupPacketInfoByTypeLookupPacketInfoByTypeAndMore(packet, c.State)
		if !ok {
			err = BadPacketTypeError
		}
		var packetBuff = bytes.NewBuffer(nil)
		err = proto_base.EncodeVarInt(packetBuff, i.PacketId)
		if err != nil {
			return
		}
		err = packet.Encode(packetBuff)
		if err != nil {
			return
		}
		err = c.SendRaw(packetBuff.Bytes())
	}
	return
}
func (c *Conn) SendRaw(rawPacketBytes []byte) (err error) {
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

func (c *Conn) ReceiveRaw() (packetBytes []byte, err error) {
	rawPacketLen, err := proto_base.DecodeVarInt(c.R)
	if err != nil {
		return
	}
	rawPacketBytes, err := io.ReadAll(io.LimitReader(c.R, int64(rawPacketLen)))
	if err != nil {
		return
	}
	if len(rawPacketBytes) != int(rawPacketLen) {
		err = errors.New("bad packet length")
		return
	}
	rawPacketBuffer := bytes.NewBuffer(rawPacketBytes)
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
	return
}

func (c *Conn) Receive() (packet proto_base.EncodeDecodeAble, err error) {
	packetBytes, err := c.ReceiveRaw()
	if err != nil {
		return
	}
	b := bytes.NewBuffer(packetBytes)
	id, err := proto_base.DecodeVarInt(b)
	if err != nil {
		return
	}
	packet, ok := LookUpTypeByPacketInfoAndCopyType(c.Actor.ReceiveDirection(), c.State, id, c.ProtocolVersion)
	if !ok {
		packet = &UnCodablePacket{Err: BadPacketIdError, Data: packetBytes, Direction: c.Actor.ReceiveDirection()}
		return
	}
	err = packet.Decode(b)
	if err != nil {
		packet = &UnCodablePacket{Err: err, Data: packetBytes, Direction: c.Actor.ReceiveDirection()}
		err = nil
	}
	if b.Len() != 0 {
		err = nil
		packet = &UnCodablePacket{Err: PacketNotFullyDecodedError, Data: packetBytes, Direction: c.Actor.ReceiveDirection()}
		return
	}
	return
}

func (c *Conn) StartOne(inst any) (err error) {
	return c.Start([]any{inst})
}

func (c *Conn) Start(insts []any) (err error) {
	_ = *c // exit early on nil connection
	handlers := event.FindHandlers(insts)
	err = event.Fire(event.OnStart{}, handlers)
	for err == nil {
		var packet any
		packet, err = c.Receive()
		if err != nil {
			break
		}
		err = event.Fire(packet, handlers)
		if err != nil {
			break
		}
		err = event.Fire(event.OnLoopCycle{}, handlers)
	}
	if errors.Is(err, event.HandlerDone) {
		err = nil
	}
	return
}
