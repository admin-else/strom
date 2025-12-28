package main

import (
	"go/ast"
	"os"
	"strconv"
	"strings"

	"github.com/admin-else/strom/data"
)

type TypesInfo struct {
	Types
	Direction string
	State     string
}

type PacketInfo struct {
	TName, Name, Direction, State, PacketId, ProtocolVersion, MinecraftVersion string
}

func CollectPacketInfo(versions []string) (packetInfos []PacketInfo, err error) {
	var versionData data.ProtocolVersion
	for _, version := range versions {
		versionData, err = data.LookUpProtocolVersionByName(version)
		if err != nil {
			return
		}
		protocol := Protocol{}
		err = data.LoadVersionedJson(version, "protocol", &protocol)
		if err != nil {
			return
		}
		prefixTypeMap := map[string]TypesInfo{
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
					var typeName string
					if strings.HasPrefix(v2, "packet_common") {
						typeName = CamelCase(v2)
					} else if v2 == "void" {
						typeName = "struct{}"
					} else {
						typeName = prefix + CamelCase(v2)
					}
					packetInfos = append(packetInfos, PacketInfo{
						TName:            typeName,
						Name:             k2,
						Direction:        types.Direction,
						State:            types.State,
						PacketId:         packetIdsRev[k2],
						ProtocolVersion:  strconv.Itoa(versionData.Version),
						MinecraftVersion: version,
					})
				}
			}
		}
	}
	return
}

func GeneratePacketInfoFile(versions []string) (err error) {
	var packetInfos []PacketInfo
	packetInfos, err = CollectPacketInfo(versions)
	if err != nil {
		return
	}

	f := NewFile("proto_generated")
	imports := []string{"reflect", "github.com/admin-else/strom/proto_base"}
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
		if typeName == "struct{}" {
			typeExpr = &ast.StructType{Fields: &ast.FieldList{}}
		} else if strings.HasPrefix(typeName, "proto_base.") {
			typeExpr = Selector("proto_base", strings.TrimPrefix(typeName, "proto_base."))
		} else {
			typeExpr = Selector(vUnderscore, typeName)
		}

		packetId, _ := strconv.ParseInt(p.PacketId, 0, 32)
		protocolVersion, _ := strconv.Atoi(p.ProtocolVersion)

		packetInfoExprs = append(packetInfoExprs, CompLit(Selector("proto_base", "PacketInfo"), []ast.Expr{
			KeyValueExpr(Ident("Type"), Call(Selector("reflect", "TypeOf"), CompLit(typeExpr, nil))),
			KeyValueExpr(Ident("Name"), StrLit(p.Name)),
			KeyValueExpr(Ident("Direction"), Selector("proto_base", p.Direction)),
			KeyValueExpr(Ident("State"), Selector("proto_base", p.State)),
			KeyValueExpr(Ident("PacketId"), NumLit(int(packetId))),
			KeyValueExpr(Ident("ProtocolVersion"), NumLit(protocolVersion)),
		}))
	}
	AppendDecl(f, VarValues("Packets", CompLit(Slice(Selector("proto_base", "PacketInfo")), packetInfoExprs)))
	out, err := os.Create("proto_generated/packet_info.go")
	if err != nil {
		return
	}
	defer out.Close()

	return PrintToFile(f, out)
}
