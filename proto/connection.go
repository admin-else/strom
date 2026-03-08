package proto

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync"

	"github.com/admin-else/strom/data"
	"github.com/admin-else/strom/event"
	"github.com/admin-else/strom/proto_base"
)

var BadPacketTypeError = errors.New("bad packet type")
var BadPacketIdError = errors.New("bad packet id")
var PacketNotFullyDecodedError = errors.New("packet not fully decoded")

// UnCodablePacket represents a packet that could not be decoded.
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
	RawConn
	*event.Loop
	state           proto_base.State
	stateMutex      sync.RWMutex
	Actor           proto_base.Actor
	Version         string
	ProtocolVersion int32
}

func (c *Conn) SetState(state proto_base.State) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()
	c.state = state
}

func (c *Conn) State() proto_base.State {
	c.stateMutex.RLock()
	defer c.stateMutex.RUnlock()
	return c.state
}

func (c *Conn) SetVersion(version string) (err error) {
	c.Version = version
	versionData, err := data.LookUpVersionByName(c.Version)
	if err != nil {
		return
	}
	c.ProtocolVersion = int32(versionData.Version)
	return
}

func (c *Conn) Send(packet proto_base.EncodeDecodeAble) (err error) {
	slog.Debug("send", "actor", c.Actor, "packet", fmt.Sprintf("%#v", packet))
	switch packet := packet.(type) {
	case *UnCodablePacket:
		err = c.SendRaw(packet.Data)
	default:
		i, ok := LookupPacketInfoByTypeLookupPacketInfoByTypeAndMore(packet, c.State())
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
	packet, ok := LookUpTypeByPacketInfoAndCopyType(c.Actor.ReceiveDirection(), c.State(), id, c.ProtocolVersion)
	if !ok {
		packet = &UnCodablePacket{Err: BadPacketIdError, Data: packetBytes, Direction: c.Actor.ReceiveDirection()}
		return
	}
	err = packet.Decode(b)
	if err != nil {
		packet = &UnCodablePacket{Err: err, Data: packetBytes, Direction: c.Actor.ReceiveDirection()}
		err = nil
		return
	}
	if b.Len() != 0 {
		err = nil
		packet = &UnCodablePacket{Err: PacketNotFullyDecodedError, Data: packetBytes, Direction: c.Actor.ReceiveDirection()}
		return
	}
	return
}

func (c *Conn) OnTick(_ event.Tick) (err error) {
	packet, err := c.Receive()
	if err != nil {
		return
	}
	slog.Debug("receive", "actor", c.Actor, "packet", fmt.Sprintf("%#v", packet))
	err = c.Loop.Fire(packet)
	return
}

func (c *Conn) Start(handlers ...any) (err error) {
	if c.Loop == nil {
		c.Loop = &event.Loop{}
	}
	c.Handlers = []any{c}
	c.Handlers = append(c.Handlers, handlers...)
	err = c.Loop.Start()
	return
}
