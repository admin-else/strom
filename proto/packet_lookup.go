package proto

import (
	"reflect"

	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/proto_generated"
)

// I know this is in-performant and i dont care until it becomes a problem for somebody it will stay like this and it is pretty easy to make caches or simular to speed this up so yeah

func LookUpTypeByPacketInfo(direction proto_base.Direction, state proto_base.State, pid, version int32) (p proto_base.PacketInfo, ok bool) {
	for _, pin := range proto_generated.Packets {
		if pin.Direction != direction || pin.State != state || pin.ProtocolVersion != version || pin.PacketId != pid {
			continue
		}
		return pin, true
	}
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
	switch packet := packet.(type) {
	case *UnCodablePacket:
		p = proto_base.PacketInfo{
			Type:            &UnCodablePacket{},
			Name:            "undecodeable",
			Direction:       packet.Direction,
			State:           -1,
			PacketId:        -1,
			ProtocolVersion: -1,
		}
		ok = true
	default:
		for _, p := range proto_generated.Packets {
			if reflect.TypeOf(p.Type) == reflect.TypeOf(packet) {
				return p, true
			}
		}
	}
	return
}

func LookupPacketInfoByTypeLookupPacketInfoByTypeAndMore(packet proto_base.EncodeDecodeAble, state proto_base.State) (p proto_base.PacketInfo, ok bool) {
	switch packet := packet.(type) {
	case *UnCodablePacket:
		p = proto_base.PacketInfo{
			Type:            &UnCodablePacket{},
			Name:            "undecodeable",
			Direction:       packet.Direction,
			State:           -1,
			PacketId:        -1,
			ProtocolVersion: -1,
		}
	default:
		for _, p := range proto_generated.Packets {
			if state == p.State && reflect.TypeOf(p.Type) == reflect.TypeOf(packet) {
				return p, true
			}
		}
	}
	return
}
