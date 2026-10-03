package nbt

import (
	"bytes"
	"errors"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNbtAccounterQuota(t *testing.T) {
	encoded := encodeTag(t, &Tag{Name: "", Value: map[string]any{}})
	accounter := NbtAccounterCreate(16)
	_, err := ReadUncompressedFileAccounter(bytes.NewReader(encoded), accounter)
	var accounterErr *NbtAccounterException
	if !errors.As(err, &accounterErr) {
		t.Fatalf("expected NbtAccounterException, got %v", err)
	}
}

func TestNbtAccounterDepth(t *testing.T) {
	encoded := encodeTag(t, &Tag{Name: "", Value: map[string]any{"nested": map[string]any{}}})
	accounter := NewNbtAccounter(NbtAccounterUNCOMPRESSED_NBT_QUOTA, 1)
	_, err := ReadUncompressedFileAccounter(bytes.NewReader(encoded), accounter)
	var accounterErr *NbtAccounterException
	if !errors.As(err, &accounterErr) {
		t.Fatalf("expected NbtAccounterException, got %v", err)
	}
}

func TestNegativeListLength(t *testing.T) {
	// List root, element type Byte, count -1.
	encoded := []byte{byte(List), byte(Byte), 0xFF, 0xFF, 0xFF, 0xFF}
	anon := &Anon{}
	err := anon.Decode(bytes.NewReader(encoded))
	var formatErr *NbtFormatException
	if !errors.As(err, &formatErr) {
		t.Fatalf("expected NbtFormatException, got %v", err)
	}
}

func TestMissingListType(t *testing.T) {
	// List root, element type End, count 1.
	encoded := []byte{byte(List), byte(End), 0x00, 0x00, 0x00, 0x01}
	anon := &Anon{}
	err := anon.Decode(bytes.NewReader(encoded))
	var formatErr *NbtFormatException
	if !errors.As(err, &formatErr) {
		t.Fatalf("expected NbtFormatException, got %v", err)
	}
}

func TestUncompressedRoundTrip(t *testing.T) {
	root := &Tag{Name: "servers", Value: map[string]any{
		"servers": []any{
			map[string]any{"name": "s", "ip": "1.2.3.4", "hidden": int8(0)},
		},
	}}
	buffer := &bytes.Buffer{}
	if err := WriteUncompressedFile(buffer, root); err != nil {
		t.Fatal(err)
	}
	got, err := ReadUncompressedFile(bytes.NewReader(buffer.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != root.Name {
		t.Errorf("name = %q", got.Name)
	}
	if !reflect.DeepEqual(got.Value, root.Value) {
		t.Errorf("value = %#v", got.Value)
	}
}

func TestUncompressedFileRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "servers.dat")
	root := &Tag{Name: "", Value: map[string]any{"servers": []any{}}}
	if err := WriteUncompressedFilePath(path, root); err != nil {
		t.Fatal(err)
	}
	got, err := ReadUncompressedFilePath(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Value, root.Value) {
		t.Errorf("value = %#v", got.Value)
	}
}

func encodeTag(t *testing.T, root *Tag) []byte {
	t.Helper()
	buffer := &bytes.Buffer{}
	if err := root.Encode(buffer); err != nil {
		t.Fatal(err)
	}
	return buffer.Bytes()
}
