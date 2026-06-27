package mapstructure

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type Inner struct {
	Value int32 `nbt:"Value"`
}

type Example struct {
	Name    string `nbt:"Name"`
	Age     int32  `nbt:"Age,omitempty"`
	Secret  string `nbt:"Secret,required"`
	Skipped string `nbt:"-"`
	NoTag   string
	Nested  Inner `nbt:"Nested"`
	Flag    bool  `nbt:"Flag"`
}

func TestBasicDecode(t *testing.T) {
	f := NewFormat("nbt")
	data := map[string]any{
		"Name":   "test",
		"Secret": "shh",
		"NoTag":  "fallback",
		"Nested": map[string]any{"Value": int32(42)},
		"Flag":   int8(1),
	}
	var got Example
	if err := f.Decode(data, &got); err != nil {
		t.Fatal(err)
	}
	want := Example{Name: "test", Secret: "shh", NoTag: "fallback", Nested: Inner{Value: 42}, Flag: true}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatal(diff)
	}
}

func TestMissingRequired(t *testing.T) {
	f := NewFormat("nbt")
	data := map[string]any{"Name": "test"}
	var got Example
	if err := f.Decode(data, &got); err == nil {
		t.Fatal("expected error for missing required field")
	}
}

func TestOmitEmpty(t *testing.T) {
	f := NewFormat("nbt")
	data := map[string]any{"Secret": "shh"}
	var got Example
	if err := f.Decode(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Age != 0 {
		t.Fatal("omitempty field should stay zero when absent")
	}
}

func TestRequireAll(t *testing.T) {
	type Simple struct {
		A int32 `nbt:"A"`
		B int32 `nbt:"B,omitempty"`
		C int32 `nbt:"C"`
	}
	f := NewFormat("nbt", WithRequireAll())
	data := map[string]any{"A": int32(1), "C": int32(3)}
	var got Simple
	if err := f.Decode(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.A != 1 || got.C != 3 {
		t.Fatal("unexpected values")
	}
}

func TestRequireAllMissing(t *testing.T) {
	type Simple struct {
		A int32 `nbt:"A"`
		B int32 `nbt:"B"`
	}
	f := NewFormat("nbt", WithRequireAll())
	data := map[string]any{"A": int32(1)}
	var got Simple
	if err := f.Decode(data, &got); err == nil {
		t.Fatal("expected error for missing field with RequireAll")
	}
}

func TestRequireAllOmitEmptyWins(t *testing.T) {
	type Simple struct {
		A int32 `nbt:"A"`
		B int32 `nbt:"B,omitempty"`
	}
	f := NewFormat("nbt", WithRequireAll())
	data := map[string]any{"A": int32(1)}
	var got Simple
	if err := f.Decode(data, &got); err != nil {
		t.Fatal("omitempty should override RequireAll")
	}
}

func TestSkipTag(t *testing.T) {
	f := NewFormat("nbt")
	data := map[string]any{"Name": "test", "Secret": "shh", "Skipped": "should not appear"}
	var got Example
	if err := f.Decode(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Skipped != "" {
		t.Fatal("skipped field should not be set")
	}
}

func TestErrorOnExtra(t *testing.T) {
	f := NewFormat("nbt", WithErrorOnExtra())
	data := map[string]any{"Name": "test", "Secret": "shh", "Unknown": "extra"}
	var got Example
	if err := f.Decode(data, &got); err == nil {
		t.Fatal("expected error for extra key")
	}
}

func TestNoErrorOnExtraByDefault(t *testing.T) {
	f := NewFormat("nbt")
	data := map[string]any{"Name": "test", "Secret": "shh", "Unknown": "extra"}
	var got Example
	if err := f.Decode(data, &got); err != nil {
		t.Fatal("extra keys should be ignored by default")
	}
}

func TestNotPointer(t *testing.T) {
	f := NewFormat("nbt")
	data := map[string]any{}
	var got Example
	if err := f.Decode(data, got); err == nil {
		t.Fatal("expected error for non-pointer target")
	}
}

func TestNotStruct(t *testing.T) {
	f := NewFormat("nbt")
	data := map[string]any{}
	var got int32
	if err := f.Decode(data, &got); err == nil {
		t.Fatal("expected error for non-struct target")
	}
}

func TestNotMap(t *testing.T) {
	f := NewFormat("nbt")
	var got Example
	if err := f.Decode("not a map", &got); err == nil {
		t.Fatal("expected error for non-map data")
	}
}

func TestTagMapping(t *testing.T) {
	type Renamed struct {
		GoName string `nbt:"MapKey"`
	}
	f := NewFormat("nbt")
	data := map[string]any{"MapKey": "hello"}
	var got Renamed
	if err := f.Decode(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.GoName != "hello" {
		t.Fatal("tag-mapped field not set correctly")
	}
}

func TestTypeConversion(t *testing.T) {
	type Conv struct {
		I8   int8    `nbt:"I8"`
		I16  int16   `nbt:"I16"`
		I32  int32   `nbt:"I32"`
		I64  int64   `nbt:"I64"`
		F32  float32 `nbt:"F32"`
		F64  float64 `nbt:"F64"`
		Bool bool    `nbt:"Bool"`
	}
	f := NewFormat("nbt")
	data := map[string]any{
		"I8":   int8(1),
		"I16":  int16(2),
		"I32":  int32(3),
		"I64":  int64(4),
		"F32":  float32(5.5),
		"F64":  float64(6.6),
		"Bool": int8(1),
	}
	var got Conv
	if err := f.Decode(data, &got); err != nil {
		t.Fatal(err)
	}
	want := Conv{I8: 1, I16: 2, I32: 3, I64: 4, F32: 5.5, F64: 6.6, Bool: true}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatal(diff)
	}
}

func TestTypeConversionWidening(t *testing.T) {
	type Wide struct {
		I32FromI8  int32 `nbt:"A"`
		I64FromI16 int64 `nbt:"B"`
	}
	f := NewFormat("nbt")
	data := map[string]any{
		"A": int8(10),
		"B": int16(20),
	}
	var got Wide
	if err := f.Decode(data, &got); err != nil {
		t.Fatal(err)
	}
	want := Wide{I32FromI8: 10, I64FromI16: 20}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatal(diff)
	}
}

func TestSliceDecode(t *testing.T) {
	type WithSlice struct {
		Items []int32 `nbt:"Items"`
	}
	f := NewFormat("nbt")
	data := map[string]any{
		"Items": []any{int32(1), int32(2), int32(3)},
	}
	var got WithSlice
	if err := f.Decode(data, &got); err != nil {
		t.Fatal(err)
	}
	want := WithSlice{Items: []int32{1, 2, 3}}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatal(diff)
	}
}

func TestNestedStructDecode(t *testing.T) {
	type Outer struct {
		Inner Inner `nbt:"Inner"`
	}
	f := NewFormat("nbt")
	data := map[string]any{
		"Inner": map[string]any{"Value": int32(99)},
	}
	var got Outer
	if err := f.Decode(data, &got); err != nil {
		t.Fatal(err)
	}
	want := Outer{Inner: Inner{Value: 99}}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatal(diff)
	}
}

