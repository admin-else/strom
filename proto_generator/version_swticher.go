package main

import (
	"go/ast"
	"io"
	"strings"
)

func GeneratePacketIdentifierToType(versions []string) (decl ast.Decl) {
	var s []ast.Stmt
	for _, version := range versions {
		s = append(s, Case(Exprs(StrLit(version)), Stmts(
			Assign121(Ident("t"), Call(SelectorExprAndStr(Ident("v"+strings.ReplaceAll(version, ".", "_")), "PacketIdentifierToType"), Ident("d"), Ident("s"), Ident("i"))),
		)))
	}
	s = append(s, Case(nil, Stmts()))
	fn := NewFunc("PacketIdentifierToType", []NameAndType{
		{"v", Ident("string")},
		{"d", Ident("proto_base.Direction")},
		{"s", Ident("proto_base.State")},
		{"i", Ident("string")},
	}, []NameAndType{
		{"t", Ident("any")},
	})
	fn.Body = NewBlockEllipsis(SwitchStmt(Ident("v"), NewBlock(s)), Return())
	decl = fn
	return
}

/*
func PacketIdentifierToType(v string, d proto_base.Direction, s proto_base.State, i string) (t any) {
	switch v {
	case "1.21.8":
		t = v1_21_8.PacketIdentifierToType(d, s, i)
	default:

	}
	return
}

func TypeToPacketIdentifier(v string, d proto_base.Direction, s proto_base.State, t any) (i string) {
	switch v {
	case "1.21.8":
		i = v1_21_8.TypeToPacketIdentifier(d, s, t)
	}
	return
}

func DecodePacket(v string, d proto_base.Direction, s proto_base.State, r io.Reader) (packet any, err error) {
	switch v {
	case "1.21.8":
		packet, err = v1_21_8.DecodePacket(d, s, r)
	}
	return
}

func EncodePacket(v string, d proto_base.Direction, s proto_base.State, i string, p any, w io.Writer) (err error) {
	switch v {
	case "1.21.8":
		err = v1_21_8.EncodePacket(d, s, i, p, w)

	}
	return
}

*/

func GenerateTypeToPacketIdentifier(versions []string) (decl ast.Decl) {
	var s []ast.Stmt
	for _, version := range versions {
		s = append(s, Case(Exprs(StrLit(version)), Stmts(
			Assign121(Ident("i"), Call(SelectorExprAndStr(Ident("v"+strings.ReplaceAll(version, ".", "_")), "TypeToPacketIdentifier"), Ident("d"), Ident("s"), Ident("t"))),
		)))
	}
	fn := NewFunc("TypeToPacketIdentifier", []NameAndType{
		{"v", Ident("string")},
		{"d", Ident("proto_base.Direction")},
		{"s", Ident("proto_base.State")},
		{"t", Ident("any")},
	}, []NameAndType{
		{"i", Ident("string")},
	})
	fn.Body = NewBlockEllipsis(SwitchStmt(Ident("v"), NewBlock(s)), Return())
	decl = fn
	return
}

func GenerateDecodePacket(versions []string) (decl ast.Decl) {
	var s []ast.Stmt
	for _, version := range versions {
		s = append(s, Case(Exprs(StrLit(version)), Stmts(
			Assign([]ast.Expr{Ident("packet"), Ident("err")}, []ast.Expr{Call(SelectorExprAndStr(Ident("v"+strings.ReplaceAll(version, ".", "_")), "DecodePacket"), Ident("d"), Ident("s"), Ident("r"))}),
		)))
	}
	fn := NewFunc("DecodePacket", []NameAndType{
		{"v", Ident("string")},
		{"d", Ident("proto_base.Direction")},
		{"s", Ident("proto_base.State")},
		{"r", Ident("io.Reader")},
	}, []NameAndType{
		{"packet", Ident("any")},
		{"err", Ident("error")},
	})
	fn.Body = NewBlockEllipsis(SwitchStmt(Ident("v"), NewBlock(s)), Return())
	decl = fn
	return
}

func GenerateEncodePacket(versions []string) (decl ast.Decl) {
	var s []ast.Stmt
	for _, version := range versions {
		s = append(s, Case(Exprs(StrLit(version)), Stmts(
			Assign121(Ident("err"), Call(SelectorExprAndStr(Ident("v"+strings.ReplaceAll(version, ".", "_")), "EncodePacket"), Ident("d"), Ident("s"), Ident("i"), Ident("p"), Ident("w"))),
		)))
	}
	fn := NewFunc("EncodePacket", []NameAndType{
		{"v", Ident("string")},
		{"d", Ident("proto_base.Direction")},
		{"s", Ident("proto_base.State")},
		{"i", Ident("string")},
		{"p", Ident("any")},
		{"w", Ident("io.Writer")},
	}, []NameAndType{
		{"err", Ident("error")},
	})
	fn.Body = NewBlockEllipsis(SwitchStmt(Ident("v"), NewBlock(s)), Return())
	decl = fn
	return
}

func GenerateVersionSwitcher(versions []string, printTo io.Writer) (err error) {
	f := NewFile("proto_generated")
	importPaths := []string{"io", "github.com/admin-else/strom/proto_base"}
	for _, version := range versions {
		importPaths = append(importPaths, "github.com/admin-else/strom/proto_generated/v"+strings.ReplaceAll(version, ".", "_"))
	}
	AppendDecl(f, Import(importPaths...))
	AppendDecl(f, GeneratePacketIdentifierToType(versions))
	AppendDecl(f, GenerateTypeToPacketIdentifier(versions))
	AppendDecl(f, GenerateDecodePacket(versions))
	AppendDecl(f, GenerateEncodePacket(versions))
	return PrintToFile(f, printTo)
}
