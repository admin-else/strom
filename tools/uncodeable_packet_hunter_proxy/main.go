package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"

	"github.com/admin-else/strom/api"
	"github.com/admin-else/strom/client/modules"
	"github.com/admin-else/strom/event"
	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/proto_generated"
	"github.com/admin-else/strom/proto_generated/v1_21_8"
	"github.com/admin-else/strom/server"
	"github.com/admin-else/strom/server/server_modules"
)

type ProxyClient struct {
	*proto.Conn
	Server *Proxy
}

func (p *ProxyClient) OnDefault(event event.Default) (err error) {
	err = p.Server.Send(event.Val)
	return
}

func (p *ProxyClient) OnCycle(_ event.OnLoopCycle) (err error) {
	return
}

func (p *ProxyClient) OnStart(_ event.OnStart) (err error) {
	return
}

func (p *ProxyClient) OnUnCodeAble(packet proto.UnCodablePacket) (err error) {
	err = p.OnDefault(event.Default{Val: packet})
	if err != nil {
		return
	}
	SaveUnCodeAbleAsTest(proto_base.ToServer, packet)
	return
}

type Proxy struct {
	*proto.Conn
	Client *ProxyClient
}

func (p *Proxy) Default(event any) (err error) {
	err = p.Client.Send(event)
	return
}

func (p *Proxy) OnCycle(_ event.OnLoopCycle) (err error) {
	return
}

func (p *Proxy) OnStart(_ event.OnStart) (err error) {
	return
}

func (p *Proxy) OnFinishConfiguration(packet v1_21_8.ConfigurationToServerPacketFinishConfiguration) (err error) {
	err = p.Default(packet)
	if err != nil {
		return
	}
	p.State = proto_base.Play
	p.Client.State = proto_base.Play
	return
}

func (p *Proxy) OnUnCodeAble(packet proto.UnCodablePacket) (err error) {
	err = p.Default(packet)
	if err != nil {
		return
	}
	SaveUnCodeAbleAsTest(proto_base.ToClient, packet)
	return
}

const TestSrcF = `package failed_packets_test

import (
	"bytes"
	"testing"

	"github.com/admin-else/strom/proto_generated/v1_21_8"
)

// Error: %v
func TestFailedPacket%10X(t *testing.T) {
	p := v1_21_8.Play%vPacket{}
	b := bytes.NewBuffer(%#v)
	p, err := p.Decode(b)
	if err != nil {
		t.Fatal(err)
	}
}
`

func SaveUnCodeAbleAsTest(d proto_base.Direction, packet proto.UnCodablePacket) {
	hUntrimmed := sha256.Sum256(packet.Data)
	h := hUntrimmed[:8]
	f, err := os.Create(fmt.Sprintf(".failed_packets/%v_%10x_test.go", len(packet.Data), h))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, TestSrcF, packet.Err, h, d.Opposite(), packet.Data)
	if err != nil {
		panic(err)
	}
	fmt.Println("Error", packet.Err, "saved to", f.Name(), "data:", packet.Data)
}

var StatusResponse = server_modules.StatusResponse{Description: struct {
	Text string `json:"text"`
}{Text: "Un-code-able packet hunter proxy"}, Version: struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}{Name: "STROM", Protocol: 772}}

func main() {
	err := server.StartServerWithOnConn(":25565", func(serveeConn *proto.Conn) (err error) {
		p := &Proxy{Conn: serveeConn, Client: nil}
		acc := api.NewOfflineAccount("sigma")
		_, err = server_modules.ServeLoginWithOtherAccount(serveeConn, acc)
		if err != nil {
			if errors.Is(err, server_modules.UnexpectedStatusRequest) {
				err = server_modules.ServeStatus(serveeConn, StatusResponse)
			}
			return
		}
		c, err := modules.ConnectAndLogin("127.0.0.1:25566", acc)
		if err != nil {
			return
		}
		pc := &ProxyClient{Conn: c, Server: p}
		p.Client = pc
		errChan := make(chan error)
		defer pc.Close()
		defer p.Close()
		go func() {
			errChan <- pc.StartOne(pc)
		}()
		go func() {
			errChan <- p.StartOne(p)
		}()
		err = <-errChan
		return
	})
	if err != nil {
		panic(err)
	}
}
