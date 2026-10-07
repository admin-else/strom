package entity

import "github.com/admin-else/strom/mc/util"

// This file ports the client-side entity movement interpolation introduced in
// 26.4: net.minecraft.world.entity.InterpolationHandler and its
// SteppedInterpolationHandler / LinearInterpolationHandler implementations,
// plus net.minecraft.core.PositionAndRotation (the Mutable variant the handlers
// inherit from). A remote entity is moved toward each server position over a
// number of client ticks; the renderer then lerps the previous tick position to
// the current one by the frame's partial tick.
//
// The reduced port has no ClientLevel, so the vanilla noCollision check in
// AbstractInterpolationHandler.adjustInterpolationPosition cannot run; deltas
// are applied unconditionally (tracked entities have no local collision).

// PositionAndRotation mirrors net.minecraft.core.PositionAndRotation.Mutable.
type PositionAndRotation struct {
	Position Vec3
	YRot     float32
	XRot     float32
}

// Set mirrors PositionAndRotation.Mutable.set.
func (p *PositionAndRotation) Set(position Vec3, yRot, xRot float32) {
	p.Position = position
	p.YRot = yRot
	p.XRot = xRot
}

// AddDelta mirrors PositionAndRotation.Mutable.addDelta.
func (p *PositionAndRotation) AddDelta(delta Vec3) {
	p.Position = p.Position.Add(delta)
}

// AddRotation mirrors PositionAndRotation.Mutable.addRotation.
func (p *PositionAndRotation) AddRotation(yRot, xRot float32) {
	p.YRot += yRot
	p.XRot += xRot
}

// Is mirrors PositionAndRotation.is.
func (p *PositionAndRotation) Is(position Vec3, yRot, xRot float32) (ret bool) {
	return p.YRot == yRot && p.XRot == xRot && p.Position == position
}

// Immutable mirrors PositionAndRotation.immutable.
func (p *PositionAndRotation) Immutable() (ret PositionAndRotation) { return *p }

// InterpolationHandler mirrors net.minecraft.world.entity.InterpolationHandler.
type InterpolationHandler interface {
	// Target mirrors target(): the active interpolation target, or ok=false
	// when no interpolation is running (Java returns null).
	Target() (target PositionAndRotation, ok bool)
	// InterpolateTo mirrors interpolateTo(PositionPath, float, float, boolean).
	// A nil position uses the entity's current position as a linear path.
	InterpolateTo(position *PositionPath, yRot, xRot float32, hasRotation bool) (ret bool)
	// Interpolate mirrors interpolate(): advance one client tick.
	Interpolate()
	// HasActiveInterpolation mirrors hasActiveInterpolation().
	HasActiveInterpolation() (ret bool)
	// Cancel mirrors cancel().
	Cancel()
}

// interpolationData is the shared mutable state of a handler; it mirrors the
// PositionAndRotation.Mutable subclass each concrete handler owns.
type interpolationData interface {
	position() Vec3
	immutable() PositionAndRotation
	is(position Vec3, yRot, xRot float32) bool
	hasActive() bool
	startInterpolating(path *PositionPath, yRot, xRot float32, interpolationSteps int)
	doInterpolate(interpolationSteps int)
	cancel()
	addDelta(delta Vec3)
	addRotation(yRot, xRot float32)
}

// NoOpInterpolationHandler mirrors InterpolationHandler.NoOpInterpolationHandler.
type NoOpInterpolationHandler struct{}

func (NoOpInterpolationHandler) Target() (PositionAndRotation, bool) {
	return PositionAndRotation{}, false
}
func (NoOpInterpolationHandler) InterpolateTo(*PositionPath, float32, float32, bool) bool {
	return false
}
func (NoOpInterpolationHandler) Interpolate()                 {}
func (NoOpInterpolationHandler) HasActiveInterpolation() bool { return false }
func (NoOpInterpolationHandler) Cancel()                      {}

// baseInterpolationHandler mirrors AbstractInterpolationHandler.
type baseInterpolationHandler struct {
	entity                  *Entity
	interpolationSteps      int
	lastPositionAndRotation PositionAndRotation
}

// target mirrors AbstractInterpolationHandler.target.
func (b *baseInterpolationHandler) target(data interpolationData) (ret PositionAndRotation, ok bool) {
	if !data.hasActive() {
		return PositionAndRotation{}, false
	}
	return data.immutable(), true
}

// interpolateTo mirrors AbstractInterpolationHandler.interpolateTo(PositionPath, float, float, boolean).
func (b *baseInterpolationHandler) interpolateTo(data interpolationData, position *PositionPath, yRot, xRot float32, hasRotation bool) (ret bool) {
	var current PositionAndRotation
	if data.hasActive() {
		current = data.immutable()
	} else {
		current = b.entity.StorePositionAndRotation()
	}
	path := position
	if path == nil {
		path = &PositionPath{EndPosition: current.Position}
	}
	toYRot, toXRot := yRot, xRot
	if !hasRotation {
		toYRot, toXRot = current.YRot, current.XRot
	}
	b.start(data, path, toYRot, toXRot)
	return true
}

