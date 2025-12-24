package main

import (
	"go/ast"
	"go/token"

	"github.com/go-viper/mapstructure/v2"
)

func DefaultEncoder(_ *Generator, varToSet ast.Expr, _ any, _ string) (s []ast.Stmt, err error) {
	s = Stmts(
		Assign121(
			Ident("err"),
			Call(SelectorExprAndStr(varToSet, "Encode"), Ident("w")),
		),
		IfErrNil(),
	)

	return
}

func ContainerEncoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
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
		statements, err = g.VisitEncoder(SelectorExprAndStr(varToSet, cName), field.Type, name+CamelCase(cName))
		if err != nil {
			return
		}
		s = append(s, statements...)
	}
	return
}

func SimpleTypeEncoder(_ *Generator, data ast.Expr, _ any, _ string) (s []ast.Stmt, err error) {
	s = Stmts(
		Assign121(
			Ident("err"),
			Call(Selector("binary", "Write"), Ident("w"), Selector("binary", "BigEndian"), data),
		),
		IfErrNil(),
	)
	return
}

func ArrayEncoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
	var data struct {
		CountType any
		Count     string
		Type      any
	}
	err = mapstructure.Decode(dataRaw, &data)
	if err != nil {
		return
	}
	if data.CountType == nil {
		data.CountType = "varint"
	}
	ct, err := g.VisitType(data.CountType)
	if err != nil {
		return
	}
	iName := "i" + name
	//s, err = g.VisitEncoder(varToSet, data.Type, name)
	var s1 []ast.Stmt
	if data.Count == "" {
		s1, err = g.VisitEncoder(Call(ct, Call(Ident("len"), varToSet)), data.CountType, name)
		if err != nil {
			return
		}
	}
	s2, err := g.VisitEncoder(Index(varToSet, Ident(iName)), data.Type, name+"Inner")
	if err != nil {
	}
	s2 = Stmts(ForRange(Ident(iName), Call(Ident("len"), varToSet), NewBlock(s2)))
	s = append(s1, s2...)
	return
}

func ToDoEncoder(_ *Generator, _ ast.Expr, _ any, _ string) (s []ast.Stmt, err error) {
	err = ToDoError
	return
}

func OptionEncoder(g *Generator, data ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
	cond := NotEquals(data, Nil())
	s1, err := g.VisitEncoder(cond, "bool", name)
	if err != nil {
		return
	}
	s2, err := g.VisitEncoder(DeRef(data), dataRaw, name)
	if err != nil {
		return
	}
	s = append(s1, If(cond, NewBlock(s2)))
	return
}

func BufferEncoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
	var data struct {
		Count     int
		CountType any
	}
	err = mapstructure.Decode(dataRaw, &data)
	if err != nil {
		return
	}
	if data.CountType != nil {
		var countType ast.Expr
		countType, err = g.VisitType(data.CountType)
		if err != nil {
			return
		}
		var s2 []ast.Stmt
		s2, err = g.VisitEncoder(Call(countType, Call(Ident("len"), varToSet)), data.CountType, name)
		if err != nil {
			return
		}
		s = append(s, s2...)
	} else {
		s = append(s, Define121(Ident("arr"), varToSet))
		varToSet = MultiIndex(Ident("arr"), nil, nil, nil)
	}
	s = append(s, Assign(Exprs(Ident("_"), Ident("err")), Exprs(Call(Selector("w", "Write"), varToSet))))
	s = append(s, IfErrNil())
	return
}

func BitFieldEncoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
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
	return g.VisitEncoder(varToSet, p, name)
}

func SwitchEncoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
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
		s1 := Define(Exprs(Ident(name), Ident("ok")), Exprs(TypeAssert(varToSet, tType)))
		s2 := If(Not(Ident("ok")), NewBlockEllipsis(Assign121(Ident("err"), Selector("proto_base", "BadTypeError")), Return()))
		var caseDecodeValueStmts []ast.Stmt
		caseDecodeValueStmts, err = g.VisitEncoder(Ident(name), fType, name)
		if err != nil {
			return
		}
		var caseStmts []ast.Stmt
		caseStmts = append(caseStmts, s1)
		caseStmts = append(caseStmts, s2)
		caseStmts = append(caseStmts, caseDecodeValueStmts...)
		s = append(s, Case(Exprs(e), caseStmts))
	}
	if data.Default != nil {
		var tType ast.Expr
		tType, err = g.VisitType(data.Default)
		if err != nil {
			return
		}
		s1 := Define(Exprs(Ident("_"), Ident("ok")), Exprs(TypeAssert(varToSet, tType)))
		s2 := If(Not(Ident("ok")), NewBlockEllipsis(Assign121(Ident("err"), Selector("proto_base", "BadTypeError")), Return()))
		defaultStmts := []ast.Stmt{s1, s2}
		var defaultEncodingStmts []ast.Stmt
		defaultEncodingStmts, err = g.VisitEncoder(TypeAssert(varToSet, tType), data.Default, name)
		if err != nil {
			return
		}
		defaultStmts = append(defaultStmts, defaultEncodingStmts...)
		s = append(s, Case(nil, defaultStmts))
	}
	if len(s) != 0 {
		s = Stmts(SwitchStmt(compareToExpr, NewBlock(s)))
	}
	return

}

func TypeForwardEncoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
	var data struct {
		Type string
	}
	if err = mapstructure.Decode(dataRaw, &data); err != nil {
		return
	}
	return g.VisitEncoder(varToSet, data.Type, name)
}

func MapperEncoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
	var data struct {
		Mappings map[string]string
		Type     any
	}
	err = mapstructure.Decode(dataRaw, &data)
	if err != nil {
		return
	}
	kType, err := g.VisitType(data.Type)
	if err != nil {
		return
	}

	mName := name + "ReverseMap"

	exprs := Exprs()
	for _, k := range OrderedKeys(data.Mappings) {
		v := data.Mappings[k]
		exprs = append(exprs, KeyValueExpr(StrLit(v), NumLitStr(k)))
	}
	g.Decl(mName, token.VAR, CompLit(MapType(Ident("string"), kType), exprs))

	vName := "v" + name
	s1 := VarStmt(vName, kType)
	s2 := Assign(Exprs(Ident(vName), Ident("err")), Exprs(Call(Selector("proto_base", "ErroringIndex"), Ident(mName), varToSet)))
	s3 := IfErrNil()
	s4, err := g.VisitEncoder(Ident(vName), data.Type, name)
	if err != nil {
		return
	}

	s = append(s, s1)
	s = append(s, s2)
	s = append(s, s3)
	s = append(s, s4...)
	return
}

func StringEncoder(_ *Generator, varToSet ast.Expr, _ any, _ string) (s []ast.Stmt, err error) {
	s = Stmts(
		Assign121(Ident("err"), Call(Selector("proto_base", "EncodeString"), Ident("w"), varToSet)), IfErrNil())
	return
}

func UUIDEncoder(_ *Generator, varToSet ast.Expr, _ any, _ string) (s []ast.Stmt, err error) {
	s = Stmts(Assign(Exprs(Ident("_"), Ident("err")), Exprs(Call(Selector("w", "Write"), MultiIndex(varToSet, nil, nil, nil)))), IfErrNil())
	return
}

func RegistryEntryHolderEncoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
	var data EntryHolderSet
	err = mapstructure.Decode(dataRaw, &data)
	otherwiseType, err := g.VisitType(data.Otherwise.Type)
	if err != nil {
		return
	}
	knownTypeName := name + "KnownType"
	otherwiseEncodeStatements, err := g.VisitEncoder(Ident(knownTypeName), data.Otherwise.Type, name+"Otherwise")
	if err != nil {
		return
	}
	baseType, err := g.VisitType("varint")
	if err != nil {
		return
	}
	baseEncodeStatements, err := g.VisitEncoder(Ident(knownTypeName), "varint", name)
	if err != nil {
		return
	}
	statements := Stmts()
	statements = append(statements, Case(Exprs(baseType), baseEncodeStatements))
	statements = append(statements, Case(Exprs(otherwiseType), otherwiseEncodeStatements))
	statements = append(statements, Case(nil, Stmts(Assign121(Ident("err"), Selector("proto_base", "BadTypeError")))))

	s = Stmts(TypeSwitch(Define121(Ident(knownTypeName), TypeAssert(varToSet, Ident("type"))), NewBlock(statements)))
	return
}

func VarIntEncoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
	s = Stmts(Assign121(Ident("err"), Call(Selector("proto_base", "EncodeVarInt"), Ident("w"), varToSet)), IfErrNil())
	return
}

func VarLongEncoder(g *Generator, varToSet ast.Expr, dataRaw any, name string) (s []ast.Stmt, err error) {
	s = Stmts(Assign121(Ident("err"), Call(Selector("proto_base", "EncodeVarLong"), Ident("w"), varToSet)), IfErrNil())
	return
}

func (g *Generator) RegisterEncoderNatives() {
	g.EncoderNatives = map[string]FunctionGeneratorFunc{
		"container": ContainerEncoder,
		"array":     ArrayEncoder,
		"option":    OptionEncoder,
		"bitflags":  TypeForwardEncoder,
		"switch":    SwitchEncoder,
		"mapper":    MapperEncoder,
		"buffer":    BufferEncoder,
		"bitfield":  BitFieldEncoder,

		"void":            DefaultEncoder,
		"varint":          VarIntEncoder,
		"varlong":         VarLongEncoder,
		"anonymousNbt":    DefaultEncoder,
		"anonOptionalNbt": DefaultEncoder, // I have no idea what the difference is between these two
		"UUID":            UUIDEncoder,
		"restBuffer":      DefaultEncoder,
		"string":          StringEncoder,

		"f32":  SimpleTypeEncoder,
		"f64":  SimpleTypeEncoder,
		"i8":   SimpleTypeEncoder,
		"i16":  SimpleTypeEncoder,
		"i32":  SimpleTypeEncoder,
		"i64":  SimpleTypeEncoder,
		"u8":   SimpleTypeEncoder,
		"u16":  SimpleTypeEncoder,
		"u32":  SimpleTypeEncoder,
		"u64":  SimpleTypeEncoder,
		"bool": SimpleTypeEncoder,

		"registryEntryHolder":      RegistryEntryHolderEncoder,
		"registryEntryHolderSet":   ToDoEncoder,
		"entityMetadataLoop":       ToDoEncoder,
		"topBitSetTerminatedArray": ToDoEncoder,
		"todo":                     ToDoEncoder,
	}
}
