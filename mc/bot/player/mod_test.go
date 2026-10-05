package player

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_11"
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

func TestDefaultClientInformation(t *testing.T) {
	got := DefaultClientInformation()
	want := ClientInformation{
		Language:       clientInformationDefaultLanguage,
		ViewDistance:   clientInformationDefaultViewDist,
		ChatVisibility: clientInformationChatVisibilityFull,
		ChatColors:     true,
		MainHand:       clientInformationRightHand,
		ParticleStatus: clientInformationParticleStatus,
	}
	if got != want {
		t.Fatalf("DefaultClientInformation = %+v, want %+v", got, want)
	}
}

func TestClampViewDistance(t *testing.T) {
	cases := []struct{ in, want int }{
		{0, minViewDistance},
		{1, minViewDistance},
		{12, 12},
		{32, 32},
		{33, maxViewDistance},
		{200, maxViewDistance},
	}
	for _, c := range cases {
		if got := clampViewDistance(c.in); got != c.want {
			t.Fatalf("clampViewDistance(%d) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestSetViewDistanceClampsAndMarksDirty(t *testing.T) {
	m := NewModule(proto.NewConn())
	if m.clientInformationDirty {
		t.Fatal("new module already dirty")
	}
	m.SetViewDistance(12)
	if got := m.ViewDistance(); got != 12 {
		t.Fatalf("ViewDistance = %d, want 12", got)
	}
	if !m.clientInformationDirty {
		t.Fatal("SetViewDistance did not mark the settings dirty")
	}

	m.clientInformationDirty = false
	m.SetViewDistance(1000)
	if got := m.ViewDistance(); got != maxViewDistance {
		t.Fatalf("ViewDistance = %d, want %d", got, maxViewDistance)
	}
}

func TestCommonSettingsPacket26_4(t *testing.T) {
	info := DefaultClientInformation()
	info.ViewDistance = 12
	packet := commonSettingsPacket(firstPositionTeleportConfirmVersion, info)

	got, ok := packet.(*v26_4_snapshot_2.PlayToServerPacketCommonSettings)
	if !ok {
		t.Fatalf("packet = %T, want *v26_4_snapshot_2.PlayToServerPacketCommonSettings", packet)
	}
	if got.Locale != info.Language || int(got.ViewDistance) != info.ViewDistance ||
		got.ChatFlags != info.ChatVisibility || got.ChatColors != info.ChatColors ||
		got.SkinParts != info.ModelCustomisation || got.MainHand != info.MainHand ||
		got.EnableTextFiltering != info.TextFilteringEnabled ||
		got.EnableServerListing != info.AllowsListing || got.ParticleStatus != info.ParticleStatus {
		t.Fatalf("packet = %+v, want fields of %+v", got, info)
	}
}

func TestCommonSettingsPacketOlderUses1_21_11(t *testing.T) {
	packet := commonSettingsPacket("26.2", DefaultClientInformation())
	if _, ok := packet.(*v1_21_11.PlayToServerPacketCommonSettings); !ok {
		t.Fatalf("packet = %T, want *v1_21_11.PlayToServerPacketCommonSettings", packet)
	}
}

// TestCommonSettingsConvertible guards the connection's SmartConversion for the
// packet this module sends on versions other than the one it was built from.
func TestCommonSettingsConvertible(t *testing.T) {
	old := reflect.TypeOf(v1_21_11.PlayToServerPacketCommonSettings{})
	newer := reflect.TypeOf(v26_4_snapshot_2.PlayToServerPacketCommonSettings{})
	if !proto.SmartConvertibleTo(old, newer) {
		t.Fatal("1.21.11 CommonSettings is not convertible to 26.4")
	}
	if !proto.SmartConvertibleTo(newer, old) {
		t.Fatal("26.4 CommonSettings is not convertible to 1.21.11")
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
