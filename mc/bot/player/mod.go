package player

import (
	"sync"
	"time"

	"github.com/admin-else/strom/mc/event"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_base"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_11"
	"github.com/admin-else/strom/mc/proto_generated/v26_4_snapshot_2"
)

const (
	// LocalPlayer.sendPosition sends a position packet once the squared move
	// distance exceeds Mth.square(2.0E-4), or every 20 ticks as a reminder.
	moveEpsilonSquared   = 4.0e-8
	positionReminderTick = 20

	// LocalPlayer.sendPosition runs once per client tick.
	tickInterval = 50 * time.Millisecond

	// ServerboundMovePlayerPacket packFlags.
	flagOnGround            = 0x01
	flagHorizontalCollision = 0x02

	// Abilities flags, net.minecraft.world.entity.player.Abilities.
	abilityFlying = 0x02
	abilityMayFly = 0x04

	// First version whose ServerboundAcceptTeleportationPacket carries the
	// accepted position, so the older teleport-id-only shape cannot be used.
	// Its CommonSettings packet also gained ParticleStatus, matching the 26.4
	// generated shape.
	firstPositionTeleportConfirmVersion = "26.4-snapshot-2"

	// Options.renderDistance IntRange(2, 32) in 26.4.
	minViewDistance = 2
	maxViewDistance = 32

	// ClientInformation.createDefault().
	clientInformationDefaultLanguage    = "en_us"
	clientInformationDefaultViewDist    = 2
	clientInformationChatVisibilityFull = 0
	// HumanoidArm.RIGHT.id, Player.DEFAULT_MAIN_HAND.
	clientInformationRightHand      = 1
	clientInformationParticleStatus = "all"
)

// ClientInformation mirrors net.minecraft.server.level.ClientInformation, the
// payload of ServerboundClientInformationPacket (generated CommonSettings).
type ClientInformation struct {
	Language             string
	ViewDistance         int
	ChatVisibility       int32
	ChatColors           bool
	ModelCustomisation   uint8
	MainHand             int32
	TextFilteringEnabled bool
	AllowsListing        bool
	ParticleStatus       string
}

// DefaultClientInformation returns ClientInformation.createDefault(). The
// vanilla client derives its settings from Options.buildPlayerInformation, but
// a headless bot has no options screen, so the record default is the baseline.
func DefaultClientInformation() (info ClientInformation) {
	return ClientInformation{
		Language:             clientInformationDefaultLanguage,
		ViewDistance:         clientInformationDefaultViewDist,
		ChatVisibility:       clientInformationChatVisibilityFull,
		ChatColors:           true,
		ModelCustomisation:   0,
		MainHand:             clientInformationRightHand,
		TextFilteringEnabled: false,
		AllowsListing:        false,
		ParticleStatus:       clientInformationParticleStatus,
	}
}

// clampViewDistance bounds a render distance to the vanilla range and the
// int8 CommonSettings field.
func clampViewDistance(distance int) (clamped int) {
	if distance < minViewDistance {
		return minViewDistance
	}
	if distance > maxViewDistance {
		return maxViewDistance
	}
	return distance
}

// commonSettingsPacket builds the ServerboundClientInformationPacket for the
// given connection version. Every version this module handles (>= 1.21.11)
// uses the CommonSettings shape; the field set is identical, so the 1.21.11
// type is converted by the connection for all versions except 26.4.
func commonSettingsPacket(version string, info ClientInformation) (packet proto_base.EncodeDecodeAble) {
	if version == firstPositionTeleportConfirmVersion {
		return &v26_4_snapshot_2.PlayToServerPacketCommonSettings{
			Locale:              info.Language,
			ViewDistance:        int8(info.ViewDistance),
			ChatFlags:           info.ChatVisibility,
			ChatColors:          info.ChatColors,
			SkinParts:           info.ModelCustomisation,
			MainHand:            info.MainHand,
			EnableTextFiltering: info.TextFilteringEnabled,
			EnableServerListing: info.AllowsListing,
			ParticleStatus:      info.ParticleStatus,
		}
	}
	return &v1_21_11.PlayToServerPacketCommonSettings{
		Locale:              info.Language,
		ViewDistance:        int8(info.ViewDistance),
		ChatFlags:           info.ChatVisibility,
		ChatColors:          info.ChatColors,
		SkinParts:           info.ModelCustomisation,
		MainHand:            info.MainHand,
		EnableTextFiltering: info.TextFilteringEnabled,
		EnableServerListing: info.AllowsListing,
		ParticleStatus:      info.ParticleStatus,
	}
}

