package data

import (
	"errors"
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
	NumValues int32 `json:"num_values"`
	Values    []string
}

var BlocksCache = make(map[string][]*Block)

func BlocksForVersion(v string) (ret []*Block) {
	var ok bool
	if ret, ok = BlocksCache[v]; ok {
		return
	}
	var b []*Block
	must(LoadVersionedJson(v, "blocks", &b))
	BlocksCache[v] = b
	return b
}

func LookupBlockByStateId(version string, stateId int32) (block *Block, ok bool) {
	for _, b := range BlocksForVersion(version) {
		if b.MinStateId <= stateId && stateId <= b.MaxStateId {
			block = b
			ok = true
		}
	}
	return
}

func LookupBlockByName(version string, stateId int32) (block *Block, ok bool) {
	for _, b := range BlocksForVersion(version) {
		if b.MinStateId <= stateId && stateId <= b.MaxStateId {
			block = b
			ok = true
		}
	}
	return
}

var BlockStateOutOfRangeErr = errors.New("stateId out of range")
var BlockNotFoundErr = errors.New("block not found")
var InvalidBoolValueErr = errors.New("invalid bool value")
var BadBlockStateTypeErr = errors.New("bad block state type")
var BadEnumErr = errors.New("bad enum")
var SupplyAllStatesErr = errors.New("supply all states")

func FromBlockState(version string, stateId int32) (b *Block, stateData map[string]any, err error) {
	b, ok := LookupBlockByStateId(version, stateId)
	if !ok {
		err = BlockNotFoundErr
	}
	stateData, err = b.StateMapFromId(stateId)
	return
}

func mcDataBlockStateTypeParse(n int32, s *BlockState) (v any, err error) {
	switch s.Type {
	case "int":
		v = n
	case "enum":
		v = s.Values[n]
	case "bool":
		if n == 0 { // MOJANK: 0 is true, 1 is false
			v = true
		} else if n == 1 {
			v = false
		} else {
			err = InvalidBoolValueErr
		}
	default:
		err = BadBlockStateTypeErr
	}
	return
}

func mcDataBlockStateTypeSerialize(v any, s *BlockState) (n int32, err error) {
	switch t := v.(type) {
	case int32:
		n = t
	case bool:
		if t {
			n = 0
		} else {
			n = 1
		}
	case string:
		for i, enumValue := range s.Values {
			if enumValue == t {
				return int32(i), nil
			}
		}
		err = BadEnumErr
	default:
		err = BadBlockStateTypeErr
	}
	return
}

func (b *Block) StateMapFromId(stateId int32) (ret map[string]any, err error) {
	if stateId < b.MinStateId || stateId > b.MaxStateId {
		err = BlockStateOutOfRangeErr
		return
	}
	ret = make(map[string]any)
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

func (b *Block) IdFromStateMap(m map[string]any) (ret int32, err error) {
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
