package data

import (
	_ "embed"
)

var Paths struct {
	Data map[string]map[string]string `json:"pc"`
}

func init() {
	must(LoadJson("minecraft-data/data/dataPaths.json", &Paths))
}
