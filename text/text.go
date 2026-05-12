package text

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/admin-else/strom/nbt"
)

// FIXME: this should use the mapstruct thing

// Component represents a Minecraft text component.
// https://minecraft.wiki/w/Text_component_format
type Component struct {
	// Common fields
	Bold          *bool       `json:"bold,omitempty" nbt:"bold,omitempty"`
	Italic        *bool       `json:"italic,omitempty" nbt:"italic,omitempty"`
	Underlined    *bool       `json:"underlined,omitempty" nbt:"underlined,omitempty"`
	Strikethrough *bool       `json:"strikethrough,omitempty" nbt:"strikethrough,omitempty"`
	Obfuscated    *bool       `json:"obfuscated,omitempty" nbt:"obfuscated,omitempty"`
	Font          string      `json:"font,omitempty" nbt:"font,omitempty"`
	Color         string      `json:"color,omitempty" nbt:"color,omitempty"`
	Insertion     string      `json:"insertion,omitempty" nbt:"insertion,omitempty"`
	ClickEvent    *ClickEvent `json:"clickEvent,omitempty" nbt:"clickEvent,omitempty"`
	HoverEvent    *HoverEvent `json:"hoverEvent,omitempty" nbt:"hoverEvent,omitempty"`
	Extra         []Component `json:"extra,omitempty" nbt:"extra,omitempty"`

	// Content fields (one of these should be set)
	Text      string      `json:"text,omitempty" nbt:"text,omitempty"`
	Translate string      `json:"translate,omitempty" nbt:"translate,omitempty"`
	With      []Component `json:"with,omitempty" nbt:"with,omitempty"`
	Score     *Score      `json:"score,omitempty" nbt:"score,omitempty"`
	Selector  string      `json:"selector,omitempty" nbt:"selector,omitempty"`
	Keybind   string      `json:"keybind,omitempty" nbt:"keybind,omitempty"`
	NBT       string      `json:"nbt,omitempty" nbt:"nbt,omitempty"`
	Interpret bool        `json:"interpret,omitempty" nbt:"interpret,omitempty"`
	Block     string      `json:"block,omitempty" nbt:"block,omitempty"`
	Entity    string      `json:"entity,omitempty" nbt:"entity,omitempty"`
	Storage   string      `json:"storage,omitempty" nbt:"storage,omitempty"`
	Separator *Component  `json:"separator,omitempty" nbt:"separator,omitempty"`
}

type ClickEvent struct {
	Action string `json:"action" nbt:"action"`
	Value  string `json:"value" nbt:"value"`
}

type HoverEvent struct {
	Action   string `json:"action" nbt:"action"`
	Contents any    `json:"contents,omitempty" nbt:"contents,omitempty"`
	Value    any    `json:"value,omitempty" nbt:"value,omitempty"` // Legacy
}

