package data

import (
	"embed"

	"github.com/admin-else/strom/nbt"
)

//go:embed registry/*
var Registry embed.FS

func LoadRegistry(name string) (n *nbt.Tag, err error) {
	r, err := Registry.Open("registry/" + name + ".nbt")
	if err != nil {
		return
	}
	defer r.Close()
	n, err = nbt.ReadFile(r)
	return
}
