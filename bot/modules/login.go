package modules

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/admin-else/strom/api"
	"github.com/admin-else/strom/bot"
	"github.com/admin-else/strom/crypto"
	"github.com/admin-else/strom/proto_base"

	"github.com/admin-else/strom/data"
	"github.com/admin-else/strom/proto_generated/v1_21_8"
)

type LoginClient struct {
	*bot.Conn
	Account api.Account

	GivenAccount v1_21_8.LoginToClientPacketSuccess
}

func (s *LoginClient) Default(event any) (err error) {
	err = fmt.Errorf("unexpected event: %#v", event)
	return
}

func (s *LoginClient) OnCycle(_ bot.OnLoopCycle) (err error) {
	return
}

func (s *LoginClient) OnStart(_ bot.OnStart) (err error) {
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
		ProtocolVersion: proto_base.VarInt(versionData.Version),
		ServerHost:      host,
		ServerPort:      uint16(port),
		NextState:       proto_base.VarInt(proto_base.Login),
	})
	if err != nil {
		return
	}
	s.State = proto_base.Login
	err = s.Send(v1_21_8.LoginToServerPacketLoginStart{Username: s.Account.Name, PlayerUUID: s.Account.Uuid})
	return
}

func (s *LoginClient) OnCompress(compress v1_21_8.LoginToClientPacketCompress) (err error) {
	s.CompressionThreshold = int32(compress.Threshold)
	return
}

func (s *LoginClient) OnEncrypt(packet v1_21_8.LoginToClientPacketEncryptionBegin) (err error) {
	sharedSecret := make([]byte, 16)
	_, _ = rand.Read(sharedSecret) //never fails

	if packet.ShouldAuthenticate {
		if s.Account.Ygg == "" {
			err = fmt.Errorf("the account has no token so we cant join online servers")
			return
		}
		serverId := crypto.AuthDigest([]byte(packet.ServerId), sharedSecret, packet.PublicKey)
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
		S: crypto.NewCFB8Decrypt(b, sharedSecret),
		R: s.Conn,
	}
	s.W = cipher.StreamWriter{
		S: crypto.NewCFB8Encrypt(b, sharedSecret),
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
	s.State = proto_base.Configuration
	err = bot.HandlerDone
	return
}

func Login(c *bot.Conn, account api.Account) (err error) {
	err = c.Start(&LoginClient{
		Conn:    c,
		Account: account,
	})
	if err != nil {
		err = errors.Join(err, errors.New("failed to login"))
	}
	return
}
