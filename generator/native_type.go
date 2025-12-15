package main

import (
	"errors"
	"fmt"
	"go/ast"

	"github.com/go-viper/mapstructure/v2"
)

func VisitContainerType(g *Generator, dataRaw any) (ast.Expr, error) {
	var data []struct {
		Name string
		Type any
		Anon bool
	}
	err := mapstructure.Decode(dataRaw, &data)
	if err != nil {
		return nil, err
	}

	g.ContainerStack = append(g.ContainerStack, ContainerStackEntry{
		Data:     data,
		VarToSet: nil,
	})
	defer func() { g.ContainerStack = g.ContainerStack[:len(g.ContainerStack)-1] }()

	s := NewStruct()
	for _, field := range data {
		t, err := g.VisitType(field.Type)
		if err != nil {
			return nil, err
		}
		if field.Anon {
			field.Name = "Anon" // lol
		}
		AddFieldToStruct(s, CamelCase(field.Name), t)
	}
	return s, nil
}

func VisitBufferType(_ *Generator, dataRaw any) (e ast.Expr, err error) {
	var data struct {
		CountType any
		Count     int
	}
	err = mapstructure.Decode(dataRaw, &data)
	if err != nil {
		return
	}
	if data.CountType != nil {
		if data.CountType == "varint" {
			e = Slice(Ident("byte"))
			return
		}
		return nil, fmt.Errorf("unsupported count type: %v", data.CountType) // TODO: old versions use uint16 as count type
	}
	e = Array(Ident("byte"), data.Count)
	return
}

func VisitArrayTypeVisitorType(g *Generator, dataRaw any) (ast.Expr, error) {
	var data struct {
		CountType any
		Type      any
	}
	must(mapstructure.Decode(dataRaw, &data))
	e, err := g.VisitType(data.Type)
	if err != nil {
		return nil, err
	}
	return Slice(e), nil
}

func VisitDontGenerateType(_ *Generator, _ any) (ast.Expr, error) {
	return nil, nil
}

func VisitToDoType(_ *Generator, _ any) (e ast.Expr, err error) {
	err = ToDoError
	return
}

func VisitOptionType(g *Generator, dataRaw any) (e ast.Expr, err error) {
	e, err = g.VisitType(dataRaw)
	if err != nil {
		return
	}
	e = Pointer(e)
	return
}

func TotalBitfieldSizeToProtodefName(totalSize int) (e string, err error) {
	if totalSize <= 8 {
		e = "u8"
	} else if totalSize <= 16 {
		e = "u16"
	} else if totalSize <= 32 {
		e = "u32"
	} else if totalSize <= 64 {
		e = "u64"
	} else {
		err = errors.New("bitfield size too large")
	}
	return
}

func VisitBitFieldType(g *Generator, dataRaw any) (e ast.Expr, err error) {
	var data []struct {
		Name   string
		Singed bool
		Size   int
	}
	err = mapstructure.Decode(dataRaw, &data)
	if err != nil {
		return
	}
	totalSize := 0
	for _, field := range data {
		totalSize += field.Size
	}
	// TODO: we will want to generate Methods and other stuff
	p, err := TotalBitfieldSizeToProtodefName(totalSize)
	if err != nil {
		return
	}
	return g.VisitType(p)
}

func VisitBitFlagsType(g *Generator, dataRaw any) (e ast.Expr, err error) {
	var data struct {
		Flags []string
		Type  any
	}
	err = mapstructure.Decode(dataRaw, &data)
	if err != nil {
		return
	}
	return g.VisitType(data.Type)
}

func MakeIdentVisitor(name string) func(g *Generator, dataRaw any) (e ast.Expr, err error) {
	return func(g *Generator, dataRaw any) (e ast.Expr, err error) {
		e = Ident(name)
		return
	}
}

func MakeSelectorVisitor(x, sel string) func(g *Generator, dataRaw any) (e ast.Expr, err error) {
	return func(g *Generator, dataRaw any) (e ast.Expr, err error) {
		e = Selector(x, sel)
		return
	}
}

func (g *Generator) RegisterNatives() {
	g.Natives = map[string]ExprGeneratorFunc{
		"container":                VisitContainerType,
		"buffer":                   VisitBufferType,
		"varint":                   MakeSelectorVisitor("queser", "VarInt"),
		"array":                    VisitArrayTypeVisitorType,
		"mapper":                   MakeIdentVisitor("string"),
		"native":                   VisitDontGenerateType,
		"anonymousNbt":             MakeSelectorVisitor("nbt", "Anon"),
		"anonOptionalNbt":          MakeSelectorVisitor("nbt", "Anon"), // I have no idea what the difference is between these two
		"void":                     MakeSelectorVisitor("queser", "Void"),
		"bool":                     MakeIdentVisitor("bool"),
		"u8":                       MakeIdentVisitor("uint8"),
		"u16":                      MakeIdentVisitor("uint16"),
		"u32":                      MakeIdentVisitor("uint32"),
		"u64":                      MakeIdentVisitor("uint64"),
		"i8":                       MakeIdentVisitor("int8"),
		"i16":                      MakeIdentVisitor("int16"),
		"i32":                      MakeIdentVisitor("int32"),
		"i64":                      MakeIdentVisitor("int64"),
		"f32":                      MakeIdentVisitor("float32"),
		"f64":                      MakeIdentVisitor("float64"),
		"registryEntryHolder":      MakeIdentVisitor("any"), // Go really fucking needs good tagged unions
		"registryEntryHolderSet":   MakeIdentVisitor("any"),
		"entityMetadataLoop":       VisitToDoType,
		"option":                   VisitOptionType,
		"switch":                   MakeIdentVisitor("any"),
		"UUID":                     MakeSelectorVisitor("uuid", "UUID"),
		"bitfield":                 VisitBitFieldType,
		"pstring":                  VisitDontGenerateType,
		"string":                   MakeIdentVisitor("string"),
		"restBuffer":               MakeSelectorVisitor("queser", "RestBuffer"),
		"bitflags":                 VisitBitFlagsType,
		"topBitSetTerminatedArray": VisitToDoType,
		"todo":                     VisitToDoType,
	}
}
