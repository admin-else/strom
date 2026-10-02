package text

import (
	"fmt"
	"reflect"
	"time"
	"unicode/utf16"

	"github.com/admin-else/strom/mc/resources"
	"github.com/google/uuid"
)

// Component mirrors net.minecraft.network.chat.Component.
type Component interface {
	FormattedText

	GetStyle() Style
	GetContents() ComponentContents
	GetString() string
	GetStringLimit(limit int) string
	GetSiblings() []Component
	TryCollapseToString() (text string, ok bool)
	PlainCopy() *MutableComponent
	Copy() *MutableComponent
	GetVisualOrderText() FormattedCharSequence
}

// ComponentGetString mirrors Component's getString() default (FormattedText.super.getString()).
func ComponentGetString(component Component) string {
	return FormattedTextGetString(component)
}

// ComponentGetStringLimit mirrors Component.getString(int).
func ComponentGetStringLimit(component Component, limit int) string {
	var units []uint16
	component.VisitContent(func(contents string) (any, bool) {
		remaining := limit - len(units)
		if remaining <= 0 {
			return FormattedTextSTOP_ITERATION, true
		}
		contentUnits := toUTF16(contents)
		if len(contentUnits) <= remaining {
			units = append(units, contentUnits...)
		} else {
			units = append(units, contentUnits[:remaining]...)
		}
		return nil, false
	})
	return string(utf16.Decode(units))
}

// ComponentVisitStyled mirrors Component's visit(StyledContentConsumer, Style) default.
func ComponentVisitStyled(component Component, output StyledContentConsumer, parentStyle Style) (any, bool) {
	selfStyle := component.GetStyle().ApplyTo(parentStyle)
	if value, present := component.GetContents().VisitStyled(output, selfStyle); present {
		return value, true
	}
	for _, sibling := range component.GetSiblings() {
		if value, present := sibling.VisitStyled(output, selfStyle); present {
			return value, true
		}
	}
	return nil, false
}

// ComponentVisitContent mirrors Component's visit(ContentConsumer) default.
func ComponentVisitContent(component Component, output ContentConsumer) (any, bool) {
	if value, present := component.GetContents().VisitContent(output); present {
		return value, true
	}
	for _, sibling := range component.GetSiblings() {
		if value, present := sibling.VisitContent(output); present {
			return value, true
		}
	}
	return nil, false
}

// ComponentToFlatList mirrors Component.toFlatList().
func ComponentToFlatList(component Component) []Component {
	return ComponentToFlatListStyle(component, StyleEMPTY)
}

// ComponentToFlatListStyle mirrors Component.toFlatList(Style).
func ComponentToFlatListStyle(component Component, rootStyle Style) []Component {
	var result []Component
	component.VisitStyled(func(style Style, contents string) (any, bool) {
		if contents != "" {
			result = append(result, ComponentLiteral(contents).WithStyle(style))
		}
		return nil, false
	}, rootStyle)
	return result
}

// ComponentContains mirrors Component.contains.
func ComponentContains(component Component, other Component) bool {
	if componentEquals(component, other) {
		return true
	}
	flat := ComponentToFlatList(component)
	otherFlat := ComponentToFlatListStyle(other, component.GetStyle())
	return indexOfSubList(flat, otherFlat) != -1
}

func componentEquals(a, b Component) bool {
	if a == nil || b == nil {
		return a == b
	}
	if mutable, ok := a.(*MutableComponent); ok {
		return mutable.Equals(b)
	}
	return reflect.DeepEqual(a, b)
}

func componentSiblingsEqual(a, b []Component) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !componentEquals(a[i], b[i]) {
			return false
		}
	}
	return true
}

// componentContentsEqual mirrors ComponentContents.equals; it compares by content
// rather than by cache state (decompose/visual-order caches are ignored).
func componentContentsEqual(a, b ComponentContents) bool {
	switch x := a.(type) {
	case PlainTextContents:
		y, ok := b.(PlainTextContents)
		return ok && x.Text() == y.Text()
	case *TranslatableContents:
		y, ok := b.(*TranslatableContents)
		return ok && x.Equals(y)
	case *KeybindContents:
		y, ok := b.(*KeybindContents)
		return ok && x.Equals(y)
	default:
		return reflect.DeepEqual(a, b)
	}
}