// Module tracks the local player and drives movement for a connection.
type Module struct {
	*proto.Conn

	mu     sync.Mutex
	joined bool
	ticker *event.Timer

	entityId int32

	dimension          int32
	worldName          string
	gamemode           string
	previousGamemode   int32
	viewDistance       int32
	simulationDistance int32
	enforcesSecureChat bool

	position            Vec3
	rotation            Rotation
	deltaMovement       Vec3
	onGround            bool
	horizontalCollision bool

	// Mirrors LocalPlayer.xLast/yLast/zLast, yRotLast/xRotLast.
	lastSentPosition        Vec3
	lastSentRotation        Rotation
	lastOnGround            bool
	lastHorizontalCollision bool
	positionReminder        int

	abilitiesFlags int8
	flyingSpeed    float32
	walkingSpeed   float32

	health         float32
	food           int32
	foodSaturation float32

	clientInformation      ClientInformation
	clientInformationDirty bool
}

// NewModule creates a player module without registering handlers.
func NewModule(c *proto.Conn) (m *Module) {
	m = &Module{Conn: c}
	m.clientInformation = DefaultClientInformation()
	return
}

// Start tracks the local player on c and begins flushing movement at the
// vanilla 20 TPS cadence.
func Start(c *proto.Conn) (m *Module) {
	m = NewModule(c)
	m.registerHandlers()
	m.startTick()
	return
}

func (m *Module) registerHandlers() {
	// ClientboundLoginPacket gained onlineMode and CommonPlayerSpawnInfo dropped
	// the seed long in 26.4-snapshot-2, so the older handler shapes cannot be
	// converted from it; register both explicitly.
	m.RegisterUntil("26.2", m.onLogin, m.onRespawn)
	m.RegisterUntil("26.4-snapshot-2", m.onLogin26_4, m.onRespawn26_4)
	m.RegisterUntilLatest(m.onAbilities)
	m.RegisterUntilLatest(m.onUpdateHealth)
	m.RegisterUntilLatest(m.onKickDisconnect)
	m.RegisterUntilLatest(m.onPosition)
	m.RegisterUntilLatest(m.onPlayerRotation)
}

func (m *Module) startTick() {
	m.ticker = &event.Timer{}
	m.ticker.Every(tickInterval, m.tick)
	m.ticker.Start(m.Conn.Loop)
}

func (m *Module) onLogin(p *v1_21_11.PlayToClientPacketLogin) (err error) {
	return m.applyJoin(p.EntityId, p.WorldState.Name, p.WorldState.Dimension, p.WorldState.Gamemode, int32(p.WorldState.PreviousGamemode), p.ViewDistance, p.SimulationDistance, p.EnforcesSecureChat)
}

func (m *Module) onLogin26_4(p *v26_4_snapshot_2.PlayToClientPacketLogin) (err error) {
	return m.applyJoin(p.EntityId, p.WorldState.Name, p.WorldState.Dimension, p.WorldState.Gamemode, int32(p.WorldState.PreviousGamemode), p.ViewDistance, p.SimulationDistance, p.EnforcesSecureChat)
}

func (m *Module) applyJoin(entityId int32, worldName string, dimension int32, gamemode string, previousGamemode int32, viewDistance, simulationDistance int32, enforcesSecureChat bool) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entityId = entityId
	m.worldName = worldName
	m.dimension = dimension
	m.gamemode = gamemode
	m.previousGamemode = previousGamemode
	m.viewDistance = viewDistance
	m.simulationDistance = simulationDistance
	m.enforcesSecureChat = enforcesSecureChat
	m.joined = true
	m.clientInformationDirty = true
	m.Log.Info("player joined", "entityId", entityId, "world", worldName, "gamemode", gamemode, "viewDistance", viewDistance, "simulationDistance", simulationDistance)
	return
}

func (m *Module) onAbilities(p *v1_21_11.PlayToClientPacketAbilities) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.abilitiesFlags = p.Flags
	m.flyingSpeed = p.FlyingSpeed
	m.walkingSpeed = p.WalkingSpeed
	return
}

func (m *Module) onUpdateHealth(p *v1_21_11.PlayToClientPacketUpdateHealth) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.health = p.Health
	m.food = p.Food
	m.foodSaturation = p.FoodSaturation
	return
}

func (m *Module) onRespawn(p *v1_21_11.PlayToClientPacketRespawn) (err error) {
	return m.applyRespawn(p.WorldState.Name, p.WorldState.Dimension, p.WorldState.Gamemode, int32(p.WorldState.PreviousGamemode))
}

