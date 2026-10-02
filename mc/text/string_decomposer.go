package text

// This file mirrors net.minecraft.util.StringDecomposer. It lives in the text
// package (rather than mc/util) because FormattedCharSequence.forward/backward
// call it and mc/util would import mc/text, creating an import cycle. Recorded
// in docs/QUESTIONS.md.

import (
	"strings"
	"unicode/utf16"
)

const stringDecomposerReplacementChar = '\uFFFD'

func isHighSurrogate(ch uint16) bool { return ch >= 0xD800 && ch <= 0xDBFF }
func isLowSurrogate(ch uint16) bool  { return ch >= 0xDC00 && ch <= 0xDFFF }

func toCodePoint(high, low uint16) int {
	return 0x10000 + (int(high)-0xD800)<<10 + (int(low) - 0xDC00)
}

func feedChar(style Style, output FormattedCharSink, pos int, ch uint16) bool {
	if isHighSurrogate(ch) || isLowSurrogate(ch) {
		return output(pos, style, 65533)
	}
	return output(pos, style, int(ch))
}

func toUTF16(s string) []uint16 { return utf16.Encode([]rune(s)) }

// StringDecomposerIterate mirrors StringDecomposer.iterate.
func StringDecomposerIterate(s string, style Style, output FormattedCharSink) bool {
	units := toUTF16(s)
	size := len(units)

	for i := 0; i < size; i++ {
		ch := units[i]
		if isHighSurrogate(ch) {
			if i+1 >= size {
				if !output(i, style, 65533) {
					return false
				}
				break
			}
			low := units[i+1]
			if isLowSurrogate(low) {
				if !output(i, style, toCodePoint(ch, low)) {
					return false
				}
				i++
			} else if !output(i, style, 65533) {
				return false
			}
		} else if !feedChar(style, output, i, ch) {
			return false
		}
	}

	return true
}

// StringDecomposerIterateBackwards mirrors StringDecomposer.iterateBackwards.
func StringDecomposerIterateBackwards(s string, style Style, output FormattedCharSink) bool {
	units := toUTF16(s)
	size := len(units)

	for i := size - 1; i >= 0; i-- {
		ch := units[i]
		if isLowSurrogate(ch) {
			if i-1 < 0 {
				if !output(0, style, 65533) {
					return false
				}
				break
			}
			high := units[i-1]
			if isHighSurrogate(high) {
				i--
				if !output(i, style, toCodePoint(high, ch)) {
					return false
				}
			} else if !output(i, style, 65533) {
				return false
			}
		} else if !feedChar(style, output, i, ch) {
			return false
		}
	}

	return true
}

// StringDecomposerIterateFormatted mirrors StringDecomposer.iterateFormatted(String, Style, FormattedCharSink).
func StringDecomposerIterateFormatted(s string, style Style, output FormattedCharSink) bool {
	return StringDecomposerIterateFormattedOffset(s, 0, style, output)
}

// StringDecomposerIterateFormattedOffset mirrors StringDecomposer.iterateFormatted(String, int, Style, FormattedCharSink).
func StringDecomposerIterateFormattedOffset(s string, offset int, style Style, output FormattedCharSink) bool {
	return StringDecomposerIterateFormattedReset(s, offset, style, style, output)
}

// StringDecomposerIterateFormattedReset mirrors StringDecomposer.iterateFormatted(String, int, Style, Style, FormattedCharSink).
func StringDecomposerIterateFormattedReset(s string, offset int, currentStyle, resetStyle Style, output FormattedCharSink) bool {
	units := toUTF16(s)
	size := len(units)
	style := currentStyle

	for i := offset; i < size; i++ {
		ch := units[i]
		if ch == 167 {
			if i+1 >= size {
				break
			}
			code := units[i+1]
			formatting := ChatFormattingGetByCode(rune(code))
			if formatting != nil {
				if *formatting == ChatFormattingRESET {
					style = resetStyle
				} else {
					style = style.ApplyLegacyFormat(*formatting)
				}
			}
			i++
		} else if isHighSurrogate(ch) {
			if i+1 >= size {
				if !output(i, style, 65533) {
					return false
				}
				break
			}
			low := units[i+1]
			if isLowSurrogate(low) {
				if !output(i, style, toCodePoint(ch, low)) {
					return false
				}
				i++
			} else if !output(i, style, 65533) {
				return false
			}
		} else if !feedChar(style, output, i, ch) {
			return false
		}
	}

	return true
}

// StringDecomposerIterateFormattedText mirrors StringDecomposer.iterateFormatted(FormattedText, Style, FormattedCharSink).
func StringDecomposerIterateFormattedText(component FormattedText, rootStyle Style, output FormattedCharSink) bool {
	value, present := component.VisitStyled(func(style Style, contents string) (any, bool) {
		if StringDecomposerIterateFormattedOffset(contents, 0, style, output) {
			return nil, false
		}
		return FormattedTextSTOP_ITERATION, true
	}, rootStyle)
	_ = value
	return !present
}

// StringDecomposerFilterBrokenSurrogates mirrors StringDecomposer.filterBrokenSurrogates.
func StringDecomposerFilterBrokenSurrogates(input string) string {
	var builder strings.Builder
	StringDecomposerIterate(input, StyleEMPTY, func(position int, style Style, codepoint int) bool {
		builder.WriteRune(rune(codepoint))
		return true
	})
	return builder.String()
}

// StringDecomposerGetPlainText mirrors StringDecomposer.getPlainText.
func StringDecomposerGetPlainText(input FormattedText) string {
	var builder strings.Builder
	StringDecomposerIterateFormattedText(input, StyleEMPTY, func(position int, style Style, codepoint int) bool {
		builder.WriteRune(rune(codepoint))
		return true
	})
	return builder.String()
}
