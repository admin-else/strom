package modules

import (
	"errors"
	"fmt"

	"github.com/admin-else/strom"

	"github.com/admin-else/queser"
	"github.com/admin-else/queser/generated/v1_21_8"
)

type ConfigIgnorer struct {
	*strom.Conn
}

func (c *ConfigIgnorer) Default(event any) (err error) {
	fmt.Printf("%#v\n", event)
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

func (c *ConfigIgnorer) OnPing(packet v1_21_8.ConfigurationToClientPacketPing) (err error) {
	err = c.Send(v1_21_8.ConfigurationToServerPacketPong(packet))
	return
}

func (c *ConfigIgnorer) OnKeepAlive(packet v1_21_8.ConfigurationToClientPacketKeepAlive) (err error) {
	err = c.Send(v1_21_8.ConfigurationToServerPacketKeepAlive(packet))
	return
}

func IgnoreConfig(c *strom.Conn) (err error) {
	err = c.Start(&ConfigIgnorer{c})
	if err != nil {
		err = errors.Join(err, errors.New("failed to ignore config"))
	}
	return
}
