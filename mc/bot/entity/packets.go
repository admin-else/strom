package entity

import (
	"encoding/binary"
	"errors"
	"io"
	"math"

	"github.com/admin-else/strom/mc/proto_base"
)

// This file holds hand-written decoders for the 26.4 entity movement family.
// The minecraft-data schema for 26.4-snapshot-2 is derived from 26.2, where
// ClientboundMoveEntityPacket still carried raw short deltas; 26.4 replaced
// those with a properties varint plus VecDelta (and ClientboundEntityPositionSync
// with a PositionPath). The generated types are therefore stale, so they are
// overridden at runtime by Mod.Start. Each type mirrors the Java packet class of
// the same name (see net/minecraft/network/protocol/game).

// Vec3 is an unsigned world position / velocity triple.
type Vec3 struct {
	X, Y, Z float64
}

// PositionStep mirrors net.minecraft.world.entity.PositionStep.
type PositionStep struct {
	Position   Vec3
	TickOffset int32
}

// PositionPath mirrors net.minecraft.world.entity.PositionPath. Stepped paths
// carry their intermediate steps; EndPosition is the final position.
type PositionPath struct {
	Stepped     bool
	EndPosition Vec3
	Steps       []PositionStep
}

// PositionPathTypeLINEAR / STEPPED mirror PositionPath.Type ids.
const (
	positionPathTypeLinear  = 0
	positionPathTypeStepped = 1
)

// PositionMoveRotation mirrors net.minecraft.world.entity.PositionMoveRotation.
type PositionMoveRotation struct {
	Position      Vec3
	DeltaMovement Vec3
	YRot          float32
	XRot          float32
}

// Relative movement flags of Relative.SET_STREAM_CODEC.
const (
	RelativeX      = 1 << 0
	RelativeY      = 1 << 1
	RelativeZ      = 1 << 2
	RelativeYRot   = 1 << 3
	RelativeXRot   = 1 << 4
	RelativeDeltaX = 1 << 5
	RelativeDeltaY = 1 << 6
	RelativeDeltaZ = 1 << 7
)

// DeltaStep mirrors VecDelta.Stepped.DeltaStep: a short delta and the ticks it
// spans.
type DeltaStep struct {
	Ticks      int32
	Xa, Ya, Za int16
}

// VecDelta mirrors net.minecraft.network.protocol.game.VecDelta: either a
// single linear short delta or a sequence of stepped deltas.
type VecDelta struct {
	Stepped    bool
	Xa, Ya, Za int16
	Steps      []DeltaStep
}

func readI8(r io.Reader) (ret int8, err error) {
	err = binary.Read(r, binary.BigEndian, &ret)
	return
}

func readI16(r io.Reader) (ret int16, err error) {
	err = binary.Read(r, binary.BigEndian, &ret)
	return
}

func readI32(r io.Reader) (ret int32, err error) {
	err = binary.Read(r, binary.BigEndian, &ret)
	return
}

func readF32(r io.Reader) (ret float32, err error) {
	var bits uint32
	err = binary.Read(r, binary.BigEndian, &bits)
	ret = math.Float32frombits(bits)
	return
}

func readF64(r io.Reader) (ret float64, err error) {
	var bits uint64
	err = binary.Read(r, binary.BigEndian, &bits)
	ret = math.Float64frombits(bits)
	return
}

func readVec3(r io.Reader) (ret Vec3, err error) {
	if ret.X, err = readF64(r); err != nil {
		return
	}
	if ret.Y, err = readF64(r); err != nil {
		return
	}
	ret.Z, err = readF64(r)
	return
}

func readVecDelta(r io.ReadSeeker, stepCount int32) (ret VecDelta, err error) {
	if stepCount <= 0 {
		ret.Stepped = false
		if ret.Xa, err = readI16(r); err != nil {
			return
		}
		if ret.Ya, err = readI16(r); err != nil {
			return
		}
		ret.Za, err = readI16(r)
		return
	}
	ret.Stepped = true
	ret.Steps = make([]DeltaStep, 0, stepCount)
	for range stepCount {
		var step DeltaStep
		if step.Ticks, err = proto_base.DecodeVarInt(r); err != nil {
			return
		}
		if step.Xa, err = readI16(r); err != nil {
			return
		}
		if step.Ya, err = readI16(r); err != nil {
			return
		}
		if step.Za, err = readI16(r); err != nil {
			return
		}
		ret.Steps = append(ret.Steps, step)
	}
	return
}

func readPositionPath(r io.ReadSeeker) (ret PositionPath, err error) {
	pathType, err := proto_base.DecodeVarInt(r)
	if err != nil {
		return
	}
	switch pathType {
	case positionPathTypeLinear:
		ret.Stepped = false
		ret.EndPosition, err = readVec3(r)
		return
	case positionPathTypeStepped:
		ret.Stepped = true
		var count int32
		count, err = proto_base.DecodeVarInt(r)
		if err != nil {
			return
		}
		ret.Steps = make([]PositionStep, 0, count)
		for range count {
			var step PositionStep
			if step.Position, err = readVec3(r); err != nil {
				return
			}
			if step.TickOffset, err = proto_base.DecodeVarInt(r); err != nil {
				return
			}
			ret.Steps = append(ret.Steps, step)
		}
		if len(ret.Steps) == 0 {
			err = errors.New("entity: empty stepped PositionPath")
			return
		}
		ret.EndPosition = ret.Steps[len(ret.Steps)-1].Position
		return
	}
	err = errors.New("entity: unknown PositionPath type")
	return
}

