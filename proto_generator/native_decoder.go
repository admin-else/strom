package main

import (
	"go/ast"
	"go/token"
	"slices"
	"strconv"

	"github.com/go-viper/mapstructure/v2"
)

func DefaultDecoder(_ *Generator, varToSet ast.Expr, _ any, _ string) (s []ast.Stmt, err error) {
	s = Stmts(
		Assign(
			Exprs(Ident("err")),
			Exprs(Call(SelectorExprAndStr(varToSet, "Decode"), Ident("r"))),
		),
		IfErrNil(),
	)

	return
}

func UUIDDecoder(_ *Generator, varToSet ast.Expr, _ any, _ string) (s []ast.Stmt, err error) {
	//_, err = io.ReadFull(r, U[:])
	s = Stmts(
		Assign(Exprs(Ident("_"), Ident("err")), Exprs(Call(Selector("io", "ReadFull"), Ident("r"), MultiIndex(varToSet, nil, nil, nil)))),
		IfErrNil(),
	)
	return
}

func ContainerDecoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
	var data []struct {
		Name string
		Type any
		Anon bool
	}
	err = mapstructure.Decode(dataRaw, &data)
	if err != nil {
		return
	}

	g.ContainerStack = append(g.ContainerStack, ContainerStackEntry{
		Data:     data,
		VarToSet: varToSet,
	})
	defer func() { g.ContainerStack = g.ContainerStack[:len(g.ContainerStack)-1] }()

	for _, field := range data {
		var statements []ast.Stmt
		cName := CamelCase(field.Name)
		if field.Anon {
			cName = "Anon"
		}
		statements, err = g.VisitDecoder(SelectorExprAndStr(varToSet, cName), field.Type, name+cName)
		if err != nil {
			return
		}
		s = append(s, statements...)
	}
	return
}

func SimpleTypeDecoder(_ *Generator, varToSet ast.Expr, _ any, _ string) (s []ast.Stmt, err error) {
	s = Stmts(
		Assign(
			Exprs(Ident("err")),
			Exprs(Call(Selector("binary", "Read"), Ident("r"), Selector("binary", "BigEndian"), AddrOf(varToSet))),
		),
		IfErrNil(),
	)
	return
}

func StringDecoder(_ *Generator, varToSet ast.Expr, _ any, _ string) (s []ast.Stmt, err error) {
	s = Stmts(
		Assign(Exprs(varToSet, Ident("err")), Exprs(Call(Selector("proto_base", "DecodeString"), Ident("r")))),
		IfErrNil(),
	)
	return
}

func OptionDecoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
	pVarName := name + "Present"
	tVarName := name + "PresentValue"
	boolStmts, err := g.VisitDecoder(Ident(pVarName), "bool", name)
	if err != nil {
		return
	}
	t, err := g.VisitType(dataRaw)
	if err != nil {
		return
	}
	tStmts1 := Stmts(VarStmt(tVarName, t))
	tStmts2, err := g.VisitDecoder(Ident(tVarName), dataRaw, name)
	if err != nil {
		return
	}

	// This looks so shit
	tStmts := append(tStmts1, tStmts2...)
	tStmts = append(tStmts, Assign121(varToSet, AddrOf(Ident(tVarName))))
	s = append(s, VarStmt(pVarName, Ident("bool")))
	s = append(s, boolStmts...)
	s = append(s, If(Ident(pVarName), NewBlock(tStmts)))
	return
}

func ArrayDecoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
	var data struct {
		CountType any
		Count     string
		Type      any
	}
	lName := "l" + CamelCase(name)
	err = mapstructure.Decode(dataRaw, &data)
	if err != nil {
		return
	}
	var s2 []ast.Stmt
	countType, err := g.VisitNameAndData("varint", nil)
	if data.CountType != nil {
		countType, err = g.VisitType(data.CountType)
		if err != nil {
			return
		}
		s2, err = g.VisitDecoder(Ident(lName), data.CountType, name)
		if err != nil {
			return
		}
	} else {
		var n int
		n, err = strconv.Atoi(data.Count)
		var e ast.Expr
		if err == nil {
			e = NumLit(n)
		} else {
			e, err = g.ParseCompareTo(data.Count)
		}
		s2 = append(s2, Assign121(Ident(lName), e))
	}

	normalType, err := g.VisitType(data.Type)
	if err != nil {
		return
	}

	arrayType, err := g.VisitNameAndData("array", dataRaw)
	if err != nil {
		return
	}

	tName := name + "Element"
	rangeStatements := Stmts(VarStmt(tName, normalType))
	s5, err := g.VisitDecoder(Ident(tName), data.Type, tName)
	rangeStatements = append(rangeStatements, s5...)
	rangeStatements = append(rangeStatements, Assign121(varToSet, Call(Ident("append"), varToSet, Ident(tName))))

	if err != nil {
		return
	}
	s1 := VarStmt(lName, countType)
	s3 := Assign121(varToSet, CompLit(arrayType, Exprs()))
	s4 := ForRange(nil, Ident(lName), NewBlock(rangeStatements))
	s = append(s, s1)
	s = append(s, s2...)
	s = append(s, s3)
	s = append(s, s4)
	return
}

func TypeForwardDecoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
	var data struct {
		Type string
	}
	if err = mapstructure.Decode(dataRaw, &data); err != nil {
		return
	}
	return g.VisitDecoder(varToSet, data.Type, name)
}

var DontStringify = []string{"true", "false"}

// CaseExprType This exists to fix this
// "switch",
//
//	{
//	 "compareTo": "itemCount",
//	 "fields": {
//	   "0": "void",
//	   "false": "void"
//	 },
//	 "default": [
//
// Protodef is best!!!
type CaseExprType int

const (
	CaseExprTypeUnset CaseExprType = iota
	CaseExprTypeNumber
	CaseExprTypeSpecial
	CaseExprTypeString
)

func MultiTypeFix(s string, t *CaseExprType) (e ast.Expr) {
	n, err := strconv.Atoi(s)
	ct := CaseExprTypeUnset
	if err == nil {
		err = nil
		e = NumLit(n)
		ct = CaseExprTypeNumber
	} else if slices.Contains(DontStringify, s) {
		e = Ident(s)
		ct = CaseExprTypeSpecial
	} else {
		e = StrLit(s)
		ct = CaseExprTypeString
	}

	if *t == CaseExprTypeUnset {
		*t = ct
		return
	}
	if *t != ct {
		e = nil
	}
	return
}

func SwitchDecoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
	var data struct {
		CompareTo string
		Fields    map[string]any
		Default   any
	}
	err = mapstructure.Decode(dataRaw, &data)
	if err != nil {
		return
	}
	compareToExpr, err := g.ParseCompareTo(data.CompareTo)
	if err != nil {
		return
	}
	cet := CaseExprTypeUnset
	for _, fName := range OrderedKeys(data.Fields) {
		fType := data.Fields[fName]
		e := MultiTypeFix(fName, &cet)
		if e == nil {
			continue
		}
		var tType ast.Expr
		tType, err = g.VisitType(fType)
		if err != nil {
			return
		}
		tName := name + CamelCase(fName) + "Tmp"
		s1 := VarStmt(tName, tType)
		s3 := Assign121(varToSet, Ident(tName))
		var caseDecodeValueStmts []ast.Stmt
		caseDecodeValueStmts, err = g.VisitDecoder(Ident(tName), fType, name)
		if err != nil {
			return
		}
		var caseStmts []ast.Stmt
		caseStmts = append(caseStmts, s1)
		caseStmts = append(caseStmts, caseDecodeValueStmts...)
		caseStmts = append(caseStmts, s3)
		s = append(s, Case(Exprs(e), caseStmts))
	}
	if data.Default != nil {
		var tType ast.Expr
		tType, err = g.VisitType(data.Default)
		tName := name + "Tmp"
		s1 := VarStmt(tName, tType)
		s3 := Assign121(varToSet, Ident(tName))
		var caseDecodeValueStmts []ast.Stmt
		caseDecodeValueStmts, err = g.VisitDecoder(Ident(tName), data.Default, name)
		if err != nil {
			return
		}
		var caseStmts []ast.Stmt
		caseStmts = append(caseStmts, s1)
		caseStmts = append(caseStmts, caseDecodeValueStmts...)
		caseStmts = append(caseStmts, s3)
		s = append(s, Case(nil, caseStmts))
	}
	if len(s) != 0 {
		s = Stmts(SwitchStmt(compareToExpr, NewBlock(s)))
	}
	return
}

var ToDoStmts = Stmts(Assign121(Ident("err"), Selector("proto_base", "ToDoError")))

