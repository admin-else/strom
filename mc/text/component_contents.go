package text

// ComponentContents mirrors net.minecraft.network.chat.ComponentContents.
type ComponentContents interface {
	VisitStyled(output StyledContentConsumer, currentStyle Style) (value any, present bool)
	VisitContent(output ContentConsumer) (value any, present bool)
	Resolve(context *ResolutionContext, recursionDepth int) (*MutableComponent, error)
}

// BaseComponentContents provides the empty-visitor defaults of ComponentContents.
// Concrete contents embed it and override the visitors they support, mirroring
// Java's default methods.
type BaseComponentContents struct{}

// VisitStyled mirrors the ComponentContents default.
func (BaseComponentContents) VisitStyled(StyledContentConsumer, Style) (any, bool) { return nil, false }

// VisitContent mirrors the ComponentContents default.
func (BaseComponentContents) VisitContent(ContentConsumer) (any, bool) { return nil, false }

// componentContentsSelfResolve mirrors the ComponentContents default resolve, which
// returns MutableComponent.create(this). Callers pass the concrete contents.
func componentContentsSelfResolve(contents ComponentContents) *MutableComponent {
	return MutableComponentCreate(contents)
}

// MutableInt mirrors org.apache.commons.lang3.mutable.MutableInt.
type MutableInt struct {
	value int
}

// NewMutableInt mirrors new MutableInt().
func NewMutableInt() *MutableInt { return &MutableInt{} }

// IntValue mirrors MutableInt.intValue().
func (m *MutableInt) IntValue() int { return m.value }

// Increment mirrors MutableInt.increment().
func (m *MutableInt) Increment() { m.value++ }

// ResolutionContextLimitBehavior mirrors ResolutionContext.LimitBehavior.
type ResolutionContextLimitBehavior int

const (
	ResolutionContextLimitBehaviorDISCARD_REMAINING ResolutionContextLimitBehavior = iota
	ResolutionContextLimitBehaviorSTOP_PROCESSING_AND_COPY_REMAINING
)

// ResolutionContext mirrors net.minecraft.network.chat.ResolutionContext. Command
// and entity types are not ported; Source and DefaultScoreboardEntity stay opaque
// and ObjectInfoValidator takes an any (see docs/QUESTIONS.md).
type ResolutionContext struct {
	Source                  any
	DefaultScoreboardEntity any
	ObjectInfoValidator     func(description any) bool
	DepthLimit              int
	DepthLimitBehavior      ResolutionContextLimitBehavior
	ResolvedComponentLimit  int
	ResolvedComponentCount  *MutableInt
}

// ResolutionContextDEFAULT_RESOLVED_COMPONENT_LIMIT mirrors the private constant.
const ResolutionContextDEFAULT_RESOLVED_COMPONENT_LIMIT = 65536

// Validate mirrors ResolutionContext.validate.
func (c *ResolutionContext) Validate(description any) any {
	if c.ObjectInfoValidator(description) {
		return description
	}
	return nil
}

// ResolutionContextBuilder mirrors ResolutionContext.Builder.
type ResolutionContextBuilder struct {
	source                  any
	defaultScoreboardEntity any
	objectInfoValidator     func(any) bool
	depthLimit              int
	depthLimitBehavior      ResolutionContextLimitBehavior
}

// ResolutionContextBuilderNew mirrors ResolutionContext.builder().
func ResolutionContextBuilderNew() *ResolutionContextBuilder {
	return &ResolutionContextBuilder{
		objectInfoValidator: func(any) bool { return true },
		depthLimit:          100,
		depthLimitBehavior:  ResolutionContextLimitBehaviorSTOP_PROCESSING_AND_COPY_REMAINING,
	}
}

// WithSource mirrors ResolutionContext.Builder.withSource. It also mirrors the
// Java builder's defaultScoreboardEntity = source.getEntity(); since the entity
// model is not ported, callers provide the entity via WithEntityOverride.
func (b *ResolutionContextBuilder) WithSource(source any) *ResolutionContextBuilder {
	b.source = source
	return b
}

// WithEntityOverride mirrors ResolutionContext.Builder.withEntityOverride.
func (b *ResolutionContextBuilder) WithEntityOverride(defaultScoreboardEntity any) *ResolutionContextBuilder {
	b.defaultScoreboardEntity = defaultScoreboardEntity
	return b
}

// WithObjectInfoValidator mirrors ResolutionContext.Builder.withObjectInfoValidator.
func (b *ResolutionContextBuilder) WithObjectInfoValidator(validator func(any) bool) *ResolutionContextBuilder {
	b.objectInfoValidator = validator
	return b
}

// SetDepthLimit mirrors ResolutionContext.Builder.setDepthLimit.
func (b *ResolutionContextBuilder) SetDepthLimit(depthLimit int) *ResolutionContextBuilder {
	b.depthLimit = depthLimit
	return b
}

// SetDepthLimitBehavior mirrors ResolutionContext.Builder.setDepthLimitBehavior.
func (b *ResolutionContextBuilder) SetDepthLimitBehavior(behavior ResolutionContextLimitBehavior) *ResolutionContextBuilder {
	b.depthLimitBehavior = behavior
	return b
}

// Build mirrors ResolutionContext.Builder.build.
func (b *ResolutionContextBuilder) Build() *ResolutionContext {
	return &ResolutionContext{
		Source:                  b.source,
		DefaultScoreboardEntity: b.defaultScoreboardEntity,
		ObjectInfoValidator:     b.objectInfoValidator,
		DepthLimit:              b.depthLimit,
		DepthLimitBehavior:      b.depthLimitBehavior,
		ResolvedComponentLimit:  ResolutionContextDEFAULT_RESOLVED_COMPONENT_LIMIT,
		ResolvedComponentCount:  NewMutableInt(),
	}
}
