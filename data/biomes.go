package data

import "errors"

// [
//
//	{
//	  "id": 0,
//	  "name": "badlands",
//	  "category": "mesa",
//	  "temperature": 2.0,
//	  "has_precipitation": false,
//	  "dimension": "overworld",
//	  "displayName": "Badlands",
//	  "color": 7254527
//	},

type Biome struct {
	Id               int32
	Name             string
	Category         string
	Temperature      float64
	HasPrecipitation bool
	Dimension        string
	DisplayName      string
	Color            int
}

var BiomesCache = make(map[string][]*Biome)

func BiomesForVersion(v string) (ret []*Biome) {
	var ok bool
	if ret, ok = BiomesCache[v]; ok {
		return
	}
	var b []*Biome
	must(LoadVersionedJson(v, "biomes", &b))
	BiomesCache[v] = b
	return b
}

func LookupBiomeById(version string, id int32) (biome *Biome, ok bool) {
	for _, b := range BiomesForVersion(version) {
		if b.Id == id {
			biome = b
			ok = true
			return
		}
	}
	return
}

var UnknownBiomeNameErr = errors.New("unknown biome Name")

func LookupBiomeByName(version string, name string) (biome *Biome, ok bool) {
	for _, b := range BiomesForVersion(version) {
		if b.Name == name {
			biome = b
			ok = true
			return
		}
	}
	return
}
