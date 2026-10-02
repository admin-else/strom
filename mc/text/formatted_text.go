package text

// Unit mirrors net.minecraft.util.Unit. Its presence as a visitor result mirrors
// Java's Optional<Unit> STOP_ITERATION sentinel.
type Unit struct{}

// FormattedTextSTOP_ITERATION mirrors FormattedText.STOP_ITERATION.
// Go cannot express Optional<T>; a visitor returns (value, present) and callers
// treat a present Unit as the stop signal.
var FormattedTextSTOP_ITERATION = Unit{}

// ContentConsumer mirrors FormattedText.ContentConsumer<T>.
type ContentConsumer func(contents string) (value any, present bool)

// StyledContentConsumer mirrors FormattedText.StyledContentConsumer<T>.
type StyledContentConsumer func(style Style, contents string) (value any, present bool)

// FormattedText mirrors net.minecraft.network.chat.FormattedText.
type FormattedText interface {
	VisitContent(output ContentConsumer) (value any, present bool)
	VisitStyled(output StyledContentConsumer, parentStyle Style) (value any, present bool)
}

// FormattedTextEMPTY mirrors FormattedText.EMPTY.
var FormattedTextEMPTY FormattedText = formattedTextEmpty{}

type formattedTextEmpty struct{}

func (formattedTextEmpty) VisitContent(ContentConsumer) (any, bool)             { return nil, false }
func (formattedTextEmpty) VisitStyled(StyledContentConsumer, Style) (any, bool) { return nil, false }

type formattedTextOf struct {
	text  string
	style Style
}

func (f formattedTextOf) VisitContent(output ContentConsumer) (any, bool) {
	return output(f.text)
}

func (f formattedTextOf) VisitStyled(output StyledContentConsumer, parentStyle Style) (any, bool) {
	return output(f.style.ApplyTo(parentStyle), f.text)
}

type formattedTextOfNoStyle struct {
	text string
}

func (f formattedTextOfNoStyle) VisitContent(output ContentConsumer) (any, bool) {
	return output(f.text)
}

func (f formattedTextOfNoStyle) VisitStyled(output StyledContentConsumer, parentStyle Style) (any, bool) {
	return output(parentStyle, f.text)
}

type formattedTextComposite struct {
	parts []FormattedText
}

func (f formattedTextComposite) VisitContent(output ContentConsumer) (any, bool) {
	for _, part := range f.parts {
		if value, present := part.VisitContent(output); present {
			return value, true
		}
	}
	return nil, false
}

func (f formattedTextComposite) VisitStyled(output StyledContentConsumer, parentStyle Style) (any, bool) {
	for _, part := range f.parts {
		if value, present := part.VisitStyled(output, parentStyle); present {
			return value, true
		}
	}
	return nil, false
}

// FormattedTextOf mirrors FormattedText.of(String).
func FormattedTextOf(text string) FormattedText {
	return formattedTextOfNoStyle{text: text}
}

// FormattedTextOfStyle mirrors FormattedText.of(String, Style).
func FormattedTextOfStyle(text string, style Style) FormattedText {
	return formattedTextOf{text: text, style: style}
}

// FormattedTextComposite mirrors FormattedText.composite(FormattedText...).
func FormattedTextComposite(parts ...FormattedText) FormattedText {
	return formattedTextComposite{parts: parts}
}

// FormattedTextCompositeList mirrors FormattedText.composite(List).
func FormattedTextCompositeList(parts []FormattedText) FormattedText {
	return formattedTextComposite{parts: parts}
}

// FormattedTextGetString mirrors FormattedText.getString().
func FormattedTextGetString(text FormattedText) string {
	builder := ""
	text.VisitContent(func(contents string) (any, bool) {
		builder += contents
		return nil, false
	})
	return builder
}
