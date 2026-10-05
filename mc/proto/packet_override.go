package proto

import "github.com/admin-else/strom/mc/proto_base"

// OverridePacketType replaces the type used to decode a packet identified by
// (direction, state, packetId, protocolVersion). It exists for hand-written
// packet layouts where the minecraft-data schema for a snapshot is stale; today
// that is the 26.4 entity movement family (VecDelta / PositionPath). It returns
// false when no packet is registered with that identity, so a caller cannot
// silently install a decoder for an unknown packet.
func OverridePacketType(direction proto_base.Direction, state proto_base.State, packetId, protocolVersion int32, name string, packetType proto_base.EncodeDecodeAble) (ok bool) {
	key := PidEtc{direction, state, packetId, protocolVersion}
	if _, exists := pidMap[key]; !exists {
		return false
	}
	pidMap[key] = &proto_base.PacketInfo{
		Type:            packetType,
		Name:            name,
		Direction:       direction,
		State:           state,
		PacketId:        packetId,
		ProtocolVersion: protocolVersion,
	}
	return true
}
