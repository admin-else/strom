package modules

import (
	"fmt"

	"github.com/admin-else/strom"
)

type LoginServer struct {
	*strom.Conn
}

func (l LoginServer) Default(event any) (err error) {
	err = fmt.Errorf("unexpected event: %#v", event)
	return
}
