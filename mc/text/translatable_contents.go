package text

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// TranslatableContentsNO_ARGS mirrors TranslatableContents.NO_ARGS.
var TranslatableContentsNO_ARGS = []any{}

var translatableContentsFormatPattern = regexp.MustCompile(`%(?:(\d+)\$)?([A-Za-z%]|$)`)

var (
	translatableContentsTextPercent = FormattedTextOf("%")
	translatableContentsTextNull    = FormattedTextOf("null")
)

// TranslatableContents mirrors net.minecraft.network.chat.contents.TranslatableContents.
type TranslatableContents struct {
	BaseComponentContents
	key             string
	fallback        *string
	args            []any
	decomposedWith  Language
	decomposedParts []FormattedText
}

// TranslatableFormatException mirrors net.minecraft.network.chat.contents.TranslatableFormatException.
type TranslatableFormatException struct {
	Message string
}

func (e *TranslatableFormatException) Error() string { return e.Message }

func newTranslatableFormatExceptionMessage(component *TranslatableContents, message string) *TranslatableFormatException {
	return &TranslatableFormatException{Message: fmt.Sprintf("Error parsing: %s: %s", component.String(), message)}
}

func newTranslatableFormatExceptionIndex(component *TranslatableContents, index int) *TranslatableFormatException {
	return &TranslatableFormatException{Message: fmt.Sprintf("Invalid index %d requested for %s", index, component.String())}
}

func newTranslatableFormatExceptionCause(component *TranslatableContents, cause error) *TranslatableFormatException {
	return &TranslatableFormatException{Message: fmt.Sprintf("Error while parsing: %s", component.String())}
}

// TranslatableContentsIsAllowedPrimitiveArgument mirrors TranslatableContents.isAllowedPrimitiveArgument.
func TranslatableContentsIsAllowedPrimitiveArgument(object any) bool {
	switch object.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool, string:
		return true
	default:
		return false
	}
}

// NewTranslatableContents mirrors the TranslatableContents constructor.
func NewTranslatableContents(key string, fallback *string, args ...any) *TranslatableContents {
	if args == nil {
		args = TranslatableContentsNO_ARGS
	}
	return &TranslatableContents{
		key:             key,
		fallback:        fallback,
		args:            args,
		decomposedParts: []FormattedText{},
	}
}

// GetKey mirrors TranslatableContents.getKey().
func (t *TranslatableContents) GetKey() string { return t.key }

// GetFallback mirrors TranslatableContents.getFallback().
func (t *TranslatableContents) GetFallback() *string { return t.fallback }

// GetArgs mirrors TranslatableContents.getArgs().
func (t *TranslatableContents) GetArgs() []any { return t.args }

func (t *TranslatableContents) decompose() {
	currentLanguage := LanguageGetInstance()
	if currentLanguage == t.decomposedWith {
		return
	}
	t.decomposedWith = currentLanguage

	var format string
	if t.fallback != nil {
		format = currentLanguage.GetOrDefaultValue(t.key, *t.fallback)
	} else {
		format = currentLanguage.GetOrDefault(t.key)
	}

	parts := []FormattedText{}
	err := t.decomposeTemplate(format, func(part FormattedText) {
		parts = append(parts, part)
	})
	if err != nil {
		t.decomposedParts = []FormattedText{FormattedTextOf(format)}
	} else {
		t.decomposedParts = parts
	}
}

