// Package resources mirrors net.minecraft.resources.
package resources

import (
	"fmt"
	"path/filepath"
	"strings"
)

// NAMESPACE_SEPARATOR mirrors Identifier.NAMESPACE_SEPARATOR.
const NAMESPACE_SEPARATOR = ':'

// DEFAULT_NAMESPACE mirrors Identifier.DEFAULT_NAMESPACE.
const DEFAULT_NAMESPACE = "minecraft"

// REALMS_NAMESPACE mirrors Identifier.REALMS_NAMESPACE.
const REALMS_NAMESPACE = "realms"

// ALLOWED_NAMESPACE_CHARACTERS mirrors Identifier.ALLOWED_NAMESPACE_CHARACTERS.
const ALLOWED_NAMESPACE_CHARACTERS = "[a-z0-9_.-]"

// Identifier mirrors net.minecraft.resources.Identifier.
type Identifier struct {
	namespace string
	path      string
}

// NewIdentifier mirrors the private Identifier(String, String) constructor.
func NewIdentifier(namespace string, path string) Identifier {
	assertValidNamespace(namespace, path)
	assertValidPath(namespace, path)
	return Identifier{namespace: namespace, path: path}
}

// createUntrusted mirrors the private Identifier.createUntrusted(String, String).
func createUntrusted(namespace string, path string) Identifier {
	assertValidNamespace(namespace, path)
	assertValidPath(namespace, path)
	return Identifier{namespace: namespace, path: path}
}

// IdentifierFromNamespaceAndPath mirrors Identifier.fromNamespaceAndPath(String, String).
func IdentifierFromNamespaceAndPath(namespace string, path string) Identifier {
	return createUntrusted(namespace, path)
}

// IdentifierParse mirrors Identifier.parse(String).
func IdentifierParse(identifier string) Identifier {
	return IdentifierBySeparator(identifier, NAMESPACE_SEPARATOR)
}

// IdentifierWithDefaultNamespace mirrors Identifier.withDefaultNamespace(String).
func IdentifierWithDefaultNamespace(path string) Identifier {
	return Identifier{namespace: DEFAULT_NAMESPACE, path: assertValidPath(DEFAULT_NAMESPACE, path)}
}

// IdentifierTryParse mirrors Identifier.tryParse(String).
func IdentifierTryParse(identifier string) *Identifier {
	return IdentifierTryBySeparator(identifier, NAMESPACE_SEPARATOR)
}

// IdentifierTryBuild mirrors Identifier.tryBuild(String, String).
func IdentifierTryBuild(namespace string, path string) *Identifier {
	if isValidNamespace(namespace) && isValidPath(path) {
		result := Identifier{namespace: namespace, path: path}
		return &result
	}
	return nil
}

// IdentifierBySeparator mirrors Identifier.bySeparator(String, char).
func IdentifierBySeparator(identifier string, separator byte) Identifier {
	separatorIndex := strings.IndexByte(identifier, separator)
	if separatorIndex >= 0 {
		path := identifier[separatorIndex+1:]
		if separatorIndex != 0 {
			namespace := identifier[:separatorIndex]
			return createUntrusted(namespace, path)
		}
		return IdentifierWithDefaultNamespace(path)
	}
	return IdentifierWithDefaultNamespace(identifier)
}

// IdentifierTryBySeparator mirrors Identifier.tryBySeparator(String, char).
func IdentifierTryBySeparator(identifier string, separator byte) *Identifier {
	separatorIndex := strings.IndexByte(identifier, separator)
	if separatorIndex >= 0 {
		path := identifier[separatorIndex+1:]
		if !isValidPath(path) {
			return nil
		} else if separatorIndex != 0 {
			namespace := identifier[:separatorIndex]
			if isValidNamespace(namespace) {
				result := Identifier{namespace: namespace, path: path}
				return &result
			}
			return nil
		}
		result := Identifier{namespace: DEFAULT_NAMESPACE, path: path}
		return &result
	}
	if isValidPath(identifier) {
		result := Identifier{namespace: DEFAULT_NAMESPACE, path: identifier}
		return &result
	}
	return nil
}

// Path mirrors Identifier.getPath().
func (i Identifier) Path() string {
	return i.path
}

// Namespace mirrors Identifier.getNamespace().
func (i Identifier) Namespace() string {
	return i.namespace
}

// WithPath mirrors Identifier.withPath(String).
func (i Identifier) WithPath(newPath string) Identifier {
	return Identifier{namespace: i.namespace, path: assertValidPath(i.namespace, newPath)}
}

