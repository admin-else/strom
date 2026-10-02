package text

import (
	"fmt"
	"reflect"
	"strings"
)

// Style mirrors net.minecraft.network.chat.Style.
type Style struct {
	color         *TextColor
	shadowColor   *int
	bold          *bool
	italic        *bool
	underlined    *bool
	strikethrough *bool
	obfuscated    *bool
	clickEvent    ClickEvent
	hoverEvent    HoverEvent
	insertion     *string
	font          FontDescription
}

// StyleEMPTY mirrors Style.EMPTY.
var StyleEMPTY = Style{}

// StyleNO_SHADOW mirrors Style.NO_SHADOW.
const StyleNO_SHADOW = 0

// newStyle mirrors the private Style constructor.
func newStyle(
	color *TextColor,
	shadowColor *int,
	bold *bool,
	italic *bool,
	underlined *bool,
	strikethrough *bool,
	obfuscated *bool,
	clickEvent ClickEvent,
	hoverEvent HoverEvent,
	insertion *string,
	font FontDescription,
) Style {
	return Style{
		color:         color,
		shadowColor:   shadowColor,
		bold:          bold,
		italic:        italic,
		underlined:    underlined,
		strikethrough: strikethrough,
		obfuscated:    obfuscated,
		clickEvent:    clickEvent,
		hoverEvent:    hoverEvent,
		insertion:     insertion,
		font:          font,
	}
}

// styleCreate mirrors Style.create; it collapses a resulting EMPTY to StyleEMPTY.
func styleCreate(s Style) Style {
	if s.Equals(StyleEMPTY) {
		return StyleEMPTY
	}
	return s
}

// GetColor mirrors Style.getColor().
func (s Style) GetColor() *TextColor { return s.color }

// GetShadowColor mirrors Style.getShadowColor().
func (s Style) GetShadowColor() *int { return s.shadowColor }

// IsBold mirrors Style.isBold().
func (s Style) IsBold() bool { return s.bold != nil && *s.bold }

// IsItalic mirrors Style.isItalic().
func (s Style) IsItalic() bool { return s.italic != nil && *s.italic }

// IsStrikethrough mirrors Style.isStrikethrough().
func (s Style) IsStrikethrough() bool { return s.strikethrough != nil && *s.strikethrough }

// IsUnderlined mirrors Style.isUnderlined().
func (s Style) IsUnderlined() bool { return s.underlined != nil && *s.underlined }

// IsObfuscated mirrors Style.isObfuscated().
func (s Style) IsObfuscated() bool { return s.obfuscated != nil && *s.obfuscated }

// IsEmpty mirrors Style.isEmpty().
func (s Style) IsEmpty() bool { return s.Equals(StyleEMPTY) }

// GetClickEvent mirrors Style.getClickEvent().
func (s Style) GetClickEvent() ClickEvent { return s.clickEvent }

// GetHoverEvent mirrors Style.getHoverEvent().
func (s Style) GetHoverEvent() HoverEvent { return s.hoverEvent }

// GetInsertion mirrors Style.getInsertion().
func (s Style) GetInsertion() *string { return s.insertion }

// GetFont mirrors Style.getFont().
func (s Style) GetFont() FontDescription {
	if s.font != nil {
		return s.font
	}
	return FontDescriptionDEFAULT
}

func checkEmptyAfterChange(newStyle Style, previousNonNil bool, nextNil bool) Style {
	if previousNonNil && nextNil && newStyle.Equals(StyleEMPTY) {
		return StyleEMPTY
	}
	return newStyle
}

// WithColor mirrors Style.withColor(TextColor).
func (s Style) WithColor(color *TextColor) Style {
	if textColorEqual(s.color, color) {
		return s
	}
	return checkEmptyAfterChange(newStyle(color, s.shadowColor, s.bold, s.italic, s.underlined, s.strikethrough, s.obfuscated, s.clickEvent, s.hoverEvent, s.insertion, s.font), s.color != nil, color == nil)
}

// WithColorFormat mirrors Style.withColor(ChatFormatting).
func (s Style) WithColorFormat(color *ChatFormatting) Style {
	if color == nil {
		return s.WithColor(nil)
	}
	return s.WithColor(FromLegacyFormat(*color))
}

// WithColorInt mirrors Style.withColor(int).
func (s Style) WithColorInt(color int) Style {
	c := FromRgb(color)
	return s.WithColor(&c)
}

