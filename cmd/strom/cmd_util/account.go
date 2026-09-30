package cmd_util

import "github.com/admin-else/strom/mc/api"

func Account(s string) (a *api.Account, err error) {
	if len(s) > 16 {
		return api.NewAccountFromYGG(s)
	}
	a = api.NewOfflineAccount(s)
	return
}
