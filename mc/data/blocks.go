package data

import (
	"errors"
	"slices"
	"strconv"

	"github.com/admin-else/strom/mc/util"
)

//  {
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
	States       []BlockState
	Drops        []int
	BoundingBox  string
}

type BlockState struct {
	Name      string
	Type      string
	MinValue  int32 // this is generated on load and used for state value calculations
	NumValues int32 `json:"num_values"`
	Values    []string
}

var BlocksCache = make(map[string][]*Block)
var BlocksCacheIdMap = make(map[string]map[int32]*Block)

// BlocksForVersion returns the block definitions for the given Minecraft version, caching them after the first load.
func BlocksForVersion(v string) (ret []*Block) {
	var ok bool
	if ret, ok = BlocksCache[v]; ok {
		return
	}
	var b []*Block
	must(LoadVersionedJson(v, "blocks", &b))
	BlocksCache[v] = b
	for _, block := range b {
		for i, state := range block.States {
			if state.Type != "int" {
				continue
			}
			values := make([]int32, state.NumValues)
			for j, stateValue := range state.Values {
				n := util.MustT(strconv.Atoi(stateValue))
				values[j] = int32(n)
			}
			block.States[i].MinValue = slices.Min(values)
		}
	}

	return b
}

// LookupBlockByStateId returns the block for the given global palette state ID in the given version.
func LookupBlockByStateId(version string, stateId int32) (block *Block, ok bool) {
	idmap, found := BlocksCacheIdMap[version]
	if !found {
		idmap = make(map[int32]*Block)
		for _, b := range BlocksForVersion(version) {
			for i := range b.MaxStateId - b.MinStateId + 1 {
				idmap[b.MinStateId+i] = b
			}
		}
		BlocksCacheIdMap[version] = idmap
	}
	block, ok = idmap[stateId]
	return
}

var UnknownBlockNameErr = errors.New("unknown block Name")

// LookupBlockByName returns the block with the given name in the given version.
func LookupBlockByName(version string, name string) (block *Block, ok bool) {
	for _, b := range BlocksForVersion(version) {
		if b.Name == name {
			block = b
			ok = true
		}
	}
	return
}

// StateIdFromBlocKAndStateMap returns the global palette state ID for the given block name and state property map in the given version.
func StateIdFromBlocKAndStateMap(version string, name string, stateMap map[string]string) (stateId int32, err error) {
	b, ok := LookupBlockByName(version, name)
	if !ok {
		err = BlockNotFoundErr
		return
	}
	if stateMap == nil {
		stateId = b.DefaultState
	} else {
		stateId, err = b.IdFromStateMap(stateMap)
	}
	if b.MinStateId > stateId || b.MaxStateId < stateId {
		err = BlockStateOutOfRangeErr
	}
	return
}

var BlockStateOutOfRangeErr = errors.New("stateId out of range")
var BlockNotFoundErr = errors.New("block not found")
var InvalidBoolValueErr = errors.New("invalid bool value")
var BadBlockStateTypeErr = errors.New("bad block state type")
var BadEnumErr = errors.New("bad enum")
var SupplyAllStatesErr = errors.New("supply all states")

// FromBlockState converts a global palette state ID back into a block definition and a map of state properties.
func FromBlockState(version string, stateId int32) (b *Block, stateData map[string]string, err error) {
	b, ok := LookupBlockByStateId(version, stateId)
	if !ok {
		err = BlockNotFoundErr
	}
	stateData, err = b.StateMapFromId(stateId)
	return
}

func mcDataBlockStateTypeParse(n int32, s *BlockState) (v string, err error) {
	switch s.Type {
	case "int":
		v = strconv.Itoa(int(n))
	case "enum":
		v = s.Values[n]
	case "bool":
		if n == 0 { // MOJANK: 0 is true, 1 is false
			v = "true"
		} else if n == 1 {
			v = "false"
		} else {
			err = InvalidBoolValueErr
		}
	default:
		err = BadBlockStateTypeErr
	}
	return
}

func mcDataBlockStateTypeSerialize(v string, s *BlockState) (n int32, err error) {
	if v == "true" {
		return 0, nil
	}
	if v == "false" {
		return 1, nil
	}
	n1, errMabye := strconv.Atoi(v)
	if errMabye == nil {
		return int32(n1) - s.MinValue, nil
	}
	for i, enumValue := range s.Values {
		if enumValue == v {
			return int32(i), nil
		}
	}
	err = BadEnumErr
	return
}

func (b *Block) StateMapFromId(stateId int32) (ret map[string]string, err error) {
	if stateId < b.MinStateId || stateId > b.MaxStateId {
		err = BlockStateOutOfRangeErr
		return
	}
	ret = make(map[string]string)
	stateId -= b.MinStateId

	a := int32(1)
	for i := range len(b.States) {
		s := b.States[len(b.States)-1-i]
		ret[s.Name], err = mcDataBlockStateTypeParse((stateId/a)%s.NumValues, &s)
		if err != nil {
			return
		}
		a *= s.NumValues
	}
	return
}

func (b *Block) IdFromStateMap(m map[string]string) (ret int32, err error) {
	a := int32(1)
	for i := range len(b.States) {
		s := b.States[len(b.States)-1-i]
		v, found := m[s.Name]
		if !found {
			err = SupplyAllStatesErr
			return
		}
		var n int32
		n, err = mcDataBlockStateTypeSerialize(v, &s)
		if err != nil {
			return
		}
		ret += n * a
		a *= s.NumValues
	}
	ret += b.MinStateId
	return
}
