package data

import (
	"embed"

	nbt2 "git.anygate.cloud/anygatecloud/strom/mc/nbt"
)

// TODO: this is stupid loginPacket.json exists next version migrate to that

//go:embed registry/*
var Registry embed.FS

// LoadRegistry loads a registry NBT file by name from the embedded registry data.
func LoadRegistry(name string) (n *nbt2.Tag, err error) {
	r, err := Registry.Open("registry/" + name + ".nbt")
	if err != nil {
		return
	}
	defer r.Close()
	n, err = nbt2.ReadUnstructuredFile(r)
	return
}
