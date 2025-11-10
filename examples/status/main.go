package main

import (
	"strom"
	"strom/modules"

	"github.com/google/uuid"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	c, err := strom.ClientConnect("127.0.0.1:25566")
	must(err)
	err = c.Start(&modules.LoginClient{
		Account: modules.Account{
			Username: "Admin_rizz",
			UUID:     uuid.UUID{},
			Token:    "",
		},
		Connection: c,
	})
	must(err)
}