// WithShadowColor mirrors Style.withShadowColor(int).
func (s Style) WithShadowColor(shadowColor int) Style {
	if s.shadowColor != nil && *s.shadowColor == shadowColor {
		return s
	}
	next := shadowColor
	return checkEmptyAfterChange(newStyle(s.color, &next, s.bold, s.italic, s.underlined, s.strikethrough, s.obfuscated, s.clickEvent, s.hoverEvent, s.insertion, s.font), s.shadowColor != nil, false)
}

// WithoutShadow mirrors Style.withoutShadow().
func (s Style) WithoutShadow() Style { return s.WithShadowColor(StyleNO_SHADOW) }

// WithBold mirrors Style.withBold(Boolean).
func (s Style) WithBold(bold *bool) Style {
	if boolPtrEqual(s.bold, bold) {
		return s
	}
	return checkEmptyAfterChange(newStyle(s.color, s.shadowColor, bold, s.italic, s.underlined, s.strikethrough, s.obfuscated, s.clickEvent, s.hoverEvent, s.insertion, s.font), s.bold != nil, bold == nil)
}

// WithItalic mirrors Style.withItalic(Boolean).
func (s Style) WithItalic(italic *bool) Style {
	if boolPtrEqual(s.italic, italic) {
		return s
	}
	return checkEmptyAfterChange(newStyle(s.color, s.shadowColor, s.bold, italic, s.underlined, s.strikethrough, s.obfuscated, s.clickEvent, s.hoverEvent, s.insertion, s.font), s.italic != nil, italic == nil)
}

// WithUnderlined mirrors Style.withUnderlined(Boolean).
func (s Style) WithUnderlined(underlined *bool) Style {
	if boolPtrEqual(s.underlined, underlined) {
		return s
	}
	return checkEmptyAfterChange(newStyle(s.color, s.shadowColor, s.bold, s.italic, underlined, s.strikethrough, s.obfuscated, s.clickEvent, s.hoverEvent, s.insertion, s.font), s.underlined != nil, underlined == nil)
}

// WithStrikethrough mirrors Style.withStrikethrough(Boolean).
func (s Style) WithStrikethrough(strikethrough *bool) Style {
	if boolPtrEqual(s.strikethrough, strikethrough) {
		return s
	}
	return checkEmptyAfterChange(newStyle(s.color, s.shadowColor, s.bold, s.italic, s.underlined, strikethrough, s.obfuscated, s.clickEvent, s.hoverEvent, s.insertion, s.font), s.strikethrough != nil, strikethrough == nil)
}

// WithObfuscated mirrors Style.withObfuscated(Boolean).
func (s Style) WithObfuscated(obfuscated *bool) Style {
	if boolPtrEqual(s.obfuscated, obfuscated) {
		return s
	}
	return checkEmptyAfterChange(newStyle(s.color, s.shadowColor, s.bold, s.italic, s.underlined, s.strikethrough, obfuscated, s.clickEvent, s.hoverEvent, s.insertion, s.font), s.obfuscated != nil, obfuscated == nil)
}

// WithClickEvent mirrors Style.withClickEvent(ClickEvent).
func (s Style) WithClickEvent(clickEvent ClickEvent) Style {
	if reflect.DeepEqual(s.clickEvent, clickEvent) {
		return s
	}
	return checkEmptyAfterChange(newStyle(s.color, s.shadowColor, s.bold, s.italic, s.underlined, s.strikethrough, s.obfuscated, clickEvent, s.hoverEvent, s.insertion, s.font), s.clickEvent != nil, clickEvent == nil)
}

// WithHoverEvent mirrors Style.withHoverEvent(HoverEvent).
func (s Style) WithHoverEvent(hoverEvent HoverEvent) Style {
	if reflect.DeepEqual(s.hoverEvent, hoverEvent) {
		return s
	}
	return checkEmptyAfterChange(newStyle(s.color, s.shadowColor, s.bold, s.italic, s.underlined, s.strikethrough, s.obfuscated, s.clickEvent, hoverEvent, s.insertion, s.font), s.hoverEvent != nil, hoverEvent == nil)
}

// WithInsertion mirrors Style.withInsertion(String).
func (s Style) WithInsertion(insertion *string) Style {
	if stringPtrEqual(s.insertion, insertion) {
		return s
	}
	return checkEmptyAfterChange(newStyle(s.color, s.shadowColor, s.bold, s.italic, s.underlined, s.strikethrough, s.obfuscated, s.clickEvent, s.hoverEvent, insertion, s.font), s.insertion != nil, insertion == nil)
}