// start mirrors the protected AbstractInterpolationHandler.interpolateTo.
func (b *baseInterpolationHandler) start(data interpolationData, position *PositionPath, yRot, xRot float32) {
	if b.interpolationSteps == 0 {
		b.entity.SnapTo(position.EndPosition, yRot, xRot)
		data.cancel()
		return
	}
	if !data.hasActive() || !data.is(position.EndPosition, yRot, xRot) {
		data.startInterpolating(position, yRot, xRot, b.interpolationSteps)
		b.setLastPositionAndRotation()
		b.entity.onInterpolationStart()
	}
}

// interpolate mirrors AbstractInterpolationHandler.interpolate.
func (b *baseInterpolationHandler) interpolate(data interpolationData) {
	if !data.hasActive() {
		data.cancel()
		return
	}
	b.adjustInterpolationTargetFromDeltas(data)
	data.doInterpolate(b.interpolationSteps)
	b.setLastPositionAndRotation()
}

// adjustInterpolationTargetFromDeltas mirrors
// AbstractInterpolationHandler.adjustInterpolationTargetFromDeltas. The vanilla
// noCollision gate in adjustInterpolationPosition is dropped (no ClientLevel in
// the reduced port), so the delta is always folded into the interpolation.
func (b *baseInterpolationHandler) adjustInterpolationTargetFromDeltas(data interpolationData) {
	delta := b.entity.position.Subtract(b.lastPositionAndRotation.Position)
	if delta.LengthSqr() > 1.0e-5 {
		data.addDelta(delta)
	}
	yRot := b.entity.yaw - b.lastPositionAndRotation.YRot
	xRot := b.entity.pitch - b.lastPositionAndRotation.XRot
	data.addRotation(yRot, xRot)
}

func (b *baseInterpolationHandler) setLastPositionAndRotation() {
	b.lastPositionAndRotation.Set(b.entity.position, b.entity.yaw, b.entity.pitch)
}

// SteppedInterpolationHandler mirrors net.minecraft.world.entity.SteppedInterpolationHandler.
type SteppedInterpolationHandler struct {
	baseInterpolationHandler
	data stepperData
}

// NewSteppedInterpolationHandler mirrors SteppedInterpolationHandler.create on
// the client side.
func NewSteppedInterpolationHandler(entity *Entity, interpolationSteps int) (ret *SteppedInterpolationHandler) {
	ret = &SteppedInterpolationHandler{
		baseInterpolationHandler: baseInterpolationHandler{entity: entity, interpolationSteps: interpolationSteps},
	}
	ret.data.entity = entity
	ret.data.interpolationSpeed = 1.0
	return
}

func (h *SteppedInterpolationHandler) Target() (PositionAndRotation, bool) {
	return h.baseInterpolationHandler.target(&h.data)
}
func (h *SteppedInterpolationHandler) InterpolateTo(position *PositionPath, yRot, xRot float32, hasRotation bool) bool {
	return h.baseInterpolationHandler.interpolateTo(&h.data, position, yRot, xRot, hasRotation)
}
func (h *SteppedInterpolationHandler) Interpolate()                 { h.baseInterpolationHandler.interpolate(&h.data) }
func (h *SteppedInterpolationHandler) HasActiveInterpolation() bool { return h.data.hasActive() }
func (h *SteppedInterpolationHandler) Cancel()                      { h.data.cancel() }

// steppedStep mirrors SteppedInterpolationHandler.Step.
type steppedStep struct {
	position   Vec3
	yRot, xRot float32
	tickOffset int32
}

func (s *steppedStep) addDelta(delta Vec3)            { s.position = s.position.Add(delta) }
func (s *steppedStep) addRotation(yRot, xRot float32) { s.yRot += yRot; s.xRot += xRot }

// stepperData mirrors SteppedInterpolationHandler.InterpolationData. entity is
// the handler's entity, so doInterpolate can move it, mirroring the outer
// SteppedInterpolationHandler's entity reference.
type stepperData struct {
	PositionAndRotation
	entity *Entity

	remainingSteps     []steppedStep
	lastStepPosRot     PositionAndRotation
	currentStepTicks   float32
	remainingTicks     float32
	interpolationSpeed float32
}

