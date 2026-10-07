package main

import (
	"fmt"
	"strings"
)

// containsTypeParameter reports whether a protodef type definition references a
// `$name` parameter placeholder anywhere inside it.
func containsTypeParameter(v any) (found bool) {
	switch t := v.(type) {
	case string:
		return strings.HasPrefix(t, "$")
	case []any:
		for _, e := range t {
			if containsTypeParameter(e) {
				return true
			}
		}
	case map[string]any:
		for _, e := range t {
			if containsTypeParameter(e) {
				return true
			}
		}
	}
	return
}

// pushTypeParameters binds the instantiation arguments of a named type. The
// returned func pops the scope and must be called by the caller. Argument values
// are resolved against the scopes already active at the call site, so a
// parameter can be forwarded through a chain of parameterized types.
func (g *Generator) pushTypeParameters(args any) (pop func()) {
	m, ok := args.(map[string]any)
	if !ok {
		return func() {}
	}
	bound := make(map[string]any, len(m))
	for k, v := range m {
		resolved, err := g.resolveTypeParameter(v)
		if err != nil {
			resolved = v
		}
		bound[k] = resolved
	}
	g.TypeParameters = append(g.TypeParameters, bound)
	return func() {
		g.TypeParameters = g.TypeParameters[:len(g.TypeParameters)-1]
	}
}

// lookupTypeParameter resolves a bare parameter name against the innermost
// instantiation scope that binds it.
func (g *Generator) lookupTypeParameter(name string) (v any, ok bool) {
	for i := len(g.TypeParameters) - 1; i >= 0; i-- {
		if v, ok := g.TypeParameters[i][name]; ok {
			return v, true
		}
	}
	return nil, false
}

// resolveTypeParameter replaces a `$name` placeholder with its bound value,
// following chains of placeholders. Non-placeholders are returned unchanged.
func (g *Generator) resolveTypeParameter(v any) (resolved any, err error) {
	for depth := 0; depth < 64; depth++ {
		s, ok := v.(string)
		if !ok || !strings.HasPrefix(s, "$") {
			return v, nil
		}
		next, found := g.lookupTypeParameter(s[1:])
		if !found {
			return nil, fmt.Errorf("unbound type parameter: %s", s)
		}
		v = next
	}
	return nil, fmt.Errorf("type parameter cycle at %v", v)
}

// instantiateNamed reports whether tName is a named parameterized type used with
// arguments, returning its definition so the caller can generate it with the
// arguments bound as the active type parameters.
func (g *Generator) instantiateNamed(tName string, args any) (def any, ok bool) {
	if args == nil || !g.ParameterizedTypes[tName] {
		return nil, false
	}
	def, found := g.Protocol.Types.Types[tName]
	if !found || def == "native" {
		return nil, false
	}
	return def, true
}
