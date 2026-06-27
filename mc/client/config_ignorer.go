package client

import (
	"errors"
	"fmt"

	"github.com/admin-else/strom/mc/event"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_base"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_8"
)

type ConfigIgnorer struct {
	*proto.Conn
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

func IgnoreConfig(c *proto.Conn) (err error) {
	ci := &ConfigIgnorer{c}
	ci.RegisterCritical(ci.Default)
	ci.RegisterUntilLatest(ci.OnKnownPacks)
	ci.RegisterUntilLatest(ci.OnFinish)
	ci.RegisterUntilLatest(ci.OnPing)
	ci.RegisterUntilLatest(ci.OnKeepAlive)

	err = ci.StartConn()
	if err != nil {
		err = errors.Join(err, errors.New("failed to ignore config"))
	}
	return
}