func (d *stepperData) position() Vec3                       { return d.Position }
func (d *stepperData) immutable() (ret PositionAndRotation) { return d.PositionAndRotation }
func (d *stepperData) is(position Vec3, yRot, xRot float32) bool {
	return d.YRot == yRot && d.XRot == xRot && d.Position == position
}
func (d *stepperData) hasActive() bool { return len(d.remainingSteps) > 0 }
func (d *stepperData) addDelta(delta Vec3) {
	d.PositionAndRotation.AddDelta(delta)
	for i := range d.remainingSteps {
		d.remainingSteps[i].addDelta(delta)
	}
	d.lastStepPosRot.AddDelta(delta)
}
func (d *stepperData) addRotation(yRot, xRot float32) {
	d.PositionAndRotation.AddRotation(yRot, xRot)
	for i := range d.remainingSteps {
		d.remainingSteps[i].addRotation(yRot, xRot)
	}
	d.lastStepPosRot.AddRotation(yRot, xRot)
}
func (d *stepperData) startInterpolating(path *PositionPath, yRot, xRot float32, interpolationSteps int) {
	if !d.hasActive() {
		d.setStartingPoint(d.entity.position, d.entity.yaw, d.entity.pitch)
	}
	endPosition := path.EndPosition
	if endPosition == d.Position {
		d.addStep(endPosition, yRot, xRot, interpolationSteps)
	} else {
		d.addSteps(path, yRot, xRot, interpolationSteps)
	}
	d.Set(endPosition, yRot, xRot)
}
func (d *stepperData) doInterpolate(interpolationSteps int) {
	target := d.getNewPositionAndRotation()
	d.entityPositionSet(target)
	d.advance(1.0, interpolationSteps)
}
func (d *stepperData) cancel() { d.reset() }

// entityPositionSet mirrors entity.setPos(target.position()) +
// entity.setRot(target.yRot(), target.xRot()) inside doInterpolate.
func (d *stepperData) entityPositionSet(target PositionAndRotation) {
	if d.entity == nil {
		return
	}
	d.entity.SetPositionAndRotation(target.Position, target.YRot, target.XRot)
}

func (d *stepperData) advance(ticks float32, interpolationSteps int) {
	targetSpeed := max(d.remainingTicks/float32(interpolationSteps), 1.0)
	d.interpolationSpeed = util.LerpFloat(1.0/float32(interpolationSteps), d.interpolationSpeed, targetSpeed)
	if ticks*d.interpolationSpeed < d.remainingTicks {
		ticks *= d.interpolationSpeed
	} else {
		ticks = d.remainingTicks
		d.interpolationSpeed = 1.0
	}
	d.currentStepTicks += ticks
	d.remainingTicks -= ticks
}

func (d *stepperData) reset() {
	d.remainingSteps = d.remainingSteps[:0]
	d.remainingTicks = 0.0
	d.interpolationSpeed = 1.0
}

func (d *stepperData) setStartingPoint(position Vec3, yRot, xRot float32) {
	d.lastStepPosRot.Set(position, yRot, xRot)
	d.currentStepTicks = 1.0
}

func (d *stepperData) addStep(position Vec3, yRot, xRot float32, interpolationSteps int) {
	d.remainingSteps = append(d.remainingSteps, steppedStep{position: position, yRot: yRot, xRot: xRot, tickOffset: int32(interpolationSteps)})
	d.remainingTicks += float32(interpolationSteps)
}

func (d *stepperData) addSteps(path *PositionPath, yRot, xRot float32, interpolationSteps int) {
	if !path.Stepped {
		d.addStep(path.EndPosition, yRot, xRot, interpolationSteps)
		return
	}
	steps := path.Steps
	if yRot == d.YRot && xRot == d.XRot {
		for _, step := range steps {
			d.addStep(step.Position, yRot, xRot, int(step.TickOffset))
		}
		return
	}
	totalInterpolationTicks := getInterpolationTicks(steps)
	offset := 0
	for _, step := range steps {
		offset += int(step.TickOffset)
		a := float32(offset) / float32(totalInterpolationTicks)
		d.addStep(step.Position, util.RotLerpFloat(a, d.YRot, yRot), util.LerpFloat(a, d.XRot, xRot), int(step.TickOffset))
	}
}

func getInterpolationTicks(steps []PositionStep) (ticks int) {
	for _, step := range steps {
		ticks += int(step.TickOffset)
	}
	return
}

// getNewPositionAndRotation mirrors
// SteppedInterpolationHandler.InterpolationData.getNewPositionAndRotation.
func (d *stepperData) getNewPositionAndRotation() (ret PositionAndRotation) {
	for len(d.remainingSteps) > 0 {
		step := d.remainingSteps[0]
		offset := step.tickOffset
		if d.currentStepTicks < float32(offset) {
			a := float64(d.currentStepTicks / float32(offset))
			return PositionAndRotation{
				Position: d.lastStepPosRot.Position.Lerp(step.position, a),
				YRot:     util.RotLerpFloat(float32(a), d.lastStepPosRot.YRot, step.yRot),
				XRot:     util.LerpFloat(float32(a), d.lastStepPosRot.XRot, step.xRot),
			}
		}
		d.currentStepTicks -= float32(offset)
		d.lastStepPosRot.Set(step.position, step.yRot, step.xRot)
		d.remainingSteps = d.remainingSteps[1:]
	}
	return d.PositionAndRotation
}

