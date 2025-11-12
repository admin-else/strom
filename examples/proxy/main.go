package main

import (
	"fmt"
	"log"
	"net"
	"strom"

	"github.com/admin-else/queser"
	"github.com/admin-else/queser/generated/v1_21_8"
)

type ProxyClient struct {
	*strom.Conn
	Proxy *Proxy
}

func (p *ProxyClient) Default(event any) (err error) {
	fmt.Printf("1 %T%v\n", event, event)
	err = p.Proxy.Send(event)
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

func (p *Proxy) Default(event any) (err error) {
	fmt.Printf("2 %T%v\n", event, event)
	err = p.ProxyClient.Send(event)
	return
}

func (p *Proxy) OnHandshake(packet v1_21_8.HandshakingToServerPacketSetProtocol) (err error) {
	err = p.ProxyClient.Send(packet)
	if err != nil {
		return
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
		errChan <- c.Start(p.ProxyClient)
	}()
	go func() {
		errChan <- s.Start(p)
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
