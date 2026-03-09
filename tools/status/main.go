package main

import (
	"flag"
	"fmt"

	"github.com/admin-else/strom/client"
)

var addrFlag = flag.String("addr", "localhost:25565", "address to connect to")

func main() {
	flag.Parse()
	status, err := client.StatusRaw(*addrFlag)
	if err != nil {
		panic(err)
	}
	fmt.Println(status.Status)
	fmt.Printf("%v ping", status.PingReceiveTime.Sub(status.PingSendTime))
}
