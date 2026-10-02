package text

import (
	"fmt"
)

// MutableComponent mirrors net.minecraft.network.chat.MutableComponent.
type MutableComponent struct {
	contents        ComponentContents
	siblings        []Component
	style           Style
	visualOrderText FormattedCharSequence
	decomposedWith  Language
}

// NewMutableComponent mirrors the package-private MutableComponent constructor.
func NewMutableComponent(contents ComponentContents, siblings []Component, style Style) *MutableComponent {
	if siblings == nil {
		siblings = []Component{}
	}
	return &MutableComponent{contents: contents, siblings: siblings, style: style}
}

// MutableComponentCreate mirrors MutableComponent.create.
func MutableComponentCreate(contents ComponentContents) *MutableComponent {
	return NewMutableComponent(contents, []Component{}, StyleEMPTY)
}

// GetContents mirrors MutableComponent.getContents().
func (m *MutableComponent) GetContents() ComponentContents { return m.contents }

// GetSiblings mirrors MutableComponent.getSiblings().
func (m *MutableComponent) GetSiblings() []Component { return m.siblings }

// SetStyle mirrors MutableComponent.setStyle.
func (m *MutableComponent) SetStyle(style Style) *MutableComponent {
	m.style = style
	return m
}

// GetStyle mirrors MutableComponent.getStyle().
func (m *MutableComponent) GetStyle() Style { return m.style }

// Append mirrors MutableComponent.append(String); an empty text returns this.
func (m *MutableComponent) Append(text string) *MutableComponent {
	if text == "" {
		return m
	}
	return m.AppendComponent(ComponentLiteral(text))
}

// AppendComponent mirrors MutableComponent.append(Component).
func (m *MutableComponent) AppendComponent(component Component) *MutableComponent {
	m.siblings = append(m.siblings, component)
	return m
}

// WithStyleFunc mirrors MutableComponent.withStyle(UnaryOperator<Style>).
func (m *MutableComponent) WithStyleFunc(updater func(Style) Style) *MutableComponent {
	m.SetStyle(updater(m.GetStyle()))
	return m
}

// WithStyle mirrors MutableComponent.withStyle(Style).
func (m *MutableComponent) WithStyle(patch Style) *MutableComponent {
	m.SetStyle(patch.ApplyTo(m.GetStyle()))
	return m
}

// WithStyleFormats mirrors MutableComponent.withStyle(ChatFormatting...).
func (m *MutableComponent) WithStyleFormats(formats ...ChatFormatting) *MutableComponent {
	m.SetStyle(m.GetStyle().ApplyFormats(formats...))
	return m
}

// WithStyleFormat mirrors MutableComponent.withStyle(ChatFormatting).
func (m *MutableComponent) WithStyleFormat(format ChatFormatting) *MutableComponent {
	m.SetStyle(m.GetStyle().ApplyFormat(format))
	return m
}

// WithColor mirrors MutableComponent.withColor(int).
func (m *MutableComponent) WithColor(color int) *MutableComponent {
	m.SetStyle(m.GetStyle().WithColorInt(color))
	return m
}

// WithTextColor mirrors MutableComponent.withColor(TextColor).
func (m *MutableComponent) WithTextColor(color *TextColor) *MutableComponent {
	m.SetStyle(m.GetStyle().WithColor(color))
	return m
}

// WithoutShadow mirrors MutableComponent.withoutShadow.
func (m *MutableComponent) WithoutShadow() *MutableComponent {
	m.SetStyle(m.GetStyle().WithoutShadow())
	return m
}

// GetVisualOrderText mirrors MutableComponent.getVisualOrderText().
func (m *MutableComponent) GetVisualOrderText() FormattedCharSequence {
	currentLanguage := LanguageGetInstance()
	if m.decomposedWith != currentLanguage {
		m.visualOrderText = currentLanguage.GetVisualOrder(m)
		m.decomposedWith = currentLanguage
	}
	return m.visualOrderText
}

// VisitContent implements FormattedText.
func (m *MutableComponent) VisitContent(output ContentConsumer) (any, bool) {
	return ComponentVisitContent(m, output)
}

// VisitStyled implements FormattedText.
func (m *MutableComponent) VisitStyled(output StyledContentConsumer, parentStyle Style) (any, bool) {
	return ComponentVisitStyled(m, output, parentStyle)
}

// GetString mirrors Component.getString().
func (m *MutableComponent) GetString() string { return ComponentGetString(m) }

// GetStringLimit mirrors Component.getString(int).
func (m *MutableComponent) GetStringLimit(limit int) string { return ComponentGetStringLimit(m, limit) }

// TryCollapseToString mirrors Component.tryCollapseToString().
func (m *MutableComponent) TryCollapseToString() (text string, ok bool) {
	if plain, isPlain := m.contents.(PlainTextContents); isPlain && len(m.siblings) == 0 && m.style.IsEmpty() {
		return plain.Text(), true
	}
	return "", false
}

// PlainCopy mirrors Component.plainCopy().
func (m *MutableComponent) PlainCopy() *MutableComponent {
	return MutableComponentCreate(m.contents)
}

// Copy mirrors Component.copy().
func (m *MutableComponent) Copy() *MutableComponent {
	siblings := make([]Component, len(m.siblings))
	copy(siblings, m.siblings)
	return NewMutableComponent(m.contents, siblings, m.style)
}

// Equals mirrors MutableComponent.equals.
func (m *MutableComponent) Equals(o Component) bool {
	if o == nil {
		return false
	}
	that, isMutable := o.(*MutableComponent)
	if !isMutable {
		return false
	}
	if m == that {
		return true
	}
	return componentContentsEqual(m.contents, that.contents) && m.style.Equals(that.style) && componentSiblingsEqual(m.siblings, that.siblings)
}

// HashCode mirrors MutableComponent.hashCode.
func (m *MutableComponent) HashCode() int {
	result := 1
	result = 31*result + stringHash(fmt.Sprintf("%v", m.contents))
	result = 31*result + m.style.HashCode()
	result = 31*result + stringHash(fmt.Sprintf("%v", m.siblings))
	return result
}

// String mirrors MutableComponent.toString.
func (m *MutableComponent) String() string {
	result := fmt.Sprintf("%v", m.contents)
	hasStyle := !m.style.IsEmpty()
	hasSiblings := len(m.siblings) > 0
	if hasStyle || hasSiblings {
		result += "["
		if hasStyle {
			result += "style="
			result += m.style.String()
		}
		if hasStyle && hasSiblings {
			result += ", "
		}
		if hasSiblings {
			result += "siblings="
			result += fmt.Sprintf("%v", m.siblings)
		}
		result += "]"
	}
	return result
}