// LinearInterpolationHandler mirrors net.minecraft.world.entity.LinearInterpolationHandler.
type LinearInterpolationHandler struct {
	baseInterpolationHandler
	data linearInterpolationData
}

// NewLinearInterpolationHandler mirrors LinearInterpolationHandler.create.
func NewLinearInterpolationHandler(entity *Entity, interpolationSteps int) (ret *LinearInterpolationHandler) {
	ret = &LinearInterpolationHandler{
		baseInterpolationHandler: baseInterpolationHandler{entity: entity, interpolationSteps: interpolationSteps},
	}
	ret.data.entity = entity
	return
}

func (h *LinearInterpolationHandler) Target() (PositionAndRotation, bool) {
	return h.baseInterpolationHandler.target(&h.data)
}
func (h *LinearInterpolationHandler) InterpolateTo(position *PositionPath, yRot, xRot float32, hasRotation bool) bool {
	return h.baseInterpolationHandler.interpolateTo(&h.data, position, yRot, xRot, hasRotation)
}
func (h *LinearInterpolationHandler) Interpolate()                 { h.baseInterpolationHandler.interpolate(&h.data) }
func (h *LinearInterpolationHandler) HasActiveInterpolation() bool { return h.data.hasActive() }
func (h *LinearInterpolationHandler) Cancel()                      { h.data.cancel() }

// linearInterpolationData mirrors LinearInterpolationHandler.InterpolationData.
type linearInterpolationData struct {
	PositionAndRotation
	entity         *Entity
	remainingSteps int
}

func (d *linearInterpolationData) position() Vec3                       { return d.Position }
func (d *linearInterpolationData) immutable() (ret PositionAndRotation) { return d.PositionAndRotation }
func (d *linearInterpolationData) is(position Vec3, yRot, xRot float32) bool {
	return d.YRot == yRot && d.XRot == xRot && d.Position == position
}
func (d *linearInterpolationData) hasActive() bool { return d.remainingSteps > 0 }
func (d *linearInterpolationData) startInterpolating(path *PositionPath, yRot, xRot float32, interpolationSteps int) {
	d.Set(path.EndPosition, yRot, xRot)
	d.remainingSteps = interpolationSteps
}
func (d *linearInterpolationData) doInterpolate(interpolationSteps int) {
	if d.entity == nil {
		return
	}
	alpha := 1.0 / float64(d.remainingSteps)
	position := d.entity.position.Lerp(d.Position, alpha)
	yRot := float32(util.RotLerpDouble(alpha, float64(d.entity.yaw), float64(d.YRot)))
	xRot := float32(util.LerpDouble(alpha, float64(d.entity.pitch), float64(d.XRot)))
	d.entity.SetPositionAndRotation(position, yRot, xRot)
	d.remainingSteps--
}
func (d *linearInterpolationData) cancel()             { d.remainingSteps = 0 }
func (d *linearInterpolationData) addDelta(delta Vec3) { d.PositionAndRotation.AddDelta(delta) }
func (d *linearInterpolationData) addRotation(yRot, xRot float32) {
	d.PositionAndRotation.AddRotation(yRot, xRot)
}

// Vec3.Add mirrors net.minecraft.world.phys.Vec3.add(Vec3).
func (v Vec3) Add(other Vec3) (ret Vec3) {
	ret.X = v.X + other.X
	ret.Y = v.Y + other.Y
	ret.Z = v.Z + other.Z
	return
}

// Vec3.Subtract mirrors net.minecraft.world.phys.Vec3.subtract(Vec3).
func (v Vec3) Subtract(other Vec3) (ret Vec3) {
	ret.X = v.X - other.X
	ret.Y = v.Y - other.Y
	ret.Z = v.Z - other.Z
	return
}

// Vec3.Lerp mirrors net.minecraft.world.phys.Vec3.lerp(Vec3, double).
func (v Vec3) Lerp(other Vec3, alpha float64) (ret Vec3) {
	ret.X = util.LerpDouble(alpha, v.X, other.X)
	ret.Y = util.LerpDouble(alpha, v.Y, other.Y)
	ret.Z = util.LerpDouble(alpha, v.Z, other.Z)
	return
}

// Vec3.LengthSqr mirrors net.minecraft.world.phys.Vec3.lengthSqr().
func (v Vec3) LengthSqr() (ret float64) {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}
