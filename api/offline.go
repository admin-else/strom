package api

import "github.com/admin-else/strom/crypto"

func NewOfflineAccount(name string) *Account {
	return &Account{
		Name: name,
		Uuid: crypto.FromOfflinePlayer(name),
	}
}
