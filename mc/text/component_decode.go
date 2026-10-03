package text

import "encoding/json"

// DecodeComponentJSON mirrors ComponentSerialization.CODEC for the subset of the
// text component format that pack metadata uses. It accepts a JSON string, an
// array of components, or an object, and returns the component.
func DecodeComponentJSON(data json.RawMessage) (component Component, err error) {
	var raw RawComponent
	if err = json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	return ComponentFromRaw(&raw), nil
}

// ComponentFromRaw converts a RawComponent tree into a Component, mirroring the
// ComponentSerialization codec's content/style/sibling split.
func ComponentFromRaw(raw *RawComponent) Component {
	if raw == nil {
		return ComponentEmpty()
	}
	var component *MutableComponent
	switch {
	case raw.Translate != "":
		args := make([]any, len(raw.With))
		for i := range raw.With {
			args[i] = ComponentFromRaw(&raw.With[i])
		}
		component = MutableComponentCreate(NewTranslatableContents(raw.Translate, nil, args...))
	case raw.Keybind != "":
		component = ComponentKeybind(raw.Keybind)
	default:
		component = ComponentLiteral(raw.Text)
	}

	component.SetStyle(componentStyleFromRaw(raw))
	for i := range raw.Extra {
		component.AppendComponent(ComponentFromRaw(&raw.Extra[i]))
	}
	return component
}

func componentStyleFromRaw(raw *RawComponent) Style {
	style := StyleEMPTY
	if raw.Bold != nil {
		style = style.WithBold(raw.Bold)
	}
	if raw.Italic != nil {
		style = style.WithItalic(raw.Italic)
	}
	if raw.Underlined != nil {
		style = style.WithUnderlined(raw.Underlined)
	}
	if raw.Strikethrough != nil {
		style = style.WithStrikethrough(raw.Strikethrough)
	}
	if raw.Obfuscated != nil {
		style = style.WithObfuscated(raw.Obfuscated)
	}
	if raw.Color != "" {
		if color, err := ParseColor(raw.Color); err == nil {
			style = style.WithColor(&color)
		}
	}
	return style
}
