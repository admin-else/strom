package modules

import (
	"errors"
	"fmt"

	"github.com/admin-else/strom/bot"
	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/proto_generated/v1_21_8"
)

type ConfigIgnorer struct {
	*bot.Conn
}

func (c *ConfigIgnorer) Default(event any) (err error) {
	fmt.Printf("%#v\n", event)
	return
}

func (c *ConfigIgnorer) OnKnownPacks(packet v1_21_8.PacketCommonSelectKnownPacks) (err error) {
	return c.Send(packet)
}

func (c *ConfigIgnorer) OnFinish(_ v1_21_8.ConfigurationToClientPacketFinishConfiguration) (err error) {
	err = c.Send(v1_21_8.ConfigurationToServerPacketFinishConfiguration{})
	if err != nil {
		return
	}
	c.State = proto_base.Play
	err = bot.HandlerDone
	return
}

func (c *ConfigIgnorer) OnPing(packet v1_21_8.ConfigurationToClientPacketPing) (err error) {
	err = c.Send(v1_21_8.ConfigurationToServerPacketPong(packet))
	return
}

func (c *ConfigIgnorer) OnKeepAlive(packet v1_21_8.ConfigurationToClientPacketKeepAlive) (err error) {
	err = c.Send(v1_21_8.ConfigurationToServerPacketKeepAlive(packet))
	return
}

func IgnoreConfig(c *bot.Conn) (err error) {
	err = c.Start(&ConfigIgnorer{c})
	if err != nil {
		err = errors.Join(err, errors.New("failed to ignore config"))
	}
	return
}
