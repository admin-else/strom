package bot

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"
	"net"

	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/proto_generated"
)

var BadPacketTypeError = errors.New("bad packet type")

// UnCodablePacket represents a packet that could not be decoded due to querser not supporting all packets.
type UnCodablePacket struct {
	Err     error
	Partial any
	Data    []byte
}

type Actor int

const (
	Server Actor = iota
	Client
)

func (a Actor) SendDirection() proto_base.Direction {
	switch a {
	case Server:
		return proto_base.ToClient
	case Client:
		return proto_base.ToServer
	default:
		panic("invalid actor")
	}
}

func (a Actor) ReceiveDirection() proto_base.Direction {
	switch a {
	case Server:
		return proto_base.ToServer
	case Client:
		return proto_base.ToClient
	default:
		panic("invalid actor")
	}
}

type Conn struct {
	net.Conn
	R                    io.Reader
	W                    io.Writer
	CompressionThreshold int32
	State                proto_base.State
	Actor                Actor
	Version              string
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
			err = proto_base.VarInt(len(packetBytes)).Encode(packetBuffer)
			if err != nil {
				return
			}
			_, err = zlib.NewWriter(packetBuffer).Write(rawPacketBytes)
			if err != nil {
				return
			}
		} else {
			err = proto_base.VarInt(0).Encode(packetBuffer)
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
	err = proto_base.VarInt(len(packetBytes)).Encode(c.W)
	if err != nil {
		return
	}
	_, err = c.W.Write(packetBytes)
	return
}

func (c *Conn) Receive() (packet any, err error) {
	var rawPacketLen proto_base.VarInt
	rawPacketLen, err = rawPacketLen.Decode(c.R)
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
		var packetLen proto_base.VarInt
		packetLen, err = packetLen.Decode(rawPacketBuffer)
		if err != nil {
			return
		}
		if packetLen == 0 {
			packetBytes, err = io.ReadAll(rawPacketBuffer)
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
		packet = UnCodablePacket{Err: err, Data: packetBytes, Partial: packet}
		err = nil
	}
	return
}

func Connect(addr string) (ret *Conn, err error) {
	ret = &Conn{}
	ret.Version = "1.21.8"
	ret.State = proto_base.Handshaking
	ret.CompressionThreshold = -1
	ret.Actor = Client
	ret.Conn, err = net.Dial("tcp", addr)
	if err != nil {
		return
	}
	ret.R = ret.Conn
	ret.W = ret.Conn
	return
}

func Servee(c net.Conn) (ret *Conn) {
	ret = &Conn{}
	ret.Version = "1.21.8"
	ret.State = proto_base.Handshaking
	ret.CompressionThreshold = -1
	ret.Actor = Server
	ret.Conn = c
	ret.R = c
	ret.W = c
	return
}
