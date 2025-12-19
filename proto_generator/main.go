package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/admin-else/strom/data"
)

var ToDoError = errors.New("to do")

type Types struct {
	Types map[string]any
}

type State struct {
	ToClient, ToServer Types
}

type Protocol struct {
	Types
	Handshaking, Status, Login, Configuration, Play State
}

type Settings struct {
	ReplaceWithTodoNames []string
}

type ContainerStackEntry struct {
	Data []struct {
		Name string
		Type any
		Anon bool
	}
	VarToSet ast.Expr
}

type Generator struct {
	// These persist the entire generate-call
	Settings         Settings
	Natives          map[string]ExprGeneratorFunc
	DecoderNatives   map[string]FunctionGeneratorFunc
	EncoderNatives   map[string]FunctionGeneratorFunc
	CompareToNatives map[string]CompareToGeneratorFunc
	Protocol         Protocol
	File             *ast.File

	// These persist in the current GenerateType call
	CurrentlyGeneratingTypes Types

	// These persist 1 type generation
	CurrentlyGeneratingTypesPrefix string
	CurrentlyGeneratingTypeName    string

	// Change during type generation
	Depth          int
	Declared       []string
	ContainerStack []ContainerStackEntry
	VtsStack       []ast.Expr
}

func (g *Generator) Decl(name string, t token.Token, e ast.Expr) {
	if slices.Contains(g.Declared, name) {
		return
	}
	g.Declared = append(g.Declared, name)
	AppendDecl(g.File, &ast.GenDecl{
		Tok: t,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{
					Ident(name),
				},
				Values: []ast.Expr{e},
			},
		},
	})
}

func ParseType(t any) (tName string, tData any, err error) {
	switch t := t.(type) {
	case string:
		tName = t
		return
	case []any:
		if len(t) != 2 {
			err = fmt.Errorf("unable to parse type %v", t)
			return
		}
		var ok bool
		tName, ok = t[0].(string)
		if !ok {
			err = fmt.Errorf("unable to parse type %v", t)
			return
		}
		tData = t[1]
		return
	}
	err = fmt.Errorf("unable to parse type %v", t)
	return
}

type ExprGeneratorFunc func(g *Generator, data any) (ast.Expr, error)

type CompareToGeneratorFunc func(g *Generator, parts []string, inExpr ast.Expr, data any) (ast.Expr, error)

type FunctionGeneratorFunc func(g *Generator, varToSet ast.Expr, data any, name string) ([]ast.Stmt, error)

func (g *Generator) VisitType(data any) (e ast.Expr, err error) {
	tName, tData, err := ParseType(data)
	if err != nil {
		return
	}
	return g.VisitNameAndData(tName, tData)
}

func (g *Generator) VisitNameAndData(tName string, tData any) (e ast.Expr, err error) {
	g.Depth += 1
	defer func() { g.Depth -= 1 }()
	n, found := g.Natives[tName]
	if found {
		return n(g, tData)
	}
	t, found := g.Protocol.Types.Types[tName]
	if t != "native" && found {
		return Ident(CamelCase(tName)), nil
	}
	t, found = g.CurrentlyGeneratingTypes.Types[tName]
	if found {
		return Ident(g.CurrentlyGeneratingTypesPrefix + CamelCase(tName)), nil
	}
	err = fmt.Errorf("unknown type %s", tName)
	return

}

func (g *Generator) VisitDecoder(varToSet ast.Expr, data any, name string) (e []ast.Stmt, err error) {
	g.Depth += 1
	defer func() { g.Depth -= 1 }()

	tName, tData, err := ParseType(data)
	if err != nil {
		return
	}
	d, found := g.DecoderNatives[tName]
	if found {
		return d(g, varToSet, tData, name)
	}
	t, found := g.Protocol.Types.Types[tName]
	if t == "native" {
		err = fmt.Errorf("decoder native not implemented %v", tName)
		return
	}
	return DefaultDecoder(g, varToSet, tData, name)
}

func (g *Generator) VisitEncoder(varToSet ast.Expr, data any, name string) (e []ast.Stmt, err error) {
	g.Depth += 1
	defer func() { g.Depth -= 1 }()

	tName, tData, err := ParseType(data)
	if err != nil {
		return
	}
	d, found := g.EncoderNatives[tName]
	if found {
		return d(g, varToSet, tData, name)
	}
	t, found := g.Protocol.Types.Types[tName]
	if t == "native" {
		err = fmt.Errorf("encoder native not implemented %v", tName)
		return
	}
	return DefaultEncoder(g, varToSet, tData, name)
}

// FixPointerReceivers
// We cannot declare a function on a ptr type as described by https://go.dev/ref/spec#Method_declarations.
// So we wrap them in a struct.
func FixPointerReceivers(e ast.Expr) (td ast.Expr, ret ast.Expr) {
	_, isStruct := e.(*ast.StructType)
	ret = Ident("ret")
	if isStruct {
		td = e
		return
	}
	ret = Selector("ret", "Val")
	s := NewStruct()
	AddFieldToStruct(s, "Val", e)
	td = s
	return
}

