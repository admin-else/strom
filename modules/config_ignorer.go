package modules

import (
	"strom"

	"github.com/admin-else/queser"
	"github.com/admin-else/queser/generated/v1_21_8"
)

type ConfigIgnorer struct {
	*strom.Conn
}

func (c *ConfigIgnorer) Default(event any) (err error) {
	return
}

func (c *ConfigIgnorer) OnKnownPacks(packet v1_21_8.PacketCommonSelectKnownPacks) (err error) {
	return c.Send(packet)
}

func (c *ConfigIgnorer) OnFinish(packet v1_21_8.ConfigurationToClientPacketFinishConfiguration) (err error) {
	err = c.Send(v1_21_8.ConfigurationToServerPacketFinishConfiguration{})
	if err != nil {
		return
	}
	c.State = queser.Play
	err = strom.HandlerDone
	return
}

func IgnoreConfig(c *strom.Conn) (err error) {
	return c.Start(&ConfigIgnorer{c})
}
