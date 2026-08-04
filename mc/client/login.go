package client

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"errors"
	"net"

	"github.com/admin-else/strom/mc/api"
	"github.com/admin-else/strom/mc/crypto"
	"github.com/admin-else/strom/mc/event"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_base"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_8"
	"github.com/admin-else/strom/mc/proto_generated/v1_8"
	"github.com/admin-else/strom/mc/proto_generated/v26_2"
	"github.com/admin-else/strom/mc/text"
	"github.com/google/uuid"
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

func (s *LoginClient) OnDisconnectV1_8(packet *v1_8.LoginToClientPacketDisconnect) (err error) {
	var reason text.Component
	err = json.Unmarshal([]byte(packet.Reason), &reason)
	if err != nil {
		return
	}
	err = KickedDuringLoginErr{reason}
	return
}

func (s *LoginClient) OnCompressV1_8(packet *v1_8.LoginToClientPacketCompress) (err error) {
	s.SetCompressionThreshold(packet.Threshold)
	return
}

func (s *LoginClient) OnSuccessV1_8(success *v1_8.LoginToClientPacketSuccess) (err error) {
	s.GivenAccount = &v1_21_8.LoginToClientPacketSuccess{
		Uuid:     uuid.MustParse(success.Uuid),
		Username: success.Username,
	}
	s.SetState(proto_base.Play)
	err = event.HandlerDoneErr{}
	return
}

func (s *LoginClient) OnSuccess(success *v1_21_8.LoginToClientPacketSuccess) (err error) {
	s.GivenAccount = success
	if s.ProtocolVersion >= 764 {
		err = s.Send(&v1_21_8.LoginToServerPacketLoginAcknowledged{})
		if err != nil {
			return
		}
		s.SetState(proto_base.Configuration)
	} else {
		s.SetState(proto_base.Play)
	}
	err = event.HandlerDoneErr{}
	return
}

func (s *LoginClient) OnSuccess26_2(success *v26_2.LoginToClientPacketSuccess) (err error) {
	return s.OnSuccess(&v1_21_8.LoginToClientPacketSuccess{
		Uuid:       success.Uuid,
		Username:   success.Username,
		Properties: success.Properties,
	})
}

func (s *LoginClient) OnStart() (err error) {
	return
}

// LoginRaw performs the login handshake on an already-established connection without DNS resolution or configuration handling.
// The handshake address is derived from the connection's RemoteAddr.
func LoginRaw(c *proto.Conn, account *api.Account) (err error) {
	return LoginRawAddr(c, account, "")
}

// LoginRawAddr performs the login handshake on an already-established connection using the given address for the handshake packet.
// If hostAddr is empty, the handshake address is derived from the connection's RemoteAddr.
func LoginRawAddr(c *proto.Conn, account *api.Account, hostAddr string) (err error) {
	lc := &LoginClient{
		Conn:    c,
		Account: account,
	}
	lc.RegisterCritical(lc.OnDefault)
	lc.RegisterCritical(lc.OnClose)
	lc.RegisterCriticalUntilLatest(lc.OnCompress)
	lc.RegisterCriticalUntilLatest(lc.OnEncrypt)
	lc.RegisterCriticalUntilLatest(lc.OnDisconnect)
	lc.RegisterCriticalUntil("26.1", lc.OnSuccess)
	lc.RegisterCriticalUntilLatest(lc.OnSuccess26_2)
	lc.RegisterCritical(lc.OnCompressV1_8)
	lc.RegisterCritical(lc.OnDisconnectV1_8)
	lc.RegisterCritical(lc.OnSuccessV1_8)

	var p proto_base.EncodeDecodeAble
	if hostAddr != "" {
		p, err = MakeHandshakePacketAddr(lc.Conn, proto_base.Login, hostAddr)
	} else {
		p, err = MakeHandshakePacket(lc.Conn, proto_base.Login)
	}
	if err != nil {
		return
	}

	err = lc.Send(p)
	if err != nil {
		return
	}
	lc.SetState(proto_base.Login)
	if lc.Conn.ProtocolVersion < 759 {
		err = lc.Send(&v1_8.LoginToServerPacketLoginStart{Username: lc.Account.Name})
	} else {
		err = lc.Send(&v1_21_8.LoginToServerPacketLoginStart{Username: lc.Account.Name, PlayerUUID: lc.Account.Uuid})
	}

	err = lc.StartConn()
	if err != nil {
		err = errors.Join(FailedToLoginErr, err)
	}
	return
}