func TestEncode(t *testing.T) {
	f := NewFormat("nbt")
	s := Example{Name: "test", Secret: "shh", NoTag: "fallback", Nested: Inner{Value: 42}, Flag: true}
	got, err := f.Encode(s)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{
		"Name":   "test",
		"Secret": "shh",
		"NoTag":  "fallback",
		"Nested": map[string]any{"Value": int32(42)},
		"Flag":   int8(1),
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatal(diff)
	}
}

func TestEncodeOmitEmpty(t *testing.T) {
	f := NewFormat("nbt")
	s := Example{Name: "test", Secret: "shh"}
	got, err := f.Encode(s)
	if err != nil {
		t.Fatal(err)
	}
	if _, exists := got["Age"]; exists {
		t.Fatal("omitempty zero field should not appear in encoded map")
	}
}

func TestEncodeSkipTag(t *testing.T) {
	f := NewFormat("nbt")
	s := Example{Name: "test", Secret: "shh", Skipped: "hidden"}
	got, err := f.Encode(s)
	if err != nil {
		t.Fatal(err)
	}
	if _, exists := got["Skipped"]; exists {
		t.Fatal("skipped field should not appear in encoded map")
	}
}

func TestEncodePointer(t *testing.T) {
	f := NewFormat("nbt")
	s := &Example{Name: "test", Secret: "shh"}
	got, err := f.Encode(s)
	if err != nil {
		t.Fatal(err)
	}
	if got["Name"] != "test" {
		t.Fatal("encode from pointer should work")
	}
}

