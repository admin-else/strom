package api

import (
	"git.anygate.cloud/anygatecloud/strom/mc/crypto"
)

// NewOfflineAccount creates an offline-mode Account with a UUID derived from the
// player name. No HTTP calls are made; the UUID is computed locally.
func NewOfflineAccount(name string) *Account {
	return &Account{
		Name: name,
		Uuid: crypto.FromOfflinePlayer(name),
	}
}