type Score struct {
	Name      string `json:"name" nbt:"name"`
	Objective string `json:"objective" nbt:"objective"`
	Value     string `json:"value,omitempty" nbt:"value,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (c *Component) MarshalJSON() ([]byte, error) {
	// If it's just a simple text component with no formatting, we can marshal it as a string
	if c.isSimpleText() {
		return json.Marshal(c.Text)
	}
	if c.hasFormattingButNoContent() {
		// just like component but text does not have omitempty
		type AliasWithText struct {
			Bold          *bool       `json:"bold,omitempty"`
			Italic        *bool       `json:"italic,omitempty"`
			Underlined    *bool       `json:"underlined,omitempty"`
			Strikethrough *bool       `json:"strikethrough,omitempty"`
			Obfuscated    *bool       `json:"obfuscated,omitempty"`
			Font          string      `json:"font,omitempty"`
			Color         string      `json:"color,omitempty"`
			Insertion     string      `json:"insertion,omitempty"`
			ClickEvent    *ClickEvent `json:"clickEvent,omitempty"`
			HoverEvent    *HoverEvent `json:"hoverEvent,omitempty"`
			Extra         []Component `json:"extra,omitempty"`

			Text string `json:"text"`

			Translate string      `json:"translate,omitempty"`
			With      []Component `json:"with,omitempty"`
			Score     *Score      `json:"score,omitempty"`
			Selector  string      `json:"selector,omitempty"`
			Keybind   string      `json:"keybind,omitempty"`
			NBT       string      `json:"nbt,omitempty"`
			Interpret bool        `json:"interpret,omitempty"`
			Block     string      `json:"block,omitempty"`
			Entity    string      `json:"entity,omitempty"`
			Storage   string      `json:"storage,omitempty"`
			Separator *Component  `json:"separator,omitempty"`
		}
		return json.Marshal(AliasWithText(*c))
	}
	type Alias Component

	return json.Marshal(Alias(*c))
}

// hasFormattingButNoContent checks if the component has formatting/styling but no content fields
func (c *Component) hasFormattingButNoContent() bool {
	hasFormatting := c.Bold != nil || c.Italic != nil || c.Underlined != nil ||
		c.Strikethrough != nil || c.Obfuscated != nil || c.Font != "" ||
		c.Color != "" || c.Insertion != "" || c.ClickEvent != nil ||
		c.HoverEvent != nil || len(c.Extra) > 0 || c.Separator != nil

	hasContent := c.Text != "" || c.Translate != "" || len(c.With) > 0 ||
		c.Score != nil || c.Selector != "" || c.Keybind != "" ||
		c.NBT != "" || c.Block != "" || c.Entity != "" || c.Storage != ""

	return hasFormatting && !hasContent
}

func (c *Component) isSimpleText() bool {
	return c.Text != "" &&
		c.Bold == nil && c.Italic == nil && c.Underlined == nil && c.Strikethrough == nil && c.Obfuscated == nil &&
		c.Font == "" && c.Color == "" && c.Insertion == "" && c.ClickEvent == nil && c.HoverEvent == nil &&
		len(c.Extra) == 0 && c.Translate == "" && c.Score == nil && c.Selector == "" && c.Keybind == "" &&
		c.NBT == "" && !c.Interpret && c.Block == "" && c.Entity == "" && c.Storage == "" && c.Separator == nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *Component) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	// Try string
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		c.Text = s
		return nil
	}

	// Try array
	var arr []Component
	if err := json.Unmarshal(data, &arr); err == nil {
		if len(arr) > 0 {
			*c = arr[0]
			if len(arr) > 1 {
				c.Extra = append(c.Extra, arr[1:]...)
			}
		}
		return nil
	}

	// Try object
	type Alias Component
	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*c = Component(aux)
	return nil
}

// ToNBT converts the component to a format suitable for the nbt package.
func (c *Component) ToNBT() any {
	m := make(map[string]any)
	if c.Bold != nil {
		m["bold"] = boolToByte(*c.Bold)
	}
	if c.Italic != nil {
		m["italic"] = boolToByte(*c.Italic)
	}
	if c.Underlined != nil {
		m["underlined"] = boolToByte(*c.Underlined)
	}
	if c.Strikethrough != nil {
		m["strikethrough"] = boolToByte(*c.Strikethrough)
	}
	if c.Obfuscated != nil {
		m["obfuscated"] = boolToByte(*c.Obfuscated)
	}
	if c.Font != "" {
		m["font"] = c.Font
	}
	if c.Color != "" {
		m["color"] = c.Color
	}
	if c.Insertion != "" {
		m["insertion"] = c.Insertion
	}
	if c.ClickEvent != nil {
		m["clickEvent"] = map[string]any{
			"action": c.ClickEvent.Action,
			"value":  c.ClickEvent.Value,
		}
	}
	if c.HoverEvent != nil {
		he := map[string]any{
			"action": c.HoverEvent.Action,
		}
		if c.HoverEvent.Contents != nil {
			he["contents"] = c.HoverEvent.Contents
		}
		if c.HoverEvent.Value != nil {
			he["value"] = c.HoverEvent.Value
		}
		m["hoverEvent"] = he
	}
	if len(c.Extra) > 0 {
		extra := make([]any, len(c.Extra))
		for i, e := range c.Extra {
			extra[i] = e.ToNBT()
		}
		m["extra"] = extra
	}

	if c.Text != "" {
		m["text"] = c.Text
	}
	if c.Translate != "" {
		m["translate"] = c.Translate
	}
	if len(c.With) > 0 {
		with := make([]any, len(c.With))
		for i, e := range c.With {
			with[i] = e.ToNBT()
		}
		m["with"] = with
	}
	if c.Score != nil {
		s := map[string]any{
			"name":      c.Score.Name,
			"objective": c.Score.Objective,
		}
		if c.Score.Value != "" {
			s["value"] = c.Score.Value
		}
		m["score"] = s
	}
	if c.Selector != "" {
		m["selector"] = c.Selector
	}
	if c.Keybind != "" {
		m["keybind"] = c.Keybind
	}
	if c.NBT != "" {
		m["nbt"] = c.NBT
	}
	if c.Interpret {
		m["interpret"] = boolToByte(c.Interpret)
	}
	if c.Block != "" {
		m["block"] = c.Block
	}
	if c.Entity != "" {
		m["entity"] = c.Entity
	}
	if c.Storage != "" {
		m["storage"] = c.Storage
	}
	if c.Separator != nil {
		m["separator"] = c.Separator.ToNBT()
	}

	if c.hasFormattingButNoContent() {
		m["text"] = ""
	}

	return m
}

func boolToByte(b bool) int8 {
	if b {
		return 1
	}
	return 0
}

// FromNBT populates the component from an NBT-like structure (map[string]any or string).
func (c *Component) FromNBT(v any) error {
	switch v := v.(type) {
	case string:
		c.Text = v
		return nil
	case map[string]any:
		if val, ok := v["bold"]; ok {
			c.Bold = ptr(byteToBool(val))
		}
		if val, ok := v["italic"]; ok {
			c.Italic = ptr(byteToBool(val))
		}
		if val, ok := v["underlined"]; ok {
			c.Underlined = ptr(byteToBool(val))
		}
		if val, ok := v["strikethrough"]; ok {
			c.Strikethrough = ptr(byteToBool(val))
		}
		if val, ok := v["obfuscated"]; ok {
			c.Obfuscated = ptr(byteToBool(val))
		}
		if val, ok := v["font"].(string); ok {
			c.Font = val
		}
		if val, ok := v["color"].(string); ok {
			c.Color = val
		}
		if val, ok := v["insertion"].(string); ok {
			c.Insertion = val
		}
		if val, ok := v["clickEvent"].(map[string]any); ok {
			c.ClickEvent = &ClickEvent{}
			if a, ok := val["action"].(string); ok {
				c.ClickEvent.Action = a
			}
			if v, ok := val["value"].(string); ok {
				c.ClickEvent.Value = v
			}
		}
		if val, ok := v["hoverEvent"].(map[string]any); ok {
			c.HoverEvent = &HoverEvent{}
			if a, ok := val["action"].(string); ok {
				c.HoverEvent.Action = a
			}
			if co, ok := val["contents"]; ok {
				c.HoverEvent.Contents = co
			}
			if v, ok := val["value"]; ok {
				c.HoverEvent.Value = v
			}
		}
		if val, ok := v["extra"].([]any); ok {
			c.Extra = make([]Component, len(val))
			for i, e := range val {
				if err := c.Extra[i].FromNBT(e); err != nil {
					return err
				}
			}
		}
		if val, ok := v["text"].(string); ok {
			c.Text = val
		}
		if val, ok := v["translate"].(string); ok {
			c.Translate = val
		}
		if val, ok := v["with"].([]any); ok {
			c.With = make([]Component, len(val))
			for i, e := range val {
				if err := c.With[i].FromNBT(e); err != nil {
					return err
				}
			}
		}
		if val, ok := v["score"].(map[string]any); ok {
			c.Score = &Score{}
			if n, ok := val["name"].(string); ok {
				c.Score.Name = n
			}
			if o, ok := val["objective"].(string); ok {
				c.Score.Objective = o
			}
			if v, ok := val["value"].(string); ok {
				c.Score.Value = v
			}
		}
		if val, ok := v["selector"].(string); ok {
			c.Selector = val
		}
		if val, ok := v["keybind"].(string); ok {
			c.Keybind = val
		}
		if val, ok := v["nbt"].(string); ok {
			c.NBT = val
		}
		if val, ok := v["interpret"]; ok {
			c.Interpret = byteToBool(val)
		}
		if val, ok := v["block"].(string); ok {
			c.Block = val
		}
		if val, ok := v["entity"].(string); ok {
			c.Entity = val
		}
		if val, ok := v["storage"].(string); ok {
			c.Storage = val
		}
		if val, ok := v["separator"]; ok {
			c.Separator = &Component{}
			if err := c.Separator.FromNBT(val); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("unsupported NBT type for component: %T", v)
	}
}

func byteToBool(v any) bool {
	switch v := v.(type) {
	case int8:
		return v != 0
	case int32:
		return v != 0
	case bool:
		return v
	default:
		return false
	}
}

func ptr[T any](v T) *T {
	return &v
}

// Pretty parses a legacy formatted string into a Component.
// It uses the § symbol for color and formatting codes.
func Pretty(s string) *Component {
	parts := strings.Split(s, "§")

	root := Component{Text: parts[0]}

	var current = &root
	for i := 1; i < len(parts); i++ {
		part := parts[i]
		if part == "" {
			continue
		}

		code := part[0]
		text := part[1:]

		newComp := Component{}
		// Inherit formatting from current if it's just a style change
		if current != &root {
			newComp.Bold = current.Bold
			newComp.Italic = current.Italic
			newComp.Underlined = current.Underlined
			newComp.Strikethrough = current.Strikethrough
			newComp.Obfuscated = current.Obfuscated
			newComp.Color = current.Color
		}

		switch code {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f':
			newComp.Color = codeToColor[code]
			// Colors reset formatting
			newComp.Bold = nil
			newComp.Italic = nil
			newComp.Underlined = nil
			newComp.Strikethrough = nil
			newComp.Obfuscated = nil
		case 'k':
			newComp.Obfuscated = ptr(true)
		case 'l':
			newComp.Bold = ptr(true)
		case 'm':
			newComp.Strikethrough = ptr(true)
		case 'n':
			newComp.Underlined = ptr(true)
		case 'o':
			newComp.Italic = ptr(true)
		case 'r':
			newComp = Component{} // Reset
		case 'x':
			// Hex color §x§r§r§g§g§b§b
			if i+6 < len(parts) {
				hex := "#"
				for j := 1; j <= 6; j++ {
					if len(parts[i+j]) > 0 {
						hex += string(parts[i+j][0])
					}
				}
				newComp.Color = hex
				i += 6
				// Hex colors also reset formatting
				newComp.Bold = nil
				newComp.Italic = nil
				newComp.Underlined = nil
				newComp.Strikethrough = nil
				newComp.Obfuscated = nil
			}
		}

		if text != "" {
			newComp.Text = text
			root.Extra = append(root.Extra, newComp)
			current = &root.Extra[len(root.Extra)-1]
		} else {
			// Just a formatting code, update current state for next parts
			current = &newComp
		}
	}

	return &root
}

// PrettyF returns a Component from a formatted string.
func PrettyF(format string, args ...any) *Component {
	s := fmt.Sprintf(format, args...)
	return Pretty(s)
}

func replacePlaceholders(s string, values map[string]any) string {
	result := s
	for key, val := range values {
		formattedPrefix := "§%" + key + ":"
		for {
			idxFormatted := strings.Index(result, formattedPrefix)
			if idxFormatted == -1 {
				break
			}
			endIdx := strings.Index(result[idxFormatted+len(formattedPrefix):], "%")
			if endIdx == -1 {
				break
			}
			formatStr := result[idxFormatted+len(formattedPrefix) : idxFormatted+len(formattedPrefix)+endIdx]
			replaceWith := fmt.Sprintf("%"+formatStr, val)
			result = result[:idxFormatted] + replaceWith + result[idxFormatted+len(formattedPrefix)+endIdx+1:]
		}
		placeholder := "§%" + key + "%"
		for {
			idx := strings.Index(result, placeholder)
			if idx == -1 {
				break
			}
			replaceWith := fmt.Sprintf("%v", val)
			result = result[:idx] + replaceWith + result[idx+len(placeholder):]
		}
	}
	return result
}

func PrettyPlaceholders(format string, placeholders map[string]any) *Component {
	s := replacePlaceholders(format, placeholders)
	return Pretty(s)
}

var colorToCode = map[string]byte{
	"black":        '0',
	"dark_blue":    '1',
	"dark_green":   '2',
	"dark_aqua":    '3',
	"dark_red":     '4',
	"dark_purple":  '5',
	"gold":         '6',
	"gray":         '7',
	"dark_gray":    '8',
	"blue":         '9',
	"green":        'a',
	"aqua":         'b',
	"red":          'c',
	"light_purple": 'd',
	"yellow":       'e',
	"white":        'f',
}

var codeToColor = map[byte]string{
	'0': "black",
	'1': "dark_blue",
	'2': "dark_green",
	'3': "dark_aqua",
	'4': "dark_red",
	'5': "dark_purple",
	'6': "gold",
	'7': "gray",
	'8': "dark_gray",
	'9': "blue",
	'a': "green",
	'b': "aqua",
	'c': "red",
	'd': "light_purple",
	'e': "yellow",
	'f': "white",
}

func (c *Component) String() string {
	var sb strings.Builder
	c.writePretty(&sb)
	return sb.String()
}

func (c *Component) writePretty(sb *strings.Builder) {
	if c.Color != "" {
		if code, ok := colorToCode[strings.ToLower(c.Color)]; ok {
			sb.WriteString("§")
			sb.WriteByte(code)
		} else if strings.HasPrefix(c.Color, "#") && len(c.Color) == 7 {
			// Hex color: §x§r§r§g§g§b§b
			sb.WriteString("§x")
			for _, char := range c.Color[1:] {
				sb.WriteString("§")
				sb.WriteRune(char)
			}
		}
	}

	if c.Bold != nil && *c.Bold {
		sb.WriteString("§l")
	}
	if c.Italic != nil && *c.Italic {
		sb.WriteString("§o")
	}
	if c.Underlined != nil && *c.Underlined {
		sb.WriteString("§n")
	}
	if c.Strikethrough != nil && *c.Strikethrough {
		sb.WriteString("§m")
	}
	if c.Obfuscated != nil && *c.Obfuscated {
		sb.WriteString("§k")
	}

	if c.Text != "" {
		sb.WriteString(c.Text)
	}

	if c.Translate != "" {
		// This is a simplified version of translation handling
		// In a real client, this would look up the translation key
		sb.WriteString(c.Translate)
		sb.WriteRune('(')
		for _, w := range c.With {
			w.writePretty(sb)
			sb.WriteByte(',')
		}
		sb.WriteRune(')')
	}

	for _, extra := range c.Extra {
		extra.writePretty(sb)
	}
}

func (c *Component) AsNBT() nbt.Anon {
	return nbt.Anon{Value: c.ToNBT()}
}

func (c *Component) AsOldNBT() nbt.Tag {
	return nbt.Tag{
		Name:  "",
		Value: c.ToNBT(),
	}
}
