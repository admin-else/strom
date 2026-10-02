package text

import (
	"cmp"
	"slices"
)

// Message mirrors com.mojang.brigadier.Message (only the members Component needs).
type Message interface {
	GetString() string
}

// ComponentUtilsDEFAULT_SEPARATOR_TEXT mirrors ComponentUtils.DEFAULT_SEPARATOR_TEXT.
const ComponentUtilsDEFAULT_SEPARATOR_TEXT = ", "

// ComponentUtilsDEFAULT_SEPARATOR mirrors ComponentUtils.DEFAULT_SEPARATOR.
var ComponentUtilsDEFAULT_SEPARATOR Component = ComponentLiteral(", ").WithStyleFormat(ChatFormattingGRAY)

// ComponentUtilsDEFAULT_NO_STYLE_SEPARATOR mirrors ComponentUtils.DEFAULT_NO_STYLE_SEPARATOR.
var ComponentUtilsDEFAULT_NO_STYLE_SEPARATOR Component = ComponentLiteral(", ")

// ComponentUtilsMergeStylesMutable mirrors ComponentUtils.mergeStyles(MutableComponent, Style).
func ComponentUtilsMergeStylesMutable(component *MutableComponent, style Style) *MutableComponent {
	if style.IsEmpty() {
		return component
	}
	inner := component.GetStyle()
	if inner.IsEmpty() {
		return component.SetStyle(style)
	}
	if inner.Equals(style) {
		return component
	}
	return component.SetStyle(inner.ApplyTo(style))
}

// ComponentUtilsMergeStyles mirrors ComponentUtils.mergeStyles(Component, Style).
func ComponentUtilsMergeStyles(component Component, style Style) Component {
	if style.IsEmpty() {
		return component
	}
	inner := component.GetStyle()
	if inner.IsEmpty() {
		return component.Copy().SetStyle(style)
	}
	if inner.Equals(style) {
		return component
	}
	return component.Copy().SetStyle(inner.ApplyTo(style))
}

// ComponentUtilsResolveOptional mirrors ComponentUtils.resolve(ResolutionContext, Optional<Component>, int).
func ComponentUtilsResolveOptional(context *ResolutionContext, component *Component, recursionDepth int) (*MutableComponent, error) {
	if component == nil {
		return nil, nil
	}
	return ComponentUtilsResolveDepth(context, *component, recursionDepth)
}

// ComponentUtilsResolve mirrors ComponentUtils.resolve(ResolutionContext, Component).
func ComponentUtilsResolve(context *ResolutionContext, component Component) (*MutableComponent, error) {
	return ComponentUtilsResolveDepth(context, component, 0)
}

// ComponentUtilsResolveDepth mirrors ComponentUtils.resolve(ResolutionContext, Component, int).
func ComponentUtilsResolveDepth(context *ResolutionContext, component Component, recursionDepth int) (*MutableComponent, error) {
	if recursionDepth > context.DepthLimit {
		switch context.DepthLimitBehavior {
		case ResolutionContextLimitBehaviorDISCARD_REMAINING:
			return CommonComponentsELLIPSIS.Copy(), nil
		default:
			return component.Copy(), nil
		}
	}
	result, err := component.GetContents().Resolve(context, recursionDepth+1)
	if err != nil {
		return nil, err
	}
	for _, sibling := range component.GetSiblings() {
		resolved, err := ComponentUtilsResolveDepth(context, sibling, recursionDepth+1)
		if err != nil {
			return nil, err
		}
		result.AppendComponent(resolved)
	}
	style, err := componentUtilsResolveStyle(context, component.GetStyle(), recursionDepth)
	if err != nil {
		return nil, err
	}
	return result.WithStyle(style), nil
}

