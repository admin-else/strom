package resources

import "testing"

func TestIdentifierParse(t *testing.T) {
	id := IdentifierParse("minecraft:textures/gui/options_background.png")
	if got := id.Namespace(); got != "minecraft" {
		t.Errorf("namespace = %q, want minecraft", got)
	}
	if got := id.Path(); got != "textures/gui/options_background.png" {
		t.Errorf("path = %q", got)
	}
}

func TestIdentifierWithDefaultNamespace(t *testing.T) {
	id := IdentifierWithDefaultNamespace("core/gui")
	if got := id.String(); got != "minecraft:core/gui" {
		t.Errorf("string = %q", got)
	}
}

func TestIdentifierBySeparatorEmptyNamespace(t *testing.T) {
	id := IdentifierBySeparator(":path/to/thing", ':')
	if got := id.Namespace(); got != "minecraft" {
		t.Errorf("namespace = %q, want minecraft", got)
	}
	if got := id.Path(); got != "path/to/thing" {
		t.Errorf("path = %q", got)
	}
}

func TestIdentifierTryParseRejectsInvalidPath(t *testing.T) {
	if id := IdentifierTryParse("minecraft:Path"); id != nil {
		t.Errorf("tryParse returned %v, want nil", id)
	}
}

func TestIdentifierTryBuildRejectsDotDotNamespace(t *testing.T) {
	if id := IdentifierTryBuild("..", "path"); id != nil {
		t.Errorf("tryBuild returned %v, want nil", id)
	}
}

func TestIdentifierToShortString(t *testing.T) {
	if got := IdentifierParse("minecraft:foo").ToShortString(); got != "foo" {
		t.Errorf("short = %q, want foo", got)
	}
	if got := IdentifierParse("custom:foo").ToShortString(); got != "custom:foo" {
		t.Errorf("short = %q, want custom:foo", got)
	}
}

func TestIdentifierCompareTo(t *testing.T) {
	a := IdentifierParse("a:foo")
	b := IdentifierParse("b:foo")
	if a.CompareTo(b) >= 0 {
		t.Error("path equal, namespace a should sort before b")
	}
	if !(IdentifierParse("minecraft:a").CompareTo(IdentifierParse("minecraft:b")) < 0) {
		t.Error("path a should sort before b")
	}
}

func TestIdentifierResolveAgainstRejectsTraversal(t *testing.T) {
	root := "/tmp/res"
	if _, err := (Identifier{namespace: "minecraft", path: "../../escape"}).ResolveAgainst(root); err == nil {
		t.Error("expected traversal above root to be rejected")
	}
}

// tryResult adapts the *Identifier returned by IdentifierTryBySeparator for
// the error check in the test.
func (i *Identifier) tryResult() *Identifier { return i }
