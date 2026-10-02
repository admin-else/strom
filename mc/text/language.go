package text

// Language mirrors net.minecraft.locale.Language. The abstract class is modelled
// as an interface plus a package-global instance, mirroring the static
// getInstance()/inject() pair. Resource loading of en_us.json is deferred
// (see docs/QUESTIONS.md); the default instance has empty storage.
type Language interface {
	GetOrDefault(elementID string) string
	GetOrDefaultValue(elementID string, defaultValue string) string
	Has(elementID string) bool
	IsDefaultRightToLeft() bool
	GetVisualOrder(logicalOrderText FormattedText) FormattedCharSequence
	GetVisualOrderLines(lines []FormattedText) []FormattedCharSequence
}

// LanguageDEFAULT mirrors Language.DEFAULT.
const LanguageDEFAULT = "en_us"

type languageDefault struct {
	storage     map[string]string
	rightToLeft bool
}

// LanguageDEFAULT_INSTANCE mirrors Language.DEFAULT_INSTANCE.
var LanguageDEFAULT_INSTANCE Language = &languageDefault{storage: map[string]string{}}

var languageInstance Language = LanguageDEFAULT_INSTANCE

// LanguageGetInstance mirrors Language.getInstance().
func LanguageGetInstance() Language { return languageInstance }

// LanguageInject mirrors Language.inject(Language).
func LanguageInject(language Language) { languageInstance = language }

// NewLanguage builds a Language from a translation map. It mirrors the anonymous
// Language created by loadDefault() once the resource store has been read.
func NewLanguage(storage map[string]string, rightToLeft bool) Language {
	copied := make(map[string]string, len(storage))
	for k, v := range storage {
		copied[k] = v
	}
	return &languageDefault{storage: copied, rightToLeft: rightToLeft}
}

func (l *languageDefault) GetOrDefault(elementID string) string {
	return l.GetOrDefaultValue(elementID, elementID)
}

func (l *languageDefault) GetOrDefaultValue(elementID string, defaultValue string) string {
	if v, ok := l.storage[elementID]; ok {
		return v
	}
	return defaultValue
}

func (l *languageDefault) Has(elementID string) bool {
	_, ok := l.storage[elementID]
	return ok
}

func (l *languageDefault) IsDefaultRightToLeft() bool { return l.rightToLeft }

func (l *languageDefault) GetVisualOrder(logicalOrderText FormattedText) FormattedCharSequence {
	return func(output FormattedCharSink) bool {
		value, present := logicalOrderText.VisitStyled(func(style Style, contents string) (any, bool) {
			if StringDecomposerIterateFormattedOffset(contents, 0, style, output) {
				return nil, false
			}
			return FormattedTextSTOP_ITERATION, true
		}, StyleEMPTY)
		_ = value
		return present
	}
}

func (l *languageDefault) GetVisualOrderLines(lines []FormattedText) []FormattedCharSequence {
	result := make([]FormattedCharSequence, len(lines))
	for i, line := range lines {
		result[i] = l.GetVisualOrder(line)
	}
	return result
}
