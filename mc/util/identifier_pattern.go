package util

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/admin-else/strom/mc/resources"
)

// IdentifierPattern mirrors net.minecraft.util.IdentifierPattern.
type IdentifierPattern struct {
	namespacePattern   *regexp.Regexp
	namespacePredicate func(string) bool
	pathPattern        *regexp.Regexp
	pathPredicate      func(string) bool
	locationPredicate  func(resources.Identifier) bool
}

// NewIdentifierPattern mirrors the private IdentifierPattern(Optional<Pattern>,
// Optional<Pattern>) constructor. Either pattern may be nil.
func NewIdentifierPattern(namespacePattern *regexp.Regexp, pathPattern *regexp.Regexp) IdentifierPattern {
	namespacePredicate := func(s string) bool { return true }
	if namespacePattern != nil {
		namespacePredicate = namespacePattern.MatchString
	}
	pathPredicate := func(s string) bool { return true }
	if pathPattern != nil {
		pathPredicate = pathPattern.MatchString
	}
	locationPredicate := func(id resources.Identifier) bool {
		return namespacePredicate(id.Namespace()) && pathPredicate(id.Path())
	}
	return IdentifierPattern{
		namespacePattern:   namespacePattern,
		namespacePredicate: namespacePredicate,
		pathPattern:        pathPattern,
		pathPredicate:      pathPredicate,
		locationPredicate:  locationPredicate,
	}
}

// ParseIdentifierPattern mirrors IdentifierPattern.CODEC for a single
// identifier pattern: an object with optional regex `namespace` and `path`
// fields. ExtraCodecs.PATTERN compiles the strings with Pattern.compile.
func ParseIdentifierPattern(raw json.RawMessage) (pattern IdentifierPattern, err error) {
	var fields struct {
		Namespace *string `json:"namespace"`
		Path      *string `json:"path"`
	}
	if err = json.Unmarshal(raw, &fields); err != nil {
		return pattern, err
	}
	var namespacePattern *regexp.Regexp
	if fields.Namespace != nil {
		namespacePattern, err = compilePattern(*fields.Namespace)
		if err != nil {
			return pattern, err
		}
	}
	var pathPattern *regexp.Regexp
	if fields.Path != nil {
		pathPattern, err = compilePattern(*fields.Path)
		if err != nil {
			return pattern, err
		}
	}
	return NewIdentifierPattern(namespacePattern, pathPattern), nil
}

func compilePattern(pattern string) (compiled *regexp.Regexp, err error) {
	compiled, err = regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("Invalid regex pattern '%s': %s", pattern, err.Error())
	}
	return compiled, nil
}

// NamespacePredicate mirrors IdentifierPattern.namespacePredicate().
func (p IdentifierPattern) NamespacePredicate() func(string) bool {
	return p.namespacePredicate
}

// PathPredicate mirrors IdentifierPattern.pathPredicate().
func (p IdentifierPattern) PathPredicate() func(string) bool {
	return p.pathPredicate
}

// LocationPredicate mirrors IdentifierPattern.locationPredicate().
func (p IdentifierPattern) LocationPredicate() func(resources.Identifier) bool {
	return p.locationPredicate
}
