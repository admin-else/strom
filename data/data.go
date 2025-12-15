package data

import (
	"embed"
	"encoding/json"
	"errors"
)

//go:embed minecraft-data/data
var MinecraftData embed.FS

func LoadJson(path string, data any) (err error) {
	b, err := MinecraftData.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, data)
	return
}

func LoadVersionedJson(version, dataName string, data any) (err error) {
	b, err := MinecraftData.ReadFile("minecraft-data/data/" + Paths.Data[version][dataName] + "/" + dataName + ".json")
	if err != nil {
		return
	}
	err = json.Unmarshal(b, data)
	return
}

var UnknownMinecraftVersionError = errors.New("unknown minecraft version")

func must(err error) {
	if err != nil {
		panic(err)
	}
}
