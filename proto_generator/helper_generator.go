package main

import (
	"go/ast"
	"os"
	"strconv"
	"strings"
)

type TypesInfo struct {
	Types
	Direction string
	State     string
}

type PacketInfo struct {
	TName, Name, Direction, State, PacketId, ProtocolVersion, MinecraftVersion string
	PacketDef                                                                  any
}

func GeneratePacketInfoFile(versions []string, packetInfos []PacketInfo) (err error) {
	f := NewFile("proto_generated")
	imports := []string{"github.com/admin-else/strom/proto_base"}
	for _, v := range versions {
		vUnderscore := strings.ReplaceAll(v, ".", "_")
		imports = append(imports, "github.com/admin-else/strom/proto_generated/v"+vUnderscore)
	}
	AppendDecl(f, Import(imports...))

	var packetInfoExprs []ast.Expr
	for _, p := range packetInfos {
		vUnderscore := "v" + strings.ReplaceAll(p.MinecraftVersion, ".", "_")
		typeName := p.TName
		var typeExpr ast.Expr
		if p.PacketDef == "void" {
			typeExpr = Selector("proto_base", "Void")
		} else {
			typeExpr = Selector(vUnderscore, typeName)
		}

		packetId, _ := strconv.ParseInt(p.PacketId, 0, 32)
		protocolVersion, _ := strconv.Atoi(p.ProtocolVersion)

		packetInfoExprs = append(packetInfoExprs, CompLit(nil, []ast.Expr{
			KeyValueExpr(Ident("Type"), AddrOf(CompLit(typeExpr, nil))),
			KeyValueExpr(Ident("Name"), StrLit(p.Name)),
			KeyValueExpr(Ident("Direction"), Selector("proto_base", p.Direction)),
			KeyValueExpr(Ident("State"), Selector("proto_base", p.State)),
			KeyValueExpr(Ident("PacketId"), NumLit(int(packetId))),
			KeyValueExpr(Ident("ProtocolVersion"), NumLit(protocolVersion)),
		}))
	}

	var supportedVersionsExprs []ast.Expr
	for _, v := range versions {
		supportedVersionsExprs = append(supportedVersionsExprs, StrLit(v))
	}
	AppendDecl(f, VarValues("SupportedVersions", CompLit(Slice(Ident("string")), supportedVersionsExprs)))
	AppendDecl(f, VarValues("Packets", CompLit(Slice(Selector("proto_base", "PacketInfo")), packetInfoExprs)))
	out, err := os.Create("proto_generated/packet_info.go")
	if err != nil {
		return
	}
	defer out.Close()

	return PrintToFile(f, out)
}
