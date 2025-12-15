package data

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

type Block struct {
	Id           int
	Name         string
	DisplayName  string
	Hardness     float64
	Resistance   float64
	StackSize    int
	Diggable     bool
	Material     string
	Transparent  bool
	EmitLight    int
	FilterLight  int
	DefaultState int32
	MinStateId   int32
	MaxStateId   int32
	States       []struct {
		Name      string
		Type      string
		NumValues int
	}
	Drops       []int
	BoundingBox string
}

var BlocksCache = make(map[string][]Block)

func BlocksForVersion(v string) (ret []Block) {
	var ok bool
	if ret, ok = BlocksCache[v]; ok {
		return
	}
	var b []Block
	must(LoadVersionedJson(v, "blocks", &b))
	BlocksCache[v] = b
	return b
}

func LookupBlockByStateId(version string, stateId int32) (block Block, ok bool) {
	for _, b := range BlocksForVersion(version) {
		if b.MinStateId <= stateId && stateId <= b.MaxStateId {
			block = b
			ok = true
		}
	}
	return
}
