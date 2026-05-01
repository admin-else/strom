package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"

	"github.com/admin-else/strom/api"
	"github.com/admin-else/strom/client"
	"github.com/admin-else/strom/data"
	"github.com/admin-else/strom/event"
	"github.com/admin-else/strom/nbt"
	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/proto_generated/v1_21_8"
	"github.com/admin-else/strom/server"
	"github.com/admin-else/strom/text"
)

func wontFail[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// THIS IS REDUNDANT USE LOGIN PACKET IN MINECRAFT DATA

var (
	// Edit these two
	Version = "1.21.11"
	SaveAs  = "data/registry/" + Version + ".nbt"

	CompatibleProtocolVersion = wontFail(data.LookUpVersionByName(Version)).Version
	FinishedConfigErr         = errors.New("finished configuration")
)

var Status = wontFail(json.Marshal(server.StatusResponse{
	Version: server.StatusResponseVersion{
		Name:     nil,
		Protocol: CompatibleProtocolVersion,
	},
	Description: text.Pretty("Registry data hunter"),
}))

type Proxy struct {
	Client, Servee *proto.Conn
	Data           nbt.Tag
}

func (p *Proxy) OnAnything(e event.Anything) (err error) {
	packet, ok := e.Val.(proto_base.EncodeDecodeAble)
	if !ok {
		return
	}
	packetInfo, ok := proto.LookupPacketInfoByType(packet)
	if !ok {
		return
	}
	switch packetInfo.Direction {
	case proto_base.ToServer:
		err = p.Client.Send(packet)
	case proto_base.ToClient:
		err = p.Servee.Send(packet)
	}
	return
}

func (p *Proxy) OnFinishConfiguration(_ *v1_21_8.ConfigurationToClientPacketFinishConfiguration) (err error) {
	err = FinishedConfigErr
	return
}

func (p *Proxy) OnSelectKnownPacks(data *v1_21_8.ConfigurationToClientPacketCommonSelectKnownPacks) (err error) {
	var entries []any
	for _, entry := range data.Packs {
		e := map[string]any{"Namespace": entry.Namespace, "Id": entry.Id, "Version": entry.Version}
		entries = append(entries, e)
	}
	p.Data.Value.(map[string]any)["KnownPacks"] = entries
	return
}

func (p *Proxy) OnData(data *v1_21_8.ConfigurationToClientPacketRegistryData) (err error) {
	d := p.Data.Value.(map[string]any)["Registries"].([]any)
	var entries []any
	for _, entry := range data.Entries {
		if entry.Value != nil {
			entries = append(entries, map[string]any{"TagType": entry.Key, "Value": entry.Value.Value})
		} else {
			entries = append(entries, map[string]any{"TagType": entry.Key})
		}
	}
	d = append(d, map[string]any{"Id": data.Id, "Entries": entries})
	p.Data.Value.(map[string]any)["Registries"] = d
	return
}

func (p *Proxy) OnTags(data *v1_21_8.ConfigurationToClientPacketTags) (err error) {
	var tags []any
	for _, tag := range data.Tags {
		var entries []any
		for _, entry := range tag.Tags.Val {
			entries = append(entries, map[string]any{"TagName": entry.TagName, "Entries": entry.Entries})
		}
		tags = append(tags, map[string]any{"TagType": tag.TagType, "Entries": entries})
	}
	p.Data.Value.(map[string]any)["Tags"] = tags
	return
}

func handleClient(c *proto.Conn) (err error) {
	acc := api.NewOfflineAccount("RegistryHunter")

	// ignore resource leak warning we own c conn so we should not close it here
	_, err = server.ServeLogin(c, server.WithOtherAccount(acc), server.WithCompatibleVersions(CompatibleProtocolVersion), server.WithRawStatus(Status))
	if err != nil {
		return
	}
	sockPuppet, err := client.Login(":25565", acc)
	if err != nil {
		return
	}
	defer sockPuppet.Close()
	errChan := make(chan error)
	p := &Proxy{Client: sockPuppet, Servee: c}
	p.Data.Value = map[string]any{"Registries": []any{}}
	p.Client.RegisterCritical(p.OnAnything)
	p.Servee.RegisterCritical(p.OnAnything)
	p.Client.RegisterCriticalUntilLatest(p.OnData)
	p.Client.RegisterCriticalUntilLatest(p.OnTags)
	p.Client.RegisterCriticalUntilLatest(p.OnSelectKnownPacks)
	p.Client.RegisterCriticalUntilLatest(p.OnFinishConfiguration)

	go func() {
		errChan <- p.Client.StartConn()
	}()
	go func() {
		errChan <- p.Servee.StartConn()
	}()
	err = <-errChan
	f, err := os.Create(SaveAs)
	if err != nil {
		return
	}
	defer f.Close()
	err = nbt.WriteUnstructuredFile(f, p.Data)
	slog.Info("saved", "file", SaveAs)
	return
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	err := server.StartServerWithOnConn(":25566", handleClient)
	if err != nil {
		panic(err)
	}
}