func (t *TranslatableContents) decomposeTemplate(template string, decomposedParts func(FormattedText)) error {
	matches := translatableContentsFormatPattern.FindAllStringSubmatchIndex(template, -1)

	replacementIndex := 0
	current := 0

	for _, match := range matches {
		start := match[0]
		end := match[1]
		if start < current {
			continue
		}
		if start > current {
			prefix := template[current:start]
			if strings.IndexByte(prefix, 37) != -1 {
				return newTranslatableFormatExceptionCause(t, fmt.Errorf("illegal format"))
			}
			decomposedParts(FormattedTextOf(prefix))
		}

		formatType := ""
		if match[4] != -1 {
			formatType = template[match[4]:match[5]]
		}
		formatString := template[start:end]
		if formatType == "%" && formatString == "%%" {
			decomposedParts(translatableContentsTextPercent)
		} else {
			if formatType != "s" {
				return newTranslatableFormatExceptionMessage(t, "Unsupported format: '"+formatString+"'")
			}

			possiblePositionIndex := ""
			if match[2] != -1 {
				possiblePositionIndex = template[match[2]:match[3]]
			}
			var index int
			if possiblePositionIndex != "" {
				parsed, parseErr := strconv.Atoi(possiblePositionIndex)
				if parseErr != nil {
					return newTranslatableFormatExceptionCause(t, parseErr)
				}
				index = parsed - 1
			} else {
				index = replacementIndex
				replacementIndex++
			}
			arg, argErr := t.getArgument(index)
			if argErr != nil {
				return argErr
			}
			decomposedParts(arg)
		}

		current = end
	}

	if current < len(template) {
		tail := template[current:]
		if strings.IndexByte(tail, 37) != -1 {
			return newTranslatableFormatExceptionCause(t, fmt.Errorf("illegal format"))
		}
		decomposedParts(FormattedTextOf(tail))
	}

	return nil
}

func (t *TranslatableContents) getArgument(index int) (FormattedText, error) {
	if index >= 0 && index < len(t.args) {
		arg := t.args[index]
		if componentArg, ok := arg.(Component); ok {
			return componentArg, nil
		}
		if arg == nil {
			return translatableContentsTextNull, nil
		}
		return FormattedTextOf(fmt.Sprintf("%v", arg)), nil
	}
	return nil, newTranslatableFormatExceptionIndex(t, index)
}

// VisitStyled mirrors TranslatableContents.visit(StyledContentConsumer, Style).
func (t *TranslatableContents) VisitStyled(output StyledContentConsumer, currentStyle Style) (any, bool) {
	t.decompose()
	for _, part := range t.decomposedParts {
		if value, present := part.VisitStyled(output, currentStyle); present {
			return value, true
		}
	}
	return nil, false
}

// VisitContent mirrors TranslatableContents.visit(ContentConsumer).
func (t *TranslatableContents) VisitContent(output ContentConsumer) (any, bool) {
	t.decompose()
	for _, part := range t.decomposedParts {
		if value, present := part.VisitContent(output); present {
			return value, true
		}
	}
	return nil, false
}

// Resolve mirrors TranslatableContents.resolve.
func (t *TranslatableContents) Resolve(context *ResolutionContext, recursionDepth int) (*MutableComponent, error) {
	argsCopy := make([]any, len(t.args))
	for i, param := range t.args {
		if component, ok := param.(Component); ok {
			resolved, err := ComponentUtilsResolveDepth(context, component, recursionDepth)
			if err != nil {
				return nil, err
			}
			argsCopy[i] = resolved
		} else {
			argsCopy[i] = param
		}
	}
	return MutableComponentCreate(NewTranslatableContents(t.key, t.fallback, argsCopy...)), nil
}

// Equals mirrors TranslatableContents.equals.
func (t *TranslatableContents) Equals(o *TranslatableContents) bool {
	if t == o {
		return true
	}
	if o == nil {
		return false
	}
	if t.key != o.key || !stringPtrEqual(t.fallback, o.fallback) {
		return false
	}
	return anySliceEqual(t.args, o.args)
}

// HashCode mirrors TranslatableContents.hashCode.
func (t *TranslatableContents) HashCode() int {
	result := stringHash(t.key)
	result = 31*result + stringPtrHash(t.fallback)
	result = 31*result + anySliceHash(t.args)
	return result
}

func (t *TranslatableContents) String() string {
	fallback := ""
	if t.fallback != nil {
		fallback = ", fallback='" + *t.fallback + "'"
	}
	return "translation{key='" + t.key + "'" + fallback + ", args=" + javaArraysToString(t.args) + "}"
}

func javaArraysToString(args []any) string {
	parts := make([]string, len(args))
	for i, arg := range args {
		if arg == nil {
			parts[i] = "null"
		} else {
			parts[i] = fmt.Sprintf("%v", arg)
		}
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

func anySliceEqual(a, b []any) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !translationArgEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func translationArgEqual(a, b any) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	if ca, ok := a.(Component); ok {
		cb, ok := b.(Component)
		return ok && componentEquals(ca, cb)
	}
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func anySliceHash(a []any) int {
	h := 1
	for _, v := range a {
		h = 31*h + stringHash(fmt.Sprintf("%v", v))
	}
	return h
}