func componentUtilsResolveStyle(context *ResolutionContext, style Style, recursionDepth int) (Style, error) {
	if showText, ok := style.GetHoverEvent().(HoverEventShowText); ok {
		resolved, err := ComponentUtilsResolveDepth(context, showText.Value, recursionDepth+1)
		if err != nil {
			return style, err
		}
		return style.WithHoverEvent(HoverEventShowText{Value: resolved}), nil
	}
	return style, nil
}

// ComponentUtilsFormatListStrings mirrors ComponentUtils.formatList(Collection<String>).
func ComponentUtilsFormatListStrings(values []string) Component {
	return ComponentUtilsFormatAndSortList(values, func(v string) Component {
		return ComponentLiteral(v).WithStyleFormat(ChatFormattingGREEN)
	})
}

// ComponentUtilsFormatAndSortList mirrors ComponentUtils.formatAndSortList.
func ComponentUtilsFormatAndSortList[T cmp.Ordered](values []T, formatter func(T) Component) Component {
	if len(values) == 0 {
		return CommonComponentsEMPTY
	}
	if len(values) == 1 {
		return formatter(values[0])
	}
	sorted := make([]T, len(values))
	copy(sorted, values)
	slices.Sort(sorted)
	return ComponentUtilsFormatList(sorted, formatter)
}

// ComponentUtilsFormatList mirrors ComponentUtils.formatList(Collection, Function).
func ComponentUtilsFormatList[T any](values []T, formatter func(T) Component) Component {
	return ComponentUtilsFormatListSeparator(values, ComponentUtilsDEFAULT_SEPARATOR, formatter)
}

// ComponentUtilsFormatListSeparator mirrors ComponentUtils.formatList(Collection, Component, Function).
func ComponentUtilsFormatListSeparator[T any](values []T, separator Component, formatter func(T) Component) *MutableComponent {
	if len(values) == 0 {
		return ComponentEmpty()
	}
	if len(values) == 1 {
		return formatter(values[0]).Copy()
	}
	result := ComponentEmpty()
	first := true
	for _, value := range values {
		if !first {
			result.AppendComponent(separator)
		}
		result.AppendComponent(formatter(value))
		first = false
	}
	return result
}

// ComponentUtilsFormatListComponents mirrors ComponentUtils.formatList(Collection<Component>, Component).
func ComponentUtilsFormatListComponents(values []Component, separator Component) Component {
	return ComponentUtilsFormatListSeparator(values, separator, func(c Component) Component { return c })
}

// ComponentUtilsWrapInSquareBrackets mirrors ComponentUtils.wrapInSquareBrackets.
func ComponentUtilsWrapInSquareBrackets(inner Component) *MutableComponent {
	return ComponentTranslatableArgs("chat.square_brackets", inner)
}

// ComponentUtilsFromMessage mirrors ComponentUtils.fromMessage.
func ComponentUtilsFromMessage(message Message) Component {
	if component, ok := message.(Component); ok {
		return component
	}
	return ComponentLiteral(message.GetString())
}

// ComponentUtilsIsTranslationResolvable mirrors ComponentUtils.isTranslationResolvable.
func ComponentUtilsIsTranslationResolvable(component Component) bool {
	if component == nil {
		return true
	}
	if translatable, ok := component.GetContents().(*TranslatableContents); ok {
		if translatable.GetFallback() != nil {
			return true
		}
		return LanguageGetInstance().Has(translatable.GetKey())
	}
	return true
}

// ComponentUtilsCopyOnClickText mirrors ComponentUtils.copyOnClickText.
func ComponentUtilsCopyOnClickText(text string) *MutableComponent {
	return ComponentUtilsWrapInSquareBrackets(
		ComponentLiteral(text).WithStyleFunc(func(s Style) Style {
			return s.WithColorFormat(&ChatFormattingGREEN).
				WithClickEvent(ClickEventCopyToClipboard{Value: text}).
				WithHoverEvent(HoverEventShowText{Value: ComponentTranslatable("chat.copy.click")}).
				WithInsertion(&text)
		}),
	)
}