func (m *Module) onRespawn26_4(p *v26_4_snapshot_2.PlayToClientPacketRespawn) (err error) {
	return m.applyRespawn(p.WorldState.Name, p.WorldState.Dimension, p.WorldState.Gamemode, p.WorldState.PreviousGamemode)
}

func (m *Module) applyRespawn(worldName string, dimension int32, gamemode string, previousGamemode int32) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.worldName = worldName
	m.dimension = dimension
	m.gamemode = gamemode
	m.previousGamemode = previousGamemode
	return
}

func (m *Module) onKickDisconnect(_ *v1_21_11.PlayToClientPacketKickDisconnect) (err error) {
	m.Log.Warn("disconnected by server")
	return event.HandlerDoneErr{}
}

func (m *Module) onPosition(p *v1_21_11.PlayToClientPacketPosition) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	source := PositionMoveRotation{Position: m.position, Delta: m.deltaMovement, Rotation: m.rotation}
	change := PositionMoveRotation{Position: Vec3{p.X, p.Y, p.Z}, Delta: Vec3{p.Dx, p.Dy, p.Dz}, Rotation: Rotation{p.Yaw, p.Pitch}}
	absolute := CalculateAbsolute(source, change, p.Flags.Val)

	m.position = absolute.Position
	m.rotation = absolute.Rotation
	m.deltaMovement = absolute.Delta
	m.lastSentPosition = absolute.Position
	m.lastSentRotation = absolute.Rotation
	m.positionReminder = 0

	return m.acceptTeleportLocked(p.TeleportId)
}

func (m *Module) onPlayerRotation(p *v1_21_11.PlayToClientPacketPlayerRotation) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	yaw := p.Yaw
	pitch := p.Pitch
	if p.RelativeYaw {
		yaw += m.rotation.Yaw
	}
	if p.RelativePitch {
		pitch += m.rotation.Pitch
	}
	m.rotation = Rotation{yaw, pitch}
	m.lastSentRotation = m.rotation

	return m.Send(&v1_21_11.PlayToServerPacketLook{Yaw: yaw, Pitch: pitch})
}

// acceptTeleportLocked answers a server teleport. The connection is assumed
// locked by the caller.
func (m *Module) acceptTeleportLocked(id int32) (err error) {
	if m.Conn.Version == firstPositionTeleportConfirmVersion {
		return m.Send(&v26_4_snapshot_2.PlayToServerPacketTeleportConfirm{
			TeleportId: id,
			X:          m.position.X,
			Y:          m.position.Y,
			Z:          m.position.Z,
			YRot:       m.rotation.Yaw,
			XRot:       m.rotation.Pitch,
		})
	}
	if err = m.Send(&v1_21_11.PlayToServerPacketTeleportConfirm{TeleportId: id}); err != nil {
		return
	}
	return m.Send(&v1_21_11.PlayToServerPacketPositionLook{
		X:     m.position.X,
		Y:     m.position.Y,
		Z:     m.position.Z,
		Yaw:   m.rotation.Yaw,
		Pitch: m.rotation.Pitch,
		Flags: m.movementFlagsLocked(),
	})
}

// tick mirrors LocalPlayer.sendPosition.
func (m *Module) tick() (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.joined {
		return nil
	}

	if m.clientInformationDirty {
		if err = m.Send(commonSettingsPacket(m.Conn.Version, m.clientInformation)); err != nil {
			return
		}
		m.clientInformationDirty = false
	}

	deltaX := m.position.X - m.lastSentPosition.X
	deltaY := m.position.Y - m.lastSentPosition.Y
	deltaZ := m.position.Z - m.lastSentPosition.Z
	deltaYRot := m.rotation.Yaw - m.lastSentRotation.Yaw
	deltaXRot := m.rotation.Pitch - m.lastSentRotation.Pitch

	m.positionReminder++
	move := deltaX*deltaX+deltaY*deltaY+deltaZ*deltaZ > moveEpsilonSquared || m.positionReminder >= positionReminderTick
	rot := deltaYRot != 0 || deltaXRot != 0

	flags := m.movementFlagsLocked()
	switch {
	case move && rot:
		err = m.Send(&v1_21_11.PlayToServerPacketPositionLook{X: m.position.X, Y: m.position.Y, Z: m.position.Z, Yaw: m.rotation.Yaw, Pitch: m.rotation.Pitch, Flags: flags})
	case move:
		err = m.Send(&v1_21_11.PlayToServerPacketPosition{X: m.position.X, Y: m.position.Y, Z: m.position.Z, Flags: flags})
	case rot:
		err = m.Send(&v1_21_11.PlayToServerPacketLook{Yaw: m.rotation.Yaw, Pitch: m.rotation.Pitch, Flags: flags})
	case m.lastOnGround != m.onGround || m.lastHorizontalCollision != m.horizontalCollision:
		err = m.Send(&v1_21_11.PlayToServerPacketFlying{Flags: flags})
	}
	if err != nil {
		return
	}

	if move {
		m.lastSentPosition = m.position
		m.positionReminder = 0
	}
	if rot {
		m.lastSentRotation = m.rotation
	}
	m.lastOnGround = m.onGround
	m.lastHorizontalCollision = m.horizontalCollision

	// Minecraft.java sends ServerboundClientTickEndPacket at the end of every
	// client tick; 26.4 uses it to reset the per-tick movement gate.
	return m.Send(&v1_21_11.PlayToServerPacketTickEnd{})
}

