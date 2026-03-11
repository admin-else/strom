package server

import (
	"fmt"
	"log/slog"

	"github.com/admin-else/strom/data"
	"github.com/admin-else/strom/event"
	"github.com/admin-else/strom/nbt"
	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/proto_generated/v1_21_8"
)

type ConfigServer struct {
	*proto.Conn
}

func (c *ConfigServer) Default(event event.Anything) (err error) {
	slog.Debug("Config server", "packet", fmt.Sprintf("%#v", event.Val))
	return
}

func (c *ConfigServer) OnStart() (err error) {
	n, err := data.LoadRegistry(c.Version)
	if err != nil {
		return
	}
	for _, p := range n.Value.(map[string]any)["Registries"].([]any) {
		packet := &v1_21_8.ConfigurationToClientPacketRegistryData{Id: p.(map[string]any)["Id"].(string)}
		for _, entry := range p.(map[string]any)["Entries"].([]any) {
			entryTyped := entry.(map[string]any)
			packet.Entries = append(packet.Entries, struct {
				Key   string
				Value *nbt.Anon
			}{entryTyped["TagType"].(string), &nbt.Anon{Value: entryTyped["Value"]}})
		}
		err = c.Send(packet)
		if err != nil {
			return
		}
	}
	tagsPacket := &v1_21_8.ConfigurationToClientPacketTags{}
	for _, p := range n.Value.(map[string]any)["Tags"].([]any) {
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

func ServeConfig(c *proto.Conn) (err error) {
	cs := &ConfigServer{c}
	err = cs.Start(cs)
	return
}
