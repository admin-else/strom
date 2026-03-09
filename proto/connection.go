package proto

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"sync"

	"github.com/admin-else/strom/data"
	"github.com/admin-else/strom/event"
	"github.com/admin-else/strom/proto_base"
)

var (
	BadPacketTypeErr                    = errors.New("bad packet type")
	BadPacketIdErr                      = errors.New("bad packet id")
	PacketNotFullyDecodedErr            = errors.New("packet not fully decoded")
	ContextDoneErr                      = errors.New("context done")
	WrongStateErr                       = errors.New("wrong state")
	WrongDirectionErr                   = errors.New("wrong direction")
	CantTranslateErr                    = errors.New("cant translate packet")
	PacketDoesntExistInFutureVersionErr = errors.New("packet doesnt exist in future version")
)

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

	receiveRoutine bool
	packetCh       chan proto_base.EncodeDecodeAble
}

func (c *Conn) ReceiveJob() {
	for {
		select {
		case <-c.Ctx.Done():
			return
		default:
			packet, err := c.Receive()
			if err != nil {
				c.ErrChan <- err
				return
			}
			c.packetCh <- packet
		}
	}
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
	versionData, err := data.LookUpVersionByName(version)
	if err != nil {
		return
	}
	c.ProtocolVersion = versionData.Version
	c.Version = version
	return
}

func (c *Conn) SetProtocolVersion(version int32) (err error) {
	versionData, err := data.LookUpVersionByProtocolVersion(version)
	if err != nil {
		return
	}
	c.ProtocolVersion = versionData.Version
	c.Version = versionData.MinecraftVersion
	return
}

// ActivateReceiveRoutine activates the receive-routine.
// It ensures that the receive-routine is only started once and creates a buffered channel for packet handling.
func (c *Conn) ActivateReceiveRoutine() {
	if c.receiveRoutine {
		return
	}
	c.receiveRoutine = true
	c.packetCh = make(chan proto_base.EncodeDecodeAble, 100) // i dunno how many to buffer 100 seems reasonable
	go c.ReceiveJob()
}

func (c *Conn) translatePacketVersion(packet proto_base.EncodeDecodeAble, packetInfo proto_base.PacketInfo) (convertedPacket proto_base.EncodeDecodeAble, err error) {
	packetInfoReal, found := LookupPacketInfoByNameProtocolVersionAndState(packetInfo.Name, c.ProtocolVersion, c.State())
	if !found {
		err = PacketDoesntExistInFutureVersionErr
		return
	}
	dstType := reflect.TypeOf(packetInfoReal.Type)
	if !reflect.TypeOf(packetInfo.Type).ConvertibleTo(dstType) {
		err = CantTranslateErr
		return
	}
	// can't fail cause pacetinfo.type is proto_base.EncodeDecodeAble
	convertedPacket = reflect.ValueOf(packet).Convert(dstType).Interface().(proto_base.EncodeDecodeAble)

	return
}

func (c *Conn) sendRegisteredPacket(packet proto_base.EncodeDecodeAble) (err error) {
	i, ok := LookupPacketInfoByTypeAndState(packet, c.State())
	if !ok {
		err = BadPacketTypeErr
		return
	}
	if i.Direction != c.Actor.SendDirection() {
		err = WrongDirectionErr
		return
	}
	if i.State != c.State() {
		err = WrongStateErr
		return
	}
	if i.ProtocolVersion != c.ProtocolVersion {
		packet, err = c.translatePacketVersion(packet, i)
		if err != nil {
			return
		}
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
	slog.Debug("send", "actor", c.Actor, "packet", fmt.Sprintf("%#v", packet))
	return c.SendRaw(packetBuff.Bytes())
}

func (c *Conn) Send(packet proto_base.EncodeDecodeAble) (err error) {
	switch packet := packet.(type) {
	case *UnCodablePacket:
		err = c.SendRaw(packet.Data)
	default:
		return c.sendRegisteredPacket(packet)
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
		packet = &UnCodablePacket{Err: BadPacketIdErr, Data: packetBytes, Direction: c.Actor.ReceiveDirection()}
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
		packet = &UnCodablePacket{Err: PacketNotFullyDecodedErr, Data: packetBytes, Direction: c.Actor.ReceiveDirection()}
		return
	}
	return
}

func (c *Conn) OnTick(_ event.Tick) (err error) {
	var packet proto_base.EncodeDecodeAble
	if c.receiveRoutine {
		select {
		case packet = <-c.packetCh:
			slog.Debug("receive", "actor", c.Actor, "packet", fmt.Sprintf("%#v", packet))
			err = c.Loop.Fire(packet)
		case <-c.Ctx.Done():
			err = ContextDoneErr
			return
		default:
		}
	} else {
		packet, err = c.Receive()
		if err != nil {
			return
		}
		slog.Debug("receive", "actor", c.Actor, "packet", fmt.Sprintf("%#v", packet))
		err = c.Loop.Fire(packet)
	}
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
