package client

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/admin-else/strom/event"
	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/proto_generated/v1_21_8"
)

type ConfigIgnorer struct {
	*proto.Conn
}

func (c *ConfigIgnorer) OnStart() (err error) {
	c.Register(c.Default)
	c.RegisterUntilLatest(c.OnKnownPacks)
	c.RegisterUntilLatest(c.OnFinish)
	c.RegisterUntilLatest(c.OnPing)
	c.RegisterUntilLatest(c.OnKeepAlive)
	return
}

func (c *ConfigIgnorer) Default(event event.Anything) (err error) {
	slog.Debug("Config ignorer", "packet", fmt.Sprintf("%#v", event))
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

func IgnoreConfig(c *proto.Conn) (err error) {
	err = c.Start(&ConfigIgnorer{c})
	if err != nil {
		err = errors.Join(err, errors.New("failed to ignore config"))
	}
	return
}
