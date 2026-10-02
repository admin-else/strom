package text

import (
	"fmt"
	"strconv"
	"strings"
)

// TextColor mirrors net.minecraft.network.chat.TextColor.
type TextColor struct {
	value int
	name  *string
}

const textColorCustomColorPrefix = "#"

var textColorNamedColors = map[string]TextColor{}

var (
	TextColorBLACK        = namedColor("black", 0)
	TextColorDARK_BLUE    = namedColor("dark_blue", 170)
	TextColorDARK_GREEN   = namedColor("dark_green", 43520)
	TextColorDARK_AQUA    = namedColor("dark_aqua", 43690)
	TextColorDARK_RED     = namedColor("dark_red", 11141120)
	TextColorDARK_PURPLE  = namedColor("dark_purple", 11141290)
	TextColorGOLD         = namedColor("gold", 16755200)
	TextColorGRAY         = namedColor("gray", 11184810)
	TextColorDARK_GRAY    = namedColor("dark_gray", 5592405)
	TextColorBLUE         = namedColor("blue", 5592575)
	TextColorGREEN        = namedColor("green", 5635925)
	TextColorAQUA         = namedColor("aqua", 5636095)
	TextColorRED          = namedColor("red", 16733525)
	TextColorLIGHT_PURPLE = namedColor("light_purple", 16733695)
	TextColorYELLOW       = namedColor("yellow", 16777045)
	TextColorWHITE        = namedColor("white", 16777215)
)

func namedColor(name string, rgb int) TextColor {
	result := TextColor{value: rgb & 16777215, name: &name}
	textColorNamedColors[name] = result
	return result
}

func newTextColorRgb(value int) TextColor {
	return TextColor{value: value & 16777215}
}

// GetValue mirrors TextColor.getValue().
func (t TextColor) GetValue() int { return t.value }

// GetName mirrors TextColor's @Nullable name (Java has no accessor; provided for parity).
func (t TextColor) GetName() *string { return t.name }

// Serialize mirrors TextColor.serialize().
func (t TextColor) Serialize() string {
	if t.name != nil {
		return *t.name
	}
	return t.FormatValue()
}

// FormatValue mirrors TextColor.formatValue().
func (t TextColor) FormatValue() string {
	return fmt.Sprintf("#%06X", t.value)
}

// Equals mirrors TextColor.equals.
func (t TextColor) Equals(o *TextColor) bool {
	if o == nil {
		return false
	}
	return t.value == o.value
}

// HashCode mirrors TextColor.hashCode().
func (t TextColor) HashCode() int {
	h := 1
	h = 31*h + t.value
	if t.name != nil {
		h = 31*h + stringHash(*t.name)
	}
	return h
}

// String mirrors TextColor.toString().
func (t TextColor) String() string { return t.Serialize() }

// FromLegacyFormat mirrors TextColor.fromLegacyFormat.
func FromLegacyFormat(format ChatFormatting) *TextColor {
	switch format {
	case ChatFormattingBLACK:
		c := TextColorBLACK
		return &c
	case ChatFormattingDARK_BLUE:
		c := TextColorDARK_BLUE
		return &c
	case ChatFormattingDARK_GREEN:
		c := TextColorDARK_GREEN
		return &c
	case ChatFormattingDARK_AQUA:
		c := TextColorDARK_AQUA
		return &c
	case ChatFormattingDARK_RED:
		c := TextColorDARK_RED
		return &c
	case ChatFormattingDARK_PURPLE:
		c := TextColorDARK_PURPLE
		return &c
	case ChatFormattingGOLD:
		c := TextColorGOLD
		return &c
	case ChatFormattingGRAY:
		c := TextColorGRAY
		return &c
	case ChatFormattingDARK_GRAY:
		c := TextColorDARK_GRAY
		return &c
	case ChatFormattingBLUE:
		c := TextColorBLUE
		return &c
	case ChatFormattingGREEN:
		c := TextColorGREEN
		return &c
	case ChatFormattingAQUA:
		c := TextColorAQUA
		return &c
	case ChatFormattingRED:
		c := TextColorRED
		return &c
	case ChatFormattingLIGHT_PURPLE:
		c := TextColorLIGHT_PURPLE
		return &c
	case ChatFormattingYELLOW:
		c := TextColorYELLOW
		return &c
	case ChatFormattingWHITE:
		c := TextColorWHITE
		return &c
	default:
		return nil
	}
}

// FromRgb mirrors TextColor.fromRgb.
func FromRgb(rgb int) TextColor { return newTextColorRgb(rgb) }

// ParseColor mirrors TextColor.parseColor; the DataResult becomes (value, error).
func ParseColor(color string) (result TextColor, err error) {
	if strings.HasPrefix(color, textColorCustomColorPrefix) {
		value, convErr := strconv.ParseInt(color[1:], 16, 64)
		if convErr != nil {
			return TextColor{}, fmt.Errorf("Invalid color value: %s", color)
		}
		if value < 0 || value > 16777215 {
			return TextColor{}, fmt.Errorf("Color value out of range: %s", color)
		}
		return FromRgb(int(value)), nil
	}
	predefined, ok := textColorNamedColors[color]
	if !ok {
		return TextColor{}, fmt.Errorf("Invalid color name: %s", color)
	}
	return predefined, nil
}

func stringHash(s string) int {
	h := 0
	for _, r := range s {
		h = 31*h + int(r)
	}
	return h
}
