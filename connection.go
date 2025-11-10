package strom

import (
	"bytes"
	"errors"
	"io"
	"net"

	"github.com/admin-else/queser"
	"github.com/admin-else/queser/generated"
)

var BadPacketTypeError = errors.New("bad packet type")

type Actor int

const (
	Server Actor = iota
	Client
)

func (a Actor) SendDirection() queser.Direction {
	switch a {
	case Server:
		return queser.ToClient
	case Client:
		return queser.ToServer
	default:
		panic("invalid actor")
	}
}

func (a Actor) ReceiveDirection() queser.Direction {
	switch a {
	case Server:
		return queser.ToServer
	case Client:
		return queser.ToClient
	default:
		panic("invalid actor")
	}
}

type Connection struct {
	io.ReadWriteCloser
	CompressionThreshold int32
	State                queser.State
	Actor                Actor
	Version              string
}

func (c *Connection) Send(packet any) (err error) {
	packetIdentifier := generated.TypeToPacketIdentifier(c.Version, c.Actor.SendDirection(), c.State, packet)
	if packetIdentifier == "" {
		err = BadPacketTypeError
		return
	}
	rawPacketBuffer := bytes.NewBuffer(nil)
	err = generated.EncodePacket(c.Version, c.Actor.SendDirection(), c.State, packetIdentifier, packet, rawPacketBuffer)
	if err != nil {
		return
	}
	var packetBytes []byte
	if c.CompressionThreshold > 0 {
		panic("todo")
	} else {
		packetBytes = rawPacketBuffer.Bytes()
	}
	err = queser.VarInt(len(packetBytes)).Encode(c)
	if err != nil {
		return
	}
	_, err = c.Write(packetBytes)
	return
}

func (c *Connection) Receive() (packet any, err error) {
	var packetLen queser.VarInt
	packetLen, err = packetLen.Decode(c)
	if err != nil {
		return
	}
	rawPacketBytes, err := io.ReadAll(io.LimitReader(c, int64(packetLen)))
	if err != nil {
		return
	}
	var packetBytes []byte
	if c.CompressionThreshold > 0 {
		panic("todo")
	} else {
		packetBytes = rawPacketBytes
	}
	packet, err = generated.DecodePacket(c.Version, c.Actor.ReceiveDirection(), c.State, bytes.NewReader(packetBytes))
	return
}

func ClientConnect(addr string) (ret *Connection, err error) {
	ret = &Connection{}
	ret.Version = "1.21.8"
	ret.State = queser.Handshaking
	ret.CompressionThreshold = -1
	ret.Actor = Client
	ret.ReadWriteCloser, err = net.Dial("tcp", addr)
	if err != nil {
		return
	}
	return
}
