package util

import (
	"reflect"
	"testing"
)

func TestDecomposePath(t *testing.T) {
	cases := []struct {
		in      string
		want    []string
		wantErr bool
	}{
		{"a", []string{"a"}, false},
		{"a/b", []string{"a", "b"}, false},
		{"a/b/c", []string{"a", "b", "c"}, false},
		{"", nil, true},
		{".", nil, true},
		{"..", nil, true},
		{"a//b", nil, true},
		{"a/../b", nil, true},
		{"a/B", nil, true},
		{"a/b/", nil, true},
	}
	for _, c := range cases {
		got, err := DecomposePath(c.in)
		if c.wantErr {
			if err == nil {
				t.Errorf("DecomposePath(%q) expected error", c.in)
			}
			continue
		}
		if err != nil {
			t.Errorf("DecomposePath(%q) error: %v", c.in, err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("DecomposePath(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestResolvePath(t *testing.T) {
	if got := ResolvePath("/root", nil); got != "/root" {
		t.Errorf("empty = %q", got)
	}
	if got := ResolvePath("/root", []string{"a"}); got != "/root/a" {
		t.Errorf("one = %q", got)
	}
	if got := ResolvePath("/root", []string{"a", "b"}); got != "/root/a/b" {
		t.Errorf("two = %q", got)
	}
}

func TestValidatePath(t *testing.T) {
	if err := ValidatePath("pack.mcmeta"); err != nil {
		t.Fatalf("valid path rejected: %v", err)
	}
	if err := ValidatePath(); err == nil {
		t.Fatal("empty path list accepted")
	}
	if err := ValidatePath(".."); err == nil {
		t.Fatal("parent segment accepted")
	}
}
