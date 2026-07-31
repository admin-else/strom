package server

import (
	"github.com/admin-else/strom/mc/data"
	"github.com/admin-else/strom/mc/event"
	"github.com/admin-else/strom/mc/nbt"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_base"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_8"
)

type ConfigServer struct {
	*proto.Conn
	registryData map[string]any
}

// ServeConfig runs the configuration phase: sends known packs, registry data, and tags, then transitions to play.
func ServeConfig(c *proto.Conn) (err error) {
	cs := &ConfigServer{Conn: c}
	n, err := data.LoadRegistry(c.Version)
	if err != nil {
		return
	}
	err = c.Send(&v1_21_8.ConfigurationToClientPacketFeatureFlags{Features: []string{"minecraft:vanilla"}})
	if err != nil {
		return
	}
	cs.registryData = n.Value.(map[string]any)
	packet := &v1_21_8.ConfigurationToClientPacketCommonSelectKnownPacks{}
	if knownPacks, ok := cs.registryData["KnownPacks"]; ok {
		for _, pAny := range knownPacks.([]any) {
			p := pAny.(map[string]any)
			packet.Packs = append(packet.Packs, struct {
				Namespace string
				Id        string
				Version   string
			}{p["Namespace"].(string), p["Id"].(string), p["Version"].(string)})
		}
	}
	err = c.Send(packet)
	if err != nil {
		return
	}
	cs.RegisterCriticalUntilLatest(cs.OnSelectKnownPacksResponse)
	cs.RegisterCriticalUntilLatest(cs.OnFinishConfiguration)
	return cs.StartConn()
}

func (c *ConfigServer) OnSelectKnownPacksResponse(_ *v1_21_8.ConfigurationToServerPacketCommonSelectKnownPacks) (err error) {
	for _, p := range c.registryData["Registries"].([]any) {
		packet := &v1_21_8.ConfigurationToClientPacketRegistryData{Id: p.(map[string]any)["Id"].(string)}
		for _, entry := range p.(map[string]any)["Entries"].([]any) {
			entryTyped := entry.(map[string]any)
			valAny, ok := entryTyped["Value"]
			var val *nbt.Anon = nil
			if ok {
				val = &nbt.Anon{Value: valAny}
			}
			packet.Entries = append(packet.Entries, struct {
				Key   string
				Value *nbt.Anon
			}{entryTyped["TagType"].(string), val})
		}
		err = c.Send(packet)
		if err != nil {
			return
		}
	}
	tagsPacket := &v1_21_8.ConfigurationToClientPacketTags{}
	for _, p := range c.registryData["Tags"].([]any) {
		p := p.(map[string]any)
		tags := v1_21_8.Tags{}
		for _, entry := range p["Entries"].([]any) {
			entry := entry.(map[string]any)
			tags.Val = append(tags.Val, struct {
				TagName string
				Entries []int32
			}{
				TagName: entry["TagName"].(string),
				Entries: entry["Entries"].([]int32),
			})
		}
		tagsPacket.Tags = append(tagsPacket.Tags, struct {
			TagType string
			Tags    v1_21_8.Tags
		}{TagType: p["TagType"].(string), Tags: tags})
	}
	err = c.Send(tagsPacket)
	if err != nil {
		return
	}
	err = c.Send(&v1_21_8.ConfigurationToClientPacketFinishConfiguration{})
	return
}

func (c *ConfigServer) OnFinishConfiguration(_ *v1_21_8.ConfigurationToServerPacketFinishConfiguration) (err error) {
	c.SetState(proto_base.Play)
	err = event.HandlerDoneErr{}
	return
}
