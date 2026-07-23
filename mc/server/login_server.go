package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"slices"

	api2 "github.com/admin-else/strom/mc/api"
	"github.com/admin-else/strom/mc/crypto"
	"github.com/admin-else/strom/mc/data"
	"github.com/admin-else/strom/mc/event"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_base"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_8"
	"github.com/admin-else/strom/mc/proto_generated/v26_2"
	"github.com/google/uuid"
)

var StatusServedErr = errors.New("status served")

var UnexpectedStatusRequest = UnexpectedNextStateError{proto_base.Status}

type NameAndUUID struct {
	Name string
	UUID uuid.UUID
}

type UnexpectedNextStateError struct {
	NextState proto_base.State
}

func (u UnexpectedNextStateError) Error() string {
	return fmt.Sprintf("unexpected next state: %v", u.NextState)
}

type LoginServer struct {
	*proto.Conn
	Encryption, OnlineMode bool
	ServerHost             string
	ServerPort             uint16
	Requested              NameAndUUID
	Given                  *NameAndUUID
	CompressionThreshold   int32
	Status                 []byte
	CompatibleVersions     []int32

	ServerId       string // always empty for now
	VerifyToken    []byte
	PrivateKey     *rsa.PrivateKey
	PublicKeyBytes []byte
}

func (l *LoginServer) OnHandshake(packet *v1_21_8.HandshakingToServerPacketSetProtocol) (err error) {
	l.SetState(proto_base.State(packet.NextState))
	l.ServerHost = packet.ServerHost
	l.ServerPort = packet.ServerPort

	if packet.NextState == int32(proto_base.Status) && l.Status != nil {
		return
	}

	if packet.NextState != int32(proto_base.Login) {
		err = UnexpectedNextStateError{proto_base.State(packet.NextState)}
		return
	}

	if len(l.CompatibleVersions) != 0 && !slices.Contains(l.CompatibleVersions, packet.ProtocolVersion) {
		err = fmt.Errorf("incompatible protocol version please use one of %v you can look up the real corresponding version at https://minecraft.wiki/w/Minecraft_Wiki:Projects/wiki.vg_merge/Protocol_version_numbers", l.CompatibleVersions)
		return
	}

	err = l.SetProtocolVersion(packet.ProtocolVersion)
	return
}

func (l *LoginServer) OnStatusRequest(_ *v1_21_8.StatusToServerPacketPingStart) (err error) {
	err = l.Send(&v1_21_8.StatusToClientPacketServerInfo{Response: string(l.Status)})
	return
}

func (l *LoginServer) OnStatusPing(packet *v1_21_8.StatusToServerPacketPing) (err error) {
	err = l.Send(&v1_21_8.StatusToClientPacketPing{Time: packet.Time})
	if err != nil {
		return
	}
	return event.HandlerDoneErr{Return: StatusServedErr}
}

func (l *LoginServer) SetCompressionThreshold(threshold int32) (err error) {
	err = l.Send(&v1_21_8.LoginToClientPacketCompress{Threshold: threshold})
	if err != nil {
		return
	}
	l.RawConn.SetCompressionThreshold(threshold)
	return
}

func (l *LoginServer) OnLoginStart(packet *v1_21_8.LoginToServerPacketLoginStart) (err error) {
	if l.CompressionThreshold < 0 {
		err = l.SetCompressionThreshold(l.CompressionThreshold)
		if err != nil {
			return
		}
	}

	l.Requested = NameAndUUID{packet.Username, packet.PlayerUUID}
	if l.Encryption {
		return l.Encrypt()
	}
	return l.FinishLogin()
}

func (l *LoginServer) OnLoginAcknowledged(_ *v1_21_8.LoginToServerPacketLoginAcknowledged) (err error) {
	l.SetState(proto_base.Configuration)
	err = event.HandlerDoneErr{}
	return
}

func (l *LoginServer) OnDefault(event event.Unhandled) (err error) {
	if _, ok := event.Val.(proto_base.EncodeDecodeAble); !ok {
		return
	}
	l.Log.Warn("unexpected packet during login serving", "packet", event.Val)
	return
}

