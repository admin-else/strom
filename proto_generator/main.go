package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"slices"
	"strconv"
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

type ContainerStackEntry struct {
	Data []struct {
		Name string
		Type any
		Anon bool
	}
	VarToSet ast.Expr
}

type ExprGeneratorFunc func(g *Generator, data any) (ast.Expr, error)

type CompareToGeneratorFunc func(g *Generator, parts []string, inExpr ast.Expr, data any) (e ast.Expr, cet CaseExprType, err error)

type FunctionGeneratorFunc func(g *Generator, varToSet ast.Expr, data any, name string) ([]ast.Stmt, error)

type Generator struct {
	packetInfos []PacketInfo

	// These persist the entire generate-call
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

func (g *Generator) GenerateTypes(prefix string, types Types) (err error) {
	g.CurrentlyGeneratingTypes = types
	g.CurrentlyGeneratingTypesPrefix = prefix
	for _, k := range OrderedKeys(types.Types) {
		v := types.Types[k]
		g.Depth = 0
		e, err2 := g.VisitType(v)
		if errors.Is(err2, ToDoError) {
			e = Selector("proto_base", "ToDo")
			err2 = nil
		}
		if err2 != nil {
			return err2
		}
		if e == nil { // This is for types we implement ourselves
			continue
		}

		tName := prefix + CamelCase(k)
		g.CurrentlyGeneratingTypeName = tName

		e, retExpr := FixPointerReceivers(e)
		AppendDecl(g.File, TypeDecl(tName, e))

		args := []NameAndType{{"r", Ident("io.Reader")}}
		rets := []NameAndType{{"err", Ident("error")}}
		decodeFunction := NewFuncWithReceiver("Decode", "ret", Pointer(Ident(tName)), args, rets)
		g.Depth = 0
		s, err2 := g.VisitDecoder(retExpr, v, tName)
		if err2 != nil {
			s = ToDoStmts
			fmt.Println("failed to make decoder for", k, err2)
		}
		decodeFunction.Body = NewBlock(append(s, Return()))
		AppendDecl(g.File, decodeFunction)

		args = []NameAndType{{"w", Ident("io.Writer")}}
		rets = []NameAndType{{"err", Ident("error")}}
		encodeFunction := NewFuncWithReceiver("Encode", "ret", Pointer(Ident(tName)), args, rets)
		g.Depth = 0
		s, err2 = g.VisitEncoder(retExpr, v, tName)
		if err2 != nil {
			s = ToDoStmts
			fmt.Println("failed to make encoder for", k, err2)
		}
		encodeFunction.Body = NewBlock(append(s, Return()))
		AppendDecl(g.File, encodeFunction)
	}
	return
}

func (g *Generator) GenerateProtocol(protocol Protocol, version string) (err error) {
	versionData, err := data.LookUpVersionByName(version)
	if err != nil {
		return
	}

	prefixTypeMap := map[string]TypesInfo{
		"":                      {protocol.Types, "", ""},
		"HandshakingToServer":   {protocol.Handshaking.ToServer, "ToServer", "Handshaking"},
		"HandshakingToClient":   {protocol.Handshaking.ToClient, "ToClient", "Handshaking"},
		"StatusToServer":        {protocol.Status.ToServer, "ToServer", "Status"},
		"StatusToClient":        {protocol.Status.ToClient, "ToClient", "Status"},
		"LoginToServer":         {protocol.Login.ToServer, "ToServer", "Login"},
		"LoginToClient":         {protocol.Login.ToClient, "ToClient", "Login"},
		"ConfigurationToServer": {protocol.Configuration.ToServer, "ToServer", "Configuration"},
		"ConfigurationToClient": {protocol.Configuration.ToClient, "ToClient", "Configuration"},
		"PlayToServer":          {protocol.Play.ToServer, "ToServer", "Play"},
		"PlayToClient":          {protocol.Play.ToClient, "ToClient", "Play"},
	}

	// Move common packet to their own areas to make packets more identifiable
	for _, k := range OrderedKeys(protocol.Types.Types) {
		v := protocol.Types.Types[k]
		if strings.Contains(k, "packet") {
			delete(protocol.Types.Types, k)
			for prefix, types := range prefixTypeMap {
				if prefix == "" {
					continue
				}
				types.Types.Types[k] = v
			}
		}
	}

	for prefix, types := range prefixTypeMap {
		for k, v := range types.Types.Types {
			if k != "packet" {
				continue
			}
			packetIds := v.([]any)[1].([]any)[0].(map[string]any)["type"].([]any)[1].(map[string]any)["mappings"].(map[string]any)
			packetIdsStrings := AssertAndConvertMapValues[string](packetIds)
			packetIdsRev := ReverseMap(packetIdsStrings)
			v := v.([]any)[1].([]any)[1].(map[string]any)["type"].([]any)[1].(map[string]any)["fields"].(map[string]any)
			for k2, v2 := range v {
				v2 := v2.(string)
				typeName := prefix + CamelCase(v2)
				g.packetInfos = append(g.packetInfos, PacketInfo{
					PacketDef:        v2,
					TName:            typeName,
					Name:             k2,
					Direction:        types.Direction,
					State:            types.State,
					PacketId:         packetIdsRev[k2],
					ProtocolVersion:  strconv.Itoa(int(versionData.Version)),
					MinecraftVersion: version,
				})
			}
		}
	}

	for prefix, types := range prefixTypeMap {
		err = g.GenerateTypes(prefix, types.Types)
		if err != nil {
			return
		}
	}
	return
}

func Generate(version string, w io.Writer) (packetInfos []PacketInfo, err error) {
	protocol := Protocol{}
	err = data.LoadVersionedJson(version, "protocol", &protocol)
	if err != nil {
		return
	}

	g := &Generator{Protocol: protocol}

	g.File = NewFile("v" + strings.ReplaceAll(version, ".", "_"))
	AppendDecl(g.File, Import("encoding/binary", "io", "github.com/admin-else/strom/proto_base", "github.com/admin-else/strom/nbt", "github.com/google/uuid"))
	g.RegisterNatives()
	g.RegisterDecoderNatives()
	g.RegisterEncoderNatives()
	g.RegisterCompareToNatives()
	err = g.GenerateProtocol(protocol, version)
	if err != nil {
		return
	}
	packetInfos = g.packetInfos
	err = PrintToFile(g.File, w)
	return
}

func generateVersion(v string) (packetInfos []PacketInfo, err error) {
	vUnderscore := strings.ReplaceAll(v, ".", "_")
	err = os.MkdirAll("proto_generated/v"+vUnderscore, 0755)
	if err != nil {
		return
	}
	f, err := os.Create("proto_generated/v" + vUnderscore + "/proto.go")
	if err != nil {
		return
	}
	defer f.Close()
	packetInfos, err = Generate(v, f)
	if err != nil {
		return
	}
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
	var packetInfos, newPacketInfos []PacketInfo
	for _, version := range versions {
		fmt.Println("generating", version)
		newPacketInfos, err = generateVersion(version)
		if err != nil {
			return
		}
		packetInfos = append(packetInfos, newPacketInfos...)
	}
	return GeneratePacketInfoFile(versions, packetInfos)
}

func main() {
	err := GenerateVersions([]string{"1.21.8", "1.21.9", "1.21.11"})
	if err != nil {
		panic(err)
	}
}
