package main

import (
	"go/ast"
	"io"
	"strings"

	_ "embed"

	"github.com/admin-else/queser/data"
)

//go:embed helper_append.txt
var helperAppend string

func GenerateHelperTypes(t Types, prefix string, f *ast.File) (err error) {
	for k, v := range t.Types {
		if k != "packet" {
			continue
		}
		v := v.([]any)[1].([]any)[1].(map[string]any)["type"].([]any)[1].(map[string]any)["fields"].(map[string]any)
		m := make(map[string]string)
		for k2, v2 := range v {
			v2 := v2.(string)
			if strings.HasPrefix(v2, "packet_common") {
				m[k2] = CamelCase(v2)
			} else if v2 == "void" {
				m[k2] = "queser.Void"
			} else {
				m[k2] = prefix + CamelCase(v2)
			}
		}
		err = PacketIdentifierToTypeGenerator(prefix, f, m)
		if err != nil {
			return
		}
		m = ReverseMap(m)
		err = TypeToPacketIdentifierGenerator(prefix, f, m)
		if err != nil {
			return
		}
	}
	return
}

func GenerateHelperState(s State, prefix string, f *ast.File) (err error) {
	err = GenerateHelperTypes(s.ToServer, prefix+"ToServer", f)
	if err != nil {
		return
	}
	return GenerateHelperTypes(s.ToClient, prefix+"ToClient", f)
}

func GenerateHelperProtocol(p Protocol, f *ast.File) (err error) {
	err = GenerateHelperState(p.Handshaking, "Handshaking", f)
	if err != nil {
		return
	}
	err = GenerateHelperState(p.Status, "Status", f)
	if err != nil {
		return
	}
	err = GenerateHelperState(p.Login, "Login", f)
	if err != nil {
		return
	}
	err = GenerateHelperState(p.Configuration, "Configuration", f)
	if err != nil {
		return
	}
	err = GenerateHelperState(p.Play, "Play", f)
	if err != nil {
		return
	}
	return
}

//	PacketIdentifierToTypeGenerator func PacketIdentifierToType(s string) any {
//		switch s {
//		case "abilities":
//			return PlayToClientPacketAbilities{}
//		default:
//			return nil
//		}
//	}
func PacketIdentifierToTypeGenerator(funcNamePrefix string, f *ast.File, m map[string]string) (err error) {
	var s []ast.Stmt
	for _, k := range OrderedKeys(m) {
		v := m[k]
		s = append(s, Case(Exprs(StrLit(k)), Stmts(Assign121(Ident("t"), CompLit(Ident(v), nil)))))
	}
	s = append(s, Case(nil, Stmts(Assign121(Ident("t"), Nil()))))
	fn := NewFunc(funcNamePrefix+"PacketIdentifierToType", []NameAndType{{"s", Ident("string")}}, []NameAndType{{"t", Ident("any")}})
	fn.Body = NewBlockEllipsis(SwitchStmt(Ident("s"), NewBlock(s)), Return())
	AppendDecl(f, fn)
	return
}

func TypeToPacketIdentifierGenerator(funcNamePrefix string, f *ast.File, m map[string]string) (err error) {
	var s []ast.Stmt
	for _, k := range OrderedKeys(m) {
		v := m[k]
		s = append(s, Case(Exprs(Ident(k)), Stmts(Assign121(Ident("s"), StrLit(v)))))
	}
	fn := NewFunc(funcNamePrefix+"TypeToPacketIdentifier", []NameAndType{{"t", Ident("any")}}, []NameAndType{{"s", Ident("string")}})
	fn.Body = NewBlockEllipsis(TypeSwitch(ExprStmt(TypeAssert(Ident("t"), Ident("type"))), NewBlock(s)), Return())
	AppendDecl(f, fn)
	return
}

func GenerateHelpers(version string, w io.Writer) (err error) {
	protocol := Protocol{}
	err = data.LoadVersionedJson(version, "protocol", &protocol)
	if err != nil {
		return
	}
	f := NewFile("v" + strings.ReplaceAll(version, ".", "_"))
	AppendDecl(f, Import("io", "github.com/admin-else/queser"))
	err = GenerateHelperProtocol(protocol, f)
	if err != nil {
		return
	}
	err = PrintToFile(f, w)
	if err != nil {
		return
	}
	_, err = w.Write([]byte(helperAppend))
	return
}
