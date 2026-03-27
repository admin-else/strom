package proto

import (
	"reflect"

	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/proto_generated"
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

func LookUpTypeByPacketInfo(direction proto_base.Direction, state proto_base.State, pid, version int32) (p proto_base.PacketInfo, ok bool) {
	pP, ok := pidMap[PidEtc{direction, state, pid, version}]
	if !ok {
		return
	}
	p = *pP
	return
}

func LookUpTypeByPacketInfoAndCopyType(direction proto_base.Direction, state proto_base.State, pid, version int32) (p proto_base.EncodeDecodeAble, ok bool) {
	i, ok := LookUpTypeByPacketInfo(direction, state, pid, version)
	if !ok {
		return
	}
	p = reflect.New(reflect.TypeOf(i.Type).Elem()).Interface().(proto_base.EncodeDecodeAble)
	return
}

func LookupPacketInfoByType(packet proto_base.EncodeDecodeAble) (p proto_base.PacketInfo, ok bool) {
	undecodedPacket, isUndecoable := packet.(*UnCodablePacket)
	if isUndecoable {
		p = proto_base.PacketInfo{Type: &UnCodablePacket{}, Name: "un_codable", Direction: undecodedPacket.Direction, State: proto_base.Handshaking, PacketId: -1, ProtocolVersion: -1}
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

func LookupPacketInfoByTypeAndState(packet proto_base.EncodeDecodeAble, state proto_base.State) (p proto_base.PacketInfo, ok bool) {
	for _, p := range proto_generated.Packets {
		if state == p.State && reflect.TypeOf(p.Type) == reflect.TypeOf(packet) {
			return p, true
		}
	}
	return
}

func LookupPacketInfoByNameProtocolVersionAndState(name string, version int32, state proto_base.State) (p proto_base.PacketInfo, ok bool) {
	for _, p := range proto_generated.Packets {
		if state == p.State && p.Name == name && p.ProtocolVersion == version {
			return p, true
		}
	}
	return
}