func TestEncodeSlice(t *testing.T) {
	type WithSlice struct {
		Items []int32 `nbt:"Items"`
	}
	f := NewFormat("nbt")
	s := WithSlice{Items: []int32{1, 2, 3}}
	got, err := f.Encode(s)
	if err != nil {
		t.Fatal(err)
	}
	items, ok := got["Items"].([]any)
	if !ok {
		t.Fatal("slice not encoded correctly")
	}
	if len(items) != 3 || items[0].(int32) != 1 {
		t.Fatal("slice values not correct")
	}
}

func TestRoundTrip(t *testing.T) {
	f := NewFormat("nbt")
	original := Example{Name: "hello", Secret: "world", Nested: Inner{Value: 99}, Flag: true}
	encoded, err := f.Encode(original)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Example
	if err := f.Decode(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(original, decoded); diff != "" {
		t.Fatal(diff)
	}
}

func TestTypeCodecDecode(t *testing.T) {
	type GameMode int32

	type Player struct {
		Mode GameMode `nbt:"Mode"`
	}

	f := NewFormat("nbt",
		WithTypeCodec(
			func(v any, tag string) (GameMode, error) {
				i, ok := v.(int32)
				if !ok {
					return 0, fmt.Errorf("expected int32, got %T", v)
				}
				return GameMode(i), nil
			},
			func(v GameMode, tag string) (any, error) {
				return int32(v), nil
			},
		),
	)

	data := map[string]any{"Mode": int32(1)}
	var got Player
	if err := f.Decode(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Mode != GameMode(1) {
		t.Fatal("type codec decode did not work")
	}
}

func TestTypeCodecEncode(t *testing.T) {
	type GameMode int32

	type Player struct {
		Mode GameMode `nbt:"Mode"`
	}

	f := NewFormat("nbt",
		WithTypeCodec(
			func(v any, tag string) (GameMode, error) {
				i, ok := v.(int32)
				if !ok {
					return 0, fmt.Errorf("expected int32, got %T", v)
				}
				return GameMode(i), nil
			},
			func(v GameMode, tag string) (any, error) {
				return int32(v), nil
			},
		),
	)

	s := Player{Mode: GameMode(2)}
	got, err := f.Encode(s)
	if err != nil {
		t.Fatal(err)
	}
	if got["Mode"].(int32) != 2 {
		t.Fatal("type codec encode did not work")
	}
}

func TestTypeCodecRoundTrip(t *testing.T) {
	type GameMode int32

	type Player struct {
		Mode GameMode `nbt:"Mode"`
	}

	f := NewFormat("nbt",
		WithTypeCodec(
			func(v any, tag string) (GameMode, error) {
				i, ok := v.(int32)
				if !ok {
					return 0, fmt.Errorf("expected int32, got %T", v)
				}
				return GameMode(i), nil
			},
			func(v GameMode, tag string) (any, error) {
				return int32(v), nil
			},
		),
	)

	original := Player{Mode: GameMode(3)}
	encoded, err := f.Encode(original)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Player
	if err := f.Decode(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(original, decoded); diff != "" {
		t.Fatal(diff)
	}
}

func TestTypeCodecReceivesTag(t *testing.T) {
	type GameMode int32

	type Player struct {
		Mode    GameMode `nbt:"Mode,strict"`
		Friends GameMode `nbt:"Friends"`
	}

	strictModes := map[string]bool{}

	f := NewFormat("nbt",
		WithTypeCodec(
			func(v any, tag string) (GameMode, error) {
				strictModes[tag] = true
				i, ok := v.(int32)
				if !ok {
					return 0, fmt.Errorf("expected int32, got %T", v)
				}
				return GameMode(i), nil
			},
			func(v GameMode, tag string) (any, error) {
				return int32(v), nil
			},
		),
	)

	data := map[string]any{"Mode": int32(1), "Friends": int32(2)}
	var got Player
	if err := f.Decode(data, &got); err != nil {
		t.Fatal(err)
	}
	if !strictModes["Mode,strict"] {
		t.Fatalf("expected tag 'Mode,strict' to be passed to codec, got keys: %v", strictModes)
	}
}

func TestMultipleTypeCodecs(t *testing.T) {
	type GameMode int32
	type Dimension int32

	type Player struct {
		Mode GameMode  `nbt:"Mode"`
		Dim  Dimension `nbt:"Dim"`
	}

	f := NewFormat("nbt",
		WithTypeCodec(
			func(v any, tag string) (GameMode, error) {
				i, _ := v.(int32)
				return GameMode(i), nil
			},
			func(v GameMode, tag string) (any, error) {
				return int32(v), nil
			},
		),
		WithTypeCodec(
			func(v any, tag string) (Dimension, error) {
				i, _ := v.(int32)
				return Dimension(i), nil
			},
			func(v Dimension, tag string) (any, error) {
				return int32(v), nil
			},
		),
	)

	data := map[string]any{"Mode": int32(1), "Dim": int32(2)}
	var got Player
	if err := f.Decode(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Mode != GameMode(1) || got.Dim != Dimension(2) {
		t.Fatal("multiple type codecs did not work")
	}
}

func TestUnexportedFieldsIgnored(t *testing.T) {
	type WithUnexported struct {
		Exported   string `nbt:"Exported"`
		unexported string
	}
	f := NewFormat("nbt")
	data := map[string]any{"Exported": "yes", "unexported": "no"}
	var got WithUnexported
	if err := f.Decode(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.unexported != "" {
		t.Fatal("unexported fields should be ignored")
	}
}

func TestEncodeMapField(t *testing.T) {
	type WithMap struct {
		Props map[string]string `nbt:"Props"`
	}
	f := NewFormat("nbt")
	s := WithMap{Props: map[string]string{"key": "val"}}
	got, err := f.Encode(s)
	if err != nil {
		t.Fatal(err)
	}
	props, ok := got["Props"].(map[string]any)
	if !ok {
		t.Fatal("map field not encoded correctly")
	}
	if props["key"] != "val" {
		t.Fatal("map field value incorrect")
	}
}

func TestDecodeMapField(t *testing.T) {
	type WithMap struct {
		Props map[string]string `nbt:"Props"`
	}
	f := NewFormat("nbt")
	data := map[string]any{
		"Props": map[string]any{"key": "val"},
	}
	var got WithMap
	if err := f.Decode(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Props["key"] != "val" {
		t.Fatal("map field value incorrect")
	}
}

func TestDecodeBoolFromInt32(t *testing.T) {
	type WithBool struct {
		Flag bool `nbt:"Flag"`
	}
	f := NewFormat("nbt")
	data := map[string]any{"Flag": int32(1)}
	var got WithBool
	if err := f.Decode(data, &got); err != nil {
		t.Fatal(err)
	}
	if !got.Flag {
		t.Fatal("bool from int32 should be true")
	}
}
