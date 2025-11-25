package main

import (
	"errors"
	"fmt"
	"strom"
	"strom/modules"

	"github.com/google/uuid"
)

var altAccount = modules.Account{
	Username: "laibergraceful70",
	Uuid:     uuid.MustParse("3a863af533e5452591e913e021bfa724"),
	Token:    "eyJraWQiOiIwNDkxODEiLCJhbGciOiJSUzI1NiJ9.eyJ4dWlkIjoiMjc3NzIxOTc3MDIzMjM1OCIsImFnZyI6IkFkdWx0Iiwic3ViIjoiMzczMjI0OTktZTAwYS0…JRnCbAsBJl5zg-GAu6NkX1JCvpWdQkgFb0H5RYi2dIUKhNpP57WG2gh562T3BoPELVXDQZQVUBzDFeDuVDGGSkKmSmhzO7l-fEPBt2CjnaLCyr1FmYFSRzDBy1w",
}

type Client struct {
	*strom.Conn
}

func (c *Client) Default(event any) (err error) {
	fmt.Printf("%T%v\n", event, event)
	return
}

func (c *Client) OnStart(_ strom.OnStart) (err error) {
	err = modules.Login(c.Conn, altAccount)
	if errors.Is(err, strom.HandlerDone) {
		return nil
	}
	return
}

func main() {
	c, err := strom.Connect("127.0.0.1:25566")
	if err != nil {
		panic(err)
	}
	defer c.Close()
	err = c.Start(&Client{c})
	if err != nil {
		panic(err)
	}
}
