package main

import (
	"bytes"
	"crypto"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strom"
	"strom/modules"
	"time"

	"github.com/admin-else/queser/generated/v1_21_8"
	"github.com/google/uuid"
)

type Client struct {
	*strom.Conn
	Ticker     *modules.Ticker
	Message    string
	AltAccount strom.Account

	SessionId  uuid.UUID
	PrivateKey *rsa.PrivateKey
}

func (c *Client) Default(event any) (err error) {
	return
}

func (c *Client) MakeChatSession() (err error) {
	c.SessionId, err = uuid.NewRandom()
	if err != nil {
		return
	}
	keys, err := c.AltAccount.FetchKeys()
	if err != nil {
		return
	}
	block, _ := pem.Decode([]byte(keys.KeyPair.PrivateKey))
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return
	}
	c.PrivateKey = privateKey.(*rsa.PrivateKey)

	x509encoded, _ := pem.Decode([]byte(keys.KeyPair.PublicKey))
	sigBytes, err := base64.StdEncoding.DecodeString(keys.PublicKeySignatureV2)
	if err != nil {
		return
	}

	packet := v1_21_8.PlayToServerPacketChatSessionUpdate{
		SessionUUID: c.SessionId,
		ExpireTime:  keys.ExpiresAt.UnixMilli(),
		PublicKey:   v1_21_8.ByteArray{Val: x509encoded.Bytes},
		Signature:   v1_21_8.ByteArray{Val: sigBytes},
	}
	err = c.Send(packet)
	return
}

func (c *Client) MakeMessageSum(salt, timestamp int64, message string) (sum []byte, err error) {
	h := sha256.New()
	err = binary.Write(h, binary.BigEndian, int32(1))
	if err != nil {
		return
	}
	_, err = h.Write(c.AltAccount.Uuid[:])
	if err != nil {
		return
	}
	_, err = h.Write(c.SessionId[:])
	if err != nil {
		return
	}
	err = binary.Write(h, binary.BigEndian, int32(0)) // message index unimplemented
	if err != nil {
		return
	}
	err = binary.Write(h, binary.BigEndian, salt)
	if err != nil {
		return
	}
	err = binary.Write(h, binary.BigEndian, timestamp)
	if err != nil {
		return
	}
	err = binary.Write(h, binary.BigEndian, int32(len(message)))
	if err != nil {
		return
	}
	_, err = h.Write([]byte(message))
	if err != nil {
		return
	}
	err = binary.Write(h, binary.BigEndian, int32(0)) // number of last seen messages unimplemented
	if err != nil {
		return
	}
	sum = h.Sum(nil)
	return
}

func (c *Client) SendReportableChat(message string) (err error) {
	salt := rand.Int64()
	timestamp := time.Now().UnixMilli()
	sum, err := c.MakeMessageSum(salt, timestamp, message)
	if err != nil {
		return
	}
	signature, err := rsa.SignPSS(cryptorand.Reader, c.PrivateKey, crypto.SHA256, sum, nil)
	return c.Send(v1_21_8.PlayToServerPacketChatMessage{Message: message, Timestamp: timestamp, Salt: salt, Signature: (*[256]byte)(signature), Offset: 0, Acknowledged: [3]uint8{0x0, 0x0, 0x0}, Checksum: 0x1})
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
	err = c.MakeChatSession()
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
	err = c.SendReportableChat(c.Message)
	return
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
	conn, err := strom.Connect("127.0.0.1:25566")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := &Client{Conn: conn, Message: "LOBA ONTOP XD"}
	t := modules.Ticker{}
	t.Interval = append(t.Interval, modules.IntervalTask{
		F:        c.Spam,
		Interval: time.Second,
	})
	c.Ticker = &t
	c.AltAccount = strom.Account{
		Name: "blossom_varJ5_11",
		Uuid: uuid.MustParse("28ea954b7a9144d192081d686df32e56"),
		Ygg:  "eyJraWQiOiIwNDkxODEiLCJhbGciOiJSUzI1NiJ9.eyJ4dWlkIjoiMjUzNTQ2NjA5NDkzNTk1NyIsImFnZyI6IkFkdWx0Iiwic3ViIjoiODRmY2JhZTgtNmZlMS00ZDc5LWI4NTMtOGMzNGIwYzI5MGM2IiwiYXV0aCI6IlhCT1giLCJucyI6ImRlZmF1bHQiLCJyb2xlcyI6W10sImlzcyI6ImF1dGhlbnRpY2F0aW9uIiwiZmxhZ3MiOlsidHdvZmFjdG9yYXV0aCIsIm11bHRpcGxheWVyIiwib3JkZXJzXzIwMjIiLCJtc2FtaWdyYXRpb25fc3RhZ2U0Il0sInByb2ZpbGVzIjp7Im1jIjoiMjhlYTk1NGItN2E5MS00NGQxLTkyMDgtMWQ2ODZkZjMyZTU2In0sInBsYXRmb3JtIjoiUENfTEFVTkNIRVIiLCJwZmQiOlt7InR5cGUiOiJtYyIsImlkIjoiMjhlYTk1NGItN2E5MS00NGQxLTkyMDgtMWQ2ODZkZjMyZTU2IiwibmFtZSI6ImJsb3Nzb21fdmFySjVfMTEifV0sIm5iZiI6MTc2NDQ0ODE4NywiZXhwIjoxNzY0NTM0NTg3LCJpYXQiOjE3NjQ0NDgxODcsImFpZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDA0YzEyYWU2ZiJ9.EC2LOt06Q4L9gyhufn97y2bqqzqJfv_kSQCdbRUNNyfhTTKonAlgi9LFv-OKN2sCcA96azTEQzpAiOTvL5jQ4H6nFT8sFz6kTMxD7vmii7ztCsCRrshJeZmW7mPCEn8WNpZukXbxVcaOhyjrNwA47RXEm-GazuPheDpWi_fvDSnY1vsLLk4ouIoW_r3kwA9bRakCJRF6r-p0feyi-bJamjLq2yDp7SEyZRVLeBhNA6pxVVU0WtyP6mgVO4cg0XeGwlMpktcs4IiYaD2C7AFjvedCdy7WJQOImLIcsc-gV4kFOn6oFROEf3EGN0RDcsrrlOJalAM01-rZTve6tUKGZQ",
	}

	keys, err := c.AltAccount.FetchKeys()
	fmt.Println(keys)
	err = conn.Start(c)
	if err != nil {
		panic(err)
	}
}
