package client

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"errors"

	"github.com/admin-else/strom/api"
	"github.com/admin-else/strom/crypto"
	"github.com/admin-else/strom/event"
	"github.com/admin-else/strom/proto"
	"github.com/admin-else/strom/proto_base"
	"github.com/admin-else/strom/text"

	"github.com/admin-else/strom/proto_generated/v1_21_8"
)

type LoginClient struct {
	*proto.Conn
	Account    *api.Account
	ServerHost string
	ServerPort uint16

	GivenAccount *v1_21_8.LoginToClientPacketSuccess
}

var (
	NoTokenErr          = errors.New("the account has no token so we cant join online servers")
	BadPublicKeyTypeErr = errors.New("the public key is not rsa")
	FailedToLoginErr    = errors.New("failed to login")
)

type KickedDuringLoginErr struct {
	text.Component
}

func (k KickedDuringLoginErr) Error() string {
	return k.String()
}

func (s *LoginClient) OnDefault(event event.Unhandled) (err error) {
	if _, ok := event.Val.(proto_base.EncodeDecodeAble); !ok {
		return
	}
	s.Log.Warn("unexpected packet during login client", "packet", event.Val)
	return
}

func (s *LoginClient) OnClose(_ event.Close) (err error) {
	return
}

func (s *LoginClient) OnCompress(compress *v1_21_8.LoginToClientPacketCompress) (err error) {
	s.SetCompressionThreshold(compress.Threshold)
	return
}

func (s *LoginClient) OnEncrypt(packet *v1_21_8.LoginToClientPacketEncryptionBegin) (err error) {
	sharedSecret := make([]byte, 16)
	_, _ = rand.Read(sharedSecret) //never fails

	if packet.ShouldAuthenticate {
		if s.Account.Ygg == "" {
			err = NoTokenErr
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
		err = BadPublicKeyTypeErr
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
	err = s.Send(&v1_21_8.LoginToServerPacketEncryptionBegin{SharedSecret: sharedSecretEnc, VerifyToken: verifyTokenEnc})
	if err != nil {
		return
	}
	err = s.SetSecret(sharedSecret)
	return

}

func (s *LoginClient) OnDisconnect(packet *v1_21_8.LoginToClientPacketDisconnect) (err error) {
	var reason text.Component
	err = json.Unmarshal([]byte(packet.Reason), &reason)
	if err != nil {
		return
	}
	err = KickedDuringLoginErr{reason}
	return
}

func (s *LoginClient) OnSuccess(success *v1_21_8.LoginToClientPacketSuccess) (err error) {
	s.GivenAccount = success
	err = s.Send(&v1_21_8.LoginToServerPacketLoginAcknowledged{})
	if err != nil {
		return
	}
	s.SetState(proto_base.Configuration)
	err = event.HandlerDoneErr{}
	return
}

func (s *LoginClient) OnStart() (err error) {
	return
}

func LoginRaw(c *proto.Conn, account *api.Account, addr string) (err error) {
	lc := &LoginClient{
		Conn:    c,
		Account: account,
	}
	lc.RegisterCritical(lc.OnDefault)
	lc.RegisterCritical(lc.OnClose)
	lc.RegisterCriticalUntilLatest(lc.OnCompress)
	lc.RegisterCriticalUntilLatest(lc.OnEncrypt)
	lc.RegisterCriticalUntilLatest(lc.OnDisconnect)
	lc.RegisterCriticalUntilLatest(lc.OnSuccess)

	p, err := MakeHandshakePacket(lc.Conn, proto_base.Login)
	if err != nil {
		return
	}

	err = lc.Send(p)
	if err != nil {
		return
	}
	lc.SetState(proto_base.Login)
	err = lc.Send(&v1_21_8.LoginToServerPacketLoginStart{Username: lc.Account.Name, PlayerUUID: lc.Account.Uuid})

	err = lc.StartConn()
	if err != nil {
		err = errors.Join(FailedToLoginErr, err)
	}
	return
}

func LoginNoSrv(connectTo string, account *api.Account) (c *proto.Conn, err error) {
	c, err = Connect(connectTo)
	if err != nil {
		return
	}
	err = LoginRaw(c, account, connectTo)
	return
}

func Login(connectTo string, account *api.Account) (c *proto.Conn, err error) {
	connectTo, err = DoDnsSimple(connectTo)
	if err != nil {
		return
	}
	c, err = Connect(connectTo)
	if err != nil {
		return
	}
	err = LoginRaw(c, account, connectTo)
	return
}
