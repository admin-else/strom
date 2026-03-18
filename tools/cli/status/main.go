package status

import (
	"flag"
	"fmt"

	"github.com/admin-else/strom/client"
)

var cmd = flag.NewFlagSet("status", flag.ContinueOnError)
var addrFlag = cmd.String("addr", "localhost:25565", "address to connect to")
var srvFlag = cmd.Bool("srv", false, "use SRV records to resolve the address")

func Run(args []string) (err error) {
	err = cmd.Parse(args)
	if err != nil {
		return
	}
	if *srvFlag {
		*addrFlag, err = client.DoDNS(*addrFlag)
		if err != nil {
			return
		}
	}

	status, err := client.StatusRaw(*addrFlag)
	if err != nil {
		return
	}
	fmt.Println(status.Status)
	return
}
