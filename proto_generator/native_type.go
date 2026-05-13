package main

import (
	"errors"
	"fmt"
	"go/ast"

	"github.com/admin-else/strom/proto_generator/protodef"
	"github.com/admin-else/strom/util"
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
		var t ast.Expr
		t, err = g.VisitType(field.Type)
		if err != nil {
			return nil, err
		}
		if field.Anon {
			field.Name = "Anon" // lol
		}
		AddFieldToStruct(s, util.CamelCase(field.Name), t)
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

func VisitArrayTypeVisitorType(g *Generator, dataRaw any) (e ast.Expr, err error) {
	var data struct {
		CountType any
		Type      any
	}
	err = mapstructure.Decode(dataRaw, &data)
	if err != nil {
		return
	}
	e, err = g.VisitType(data.Type)
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

func BitSizeToProtodef(totalSize int, singed bool) (e string, err error) {
	if !singed && totalSize == 1 {
		e = "bool"
		return
	}
	if singed {
		e = "i"
	} else {
		e = "u"
	}
	if totalSize <= 8 {
		e += "8"
	} else if totalSize <= 16 {
		e += "16"
	} else if totalSize <= 32 {
		e += "32"
	} else if totalSize <= 64 {
		e += "64"
	} else {
		err = errors.New("bitfield size too large")
	}
	return
}

func VisitBitFieldType(g *Generator, dataRaw any) (e ast.Expr, err error) {
	var data protodef.BitField
	err = mapstructure.Decode(dataRaw, &data)
	if err != nil {
		return
	}
	ret := NewStruct()
	for _, field := range data {
		var p string
		p, err = BitSizeToProtodef(field.Size, field.Signed)
		if err != nil {
			return
		}
		var fieldType ast.Expr
		fieldType, err = g.VisitType(p)
		if err != nil {
			return
		}
		AddFieldToStruct(ret, util.CamelCase(field.Name), fieldType)
	}
	e = ret
	return
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

func VoidType(_ *Generator, _ any) (e ast.Expr, err error) {
	e = NewStruct()
	return
}

func EntityMetadataLoopType(g *Generator, v any) (e ast.Expr, err error) {
	var data protodef.EntityMetadataLoop
	err = mapstructure.Decode(v, &data)
	if err != nil {
		return
	}
	innerType, err := g.VisitType(data.Type)
	if err != nil {
		return
	}
	e = MapType(Ident("byte"), innerType)
	return
}

func (g *Generator) RegisterNatives() {
	g.Natives = map[string]ExprGeneratorFunc{
		"container":                VisitContainerType,
		"buffer":                   VisitBufferType,
		"varint":                   MakeIdentVisitor("int32"),
		"varlong":                  MakeIdentVisitor("int64"),
		"array":                    VisitArrayTypeVisitorType,
		"mapper":                   MakeIdentVisitor("string"),
		"native":                   VisitDontGenerateType,
		"anonymousNbt":             MakeSelectorVisitor("nbt", "Anon"),
		"anonOptionalNbt":          MakeSelectorVisitor("nbt", "Anon"), // I have no idea what the difference is between these two
		"void":                     VoidType,
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
		"entityMetadataLoop":       EntityMetadataLoopType,
		"option":                   VisitOptionType,
		"switch":                   MakeIdentVisitor("any"),
		"UUID":                     MakeSelectorVisitor("uuid", "UUID"),
		"bitfield":                 VisitBitFieldType,
		"pstring":                  VisitDontGenerateType,
		"string":                   MakeIdentVisitor("string"),
		"restBuffer":               MakeSelectorVisitor("proto_base", "RestBuffer"),
		"bitflags":                 VisitBitFlagsType,
		"topBitSetTerminatedArray": VisitToDoType,
		"todo":                     VisitToDoType,
		"lpVec3":                   VisitToDoType,
	}
}
