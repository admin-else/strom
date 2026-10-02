package text

// FormattedCharSink mirrors net.minecraft.util.FormattedCharSink.
type FormattedCharSink func(position int, style Style, codepoint int) bool

// FormattedCharSequence mirrors net.minecraft.util.FormattedCharSequence.
type FormattedCharSequence func(output FormattedCharSink) bool

// FormattedCharSequenceEMPTY mirrors FormattedCharSequence.EMPTY.
var FormattedCharSequenceEMPTY FormattedCharSequence = func(output FormattedCharSink) bool { return true }

// FormattedCharSequenceCodepoint mirrors FormattedCharSequence.codepoint.
func FormattedCharSequenceCodepoint(codepoint int, style Style) FormattedCharSequence {
	return func(output FormattedCharSink) bool { return output(0, style, codepoint) }
}

// FormattedCharSequenceForward mirrors FormattedCharSequence.forward(String, Style).
func FormattedCharSequenceForward(plainText string, style Style) FormattedCharSequence {
	if plainText == "" {
		return FormattedCharSequenceEMPTY
	}
	return func(output FormattedCharSink) bool { return StringDecomposerIterate(plainText, style, output) }
}

// FormattedCharSequenceForwardModifier mirrors FormattedCharSequence.forward(String, Style, Int2IntFunction).
func FormattedCharSequenceForwardModifier(plainText string, style Style, modifier func(int) int) FormattedCharSequence {
	if plainText == "" {
		return FormattedCharSequenceEMPTY
	}
	return func(output FormattedCharSink) bool {
		return StringDecomposerIterate(plainText, style, FormattedCharSequenceDecorateOutput(output, modifier))
	}
}

// FormattedCharSequenceBackward mirrors FormattedCharSequence.backward(String, Style).
func FormattedCharSequenceBackward(plainText string, style Style) FormattedCharSequence {
	if plainText == "" {
		return FormattedCharSequenceEMPTY
	}
	return func(output FormattedCharSink) bool { return StringDecomposerIterateBackwards(plainText, style, output) }
}

// FormattedCharSequenceBackwardModifier mirrors FormattedCharSequence.backward(String, Style, Int2IntFunction).
func FormattedCharSequenceBackwardModifier(plainText string, style Style, modifier func(int) int) FormattedCharSequence {
	if plainText == "" {
		return FormattedCharSequenceEMPTY
	}
	return func(output FormattedCharSink) bool {
		return StringDecomposerIterateBackwards(plainText, style, FormattedCharSequenceDecorateOutput(output, modifier))
	}
}

// FormattedCharSequenceDecorateOutput mirrors FormattedCharSequence.decorateOutput.
func FormattedCharSequenceDecorateOutput(output FormattedCharSink, modifier func(int) int) FormattedCharSink {
	return func(position int, style Style, codepoint int) bool {
		return output(position, style, modifier(codepoint))
	}
}

// FormattedCharSequenceComposite mirrors FormattedCharSequence.composite().
func FormattedCharSequenceComposite() FormattedCharSequence { return FormattedCharSequenceEMPTY }

// FormattedCharSequenceCompositeList mirrors FormattedCharSequence.composite(FormattedCharSequence...).
func FormattedCharSequenceCompositeList(parts ...FormattedCharSequence) FormattedCharSequence {
	switch len(parts) {
	case 0:
		return FormattedCharSequenceEMPTY
	case 1:
		return parts[0]
	case 2:
		return FormattedCharSequenceFromPair(parts[0], parts[1])
	default:
		return FormattedCharSequenceFromList(parts)
	}
}

// FormattedCharSequenceFromPair mirrors FormattedCharSequence.fromPair.
func FormattedCharSequenceFromPair(first, second FormattedCharSequence) FormattedCharSequence {
	return func(output FormattedCharSink) bool { return first(output) && second(output) }
}

// FormattedCharSequenceFromList mirrors FormattedCharSequence.fromList.
func FormattedCharSequenceFromList(partCopy []FormattedCharSequence) FormattedCharSequence {
	return func(output FormattedCharSink) bool {
		for _, part := range partCopy {
			if !part(output) {
				return false
			}
		}
		return true
	}
}