// WithFont mirrors Style.withFont(FontDescription).
func (s Style) WithFont(font FontDescription) Style {
	if reflect.DeepEqual(s.font, font) {
		return s
	}
	return checkEmptyAfterChange(newStyle(s.color, s.shadowColor, s.bold, s.italic, s.underlined, s.strikethrough, s.obfuscated, s.clickEvent, s.hoverEvent, s.insertion, font), s.font != nil, font == nil)
}

// ApplyFormat mirrors Style.applyFormat.
func (s Style) ApplyFormat(format ChatFormatting) Style {
	color := s.color
	bold := s.bold
	italic := s.italic
	strikethrough := s.strikethrough
	underlined := s.underlined
	obfuscated := s.obfuscated

	switch format {
	case ChatFormattingOBFUSCATED:
		t := true
		obfuscated = &t
	case ChatFormattingBOLD:
		t := true
		bold = &t
	case ChatFormattingSTRIKETHROUGH:
		t := true
		strikethrough = &t
	case ChatFormattingUNDERLINE:
		t := true
		underlined = &t
	case ChatFormattingITALIC:
		t := true
		italic = &t
	case ChatFormattingRESET:
		return StyleEMPTY
	default:
		color = FromLegacyFormat(format)
	}

	return newStyle(color, s.shadowColor, bold, italic, underlined, strikethrough, obfuscated, s.clickEvent, s.hoverEvent, s.insertion, s.font)
}

// ApplyLegacyFormat mirrors Style.applyLegacyFormat.
func (s Style) ApplyLegacyFormat(format ChatFormatting) Style {
	color := s.color
	bold := s.bold
	italic := s.italic
	strikethrough := s.strikethrough
	underlined := s.underlined
	obfuscated := s.obfuscated

	switch format {
	case ChatFormattingOBFUSCATED:
		t := true
		obfuscated = &t
	case ChatFormattingBOLD:
		t := true
		bold = &t
	case ChatFormattingSTRIKETHROUGH:
		t := true
		strikethrough = &t
	case ChatFormattingUNDERLINE:
		t := true
		underlined = &t
	case ChatFormattingITALIC:
		t := true
		italic = &t
	case ChatFormattingRESET:
		return StyleEMPTY
	default:
		f := false
		obfuscated = &f
		bold = &f
		strikethrough = &f
		underlined = &f
		italic = &f
		color = FromLegacyFormat(format)
	}

	return newStyle(color, s.shadowColor, bold, italic, underlined, strikethrough, obfuscated, s.clickEvent, s.hoverEvent, s.insertion, s.font)
}

// ApplyFormats mirrors Style.applyFormats.
func (s Style) ApplyFormats(formats ...ChatFormatting) Style {
	color := s.color
	bold := s.bold
	italic := s.italic
	strikethrough := s.strikethrough
	underlined := s.underlined
	obfuscated := s.obfuscated

	for _, format := range formats {
		switch format {
		case ChatFormattingOBFUSCATED:
			t := true
			obfuscated = &t
		case ChatFormattingBOLD:
			t := true
			bold = &t
		case ChatFormattingSTRIKETHROUGH:
			t := true
			strikethrough = &t
		case ChatFormattingUNDERLINE:
			t := true
			underlined = &t
		case ChatFormattingITALIC:
			t := true
			italic = &t
		case ChatFormattingRESET:
			return StyleEMPTY
		default:
			color = FromLegacyFormat(format)
		}
	}

	return newStyle(color, s.shadowColor, bold, italic, underlined, strikethrough, obfuscated, s.clickEvent, s.hoverEvent, s.insertion, s.font)
}

// ApplyTo mirrors Style.applyTo.
func (s Style) ApplyTo(other Style) Style {
	if s.IsEmpty() {
		return other
	}
	if other.IsEmpty() {
		return s
	}
	return newStyle(
		firstNonNilTextColor(s.color, other.color),
		firstNonNilInt(s.shadowColor, other.shadowColor),
		firstNonNilBool(s.bold, other.bold),
		firstNonNilBool(s.italic, other.italic),
		firstNonNilBool(s.underlined, other.underlined),
		firstNonNilBool(s.strikethrough, other.strikethrough),
		firstNonNilBool(s.obfuscated, other.obfuscated),
		firstNonNilClick(s.clickEvent, other.clickEvent),
		firstNonNilHover(s.hoverEvent, other.hoverEvent),
		firstNonNilString(s.insertion, other.insertion),
		firstNonNilFont(s.font, other.font),
	)
}