func (l *LoginServer) Encrypt() (err error) {
	l.PrivateKey, err = rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return
	}
	l.VerifyToken = make([]byte, 4)
	_, _ = rand.Read(l.VerifyToken)
	l.PublicKeyBytes, err = x509.MarshalPKIXPublicKey(l.PrivateKey.Public())
	if err != nil {
		return
	}

	return l.Send(&v1_21_8.LoginToClientPacketEncryptionBegin{
		ServerId:           "",
		PublicKey:          l.PublicKeyBytes,
		VerifyToken:        l.VerifyToken,
		ShouldAuthenticate: l.OnlineMode,
	})
}

var ErrEncryptionFailedVerifyToken = errors.New("clients verify token does not match server verify token")

func (l *LoginServer) OnEncryptionResponse(packet *v1_21_8.LoginToServerPacketEncryptionBegin) (err error) {
	plain, err := l.PrivateKey.Decrypt(rand.Reader, packet.VerifyToken, &rsa.PKCS1v15DecryptOptions{}) // Maybe supply 1024????
	if err != nil {
		return
	}
	if !slices.Equal(plain, l.VerifyToken) {
		err = ErrEncryptionFailedVerifyToken
		return
	}
	sharedSecret, err := l.PrivateKey.Decrypt(rand.Reader, packet.SharedSecret, &rsa.PKCS1v15DecryptOptions{}) // Maybe supply 1024????
	if err != nil {
		return
	}

	err = l.SetSecret(sharedSecret)
	if err != nil {
		return
	}

	serverId := crypto.AuthDigest([]byte(l.ServerId), sharedSecret, l.PublicKeyBytes)
	if l.OnlineMode {
		_, err = api2.HasJoined(l.Requested.Name, serverId, "")
		if err != nil {
			return
		}
	}
	err = l.FinishLogin()
	return
}

func (l *LoginServer) FinishLogin() (err error) {
	if l.Given == nil {
		l.Given = &l.Requested
	}
	if l.ProtocolVersion >= data.MustLookupProtocolVersion("26.2") {
		err = l.Send(&v26_2.LoginToClientPacketSuccess{
			Uuid:       l.Given.UUID,
			Username:   l.Given.Name,
			Properties: nil,
			SessionId:  uuid.New(),
		})
	} else {
		err = l.Send(&v1_21_8.LoginToClientPacketSuccess{
			Uuid:       l.Given.UUID,
			Username:   l.Given.Name,
			Properties: nil,
		})
	}
	return
}

type LoginServerSetting func(*LoginServer)

func WithOtherAccount(a *api2.Account) LoginServerSetting {
	return func(s *LoginServer) {
		s.Given = &NameAndUUID{
			Name: a.Name,
			UUID: a.Uuid,
		}
	}
}

func WithRawStatus(status []byte) LoginServerSetting {
	return func(s *LoginServer) {
		s.Status = status
	}
}

func WithStatus(response StatusResponse) LoginServerSetting {
	var status []byte
	status, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	slog.Debug("With status", "status", status)
	return WithRawStatus(status)
}

func WithCompatibleVersionsRange(start, end int32) LoginServerSetting {
	var versions []int32
	for i := start; i <= end; i++ {
		versions = append(versions, i)
	}
	return WithCompatibleVersions(versions...)
}

func WithCompatibleVersions(versions ...int32) LoginServerSetting {
	return func(s *LoginServer) {
		s.CompatibleVersions = versions
	}
}

func WithoutOnlineMode() LoginServerSetting {
	return func(s *LoginServer) {
		s.OnlineMode = false
	}
}

func WithoutEncryption() LoginServerSetting {
	return func(s *LoginServer) {
		s.Encryption = false
	}
}

func ServeLogin(c *proto.Conn, settings ...LoginServerSetting) (ret *LoginServer, err error) {
	ret = &LoginServer{Conn: c}
	ret.OnlineMode = true
	ret.Encryption = true

	ret.RegisterCritical(ret.OnDefault)
	ret.RegisterCriticalUntilLatest(ret.OnHandshake)
	ret.RegisterCriticalUntilLatest(ret.OnLoginStart)
	ret.RegisterCriticalUntilLatest(ret.OnLoginAcknowledged)
	ret.RegisterCriticalUntilLatest(ret.OnStatusRequest)
	ret.RegisterCriticalUntilLatest(ret.OnStatusPing)
	ret.RegisterCriticalUntilLatest(ret.OnEncryptionResponse)
	for _, s := range settings {
		s(ret)
	}
	if ret.OnlineMode && !ret.Encryption {
		err = fmt.Errorf("online mode requires encryption")
		return
	}

	ret.CompressionThreshold = 256
	err = ret.StartConn()
	return
}