func (m *Module) movementFlagsLocked() (flags v1_21_11.PlayToServerMovementFlags) {
	if m.onGround {
		flags.Val |= flagOnGround
	}
	if m.horizontalCollision {
		flags.Val |= flagHorizontalCollision
	}
	return
}

// SetPosition sets the client position flushed on the next tick.
func (m *Module) SetPosition(x, y, z float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.position = Vec3{x, y, z}
}

// SetRotation sets the client yaw and pitch flushed on the next tick.
func (m *Module) SetRotation(yaw, pitch float32) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rotation = Rotation{yaw, pitch}
}

// Look sets the client rotation, matching the vanilla camera look direction.
func (m *Module) Look(yaw, pitch float32) {
	m.SetRotation(yaw, pitch)
}

// SetViewDistance sets the render distance, in chunks, reported to the server
// in ServerboundClientInformationPacket. The value is clamped to the vanilla
// Options.renderDistance range and flushed on the next tick.
func (m *Module) SetViewDistance(distance int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.setViewDistanceLocked(distance)
}

func (m *Module) setViewDistanceLocked(distance int) {
	m.clientInformation.ViewDistance = clampViewDistance(distance)
	m.clientInformationDirty = true
}

// SetClientInformation replaces the client settings sent to the server and
// flushes them on the next tick.
func (m *Module) SetClientInformation(info ClientInformation) {
	m.mu.Lock()
	defer m.mu.Unlock()
	info.ViewDistance = clampViewDistance(info.ViewDistance)
	m.clientInformation = info
	m.clientInformationDirty = true
}

// ClientInformation returns the settings queued for the server.
func (m *Module) ClientInformation() (info ClientInformation) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.clientInformation
}

// ViewDistance returns the queued render distance in chunks.
func (m *Module) ViewDistance() (distance int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.clientInformation.ViewDistance
}

// Move adds a delta to the client position flushed on the next tick.
func (m *Module) Move(dx, dy, dz float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.position.X += dx
	m.position.Y += dy
	m.position.Z += dz
}

// Fly moves the client by a delta while marking the player as flying, which the
// spectator camera relies on to ignore collision.
func (m *Module) Fly(dx, dy, dz float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.abilitiesFlags |= abilityFlying | abilityMayFly
	m.position.X += dx
	m.position.Y += dy
	m.position.Z += dz
}

// Position returns the current client position.
func (m *Module) Position() (x, y, z float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.position.X, m.position.Y, m.position.Z
}

// Rotation returns the current client yaw and pitch.
func (m *Module) Rotation() (yaw, pitch float32) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.rotation.Yaw, m.rotation.Pitch
}

// EntityId returns the local player's entity id.
func (m *Module) EntityId() (id int32) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.entityId
}

// Joined reports whether the play login packet has been received.
func (m *Module) Joined() (joined bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.joined
}

// WorldName returns the current dimension identifier.
func (m *Module) WorldName() (name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.worldName
}

// Gamemode returns the current game mode name.
func (m *Module) Gamemode() (gamemode string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.gamemode
}

// IsSpectator reports whether the player is in spectator mode.
func (m *Module) IsSpectator() (ok bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.gamemode == "spectator"
}

// Abilities returns the raw ability flags.
func (m *Module) Abilities() (flags int8) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.abilitiesFlags
}

// Flying reports whether the player may fly and is currently flying.
func (m *Module) Flying() (ok bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.abilitiesFlags&abilityFlying != 0
}

// Health returns the last health, food and saturation values.
func (m *Module) Health() (health float32, food int32, saturation float32) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.health, m.food, m.foodSaturation
}