func ToDoDecoder(_ *Generator, _ ast.Expr, _ any, _ string) (s []ast.Stmt, err error) {
	err = ToDoError
	return
}

func MapperDecoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
	var data struct {
		Mappings map[string]string
		Type     any
	}
	err = mapstructure.Decode(dataRaw, &data)
	if err != nil {
		return
	}
	kName := CamelCase(name) + "Key"
	kType, err := g.VisitType(data.Type)
	if err != nil {
		return
	}
	s1 := VarStmt(kName, kType)
	s2, err := g.VisitDecoder(Ident(kName), data.Type, name)
	if err != nil {
		return
	}

	mName := name + "Map"

	expressions := Exprs()
	for _, k := range OrderedKeys(data.Mappings) {
		v := data.Mappings[k]
		expressions = append(expressions, KeyValueExpr(NumLitStr(k), StrLit(v)))
	}
	g.Decl(mName, token.VAR, CompLit(MapType(kType, Ident("string")), expressions))
	s3 := Assign(Exprs(varToSet, Ident("err")), Exprs(Call(Selector("proto_base", "ErroringIndex"), Ident(mName), Ident(kName))))
	s = append(s, s1)
	s = append(s, s2...)
	s = append(s, s3)
	s = append(s, IfErrNil())
	return
}

func BufferDecoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
	var data struct {
		Count     int
		CountType any
	}
	err = mapstructure.Decode(dataRaw, &data)
	if err != nil {
		return
	}

	// 	rawString, err := io.ReadAll(io.LimitReader(r, int64(l)))

	if data.CountType != nil {
		lName := "l" + CamelCase(name)
		var countType ast.Expr
		countType, err = g.VisitType(data.CountType)
		if err != nil {
			return
		}
		var s2 []ast.Stmt
		s2, err = g.VisitDecoder(Ident(lName), data.CountType, name)
		if err != nil {
			return
		}
		s1 := VarStmt(lName, countType)
		s3 := Assign(Exprs(varToSet, Ident("err")), Exprs(Call(Selector("io", "ReadAll"), Call(Selector("io", "LimitReader"), Ident("r"), Call(Ident("int64"), Ident(lName))))))

		s = append(s, s1)
		s = append(s, s2...)
		s = append(s, s3)
	} else {
		s = Stmts(Assign(Exprs(Ident("_"), Ident("err")), Exprs(Call(Selector("r", "Read"), MultiIndex(varToSet, nil, nil, nil)))))
	}

	s = append(s, IfErrNil())
	return
}

func BitFieldDecoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
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
	p, err := TotalBitfieldSizeToProtodefName(totalSize)
	if err != nil {
		return
	}
	return g.VisitDecoder(varToSet, p, name)
}

func RegistryEntryHolderSetDecoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
	var data struct {
		Base struct {
			Name string
			Type any
		}
		Otherwise struct {
			Name string
			Type any
		}
	}
	err = mapstructure.Decode(dataRaw, &data)
	if err != nil {
		return
	}
	var countType ast.Expr
	countType, err = g.VisitType("varint")
	if err != nil {
		return
	}
	baseType, err := g.VisitType(data.Base.Type)
	if err != nil {
		return
	}
	otherwiseArrayType, err := g.VisitNameAndData("array", map[string]any{"type": data.Otherwise.Type})
	if err != nil {
		return
	}
	otherwiseType, err := g.VisitType(data.Otherwise.Type)
	if err != nil {
		return
	}
	lName := "l" + CamelCase(name)
	s1 := VarStmt(lName, countType)
	s2, err := g.VisitDecoder(Ident(lName), "varint", name)
	if err != nil {
		return
	}
	s3 := Assign121(Ident(lName), Sub(Ident(lName), NumLit(1)))
	resultName := name + "Result"
	b1s1 := VarStmt(resultName, baseType)
	b1s2, err := g.VisitDecoder(Ident(resultName), data.Base.Type, resultName)
	b1s3 := Assign121(varToSet, Ident(resultName))
	b1s4 := Return()
	var b1 []ast.Stmt
	b1 = append(b1, b1s1)
	b1 = append(b1, b1s2...)
	b1 = append(b1, b1s3)
	b1 = append(b1, b1s4)

	tName := name + "Element"
	rangeStatements := Stmts(VarStmt(tName, otherwiseType))
	rangeStatementsDecode, err := g.VisitDecoder(Ident(tName), data.Otherwise.Type, tName)
	rangeStatements = append(rangeStatements, rangeStatementsDecode...)
	rangeStatements = append(rangeStatements, Assign121(Ident(resultName), Call(Ident("append"), Ident(resultName), Ident(tName))))

	s5 := If(Equals(Ident(lName), NumLit(-1)), NewBlock(b1))
	s6 := VarStmt(resultName, otherwiseArrayType)
	s7 := ForRange(nil, Ident(lName), NewBlock(rangeStatements))
	s = append(s, s1)
	s = append(s, s2...)
	s = append(s, s3)
	s = append(s, s5)
	s = append(s, s6)
	s = append(s, s7)
	return
}

func RegistryEntryHolderDecoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
	var data EntryHolderSet
	err = mapstructure.Decode(dataRaw, &data)
	if err != nil {
		return
	}
	idName := name + "Id"
	varIntType, err := g.VisitNameAndData("varint", nil)
	if err != nil {
		return
	}
	IdVar := VarStmt(idName, varIntType)
	idDecodeStatements, err := g.VisitDecoder(Ident(idName), "varint", name)
	if err != nil {
		return
	}

	otherwiseType, err := g.VisitType(data.Otherwise.Type)
	if err != nil {
		return
	}
	resName := name + "Result"
	resNameIdent := Ident(resName)
	var otherwisePath []ast.Stmt
	otherwisePath = append(otherwisePath, VarStmt(resName, otherwiseType))

	otherwiseDecodeStatements, err := g.VisitDecoder(resNameIdent, data.Otherwise.Type, name+"Otherwise")
	if err != nil {
		return
	}
	otherwisePath = append(otherwisePath, otherwiseDecodeStatements...)
	otherwisePath = append(otherwisePath, Assign121(varToSet, resNameIdent))
	idPath := IfElse(NotEquals(Ident(idName), NumLit(0)), NewBlockEllipsis(Assign121(varToSet, Ident(idName))), NewBlock(otherwisePath))

	s = append(s, IdVar)
	s = append(s, idDecodeStatements...)
	s = append(s, idPath)
	return
}

func VarIntDecoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
	s = Stmts(Assign(Exprs(varToSet, Ident("err")), Exprs(Call(Selector("proto_base", "DecodeVarInt"), Ident("r")))), IfErrNil())
	return
}

func VarLongDecoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
	s = Stmts(Assign(Exprs(varToSet, Ident("err")), Exprs(Call(Selector("proto_base", "DecodeVarLong"), Ident("r")))), IfErrNil())
	return
}

func VoidDecoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
	return
}

func (g *Generator) RegisterDecoderNatives() {
	g.DecoderNatives = map[string]FunctionGeneratorFunc{
		"container": ContainerDecoder,
		"array":     ArrayDecoder,
		"option":    OptionDecoder,
		"bitflags":  TypeForwardDecoder,
		"switch":    SwitchDecoder,
		"mapper":    MapperDecoder,
		"buffer":    BufferDecoder,
		"bitfield":  BitFieldDecoder,
		"string":    StringDecoder,

		"void":            VoidDecoder,
		"varint":          VarIntDecoder,
		"varlong":         VarLongDecoder,
		"anonymousNbt":    DefaultDecoder,
		"anonOptionalNbt": DefaultDecoder, // I have no idea what the difference is between these two
		"UUID":            UUIDDecoder,
		"restBuffer":      DefaultDecoder,

		"f32":  SimpleTypeDecoder,
		"f64":  SimpleTypeDecoder,
		"i8":   SimpleTypeDecoder,
		"i16":  SimpleTypeDecoder,
		"i32":  SimpleTypeDecoder,
		"i64":  SimpleTypeDecoder,
		"u8":   SimpleTypeDecoder,
		"u16":  SimpleTypeDecoder,
		"u32":  SimpleTypeDecoder,
		"u64":  SimpleTypeDecoder,
		"bool": SimpleTypeDecoder,

		"registryEntryHolder":      RegistryEntryHolderDecoder,
		"registryEntryHolderSet":   ToDoDecoder, // RegistryEntryHolderSetDecoder
		"entityMetadataLoop":       ToDoDecoder,
		"topBitSetTerminatedArray": ToDoDecoder,
		"todo":                     ToDoDecoder,
	}
}
