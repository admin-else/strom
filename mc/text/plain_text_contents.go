package text

// PlainTextContents mirrors net.minecraft.network.chat.contents.PlainTextContents.
type PlainTextContents interface {
	ComponentContents
	Text() string
}

// PlainTextContentsEMPTY mirrors PlainTextContents.EMPTY.
var PlainTextContentsEMPTY PlainTextContents = plainTextContentsEmpty{}

type plainTextContentsEmpty struct {
	BaseComponentContents
}

func (plainTextContentsEmpty) Text() string { return "" }

func (plainTextContentsEmpty) String() string { return "empty" }

func (p plainTextContentsEmpty) Resolve(context *ResolutionContext, recursionDepth int) (*MutableComponent, error) {
	return MutableComponentCreate(p), nil
}

// PlainTextContentsCreate mirrors PlainTextContents.create.
func PlainTextContentsCreate(text string) PlainTextContents {
	if text == "" {
		return PlainTextContentsEMPTY
	}
	return &PlainTextContentsLiteralContents{TextValue: text}
}

// PlainTextContentsLiteralContents mirrors PlainTextContents.LiteralContents.
type PlainTextContentsLiteralContents struct {
	BaseComponentContents
	TextValue string
}

func (p *PlainTextContentsLiteralContents) Text() string { return p.TextValue }

func (p *PlainTextContentsLiteralContents) VisitContent(output ContentConsumer) (any, bool) {
	return output(p.TextValue)
}

func (p *PlainTextContentsLiteralContents) VisitStyled(output StyledContentConsumer, currentStyle Style) (any, bool) {
	return output(currentStyle, p.TextValue)
}

func (p *PlainTextContentsLiteralContents) Resolve(context *ResolutionContext, recursionDepth int) (*MutableComponent, error) {
	return MutableComponentCreate(p), nil
}

func (p *PlainTextContentsLiteralContents) String() string {
	return "literal{" + p.TextValue + "}"
}
