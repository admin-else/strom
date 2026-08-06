package cmd_util

import "git.anygate.cloud/anygatecloud/strom/mc/api"

func Account(s string) (a *api.Account, err error) {
	if len(s) > 16 {
		return api.NewAccountFromYGG(s)
	}
	a = api.NewOfflineAccount(s)
	return
}