// String mirrors Style.toString().
func (s Style) String() string {
	var result strings.Builder
	result.WriteString("{")
	first := true
	prependSeparator := func() {
		if first {
			first = false
			return
		}
		result.WriteString(",")
	}
	addFlagString := func(name string, value *bool) {
		if value != nil {
			prependSeparator()
			if !*value {
				result.WriteString("!")
			}
			result.WriteString(name)
		}
	}
	addValueString := func(name string, present bool, value string) {
		if present {
			prependSeparator()
			result.WriteString(name)
			result.WriteString("=")
			result.WriteString(value)
		}
	}

	if s.color != nil {
		addValueString("color", true, s.color.String())
	}
	if s.shadowColor != nil {
		addValueString("shadowColor", true, fmt.Sprintf("%d", *s.shadowColor))
	}
	addFlagString("bold", s.bold)
	addFlagString("italic", s.italic)
	addFlagString("underlined", s.underlined)
	addFlagString("strikethrough", s.strikethrough)
	addFlagString("obfuscated", s.obfuscated)
	if s.clickEvent != nil {
		addValueString("clickEvent", true, fmt.Sprintf("%v", s.clickEvent))
	}
	if s.hoverEvent != nil {
		addValueString("hoverEvent", true, fmt.Sprintf("%v", s.hoverEvent))
	}
	if s.insertion != nil {
		addValueString("insertion", true, *s.insertion)
	}
	if s.font != nil {
		addValueString("font", true, fmt.Sprintf("%v", s.font))
	}
	result.WriteString("}")
	return result.String()
}

// Equals mirrors Style.equals.
func (s Style) Equals(o Style) bool {
	return boolPtrEqual(s.bold, o.bold) &&
		textColorEqual(s.color, o.color) &&
		intPtrEqual(s.shadowColor, o.shadowColor) &&
		boolPtrEqual(s.italic, o.italic) &&
		boolPtrEqual(s.obfuscated, o.obfuscated) &&
		boolPtrEqual(s.strikethrough, o.strikethrough) &&
		boolPtrEqual(s.underlined, o.underlined) &&
		reflect.DeepEqual(s.clickEvent, o.clickEvent) &&
		reflect.DeepEqual(s.hoverEvent, o.hoverEvent) &&
		stringPtrEqual(s.insertion, o.insertion) &&
		reflect.DeepEqual(s.font, o.font)
}

// HashCode mirrors Style.hashCode.
func (s Style) HashCode() int {
	h := 1
	h = 31*h + textColorHash(s.color)
	h = 31*h + intPtrHash(s.shadowColor)
	h = 31*h + boolPtrHash(s.bold)
	h = 31*h + boolPtrHash(s.italic)
	h = 31*h + boolPtrHash(s.underlined)
	h = 31*h + boolPtrHash(s.strikethrough)
	h = 31*h + boolPtrHash(s.obfuscated)
	h = 31*h + interfaceHash(s.clickEvent)
	h = 31*h + interfaceHash(s.hoverEvent)
	h = 31*h + stringPtrHash(s.insertion)
	return h
}

func boolPtrEqual(a, b *bool) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func intPtrEqual(a, b *int) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func stringPtrEqual(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func textColorEqual(a, b *TextColor) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Equals(b)
}

func boolPtrHash(a *bool) int {
	if a == nil {
		return 0
	}
	if *a {
		return 1231
	}
	return 1237
}

func intPtrHash(a *int) int {
	if a == nil {
		return 0
	}
	return *a
}

func stringPtrHash(a *string) int {
	if a == nil {
		return 0
	}
	return stringHash(*a)
}

func textColorHash(a *TextColor) int {
	if a == nil {
		return 0
	}
	return a.HashCode()
}

func interfaceHash(v any) int {
	if v == nil {
		return 0
	}
	return stringHash(fmt.Sprintf("%v", v))
}

func firstNonNilTextColor(a, b *TextColor) *TextColor {
	if a != nil {
		return a
	}
	return b
}

func firstNonNilInt(a, b *int) *int {
	if a != nil {
		return a
	}
	return b
}

func firstNonNilBool(a, b *bool) *bool {
	if a != nil {
		return a
	}
	return b
}

func firstNonNilString(a, b *string) *string {
	if a != nil {
		return a
	}
	return b
}

func firstNonNilClick(a, b ClickEvent) ClickEvent {
	if a != nil {
		return a
	}
	return b
}

func firstNonNilHover(a, b HoverEvent) HoverEvent {
	if a != nil {
		return a
	}
	return b
}

func firstNonNilFont(a, b FontDescription) FontDescription {
	if a != nil {
		return a
	}
	return b
}
