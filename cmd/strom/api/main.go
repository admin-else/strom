package api

import (
	"encoding/json"
	"errors"
	"flag"
	"log/slog"
	"os"

	"git.anygate.cloud/anygatecloud/strom/mc/api"
)

var cmd = flag.NewFlagSet("api", flag.ContinueOnError)
var tokenFlag = cmd.String("token", os.Getenv("MC_TOKEN"), "token to use")

var renameFlag = cmd.String("rename", "", "rename the player")

var uuidFlag = cmd.String("uuid", "", "get data about player with uuid")
var nameFlag = cmd.String("name", "", "get data about player with name")
var unsignedFlag = cmd.Bool("unsigned", true, "use unsigned mode for properties signatures")

var PleaseProvideArgumentErr = errors.New("please provide an argument")

func Run(args []string) (err error) {
	err = cmd.Parse(args)
	if err != nil {
		return
	}
	a := &api.Account{}
	if *tokenFlag != "" {
		a, err = api.NewAccountFromYGG(*tokenFlag)
		if err != nil {
			return
		}
	}
	var data any
	if *renameFlag != "" {
		data, err = a.ChangeName(*renameFlag)
		if err != nil {
			return
		}
	} else if *uuidFlag != "" || *nameFlag != "" {
		if *nameFlag != "" {
			var idName api.IdName
			idName, err = api.NameToId(*nameFlag)
			if err != nil {
				return
			}
			*uuidFlag = idName.Id
		}
		data, err = api.Profile(*uuidFlag, *unsignedFlag)
		if err != nil {
			return
		}
	} else {
		slog.Info("account", "name", a.Name, "uuid", a.Uuid)
		return
	}
	err = json.NewEncoder(os.Stdout).Encode(data)
	return
}
