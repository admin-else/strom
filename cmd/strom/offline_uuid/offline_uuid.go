package offline_uuid

import (
	"flag"
	"fmt"

	"github.com/admin-else/strom/crypto"
)

var cmd = flag.NewFlagSet("offline-uuid", flag.ContinueOnError)
var nameFlag = cmd.String("name", "Notch", "name of the player")

func Run(args []string) (err error) {
	err = cmd.Parse(args)
	if err != nil {
		return
	}
	fmt.Println(crypto.FromOfflinePlayer(*nameFlag))
	return
}