func indexOfSubList(list, target []Component) int {
	targetSize := len(target)
	if targetSize == 0 {
		return 0
	}
	max := len(list) - targetSize
	for i := 0; i <= max; i++ {
		matches := true
		for j := 0; j < targetSize; j++ {
			if !componentEquals(list[i+j], target[j]) {
				matches = false
				break
			}
		}
		if matches {
			return i
		}
	}
	return -1
}

// NullToEmpty mirrors Component.nullToEmpty.
func ComponentNullToEmpty(text *string) Component {
	if text != nil {
		return ComponentLiteral(*text)
	}
	return CommonComponentsEMPTY
}

// ComponentLiteral mirrors Component.literal.
func ComponentLiteral(text string) *MutableComponent {
	return MutableComponentCreate(PlainTextContentsCreate(text))
}

// ComponentTranslatable mirrors Component.translatable(String).
func ComponentTranslatable(key string) *MutableComponent {
	return MutableComponentCreate(NewTranslatableContents(key, nil, TranslatableContentsNO_ARGS...))
}

// ComponentTranslatableArgs mirrors Component.translatable(String, Object...).
func ComponentTranslatableArgs(key string, args ...any) *MutableComponent {
	return MutableComponentCreate(NewTranslatableContents(key, nil, args...))
}

// ComponentTranslatableEscape mirrors Component.translatableEscape.
func ComponentTranslatableEscape(key string, args ...any) *MutableComponent {
	escaped := make([]any, len(args))
	copy(escaped, args)
	for i, arg := range escaped {
		if !TranslatableContentsIsAllowedPrimitiveArgument(arg) {
			if _, ok := arg.(Component); !ok {
				escaped[i] = stringifyArg(arg)
			}
		}
	}
	return ComponentTranslatableArgs(key, escaped...)
}

// ComponentTranslatableWithFallback mirrors Component.translatableWithFallback(String, String).
func ComponentTranslatableWithFallback(key string, fallback *string) *MutableComponent {
	return MutableComponentCreate(NewTranslatableContents(key, fallback, TranslatableContentsNO_ARGS...))
}

// ComponentTranslatableWithFallbackArgs mirrors Component.translatableWithFallback(String, String, Object...).
func ComponentTranslatableWithFallbackArgs(key string, fallback *string, args ...any) *MutableComponent {
	return MutableComponentCreate(NewTranslatableContents(key, fallback, args...))
}

// ComponentEmpty mirrors Component.empty.
func ComponentEmpty() *MutableComponent {
	return MutableComponentCreate(PlainTextContentsEMPTY)
}

// ComponentKeybind mirrors Component.keybind.
func ComponentKeybind(name string) *MutableComponent {
	return MutableComponentCreate(NewKeybindContents(name))
}

// ComponentTranslationArgDate mirrors Component.translationArg(Date).
func ComponentTranslationArgDate(date time.Time) Component {
	return ComponentLiteral(date.Format("Mon Jan 02 15:04:05 MST 2006"))
}

// ComponentTranslationArgMessage mirrors Component.translationArg(Message).
func ComponentTranslationArgMessage(message Message) Component {
	if component, ok := message.(Component); ok {
		return component
	}
	return ComponentLiteral(message.GetString())
}

// ComponentTranslationArgUUID mirrors Component.translationArg(UUID).
func ComponentTranslationArgUUID(id uuid.UUID) Component {
	return ComponentLiteral(id.String())
}

// ComponentTranslationArgIdentifier mirrors Component.translationArg(Identifier).
func ComponentTranslationArgIdentifier(id resources.Identifier) Component {
	return ComponentLiteral(id.String())
}

// ComponentTranslationArgURI mirrors Component.translationArg(URI).
func ComponentTranslationArgURI(uri string) Component {
	return ComponentLiteral(uri)
}

func stringifyArg(arg any) string {
	if arg == nil {
		return "null"
	}
	if s, ok := arg.(interface{ String() string }); ok {
		return s.String()
	}
	return fmt.Sprintf("%v", arg)
}

func optionalComponentEquals(a, b *Component) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return reflect.DeepEqual(*a, *b)
}

func optionalComponentHash(c *Component) int {
	if c == nil {
		return 0
	}
	return stringHash(fmt.Sprintf("%v", *c))
}
