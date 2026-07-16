package chat

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/admin-else/strom/mc/api"
	"github.com/admin-else/strom/mc/proto"
	"github.com/admin-else/strom/mc/proto_generated/v1_21_11"
	"github.com/google/uuid"
)

type Module struct {
	*proto.Conn

	Account *api.Account

	// Chat signing state
	canSign     bool
	privKey     *rsa.PrivateKey
	sessionUUID uuid.UUID
	msgIndex    int32
}

// Start creates a chat module for the connection. If acc is an online account,
// it fetches the profile key pair from Mojang and enables signed chat.
func Start(c *proto.Conn, acc *api.Account) (*Module, error) {
	m := &Module{Conn: c, Account: acc}
	if acc != nil && acc.Ygg != "" {
		if err := m.initSigning(acc); err != nil {
			return nil, err
		}
	}
	return m, nil
}

func (m *Module) initSigning(acc *api.Account) (err error) {
	keys, err := acc.FetchKeys()
	if err != nil {
		return fmt.Errorf("fetch profile keys: %w", err)
	}

	priv, err := parsePrivateKey(keys.KeyPair.PrivateKey)
	if err != nil {
		return fmt.Errorf("parse private key: %w", err)
	}

	m.privKey = priv
	m.sessionUUID = uuid.New()

	pubDer, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return fmt.Errorf("marshal public key: %w", err)
	}

	sigBytes, err := base64.StdEncoding.DecodeString(keys.PublicKeySignatureV2)
	if err != nil {
		return fmt.Errorf("decode key signature: %w", err)
	}

	if err = m.Conn.Send(&v1_21_11.PlayToServerPacketChatSessionUpdate{
		SessionUUID: m.sessionUUID,
		ExpireTime:  keys.ExpiresAt.UnixMilli(),
		PublicKey:   v1_21_11.ByteArray{Val: pubDer},
		Signature:   v1_21_11.ByteArray{Val: sigBytes},
	}); err != nil {
		return fmt.Errorf("send chat session update: %w", err)
	}

	m.canSign = true
	return nil
}

// SendMessage sends a chat message. If signing is enabled, the message is signed;
// otherwise it is sent unsigned.
func (m *Module) SendMessage(message string) (err error) {
	if m.canSign {
		return m.sendSigned(message)
	}
	return m.Conn.Send(&v1_21_11.PlayToServerPacketChatMessage{Message: message})
}

func (m *Module) sendSigned(message string) (err error) {
	now := time.Now()

	saltBytes := make([]byte, 8)
	if _, err = rand.Read(saltBytes); err != nil {
		return
	}
	salt := int64(binary.BigEndian.Uint64(saltBytes))

	timestamp := now.UnixMilli()

	payload, err := m.buildSignaturePayload(message, salt, timestamp)
	if err != nil {
		return
	}

	hashed := sha256.Sum256(payload)
	sigBytes, err := rsa.SignPKCS1v15(rand.Reader, m.privKey, crypto.SHA256, hashed[:])
	if err != nil {
		return
	}
	if len(sigBytes) != 256 {
		return fmt.Errorf("unexpected signature length: %d", len(sigBytes))
	}
	var signature [256]byte
	copy(signature[:], sigBytes)

	m.msgIndex++

	return m.Conn.Send(&v1_21_11.PlayToServerPacketChatMessage{
		Message:      message,
		Timestamp:    timestamp,
		Salt:         salt,
		Signature:    &signature,
		Offset:       0,
		Acknowledged: [3]byte{},
		Checksum:     0,
	})
}

func (m *Module) buildSignaturePayload(content string, salt, timestampMillis int64) ([]byte, error) {
	buf := make([]byte, 0, 64+len(content))
	w := func(b []byte) { buf = append(buf, b...) }

	// Header
	w(int32Bytes(1))

	// SignedMessageLink: sender UUID, session UUID, index
	w(uuidBytes(m.Account.Uuid))
	w(uuidBytes(m.sessionUUID))
	w(int32Bytes(m.msgIndex))

	// SignedMessageBody: salt, timestamp seconds, content
	w(int64Bytes(uint64(salt)))
	w(int64Bytes(uint64(timestampMillis / 1000)))
	contentBytes := []byte(content)
	w(int32Bytes(int32(len(contentBytes))))
	w(contentBytes)

	// LastSeenMessages: empty
	w(int32Bytes(0))

	return buf, nil
}

func parsePrivateKey(s string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(s))
	if block != nil {
		k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			k, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		}
		if err != nil {
			return nil, err
		}
		if key, ok := k.(*rsa.PrivateKey); ok {
			return key, nil
		}
		return nil, fmt.Errorf("not an RSA private key")
	}

	// Some implementations return raw base64 DER instead of PEM
	raw, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	k, err := x509.ParsePKCS8PrivateKey(raw)
	if err != nil {
		k, err = x509.ParsePKCS1PrivateKey(raw)
	}
	if err != nil {
		return nil, err
	}
	key, ok := k.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA private key")
	}
	return key, nil
}

func int32Bytes(v int32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(v))
	return b
}

func int64Bytes(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func uuidBytes(u uuid.UUID) []byte {
	b := make([]byte, 16)
	copy(b, u[:])
	return b
}
