package proto

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/admin-else/strom/data"
	"github.com/admin-else/strom/event"
	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/proto_generated"
)

const (
	BufferNPackets = 100
)

var (
	BadPacketTypeErr                    = errors.New("bad packet type")
	BadPacketIdErr                      = errors.New("bad packet id")
	PacketNotFullyDecodedErr            = errors.New("packet not fully decoded")
	WrongStateErr                       = errors.New("wrong state")
	WrongDirectionErr                   = errors.New("wrong direction")
	CantTranslateErr                    = errors.New("cant translate packet")
	PacketDoesntExistInFutureVersionErr = errors.New("packet doesnt exist in future version")
	NoHandlerRegisteredErr              = errors.New("no handler registered")
)

// LoginTick is used as a subsite for the old event system because it is very useful in login contexts
type LoginTick struct{}

var LatestVersion = proto_generated.SupportedVersions[len(proto_generated.SupportedVersions)-1]
var EarliestVersion = proto_generated.SupportedVersions[len(proto_generated.SupportedVersions)-1]

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
	state             proto_base.State
	stateMutex        sync.RWMutex
	Actor             proto_base.Actor
	Version           string
	ProtocolVersion   int32
	DebugPrintPackets []string

	asyncPackets      bool
	asyncPacketsMutex sync.RWMutex
}

func (c *Conn) StartConn() (err error) {
	c.RegisterEventSource(c.ActivateReceiveRoutine())
	c.RegisterCritical(c.OnLoginTick)
	err = c.Loop.StartLoop()
	c.Loop = event.NewLoop()
	return
}

func (c *Conn) OnLoginTick(_ LoginTick) (err error) {
	packet, err := c.Receive()
	if err != nil {
		return
	}
	err = c.Fire(packet)
	return
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

func (c *Conn) SetAsyncPackets(asyncPackets bool) {
	c.asyncPacketsMutex.Lock()
	defer c.asyncPacketsMutex.Unlock()
	c.asyncPackets = asyncPackets
}

func (c *Conn) GetAsyncPackets() bool {
	c.asyncPacketsMutex.RLock()
	defer c.asyncPacketsMutex.RUnlock()
	return c.asyncPackets
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
func (c *Conn) ActivateReceiveRoutine() (chSending <-chan any) {
	// This expects you to start a new routine when switching into play state
	ch := make(chan any, BufferNPackets)
	go func() {
		ctx := c.Ctx // this prevents c.Ctx being reassigned
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if !c.GetAsyncPackets() {
					ch <- LoginTick{}
				} else {
					packet, err := c.Receive()
					if err != nil {
						ch <- err
						return
					}
					ch <- packet
				}
			}
		}
	}()
	chSending = ch
	return
}

func (c *Conn) translatePacketVersion(packet proto_base.EncodeDecodeAble, packetInfo proto_base.PacketInfo) (convertedPacket proto_base.EncodeDecodeAble, packetInfoReal proto_base.PacketInfo, err error) {
	packetInfoReal, found := LookupPacketInfoByNameProtocolVersionAndState(packetInfo.Name, c.ProtocolVersion, c.State())
	if !found {
		err = PacketDoesntExistInFutureVersionErr
		return
	}

	dstType := reflect.TypeOf(packetInfoReal.Type)
	if !SmartConvertibleTo(reflect.TypeOf(packetInfo.Type), dstType) {
		err = CantTranslateErr
		return
	}
	// can't fail cause pacetinfo.type is proto_base.EncodeDecodeAble
	convertedPacket = SmartConvert(reflect.ValueOf(packet), dstType).Interface().(proto_base.EncodeDecodeAble)
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
		packet, i, err = c.translatePacketVersion(packet, i)
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
	debugPrint := slog.Default().Enabled(context.Background(), slog.LevelDebug) &&
		slices.ContainsFunc(c.DebugPrintPackets, func(s string) bool { return strings.Contains(i.Name, s) })
	if debugPrint {
		slog.Debug("send", "actor", c.Actor, "packet", fmt.Sprintf("%#v", packet))
	}
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
	if c.R == nil {
		return
	}
	packetBytes, err := c.ReceiveRaw()
	if err != nil {
		return
	}
	b := bytes.NewBuffer(packetBytes)
	id, err := proto_base.DecodeVarInt(b)
	if err != nil {
		return
	}

	i, ok := LookUpTypeByPacketInfo(c.Actor.ReceiveDirection(), c.State(), id, c.ProtocolVersion)
	if !ok {
		packet = &UnCodablePacket{Err: BadPacketIdErr, Data: packetBytes, Direction: c.Actor.ReceiveDirection()}
		return
	}
	packet = reflect.New(reflect.TypeOf(i.Type).Elem()).Interface().(proto_base.EncodeDecodeAble)
	debugPrint := slog.Default().Enabled(context.Background(), slog.LevelDebug) &&
		slices.ContainsFunc(c.DebugPrintPackets, func(s string) bool { return strings.Contains(i.Name, s) })
	isListenerRegistered := c.IsTypeRegistered(reflect.TypeOf(packet))
	if !isListenerRegistered && !debugPrint {
		packet = &UnCodablePacket{Err: NoHandlerRegisteredErr, Data: packetBytes, Direction: c.Actor.ReceiveDirection()}
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
	if debugPrint {
		slog.Debug("recv", "packet", fmt.Sprintf("%#v", packet), "isListenerRegistered", isListenerRegistered)
	}
	return
}

func NewConn() (ret *Conn) {
	ret = &Conn{}
	ret.SetState(proto_base.Handshaking)
	ret.SetCompressionThreshold(-1)
	err := ret.SetVersion(EarliestVersion)
	if err != nil {
		panic(err)
	}
	ret.Loop = event.NewLoop()
	return
}
