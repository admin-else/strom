package data

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"

	mcdata "github.com/admin-else/strom/mc/data"
	"github.com/admin-else/strom/mc/proto_generated"
	"github.com/admin-else/strom/mc/util"
)

// strom data 1.21.11 blocks id 0 name => air
//   {
//    "id": 0,
//    "name": "air",
//    "displayName": "Air",
//    "hardness": 0.0,
//    "resistance": 0.0,
//    "stackSize": 64,
//    "diggable": false,
//    "material": "default",
//    "transparent": true,
//    "emitLight": 0,
//    "filterLight": 0,
//    "defaultState": 0,
//    "minStateId": 0,
//    "maxStateId": 0,
//    "states": [],
//    "drops": [],
//    "boundingBox": "empty"
//  },

func Run(args []string) (err error) {
	var dataValues any
	switch len(args) {
	case 0: // print available versions
		dataValues = proto_generated.SupportedVersions
	case 1: // print available data for version
		dataValues = slices.Collect(maps.Keys(mcdata.Paths.Data[args[0]]))
	case 2:
		err = mcdata.LoadVersionedJson(args[0], args[1], &dataValues)
	case 3:
		err = mcdata.LoadVersionedJson(args[0], args[1], &dataValues)
		if err != nil {
			return
		}
		var dataSlices []any
		dataSlices, err = util.AssertTypeError[[]any](dataValues)
		if err != nil {
			return
		}
		var result = make([]any, len(dataSlices))
		for _, v := range dataSlices {
			var dataMap map[string]any
			dataMap, err = util.AssertTypeError[map[string]any](v)
			if err != nil {
				return
			}
			var vToAppend any
			vToAppend, err = util.MapGetError(dataMap, args[2])
			if err != nil {
				return
			}
			result = append(result, vToAppend)
		}
		dataValues = result
	case 4, 5:
		err = mcdata.LoadVersionedJson(args[0], args[1], &dataValues)
		if err != nil {
			return
		}
		var dataSlices []any
		dataSlices, err = util.AssertTypeError[[]any](dataValues)
		if err != nil {
			return
		}
		found := false
		for _, v := range dataSlices {
			var dataMap map[string]any
			dataMap, err = util.AssertTypeError[map[string]any](v)
			if err != nil {
				return
			}
			if v2, ok := dataMap[args[2]]; ok {
				if args[3] != fmt.Sprint(v2) {
					continue
				}

				dataValues = dataMap
				found = true
				if len(args) == 5 {
					dataValues, err = util.MapGetError(dataMap, args[4])
					if err != nil {
						return
					}
				}
				break
			}
		}
		if !found {
			err = fmt.Errorf("not found")
			return
		}
	default:
		err = fmt.Errorf("too many arguments")
	}

	if err == nil {
		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "  ")
		err = e.Encode(dataValues)
		return
	}
	return
}
