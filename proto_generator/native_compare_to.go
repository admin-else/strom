package main

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"

	"github.com/go-viper/mapstructure/v2"
)

func (g *Generator) ParseCompareTo(compareTo string) (e ast.Expr, err error) {
	parts := strings.Split(compareTo, "/")
	downPrefixCount := 0
	for _, part := range parts {
		if part == ".." {
			downPrefixCount++
		}
	}
	parts = parts[downPrefixCount:]
	startingPoint := g.ContainerStack[len(g.ContainerStack)-downPrefixCount-1]
	e, err = g.VisitCompareTo(parts, startingPoint.VarToSet, CombineNamAndData("container", startingPoint.Data))
	//	return ParseCompareToLegacy(compareTo, varToSet)
	return
}

func (g *Generator) VisitCompareTo(parts []string, inExpr ast.Expr, data any) (e ast.Expr, err error) {
	tName, tData, err := ParseType(data)
	if err != nil {
		return
	}
	d, found := g.CompareToNatives[tName]
	if found {
		return d(g, parts, inExpr, tData)
	}
	err = fmt.Errorf("native compare to not implemented for %v", tName)
	return
}

func ContainerCompareTo(g *Generator, parts []string, inExpr ast.Expr, dataRaw any) (e ast.Expr, err error) {
	var data []struct {
		Name string
		Type any
		Anon bool
	}
	err = mapstructure.Decode(dataRaw, &data)
	if err != nil {
		return
	}

	name, parts, ok := PopFront(parts)
	if !ok {
		err = errors.New("expected name part")
		return
	}
	for _, field := range data {
		if field.Name == name {
			return g.VisitCompareTo(parts, SelectorExprAndStr(inExpr, CamelCase(field.Name)), field.Type)
		}
	}
	err = errors.New("field not found")
	return
}

func ReturnInputCompareTo(_ *Generator, _ []string, inExpr ast.Expr, _ any) (e ast.Expr, err error) {
	return inExpr, nil
}

func BitflagsCompareTo(_ *Generator, parts []string, inExpr ast.Expr, dataRaw any) (e ast.Expr, err error) {
	var data struct {
		Flags []string
		Type  any
	}
	err = mapstructure.Decode(dataRaw, &data)
	if err != nil {
		return
	}
	name, parts, ok := PopFront(parts)
	if !ok {
		err = errors.New("expected name part")
		return
	}

	for i, flag := range data.Flags {
		if flag != name {
			continue
		}
		e = NotEquals(BinAnd(inExpr, NumLitHex(1<<i)), NumLit(0))
		return
	}
	err = errors.New("flag not found")
	return
}

func (g *Generator) RegisterCompareToNatives() {
	g.CompareToNatives = map[string]CompareToGeneratorFunc{
		"container": ContainerCompareTo,
		"bool":      ReturnInputCompareTo,
		"mapper":    ReturnInputCompareTo,
		"varint":    ReturnInputCompareTo,
		"bitfield":  ReturnInputCompareTo,
		"u8":        ReturnInputCompareTo,
		"i8":        ReturnInputCompareTo,
		"bitflags":  BitflagsCompareTo,
		"switch":    ReturnInputCompareTo,
		"varlong":   ReturnInputCompareTo,
	}
}