func (g *Generator) GenerateTypes(prefix string, types Types) error {
	g.CurrentlyGeneratingTypes = types
	g.CurrentlyGeneratingTypesPrefix = prefix
	for _, k := range OrderedKeys(types.Types) {
		v := types.Types[k]
		if slices.Contains(g.Settings.ReplaceWithTodoNames, k) {
			v = "todo"
		}
		g.Depth = 0
		e, err := g.VisitType(v)
		if errors.Is(err, ToDoError) {
			e = Selector("proto_base", "ToDo")
			err = nil
		}
		if err != nil {
			return err
		}
		if e == nil { // This is for types we implement ourselves
			continue
		}

		tName := prefix + CamelCase(k)
		g.CurrentlyGeneratingTypeName = tName

		e, retExpr := FixPointerReceivers(e)
		AppendDecl(g.File, TypeDecl(tName, e))

		args := []NameAndType{{"r", Ident("io.Reader")}}
		rets := []NameAndType{{"ret", Ident(tName)}, {"err", Ident("error")}}
		decodeFunction := NewFuncWithReceiver("Decode", "_", tName, args, rets)
		g.Depth = 0
		s, err := g.VisitDecoder(retExpr, v, tName)
		if err != nil {
			s = ToDoStmts
			if !errors.Is(err, ToDoError) {
				fmt.Println("failed to make decoder for", k, err)
			}
		}
		decodeFunction.Body = NewBlock(append(s, Return()))
		AppendDecl(g.File, decodeFunction)

		args = []NameAndType{{"w", Ident("io.Writer")}}
		rets = []NameAndType{{"err", Ident("error")}}
		encodeFunction := NewFuncWithReceiver("Encode", "ret", tName, args, rets)
		g.Depth = 0
		s, err = g.VisitEncoder(retExpr, v, tName)
		if err != nil {
			s = ToDoStmts
			if !errors.Is(err, ToDoError) {
				fmt.Println("failed to make encoder for", k, err)
			}
		}
		encodeFunction.Body = NewBlock(append(s, Return()))
		AppendDecl(g.File, encodeFunction)

	}
	return nil
}

func (g *Generator) GenerateState(prefix string, state State) (err error) {
	err = g.GenerateTypes(prefix+"ToServer", state.ToServer)
	if err != nil {
		return
	}
	err = g.GenerateTypes(prefix+"ToClient", state.ToClient)
	return
}

func (g *Generator) GenerateProtocol(protocol Protocol) (err error) {
	err = g.GenerateTypes("", protocol.Types)
	if err != nil {
		return
	}
	err = g.GenerateState("Handshaking", protocol.Handshaking)
	if err != nil {
		return
	}
	err = g.GenerateState("Status", protocol.Status)
	if err != nil {
		return
	}
	err = g.GenerateState("Login", protocol.Login)
	if err != nil {
		return
	}
	err = g.GenerateState("Configuration", protocol.Configuration)
	if err != nil {
		return
	}
	err = g.GenerateState("Play", protocol.Play)
	return
}

func Generate(version string, w io.Writer) (err error) {
	protocol := Protocol{}
	err = data.LoadVersionedJson(version, "protocol", &protocol)
	if err != nil {
		return
	}

	g := &Generator{Protocol: protocol}

	// These have annoying edge cases that would require me to implement a whole new compareTo parser
	g.Settings.ReplaceWithTodoNames = []string{"packet_scoreboard_score", "packet_advancements"}

	g.File = NewFile("v" + strings.ReplaceAll(version, ".", "_"))
	AppendDecl(g.File, Import("encoding/binary", "io", "github.com/admin-else/strom/proto_base", "github.com/admin-else/strom/nbt", "github.com/google/uuid"))
	g.RegisterNatives()
	g.RegisterDecoderNatives()
	g.RegisterEncoderNatives()
	g.RegisterCompareToNatives()
	err = g.GenerateProtocol(protocol)
	if err != nil {
		return
	}
	return PrintToFile(g.File, w)
}

func generateVersion(v string) (err error) {
	vUnderscore := strings.ReplaceAll(v, ".", "_")
	err = os.MkdirAll("proto_generated/v"+vUnderscore, 0755)
	if err != nil {
		return
	}
	f, err := os.Create("proto_generated/v" + vUnderscore + "/proto_base.go")
	if err != nil {
		return
	}
	defer f.Close()
	err = Generate(v, f)
	if err != nil {
		return
	}
	helperF, err := os.Create("proto_generated/v" + vUnderscore + "/helpers.go")
	if err != nil {
		return
	}
	defer f.Close()
	err = GenerateHelpers(v, helperF)
	return
}

func GenerateVersions(versions []string) (err error) {
	err = os.RemoveAll("proto_generated/")
	if err != nil {
		return
	}
	err = os.MkdirAll("proto_generated/", 0755)
	if err != nil {
		return
	}
	f, err := os.Create("proto_generated/version_switcher.go")
	if err != nil {
		return
	}
	defer f.Close()
	err = GenerateVersionSwitcher(versions, f)
	for _, version := range versions {
		fmt.Println("generating", version)
		err = generateVersion(version)
		if err != nil {
			return
		}
	}
	return
}

func main() {
	err := GenerateVersions([]string{"1.21.8"})
	if err != nil {
		panic(err)
	}
}
