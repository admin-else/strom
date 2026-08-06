package proto

import (
	"reflect"

	"git.anygate.cloud/anygatecloud/strom/mc/proto_base"
	"git.anygate.cloud/anygatecloud/strom/mc/proto_generated"
)

type PidEtc struct {
	Direction       proto_base.Direction
	State           proto_base.State
	PacketId        int32
	ProtocolVersion int32
}

var pidMap = make(map[PidEtc]*proto_base.PacketInfo)

func init() {
	for _, pin := range proto_generated.Packets {
		pidMap[PidEtc{pin.Direction, pin.State, pin.PacketId, pin.ProtocolVersion}] = &pin
	}
}

// LookUpTypeByPacketInfo looks up a packet's type information by direction, state, packet id, and protocol version.
func LookUpTypeByPacketInfo(direction proto_base.Direction, state proto_base.State, pid, version int32) (p proto_base.PacketInfo, ok bool) {
	pP, ok := pidMap[PidEtc{direction, state, pid, version}]
	if !ok {
		p = proto_base.PacketInfo{
			Type:            nil,
			Name:            "",
			Direction:       direction,
			State:           state,
			PacketId:        pid,
			ProtocolVersion: version,
		}
		return
	}
	p = *pP
	return
}

// LookUpTypeByPacketInfoAndCopyType looks up a packet and returns a newly allocated instance of its type.
func LookUpTypeByPacketInfoAndCopyType(direction proto_base.Direction, state proto_base.State, pid, version int32) (p proto_base.EncodeDecodeAble, ok bool) {
	i, ok := LookUpTypeByPacketInfo(direction, state, pid, version)
	if !ok {
		return
	}
	p = reflect.New(reflect.TypeOf(i.Type).Elem()).Interface().(proto_base.EncodeDecodeAble)
	return
}

// LookupPacketInfoByType looks up a packet's type information by matching against the type of the given packet value.
func LookupPacketInfoByType(packet proto_base.EncodeDecodeAble) (p proto_base.PacketInfo, ok bool) {
	undecodedPacket, isUndecoable := packet.(*UnCodablePacket)
	if isUndecoable {
		p = undecodedPacket.Info
		ok = true
		return
	}

	for _, p := range proto_generated.Packets {
		if reflect.TypeOf(p.Type) == reflect.TypeOf(packet) {
			return p, true
		}
	}
	return
}

// LookupPacketInfoByTypeAndState looks up a packet's type information by type and state.
func LookupPacketInfoByTypeAndState(packet proto_base.EncodeDecodeAble, state proto_base.State) (p proto_base.PacketInfo, ok bool) {
	for _, p := range proto_generated.Packets {
		if state == p.State && reflect.TypeOf(p.Type) == reflect.TypeOf(packet) {
			return p, true
		}
	}
	return
}

// LookupPacketInfoByNameProtocolVersionAndState looks up a packet's type information by name, protocol version, and state.
func LookupPacketInfoByNameProtocolVersionAndState(name string, version int32, state proto_base.State) (p proto_base.PacketInfo, ok bool) {
	for _, p := range proto_generated.Packets {
		if state == p.State && p.Name == name && p.ProtocolVersion == version {
			return p, true
		}
	}
	return
}

// LookupPacketInfoByNameProtocolVersionStateAndDirection looks up a packet's type information by name, protocol version, state, and direction.
func LookupPacketInfoByNameProtocolVersionStateAndDirection(name string, version int32, state proto_base.State, direction proto_base.Direction) (p proto_base.PacketInfo, ok bool) {
	for _, p := range proto_generated.Packets {
		if state == p.State && p.Name == name && p.ProtocolVersion == version && p.Direction == direction {
			return p, true
		}
	}
	return
}
