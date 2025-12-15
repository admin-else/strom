package bot

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/admin-else/strom/api"
	"github.com/google/uuid"
)

type SecureChatSession struct {
	Account    *api.Account
	PlayerKeys api.PlayerKeys
	Private    *rsa.PrivateKey
	SessionId  uuid.UUID
}

func NewSecureChat(a *api.Account) (s *SecureChatSession, err error) {
	s = &SecureChatSession{Account: a}
	s.PlayerKeys, err = a.FetchKeys()

	if err != nil {
		return
	}
	block, _ := pem.Decode([]byte(s.PlayerKeys.KeyPair.PrivateKey))

	k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return
	}

	s.Private = k.(*rsa.PrivateKey)

	s.SessionId, err = uuid.NewRandom()
	return
}