func readPositionMoveRotation(r io.Reader) (ret PositionMoveRotation, err error) {
	if ret.Position, err = readVec3(r); err != nil {
		return
	}
	if ret.DeltaMovement, err = readVec3(r); err != nil {
		return
	}
	if ret.YRot, err = readF32(r); err != nil {
		return
	}
	ret.XRot, err = readF32(r)
	return
}

// PlayToClientPacketMoveEntityPos mirrors ClientboundMoveEntityPacket.Pos.
type PlayToClientPacketMoveEntityPos struct {
	EntityId int32
	Delta    VecDelta
	OnGround bool
}

func (p *PlayToClientPacketMoveEntityPos) Decode(r io.ReadSeeker) (err error) {
	if p.EntityId, err = proto_base.DecodeVarInt(r); err != nil {
		return
	}
	var properties int32
	if properties, err = proto_base.DecodeVarInt(r); err != nil {
		return
	}
	p.OnGround = properties&1 != 0
	p.Delta, err = readVecDelta(r, int32(uint32(properties)>>1))
	return
}

func (p *PlayToClientPacketMoveEntityPos) Encode(w io.Writer) (err error) {
	err = errors.New("entity: encoding ClientboundMoveEntityPacket.Pos is not supported")
	return
}

// PlayToClientPacketMoveEntityPosRot mirrors ClientboundMoveEntityPacket.PosRot.
type PlayToClientPacketMoveEntityPosRot struct {
	EntityId int32
	Delta    VecDelta
	YRot     int8
	XRot     int8
	OnGround bool
}

func (p *PlayToClientPacketMoveEntityPosRot) Decode(r io.ReadSeeker) (err error) {
	if p.EntityId, err = proto_base.DecodeVarInt(r); err != nil {
		return
	}
	var properties int32
	if properties, err = proto_base.DecodeVarInt(r); err != nil {
		return
	}
	p.OnGround = properties&1 != 0
	if p.Delta, err = readVecDelta(r, int32(uint32(properties)>>1)); err != nil {
		return
	}
	if p.YRot, err = readI8(r); err != nil {
		return
	}
	p.XRot, err = readI8(r)
	return
}

func (p *PlayToClientPacketMoveEntityPosRot) Encode(w io.Writer) (err error) {
	err = errors.New("entity: encoding ClientboundMoveEntityPacket.PosRot is not supported")
	return
}

// PlayToClientPacketMoveEntityRot mirrors ClientboundMoveEntityPacket.Rot.
type PlayToClientPacketMoveEntityRot struct {
	EntityId int32
	OnGround bool
	YRot     int8
	XRot     int8
}

func (p *PlayToClientPacketMoveEntityRot) Decode(r io.ReadSeeker) (err error) {
	if p.EntityId, err = proto_base.DecodeVarInt(r); err != nil {
		return
	}
	if p.OnGround, err = readBool(r); err != nil {
		return
	}
	if p.YRot, err = readI8(r); err != nil {
		return
	}
	p.XRot, err = readI8(r)
	return
}

func (p *PlayToClientPacketMoveEntityRot) Encode(w io.Writer) (err error) {
	err = errors.New("entity: encoding ClientboundMoveEntityPacket.Rot is not supported")
	return
}

// PlayToClientPacketEntityPositionSync mirrors ClientboundEntityPositionSyncPacket.
type PlayToClientPacketEntityPositionSync struct {
	EntityId int32
	Position PositionPath
	YRot     float32
	XRot     float32
	OnGround bool
}

func (p *PlayToClientPacketEntityPositionSync) Decode(r io.ReadSeeker) (err error) {
	if p.EntityId, err = proto_base.DecodeVarInt(r); err != nil {
		return
	}
	if p.Position, err = readPositionPath(r); err != nil {
		return
	}
	if p.YRot, err = readF32(r); err != nil {
		return
	}
	if p.XRot, err = readF32(r); err != nil {
		return
	}
	p.OnGround, err = readBool(r)
	return
}

func (p *PlayToClientPacketEntityPositionSync) Encode(w io.Writer) (err error) {
	err = errors.New("entity: encoding ClientboundEntityPositionSyncPacket is not supported")
	return
}

// PlayToClientPacketTeleportEntity mirrors ClientboundTeleportEntityPacket.
type PlayToClientPacketTeleportEntity struct {
	EntityId  int32
	Change    PositionMoveRotation
	Relatives int32
	OnGround  bool
}

func (p *PlayToClientPacketTeleportEntity) Decode(r io.ReadSeeker) (err error) {
	if p.EntityId, err = proto_base.DecodeVarInt(r); err != nil {
		return
	}
	if p.Change, err = readPositionMoveRotation(r); err != nil {
		return
	}
	if p.Relatives, err = readI32(r); err != nil {
		return
	}
	p.OnGround, err = readBool(r)
	return
}

func (p *PlayToClientPacketTeleportEntity) Encode(w io.Writer) (err error) {
	err = errors.New("entity: encoding ClientboundTeleportEntityPacket is not supported")
	return
}

func readBool(r io.Reader) (ret bool, err error) {
	var b uint8
	err = binary.Read(r, binary.BigEndian, &b)
	ret = b != 0
	return
}
