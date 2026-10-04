package player

import (
	"bytes"
	"testing"

	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_generated/v26_4_snapshot_2"
)

// TestRegisterHandlersLatest guards the 26.4 handler shapes: SmartConversion
// panics at registration time when a packet field set diverges.
func TestRegisterHandlersLatest(t *testing.T) {
	c := proto.NewConn()
	m := NewModule(c)
	m.registerHandlers()
}

func TestCalculateAbsoluteAbsolute(t *testing.T) {
	source := PositionMoveRotation{Position: Vec3{10, 20, 30}, Rotation: Rotation{90, 10}}
	change := PositionMoveRotation{Position: Vec3{1, 2, 3}, Rotation: Rotation{5, 5}}

	got := CalculateAbsolute(source, change, 0)
	want := PositionMoveRotation{Position: Vec3{1, 2, 3}, Rotation: Rotation{5, 5}}
	if got.Position != want.Position || got.Rotation != want.Rotation {
		t.Fatalf("CalculateAbsolute = %+v, want %+v", got, want)
	}
}

func TestCalculateAbsoluteRelativePositionAndRotation(t *testing.T) {
	source := PositionMoveRotation{Position: Vec3{10, 20, 30}, Rotation: Rotation{90, 10}}
	change := PositionMoveRotation{Position: Vec3{1, 2, 3}, Rotation: Rotation{5, -5}}

	got := CalculateAbsolute(source, change, relativeX|relativeY|relativeZ|relativeYRot|relativeXRot)
	want := PositionMoveRotation{Position: Vec3{11, 22, 33}, Rotation: Rotation{95, 5}}
	if got.Position != want.Position || got.Rotation != want.Rotation {
		t.Fatalf("CalculateAbsolute = %+v, want %+v", got, want)
	}
}

func TestCalculateAbsoluteClampsPitch(t *testing.T) {
	source := PositionMoveRotation{Rotation: Rotation{Pitch: 80}}
	change := PositionMoveRotation{Rotation: Rotation{Pitch: 30}}

	got := CalculateAbsolute(source, change, relativeXRot)
	if got.Rotation.Pitch != 90 {
		t.Fatalf("pitch = %v, want 90", got.Rotation.Pitch)
	}
}

func TestMovementFlags(t *testing.T) {
	m := NewModule(proto.NewConn())
	m.onGround = true
	m.horizontalCollision = true

	got := m.movementFlagsLocked()
	if got.Val != flagOnGround|flagHorizontalCollision {
		t.Fatalf("flags = %#x, want %#x", got.Val, flagOnGround|flagHorizontalCollision)
	}
}

// TestTeleportConfirm26_4RoundTrip locks the 26.4 shape of
// ServerboundAcceptTeleportationPacket, which carries the accepted position.
func TestTeleportConfirm26_4RoundTrip(t *testing.T) {
	want := &v26_4_snapshot_2.PlayToServerPacketTeleportConfirm{
		TeleportId: 7,
		X:          1.5,
		Y:          2.5,
		Z:          3.5,
		YRot:       90,
		XRot:       45,
	}
	var buf bytes.Buffer
	if err := want.Encode(&buf); err != nil {
		t.Fatal(err)
	}
	got := &v26_4_snapshot_2.PlayToServerPacketTeleportConfirm{}
	if err := got.Decode(bytes.NewReader(buf.Bytes())); err != nil {
		t.Fatal(err)
	}
	if *got != *want {
		t.Fatalf("round trip = %+v, want %+v", got, want)
	}
}

// TestLogin26_4OnlineModeRoundTrip locks the onlineMode field added to
// ClientboundLoginPacket in 26.4-snapshot-2.
func TestLogin26_4OnlineModeRoundTrip(t *testing.T) {
	want := &v26_4_snapshot_2.PlayToClientPacketLogin{
		EntityId:            42,
		WorldNames:          []string{"minecraft:overworld"},
		MaxPlayers:          10,
		ViewDistance:        10,
		SimulationDistance:  10,
		EnableRespawnScreen: true,
		WorldState: v26_4_snapshot_2.PlayToClientSpawnInfo{
			Name:           "minecraft:overworld",
			Gamemode:       "survival",
			PortalCooldown: 0,
		},
		OnlineMode:         true,
		EnforcesSecureChat: false,
	}
	var buf bytes.Buffer
	if err := want.Encode(&buf); err != nil {
		t.Fatal(err)
	}
	got := &v26_4_snapshot_2.PlayToClientPacketLogin{}
	if err := got.Decode(bytes.NewReader(buf.Bytes())); err != nil {
		t.Fatal(err)
	}
	if !got.OnlineMode {
		t.Fatalf("OnlineMode = false, want true (decoded %+v)", got)
	}
	if got.EntityId != want.EntityId || got.WorldState.Name != want.WorldState.Name {
		t.Fatalf("round trip = %+v, want %+v", got, want)
	}
}
