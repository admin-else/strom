package main

import (
	"bytes"
	"strings"
	"testing"
)

func protoField(name string, t any) map[string]any {
	return map[string]any{"name": name, "type": t}
}

// parameterizedTestProtocol models the protodef parameterization patterns the
// generator has to support: a `$compareTo` switch instantiated from a field of
// the enclosing container, a parameter forwarded through another parameterized
// type, and option/array wrappers around a parameterized type.
func parameterizedTestProtocol() Protocol {
	return Protocol{Types: Types{Types: map[string]any{
		"entityMetadataItem": []any{"switch", map[string]any{
			"compareTo": "$compareTo",
			"fields": map[string]any{
				"byte":  "i8",
				"int":   "varint",
				"float": "f32",
			},
		}},
		"entityMetadata": []any{"entityMetadataLoop", map[string]any{
			"endVal": 255,
			"type": []any{"container", []any{
				map[string]any{"anon": true, "type": []any{"container", []any{
					protoField("key", "u8"),
					protoField("type", []any{"mapper", map[string]any{
						"type":     "varint",
						"mappings": map[string]any{"0": "byte", "1": "int", "2": "float"},
					}}),
				}}},
				protoField("value", []any{"entityMetadataItem", map[string]any{"compareTo": "type"}}),
			}},
		}},
		"particleData": []any{"switch", map[string]any{
			"compareTo": "$compareTo",
			"fields": map[string]any{
				"1": "varint",
				"14": []any{"container", []any{
					protoField("red", "f32"),
					protoField("green", "f32"),
				}},
			},
			"default": "void",
		}},
		"particle": []any{"container", []any{
			protoField("particleId", "varint"),
			protoField("data", []any{"particleData", map[string]any{"compareTo": "particleId"}}),
		}},
		"item": []any{"switch", map[string]any{
			"compareTo": "$compareTo",
			"fields":    map[string]any{"1": "varint", "2": "i8"},
		}},
		"inner": []any{"container", []any{
			protoField("value", []any{"item", map[string]any{"compareTo": "$compareTo"}}),
		}},
		"outer": []any{"container", []any{
			protoField("kind", "varint"),
			protoField("wrapped", []any{"inner", map[string]any{"compareTo": "kind"}}),
			protoField("maybe", []any{"option", []any{"item", map[string]any{"compareTo": "kind"}}}),
			protoField("list", []any{"array", map[string]any{
				"countType": "varint",
				"type":      []any{"item", map[string]any{"compareTo": "kind"}},
			}}),
		}},
	}}}
}

func TestParameterizedTypesGenerateConcreteDecoders(t *testing.T) {
	var buf bytes.Buffer
	if _, err := GenerateFromProtocol(parameterizedTestProtocol(), "1.20.4", &buf, "test"); err != nil {
		t.Fatalf("GenerateFromProtocol: %v", err)
	}
	out := buf.String()

	if strings.Contains(out, "ToDoError") {
		t.Fatalf("generated code fell back to ToDoError:\n%s", out)
	}

	// Abstract parameterized types must not be emitted standalone.
	for _, decl := range []string{
		"type EntityMetadataItem struct",
		"type EntityMetadataItem ",
		"type ParticleData struct",
		"type Item struct",
		"type Inner struct",
	} {
		if strings.Contains(out, decl) {
			t.Errorf("standalone parameterized type was emitted: %q", decl)
		}
	}

	// $compareTo resolves against the enclosing container's anon field.
	if !strings.Contains(out, "switch entry.Anon.Type {") {
		t.Errorf("entity metadata switch did not resolve $compareTo to entry.Anon.Type")
	}
	if !strings.Contains(out, "case \"byte\":") {
		t.Errorf("entity metadata switch is missing its string cases")
	}

	// $compareTo resolves against the particle id field.
	if !strings.Contains(out, "switch ret.ParticleId {") {
		t.Errorf("particle switch did not resolve $compareTo to ret.ParticleId")
	}
	if !strings.Contains(out, "case 14:") {
		t.Errorf("particle switch is missing its numeric cases")
	}

	// The parameter forwarded through `inner` resolves in `outer`.
	if count := strings.Count(out, "switch ret.Kind {"); count < 3 {
		t.Errorf("expected the forwarded parameter to resolve in wrapped/maybe/list, got %d switches on ret.Kind", count)
	}
}

func TestContainsTypeParameter(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want bool
	}{
		{"plain string", "varint", false},
		{"placeholder", "$compareTo", true},
		{"nested", []any{"switch", map[string]any{"compareTo": "$compareTo"}}, true},
		{"no placeholder nested", []any{"container", []any{map[string]any{"name": "x", "type": "u8"}}}, false},
	}
	for _, c := range cases {
		if got := containsTypeParameter(c.v); got != c.want {
			t.Errorf("%s: containsTypeParameter = %v, want %v", c.name, got, c.want)
		}
	}
}
