package data

import (
	"embed"

	"github.com/admin-else/strom/nbt"
)

// TODO: this is stupid loginPacket.json exists next version migrate to that

//go:embed registry/*
var Registry embed.FS

func LoadRegistry(name string) (n *nbt.Tag, err error) {
	r, err := Registry.Open("registry/" + name + ".nbt")
	if err != nil {
		return
	}
	defer r.Close()
	n, err = nbt.ReadUnstructuredFile(r)
	return
}
