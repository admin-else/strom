package util

import (
	"encoding/json"
	"testing"

	"github.com/admin-else/strom/mc/resources"
)

func TestParseIdentifierPattern(t *testing.T) {
	pattern, err := ParseIdentifierPattern(json.RawMessage(`{"namespace":"minecraft","path":"textures/.*"}`))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if !pattern.NamespacePredicate()("minecraft") || pattern.NamespacePredicate()("other") {
		t.Fatal("namespace predicate wrong")
	}
	if !pattern.PathPredicate()("textures/block/stone.png") || pattern.PathPredicate()("models/block/stone.json") {
		t.Fatal("path predicate wrong")
	}
	id := resources.NewIdentifier("minecraft", "textures/block/stone.png")
	if !pattern.LocationPredicate()(id) {
		t.Fatal("location predicate should match")
	}
}

func TestParseIdentifierPatternEmpty(t *testing.T) {
	pattern, err := ParseIdentifierPattern(json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if !pattern.NamespacePredicate()("anything") || !pattern.PathPredicate()("anything") {
		t.Fatal("omitted patterns must accept everything")
	}
}

func TestParseIdentifierPatternInvalidRegex(t *testing.T) {
	if _, err := ParseIdentifierPattern(json.RawMessage(`{"path":"["}`)); err == nil {
		t.Fatal("expected invalid regex error")
	}
}
