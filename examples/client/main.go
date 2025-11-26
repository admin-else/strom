package main

import (
	"bytes"
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"strom"
	"strom/modules"
	"sync"
	"time"

	"github.com/admin-else/queser/generated/v1_21_8"
)

type Client struct {
	*strom.Conn
	Ticker     *modules.Ticker
	Message    string
	AltAccount strom.Account
}

func (c *Client) Default(event any) (err error) {
	return
}

func (c *Client) SendUnReportableChat(message string) (err error) {
	return c.Send(v1_21_8.PlayToServerPacketChatMessage{Message: message, Timestamp: time.Now().Unix(), Salt: rand.Int64(), Signature: nil, Offset: 0, Acknowledged: [3]uint8{0x0, 0x0, 0x0}, Checksum: 0x1})
}

func (c *Client) OnStart(_ strom.OnStart) (err error) {
	err = modules.Login(c.Conn, c.AltAccount)
	if err != nil {
		return
	}
	err = modules.IgnoreConfig(c.Conn)
	if err != nil {
		return
	}
	c.Ticker.Active = true
	return
}

func (c *Client) OnLoopCycle(strom.OnLoopCycle) (err error) {
	err = c.Ticker.Tick()
	return
}

func (c *Client) KeepAlive(packet v1_21_8.PlayToClientPacketKeepAlive) (err error) {
	return c.Send(v1_21_8.PlayToServerPacketKeepAlive(packet))
}

func (c *Client) Spam() (err error) {
	err = c.SendUnReportableChat(c.Message)
	return
}

func spawnSpammer(message, connectTo string, account strom.Account, wg *sync.WaitGroup) {
	conn, err := strom.Connect(connectTo)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := &Client{Conn: conn, Message: message}
	t := modules.Ticker{}
	t.Interval = append(t.Interval, modules.IntervalTask{
		F:        c.Spam,
		Interval: time.Second,
	})
	c.Ticker = &t
	c.AltAccount = account
	err = conn.Start(c)
	wg.Done()
}

func DownloadAccounts(token string) (accounts []strom.Account, err error) {
	r, err := http.NewRequest("GET", "https://griefing.homes/api/list/active", nil)
	if err != nil {
		return
	}
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b := bytes.NewBuffer(nil)
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(b.Bytes(), &accounts)
	return
}

func main() {
	accounts, err := DownloadAccounts("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIzNDBjNWIzMS02MmEzLTRjODUtYTE5ZC1lZjcxMWVkYzY3ZWQifQ.M_gwEK2HyyY6d4RFY7JzyeuqGVvmqmPK-13wNP0AxGc")
	if err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(accounts))
	for _, account := range accounts {
		go spawnSpammer("LOBA ONTOP XD", "127.0.0.1:25566", account, wg)
		time.Sleep(time.Second * 6)
	}
	wg.Wait()
}
