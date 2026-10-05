package entity

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"

	"github.com/admin-else/strom/mc/proto_base"
)

func putI16(t *testing.T, b *bytes.Buffer, value int16) {
	t.Helper()
	if err := binary.Write(b, binary.BigEndian, value); err != nil {
		t.Fatal(err)
	}
}

func putI8(t *testing.T, b *bytes.Buffer, value int8) {
	t.Helper()
	if err := binary.Write(b, binary.BigEndian, value); err != nil {
		t.Fatal(err)
	}
}

func putI32(t *testing.T, b *bytes.Buffer, value int32) {
	t.Helper()
	if err := binary.Write(b, binary.BigEndian, value); err != nil {
		t.Fatal(err)
	}
}

func putF32(t *testing.T, b *bytes.Buffer, value float32) {
	t.Helper()
	if err := binary.Write(b, binary.BigEndian, math.Float32bits(value)); err != nil {
		t.Fatal(err)
	}
}

func putF64(t *testing.T, b *bytes.Buffer, value float64) {
	t.Helper()
	if err := binary.Write(b, binary.BigEndian, math.Float64bits(value)); err != nil {
		t.Fatal(err)
	}
}

func putVarInt(t *testing.T, b *bytes.Buffer, value int32) {
	t.Helper()
	if err := proto_base.EncodeVarInt(b, value); err != nil {
		t.Fatal(err)
	}
}

func TestMoveEntityPosLinear(t *testing.T) {
	var b bytes.Buffer
	putVarInt(t, &b, 7)
	putVarInt(t, &b, 1) // stepCount 0, onGround true
	putI16(t, &b, 4096)
	putI16(t, &b, -4096)
	putI16(t, &b, 0)
	var packet PlayToClientPacketMoveEntityPos
	if err := packet.Decode(bytes.NewReader(b.Bytes())); err != nil {
		t.Fatal(err)
	}
	if packet.EntityId != 7 || !packet.OnGround || packet.Delta.Stepped {
		t.Fatalf("unexpected packet %+v", packet)
	}
	entity := &Entity{Id: 7, Position: Vec3{}}
	path := entity.codec.decodePath(packet.Delta)
	if path.EndPosition.X != 1.0 || path.EndPosition.Y != -1.0 || path.EndPosition.Z != 0 {
		t.Fatalf("unexpected end %+v", path.EndPosition)
	}
}

func TestMoveEntityPosRotSteppedAndRotation(t *testing.T) {
	var b bytes.Buffer
	putVarInt(t, &b, 9)
	putVarInt(t, &b, 4) // stepCount 2, onGround false
	for range 2 {
		putVarInt(t, &b, 1)
		putI16(t, &b, 4096)
		putI16(t, &b, 0)
		putI16(t, &b, 0)
	}
	putI8(t, &b, 64)
	putI8(t, &b, -32)
	var packet PlayToClientPacketMoveEntityPosRot
	if err := packet.Decode(bytes.NewReader(b.Bytes())); err != nil {
		t.Fatal(err)
	}
	if len(packet.Delta.Steps) != 2 || packet.YRot != 64 || packet.XRot != -32 || packet.OnGround {
		t.Fatalf("unexpected packet %+v", packet)
	}
	entity := &Entity{Id: 9}
	path := entity.codec.decodePath(packet.Delta)
	if path.EndPosition.X != 2.0 || path.EndPosition.Y != 0 || path.EndPosition.Z != 0 {
		t.Fatalf("unexpected end %+v", path.EndPosition)
	}
	if got := unpackDegrees(packet.YRot); got != 90 {
		t.Fatalf("yaw = %v, want 90", got)
	}
}

func TestMoveEntityRotWireOrder(t *testing.T) {
	// entityId, onGround, yRot, xRot (the 26.4 Pos/Rot order puts onGround first).
	var b bytes.Buffer
	putVarInt(t, &b, 5)
	putI8(t, &b, 1)
	putI8(t, &b, 32)
	putI8(t, &b, 96)
	var packet PlayToClientPacketMoveEntityRot
	if err := packet.Decode(bytes.NewReader(b.Bytes())); err != nil {
		t.Fatal(err)
	}
	if !packet.OnGround || packet.YRot != 32 || packet.XRot != 96 {
		t.Fatalf("unexpected packet %+v", packet)
	}
}

