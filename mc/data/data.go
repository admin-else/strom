package data

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
)

//go:generate ./download_minecraft_data.sh

//go:embed minecraft-data/pc minecraft-data/dataPaths.json
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
	b, err := MinecraftData.ReadFile("minecraft-data/" + Paths.Data[version][dataName] + "/" + dataName + ".json")
	if err != nil {
		// If you get this error, you need to probably need to update the download script
		return fmt.Errorf("minecraft-data/%s/%s/%s.json: %w", version, Paths.Data[version][dataName], dataName, err)
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
