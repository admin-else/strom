package client

import (
	"errors"
	"fmt"

	"github.com/admin-else/strom/mc/event"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_base"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_8"
	"github.com/admin-else/strom/mc/registry"
)

type ConfigIgnorer struct {
	*proto.Conn
	// registries, when non-nil, receives every configuration RegistryData
	// packet in wire (id) order.
	registries *registry.Store
}

// OnRegistryData records one dynamic registry sent during configuration. The
// optional element payload (NBT) is consumed by the packet decoder but not
// stored; only the ordered resource ids matter for id mapping.
func (c *ConfigIgnorer) OnRegistryData(packet *v1_21_8.ConfigurationToClientPacketRegistryData) (err error) {
	if c.registries == nil {
		return
	}
	r := registry.NewRegistry(packet.Id)
	for _, entry := range packet.Entries {
		r.Add(entry.Key)
	}
	c.registries.Set(r)
	return
}

func (c *ConfigIgnorer) OnStart() (err error) {
	return
}

func (c *ConfigIgnorer) Default(e event.Anything) (err error) {
	c.Log.Debug("Config ignorer", "packet", fmt.Sprintf("%#v", e))
	return
}

func (c *ConfigIgnorer) OnKnownPacks(packet *v1_21_8.ConfigurationToClientPacketCommonSelectKnownPacks) (err error) {
	return c.Send(&v1_21_8.ConfigurationToServerPacketCommonSelectKnownPacks{Packs: packet.Packs})
}

func (c *ConfigIgnorer) OnFinish(_ *v1_21_8.ConfigurationToClientPacketFinishConfiguration) (err error) {
	err = c.Send(&v1_21_8.ConfigurationToServerPacketFinishConfiguration{})
	if err != nil {
		return
	}
	c.SetState(proto_base.Play)
	err = event.HandlerDoneErr{}
	return
}

func (c *ConfigIgnorer) OnPing(packet *v1_21_8.ConfigurationToClientPacketPing) (err error) {
	err = c.Send(&v1_21_8.ConfigurationToServerPacketPong{Id: packet.Id})
	return
}

func (c *ConfigIgnorer) OnKeepAlive(packet *v1_21_8.ConfigurationToClientPacketKeepAlive) (err error) {
	err = c.Send(&v1_21_8.ConfigurationToServerPacketKeepAlive{KeepAliveId: packet.KeepAliveId})
	return
}

// IgnoreConfig acknowledges configuration packets on the connection and transitions it to the Play state once configuration is complete.
func IgnoreConfig(c *proto.Conn) (err error) {
	return IgnoreConfigStore(c, nil)
}

// IgnoreConfigStore is IgnoreConfig plus dynamic-registry capture: every
// configuration RegistryData packet is added to store before the connection
// transitions to the Play state. A nil store discards the registries.
func IgnoreConfigStore(c *proto.Conn, store *registry.Store) (err error) {
	ci := &ConfigIgnorer{Conn: c, registries: store}
	ci.RegisterCritical(ci.Default)
	ci.RegisterUntilLatest(ci.OnKnownPacks)
	ci.RegisterUntilLatest(ci.OnRegistryData)
	ci.RegisterUntilLatest(ci.OnFinish)
	ci.RegisterUntilLatest(ci.OnPing)
	ci.RegisterUntilLatest(ci.OnKeepAlive)

	err = ci.StartConn()
	if err != nil {
		err = errors.Join(err, errors.New("failed to ignore config"))
	}
	return
}