// WithPrefix mirrors Identifier.withPrefix(String).
func (i Identifier) WithPrefix(prefix string) Identifier {
	return i.WithPath(prefix + i.path)
}

// WithSuffix mirrors Identifier.withSuffix(String).
func (i Identifier) WithSuffix(suffix string) Identifier {
	return i.WithPath(i.path + suffix)
}

// String mirrors Identifier.toString().
func (i Identifier) String() string {
	return i.namespace + ":" + i.path
}

// CompareTo mirrors Identifier.compareTo(Identifier).
func (i Identifier) CompareTo(other Identifier) int {
	result := strings.Compare(i.path, other.path)
	if result == 0 {
		result = strings.Compare(i.namespace, other.namespace)
	}
	return result
}

// ResolveAgainst mirrors Identifier.resolveAgainst(Path).
func (i Identifier) ResolveAgainst(root string) (string, error) {
	resultingPath := filepath.Join(root, i.Namespace(), i.Path())
	normalizedPath := filepath.Clean(resultingPath)
	normalizedRoot := filepath.Clean(root)
	if !strings.HasPrefix(normalizedPath, normalizedRoot) {
		return "", fmt.Errorf(
			"Identifier %q tried to access path %q from root %q",
			i.String(), normalizedPath, normalizedRoot,
		)
	}
	return resultingPath, nil
}

// ToDebugFileName mirrors Identifier.toDebugFileName().
func (i Identifier) ToDebugFileName() string {
	return strings.NewReplacer("/", "_", ":", "_").Replace(i.String())
}

// ToLanguageKey mirrors Identifier.toLanguageKey().
func (i Identifier) ToLanguageKey() string {
	return i.namespace + "." + i.path
}

// ToShortLanguageKey mirrors Identifier.toShortLanguageKey().
func (i Identifier) ToShortLanguageKey() string {
	if i.namespace == DEFAULT_NAMESPACE {
		return i.path
	}
	return i.ToLanguageKey()
}

// ToShortString mirrors Identifier.toShortString().
func (i Identifier) ToShortString() string {
	if i.namespace == DEFAULT_NAMESPACE {
		return i.path
	}
	return i.String()
}

// ToLanguageKeyWithPrefix mirrors Identifier.toLanguageKey(String).
func (i Identifier) ToLanguageKeyWithPrefix(prefix string) string {
	return prefix + "." + i.ToLanguageKey()
}

// ToLanguageKeyFull mirrors Identifier.toLanguageKey(String, String).
func (i Identifier) ToLanguageKeyFull(prefix string, suffix string) string {
	return prefix + "." + i.ToLanguageKey() + "." + suffix
}

// IsAllowedInIdentifier mirrors Identifier.isAllowedInIdentifier(char).
func IsAllowedInIdentifier(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c == '_' || c == ':' || c == '/' || c == '.' || c == '-'
}

// IsValidPath mirrors Identifier.isValidPath(String).
func IsValidPath(path string) bool {
	for index := 0; index < len(path); index++ {
		if !ValidPathChar(path[index]) {
			return false
		}
	}
	return true
}

// isValidPath mirrors the package-private Identifier.isValidPath(String) used
// internally.
func isValidPath(path string) bool {
	return IsValidPath(path)
}

// IsValidNamespace mirrors Identifier.isValidNamespace(String).
func IsValidNamespace(namespace string) bool {
	if namespace == ".." {
		return false
	}
	for index := 0; index < len(namespace); index++ {
		if !validNamespaceChar(namespace[index]) {
			return false
		}
	}
	return true
}

// isValidNamespace mirrors the package-private Identifier.isValidNamespace(String).
func isValidNamespace(namespace string) bool {
	return IsValidNamespace(namespace)
}

func assertValidNamespace(namespace string, path string) string {
	if !isValidNamespace(namespace) {
		panic(&IdentifierException{Message: "Non [a-z0-9_.-] character in namespace of identifier: " + namespace + ":" + path})
	}
	return namespace
}

// ValidPathChar mirrors Identifier.validPathChar(char).
func ValidPathChar(c byte) bool {
	return c == '_' || c == '-' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '/' || c == '.'
}

func validNamespaceChar(c byte) bool {
	return c == '_' || c == '-' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '.'
}

func assertValidPath(namespace string, path string) string {
	if !isValidPath(path) {
		panic(&IdentifierException{Message: "Non [a-z0-9/._-] character in path of location: " + namespace + ":" + path})
	}
	return path
}
