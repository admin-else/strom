package main

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"

	"github.com/admin-else/strom/proto_generator/protodef"
	"github.com/admin-else/strom/util"
	"github.com/go-viper/mapstructure/v2"
)

func (g *Generator) ParseCompareTo(compareTo string) (e ast.Expr, cet CaseExprType, err error) {
	parts := strings.Split(compareTo, "/")
	downPrefixCount := 0
	for _, part := range parts {
		if part == ".." {
			downPrefixCount++
		}
	}
	parts = parts[downPrefixCount:]
	startingPoint := g.ContainerStack[len(g.ContainerStack)-downPrefixCount-1]
	e, cet, err = g.VisitCompareTo(parts, startingPoint.VarToSet, CombineNamAndData("container", startingPoint.Data))
	//	return ParseCompareToLegacy(compareTo, varToSet)
	return
}

func (g *Generator) VisitCompareTo(parts []string, inExpr ast.Expr, data any) (e ast.Expr, cet CaseExprType, err error) {
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

func ContainerCompareTo(g *Generator, parts []string, inExpr ast.Expr, dataRaw any) (e ast.Expr, cet CaseExprType, err error) {
	var data protodef.Container
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
			return g.VisitCompareTo(parts, SelectorExprAndStr(inExpr, util.CamelCase(field.Name)), field.Type)
		}
	}
	err = errors.New("field not found")
	return
}

func ReturnInputCompareTo(_ *Generator, _ []string, inExpr ast.Expr, _ any) (e ast.Expr, cet CaseExprType, err error) {
	return inExpr, CaseExprTypeUnset, nil
}

func BitfieldCompareTo(_ *Generator, parts []string, inExpr ast.Expr, dataRaw any) (e ast.Expr, cet CaseExprType, err error) {
	var data protodef.BitField
	err = mapstructure.Decode(dataRaw, &data)
	if err != nil {
		return
	}
	fieldName, parts, found := PopFront(parts)
	if !found {
		err = errors.New("expected name part")
		return
	}
	var field protodef.BitFieldEntry
	for _, f := range data {
		if f.Name == fieldName {
			field = f
			break
		}
	}
	if field.Size == 1 {
		cet = CaseExprTypeOnlyBool
	} else {
		cet = CaseExprTypeOnlyNumber
	}

	e = SelectorExprAndStr(inExpr, util.CamelCase(fieldName))
	return
}

func BitflagsCompareTo(_ *Generator, parts []string, inExpr ast.Expr, dataRaw any) (e ast.Expr, cet CaseExprType, err error) {
	var data protodef.BitFlags
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
		"bitfield":  BitfieldCompareTo,
		"u8":        ReturnInputCompareTo,
		"i8":        ReturnInputCompareTo,
		"bitflags":  BitflagsCompareTo,
		"switch":    ReturnInputCompareTo,
		"varlong":   ReturnInputCompareTo,
	}
}
