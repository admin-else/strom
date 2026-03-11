package main

import (
	"flag"
	"fmt"

	"github.com/admin-else/strom/client"
)

var addrFlag = flag.String("addr", "localhost:25565", "address to connect to")
var srvFlag = flag.Bool("srv", false, "use SRV records to resolve the address")

func main() {
	flag.Parse()
	var err error
	if *srvFlag {
		*addrFlag, err = client.DoDNS(*addrFlag)
		if err != nil {
			panic(err)
		}
	}

	status, err := client.StatusRaw(*addrFlag)
	if err != nil {
		panic(err)
	}
	fmt.Println(status.Status)
	fmt.Printf("%v ping", status.PingReceiveTime.Sub(status.PingSendTime))
}
