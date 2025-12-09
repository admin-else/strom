package main

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/admin-else/strom"

	"github.com/admin-else/queser"
	"github.com/admin-else/queser/generated/v1_21_8"
)

type ProxyClient struct {
	*strom.Conn
	Proxy *Proxy
}

func (p *ProxyClient) OnLoopCycle(event strom.OnLoopCycle) (err error) {
	return
}

func (p *ProxyClient) Default(event any) (err error) {
	fmt.Printf("S2C %#v\n", event)
	err = p.Proxy.Send(event)
	if err != nil {
		return
	}
	return
}

func (p *ProxyClient) Compress(packet v1_21_8.LoginToClientPacketCompress) (err error) {
	err = p.Proxy.Send(packet)
	if err != nil {
		return
	}
	p.CompressionThreshold = int32(packet.Threshold)
	p.Proxy.CompressionThreshold = int32(packet.Threshold)
	return
}

func (p *ProxyClient) OnStart(_ strom.OnStart) (err error) {
	return
}

type Proxy struct {
	*strom.Conn
	ProxyClient *ProxyClient
}

func (p *Proxy) OnLastTrigger(event strom.OnLoopCycle) (err error) {
	return
}

func (p *Proxy) Default(event any) (err error) {
	fmt.Printf("C2S %#v\n", event)
	err = p.ProxyClient.Send(event)
	return
}

func (p *Proxy) OnHandshake(packet v1_21_8.HandshakingToServerPacketSetProtocol) (err error) {
	err = p.ProxyClient.Send(packet)
	if err != nil {
		return
	}
	if packet.NextState != queser.VarInt(queser.Status) && packet.NextState != queser.VarInt(queser.Login) {
		err = errors.New("invalid next state")
	}
	p.State = queser.State(packet.NextState)
	p.ProxyClient.State = queser.State(packet.NextState)
	return
}

func (p *Proxy) OnLoginAcknowledged(packet v1_21_8.LoginToServerPacketLoginAcknowledged) (err error) {
	err = p.ProxyClient.Send(packet)
	if err != nil {
		return
	}
	p.State = queser.Configuration
	p.ProxyClient.State = queser.Configuration
	return
}

func (p *Proxy) OnConfigFinish(packet v1_21_8.ConfigurationToServerPacketFinishConfiguration) (err error) {
	err = p.ProxyClient.Send(packet)
	if err != nil {
		return
	}
	p.State = queser.Play
	p.ProxyClient.State = queser.Play
	return
}

func (p *Proxy) OnStart(_ strom.OnStart) (err error) {
	return
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	s := strom.Servee(conn)
	p := &Proxy{Conn: s}
	c, err := strom.Connect("127.0.0.1:25566")
	if err != nil {
		return
	}
	p.ProxyClient = &ProxyClient{Conn: c, Proxy: p}
	errChan := make(chan error)
	go func() {
		errChan <- fmt.Errorf("client error: %v", c.Start(p.ProxyClient))
	}()
	go func() {
		errChan <- fmt.Errorf("server error: %v", s.Start(p))
	}()
	err = <-errChan
	if err != nil {
		log.Println(err)
	}
	return
}

func main() {
	l, err := net.Listen("tcp", ":25565")
	if err != nil {
		panic(err)
	}
	for {
		var c net.Conn
		c, err = l.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(c)
	}
}