// WithoutDns returns a login setting that disables SRV DNS resolution during login.
func WithoutDns() func(ls *loginSettings) {
	return func(ls *loginSettings) {
		ls.NoDns = true
	}
}

// WithContext returns a login setting that sets the context used for the connection dial and DNS resolution.
func WithContext(ctx context.Context) func(ls *loginSettings) {
	return func(ls *loginSettings) {
		ls.Ctx = ctx
	}
}

// WithVersion returns a login setting that forces a specific Minecraft version string for the connection.
func WithVersion(version string) func(ls *loginSettings) {
	return func(ls *loginSettings) {
		ls.Version = version
	}
}

func WithIgnoreConfig() func(ls *loginSettings) {
	return func(ls *loginSettings) {
		ls.IgnoreConfig = true
	}
}

// DialFunc is a context-aware dial function used by Login to establish a TCP connection.
type DialFunc func(ctx context.Context, network, addr string) (net.Conn, error)

// WithConn returns a login setting that uses an already-established net.Conn instead of dialling.
// WithVersion must also be set; auto-detection is not possible with a pre-established connection.
func WithConn(conn net.Conn) func(ls *loginSettings) {
	return func(ls *loginSettings) {
		ls.Conn = conn
	}
}

// WithDialer returns a login setting that uses the given dial function instead of net.Dialer.
// Auto version detection via status ping is still supported when no explicit version is set.
func WithDialer(d DialFunc) func(ls *loginSettings) {
	return func(ls *loginSettings) {
		ls.Dialer = d
	}
}

var (
	ErrWithConnRequiresVersion = errors.New("WithVersion is required when using WithConn")
	ErrWithConnNoStatus        = errors.New("auto version detection is not supported with WithConn")
)

type loginSettings struct {
	NoDns        bool
	Ctx          context.Context
	Version      string
	IgnoreConfig bool
	Conn         net.Conn
	Dialer       DialFunc
}

// Login resolves the address with SRV DNS, connects, performs the login handshake, and handles the configuration phase, returning an authenticated connection ready for play.
func Login(connectTo string, account *api.Account, settings ...func(ls *loginSettings)) (c *proto.Conn, err error) {
	ls := &loginSettings{}

	for _, setting := range settings {
		setting(ls)
	}

	if ls.Conn != nil {
		if ls.Version == "" {
			return nil, ErrWithConnRequiresVersion
		}
		c, err = ConnectWithConn(ls.Ctx, ls.Conn, ls.Version)
		if err != nil {
			return
		}
	} else if ls.Dialer != nil {
		if ls.Version == "" {
			return nil, ErrWithConnRequiresVersion
		}
		if !ls.NoDns {
			connectTo, _, _, err = DoDNSChecked(ls.Ctx, connectTo, nil)
			if err != nil {
				return
			}
		}
		var conn net.Conn
		conn, err = ls.Dialer(ls.Ctx, "tcp", connectTo)
		if err != nil {
			return
		}
		c, err = ConnectWithConn(ls.Ctx, conn, ls.Version)
		if err != nil {
			conn.Close()
			return
		}
	} else {
		if !ls.NoDns {
			connectTo, _, _, err = DoDNSChecked(ls.Ctx, connectTo, nil)
			if err != nil {
				return
			}
		}
		c, err = Connect(ls.Ctx, connectTo, ls.Version)
		if err != nil {
			return
		}
	}

	err = LoginRawAddr(c, account, connectTo)
	if err != nil {
		return
	}
	if !ls.IgnoreConfig && c.ProtocolVersion >= 764 {
		err = IgnoreConfig(c)
	}
	return
}
