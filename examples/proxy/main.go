package main

import (
	"fmt"

	"github.com/admin-else/strom"
)

type Proxy struct{ *strom.Conn }

func (p *Proxy) Default(event any) (err error) {
	fmt.Printf("%#v\n", event)
	return
}

func (p *Proxy) OnStart(_ strom.OnStart) (err error) {
	return
}

func main() {
	err := strom.StartServer(":25565", func(c *strom.Conn) (h strom.HandlerInst, err error) {
		return &Proxy{c}, nil
	})
	if err != nil {
		panic(err)
	}
}
