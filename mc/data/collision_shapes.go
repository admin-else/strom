package data

import (
	"encoding/json"
	"strconv"
	"sync"
)

// BlockCollisionShapeBoxes mirrors minecraft-data's blockCollisionShapes.json:
// every block maps to a collision shape id, or per-state-id array of shape ids;
// each shape is a list of axis-aligned boxes [x1, y1, z1, x2, y2, z2] in block
// space.
type blockCollisionShapesFile struct {
	Blocks map[string]json.RawMessage `json:"blocks"`
	Shapes map[string][][6]float64    `json:"shapes"`
}

var (
	collisionShapesCache     = make(map[string]*blockCollisionShapesFile)
	collisionShapesCacheLock sync.Mutex
)

func loadBlockCollisionShapes(version string) (file *blockCollisionShapesFile, err error) {
	collisionShapesCacheLock.Lock()
	defer collisionShapesCacheLock.Unlock()
	if cached, ok := collisionShapesCache[version]; ok {
		return cached, nil
	}
	file = &blockCollisionShapesFile{}
	if err = LoadVersionedJson(version, "blockCollisionShapes", file); err != nil {
		return nil, err
	}
	collisionShapesCache[version] = file
	return file, nil
}

// CollisionShapeBoxesAt returns the collision boxes for the given global block
// state id, in block-local space. ok is false when the block has no collision
// data or an empty shape.
func CollisionShapeBoxesAt(version string, stateId int32) (boxes [][6]float64, ok bool) {
	file, err := loadBlockCollisionShapes(version)
	if err != nil {
		return nil, false
	}
	block, found := LookupBlockByStateId(version, stateId)
	if !found {
		return nil, false
	}
	raw, found := file.Blocks[block.Name]
	if !found {
		return nil, false
	}

	shapeId := -1
	var single int
	if err := json.Unmarshal(raw, &single); err == nil {
		shapeId = single
	} else {
		var list []int
		if err := json.Unmarshal(raw, &list); err != nil {
			return nil, false
		}
		index := int(stateId - block.MinStateId)
		if index < 0 || index >= len(list) {
			return nil, false
		}
		shapeId = list[index]
	}

	shape, found := file.Shapes[strconv.Itoa(shapeId)]
	if !found || len(shape) == 0 {
		return nil, false
	}
	return shape, true
}