func TestEntityPositionSyncLinear(t *testing.T) {
	var b bytes.Buffer
	putVarInt(t, &b, 11)
	putVarInt(t, &b, positionPathTypeLinear)
	putF64(t, &b, 1.5)
	putF64(t, &b, 64.0)
	putF64(t, &b, -2.5)
	putF32(t, &b, 45)
	putF32(t, &b, 10)
	putI8(t, &b, 1)
	var packet PlayToClientPacketEntityPositionSync
	if err := packet.Decode(bytes.NewReader(b.Bytes())); err != nil {
		t.Fatal(err)
	}
	if packet.Position.Stepped || packet.Position.EndPosition != (Vec3{1.5, 64, -2.5}) {
		t.Fatalf("unexpected position %+v", packet.Position)
	}
	if packet.YRot != 45 || packet.XRot != 10 || !packet.OnGround {
		t.Fatalf("unexpected packet %+v", packet)
	}
}

func TestEntityPositionSyncStepped(t *testing.T) {
	var b bytes.Buffer
	putVarInt(t, &b, 12)
	putVarInt(t, &b, positionPathTypeStepped)
	putVarInt(t, &b, 2)
	putF64(t, &b, 1.0)
	putF64(t, &b, 0)
	putF64(t, &b, 0)
	putVarInt(t, &b, 1)
	putF64(t, &b, 2.0)
	putF64(t, &b, 0)
	putF64(t, &b, 0)
	putVarInt(t, &b, 1)
	putF32(t, &b, 0)
	putF32(t, &b, 0)
	putI8(t, &b, 0)
	var packet PlayToClientPacketEntityPositionSync
	if err := packet.Decode(bytes.NewReader(b.Bytes())); err != nil {
		t.Fatal(err)
	}
	if !packet.Position.Stepped || len(packet.Position.Steps) != 2 || packet.Position.EndPosition.X != 2.0 {
		t.Fatalf("unexpected position %+v", packet.Position)
	}
}

func TestTeleportEntityDecode(t *testing.T) {
	var b bytes.Buffer
	putVarInt(t, &b, 13)
	// PositionMoveRotation: position, deltaMovement, yRot, xRot.
	for _, value := range []float64{10, 20, 30, 1, 2, 3} {
		putF64(t, &b, value)
	}
	putF32(t, &b, 90)
	putF32(t, &b, 15)
	putI32(t, &b, RelativeX|RelativeYRot) // 1 | 8
	putI8(t, &b, 1)
	var packet PlayToClientPacketTeleportEntity
	if err := packet.Decode(bytes.NewReader(b.Bytes())); err != nil {
		t.Fatal(err)
	}
	if packet.Change.Position != (Vec3{10, 20, 30}) || packet.Relatives != RelativeX|RelativeYRot || !packet.OnGround {
		t.Fatalf("unexpected packet %+v", packet)
	}
	absolute := calculateAbsolute(PositionMoveRotation{
		Position: Vec3{100, 200, 300},
		YRot:     45,
	}, packet.Change, packet.Relatives)
	if absolute.Position != (Vec3{110, 20, 30}) {
		t.Fatalf("position = %+v, want {110 20 30}", absolute.Position)
	}
	if absolute.YRot != 135 {
		t.Fatalf("yaw = %v, want 135", absolute.YRot)
	}
	if absolute.XRot != 15 {
		t.Fatalf("pitch = %v, want 15", absolute.XRot)
	}
}

func TestVecDeltaCodecRounding(t *testing.T) {
	codec := &VecDeltaCodec{}
	codec.SetBase(Vec3{X: 0.5})
	got := codec.Decode(4096, 0, 0)
	if got.X != 1.5 {
		t.Fatalf("x = %v, want 1.5", got.X)
	}
}
