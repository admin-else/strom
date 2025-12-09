package modules

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/admin-else/strom"

	"github.com/admin-else/queser"
	"github.com/admin-else/queser/data"
	"github.com/admin-else/queser/generated/v1_21_8"
)

type LoginClient struct {
	*strom.Conn
	Account strom.Account

	GivenAccount v1_21_8.LoginToClientPacketSuccess
}

func (s *LoginClient) Default(event any) (err error) {
	err = fmt.Errorf("unexpected event: %#v", event)
	return
}

func (s *LoginClient) OnCycle(_ strom.OnLoopCycle) (err error) {
	return
}

func (s *LoginClient) OnStart(_ strom.OnStart) (err error) {
	host, portStr, err := net.SplitHostPort(s.RemoteAddr().String())
	if err != nil {
		return
	}
	versionData, err := data.LookUpProtocolVersionByName(s.Version)
	if err != nil {
		return
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return
	}
	err = s.Send(v1_21_8.HandshakingToServerPacketSetProtocol{
		ProtocolVersion: queser.VarInt(versionData.Version),
		ServerHost:      host,
		ServerPort:      uint16(port),
		NextState:       queser.VarInt(queser.Login),
	})
	if err != nil {
		return
	}
	s.State = queser.Login
	err = s.Send(v1_21_8.LoginToServerPacketLoginStart{Username: s.Account.Name, PlayerUUID: s.Account.Uuid})
	return
}

func (s *LoginClient) OnCompress(compress v1_21_8.LoginToClientPacketCompress) (err error) {
	s.CompressionThreshold = int32(compress.Threshold)
	return
}

func twosComplement(p []byte) []byte {
	carry := true
	for i := len(p) - 1; i >= 0; i-- {
		p[i] = byte(^p[i])
		if carry {
			carry = p[i] == 0xff
			p[i]++
		}
	}
	return p
}

// AuthDigest stolen from https://gist.github.com/toqueteos/5372776
func AuthDigest(elems ...[]byte) string {
	h := sha1.New()
	for _, elem := range elems {
		h.Write(elem)
	}
	hash := h.Sum(nil)

	negative := (hash[0] & 0x80) == 0x80
	if negative {
		hash = twosComplement(hash)
	}

	res := strings.TrimLeft(hex.EncodeToString(hash), "0")
	if negative {
		res = "-" + res
	}

	return res
}

func (s *LoginClient) OnEncrypt(packet v1_21_8.LoginToClientPacketEncryptionBegin) (err error) {
	sharedSecret := make([]byte, 16)
	_, _ = rand.Read(sharedSecret) //never fails

	if packet.ShouldAuthenticate {
		if s.Account.Ygg == "" {
			err = fmt.Errorf("the account has no token so we cant join online servers")
			return
		}
		serverId := AuthDigest([]byte(packet.ServerId), sharedSecret, packet.PublicKey)
		err = s.Account.JoinServer(serverId)
		if err != nil {
			return
		}
	}

	pubAny, err := x509.ParsePKIXPublicKey(packet.PublicKey)
	if err != nil {
		return
	}
	pub, ok := pubAny.(*rsa.PublicKey)
	if !ok {
		err = fmt.Errorf("public key is not rsa")
		return
	}
	verifyTokenEnc, err := rsa.EncryptPKCS1v15(rand.Reader, pub, packet.VerifyToken)
	if err != nil {
		return
	}

	sharedSecretEnc, err := rsa.EncryptPKCS1v15(rand.Reader, pub, sharedSecret)
	if err != nil {
		return
	}
	err = s.Send(v1_21_8.LoginToServerPacketEncryptionBegin{SharedSecret: sharedSecretEnc, VerifyToken: verifyTokenEnc})
	if err != nil {
		return
	}

	var b cipher.Block
	b, err = aes.NewCipher(sharedSecret)
	if err != nil {
		return
	}
	s.R = cipher.StreamReader{
		S: strom.NewCFB8Decrypt(b, sharedSecret),
		R: s.Conn,
	}
	s.W = cipher.StreamWriter{
		S: strom.NewCFB8Encrypt(b, sharedSecret),
		W: s.Conn,
	}
	return

}

func (s *LoginClient) OnSuccess(success v1_21_8.LoginToClientPacketSuccess) (err error) {
	s.GivenAccount = success
	err = s.Send(v1_21_8.LoginToServerPacketLoginAcknowledged{})
	if err != nil {
		return
	}
	s.State = queser.Configuration
	err = strom.HandlerDone
	return
}

func Login(c *strom.Conn, account strom.Account) (err error) {
	err = c.Start(&LoginClient{
		Conn:    c,
		Account: account,
	})
	if err != nil {
		err = errors.Join(err, errors.New("failed to login"))
	}
	return
}
