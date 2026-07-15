package data

import (
	"fmt"
	"math"
	"sync"
)

// MultiNoisePreset names the two built-in multi-noise parameter lists.
const (
	MultiNoiseOverworld = "minecraft:overworld"
	MultiNoiseNether    = "minecraft:nether"
)

// ParamRange is a single quantized parameter range [Min, Max].
type ParamRange struct {
	Min, Max int64
}

// ParameterPoint is a 6-D parameter box plus an offset used by Minecraft's
// climate distance function.
type ParameterPoint struct {
	Biome           string
	Temperature     ParamRange
	Humidity        ParamRange
	Continentalness ParamRange
	Erosion         ParamRange
	Depth           ParamRange
	Weirdness       ParamRange
	Offset          int64
}

// MultiNoiseData holds the parameter lists for each preset.
type MultiNoiseData struct {
	Presets map[string][]ParameterPoint

	bvls    map[string]*bvhNode
	bvhOnce sync.Once
	bvhErr  error
}

// raw JSON shapes.
type rawParameter struct {
	Temperature     [2]int64 `json:"temperature"`
	Humidity        [2]int64 `json:"humidity"`
	Continentalness [2]int64 `json:"continentalness"`
	Erosion         [2]int64 `json:"erosion"`
	Depth           [2]int64 `json:"depth"`
	Weirdness       [2]int64 `json:"weirdness"`
	Offset          float64  `json:"offset"`
}

type rawEntry struct {
	Biome      string       `json:"biome"`
	Parameters rawParameter `json:"parameters"`
}

type rawMultiNoise map[string][]rawEntry

// LoadMultiNoise loads the expanded multi-noise parameter list for a Minecraft
// version. The data is expected to be present under the
// "multi_noise_biome_source_parameter_list" key in minecraft-data.
func LoadMultiNoise(version string) (*MultiNoiseData, error) {
	var raw rawMultiNoise
	if err := LoadVersionedJson(version, "multi_noise_biome_source_parameter_list", &raw); err != nil {
		return nil, fmt.Errorf("load multi_noise for %s: %w", version, err)
	}

	out := &MultiNoiseData{Presets: make(map[string][]ParameterPoint, len(raw))}
	for preset, entries := range raw {
		points := make([]ParameterPoint, len(entries))
		for i, e := range entries {
			points[i] = ParameterPoint{
				Biome:           e.Biome,
				Temperature:     ParamRange{Min: e.Parameters.Temperature[0], Max: e.Parameters.Temperature[1]},
				Humidity:        ParamRange{Min: e.Parameters.Humidity[0], Max: e.Parameters.Humidity[1]},
				Continentalness: ParamRange{Min: e.Parameters.Continentalness[0], Max: e.Parameters.Continentalness[1]},
				Erosion:         ParamRange{Min: e.Parameters.Erosion[0], Max: e.Parameters.Erosion[1]},
				Depth:           ParamRange{Min: e.Parameters.Depth[0], Max: e.Parameters.Depth[1]},
				Weirdness:       ParamRange{Min: e.Parameters.Weirdness[0], Max: e.Parameters.Weirdness[1]},
				Offset:          int64(math.Round(e.Parameters.Offset * 10000.0)),
			}
		}
		out.Presets[preset] = points
	}

	return out, nil
}

// Lookup finds the biome whose parameter box is closest to the supplied
// quantized climate parameters. The preset argument should be one of the
// MultiNoise* constants.
func (m *MultiNoiseData) Lookup(preset string, t, h, c, e, d, w int64) (string, error) {
	points, ok := m.Presets[preset]
	if !ok {
		return "", fmt.Errorf("unknown multi-noise preset: %s", preset)
	}

	if len(points) == 0 {
		return "", fmt.Errorf("multi-noise preset %s has no entries", preset)
	}

	m.bvhOnce.Do(func() { m.bvhErr = m.buildBVHs() })
	if m.bvhErr != nil {
		return "", m.bvhErr
	}

	root := m.bvls[preset]
	if root == nil {
		return "", fmt.Errorf("no BVH for preset: %s", preset)
	}

	q := point6{t, h, c, e, d, w}
	best := result{dist: int64(math.MaxInt64)}
	root.lookup(q, points, &best)
	return best.biome, nil
}

func distanceToRange(value int64, r ParamRange) int64 {
	if value < r.Min {
		d := r.Min - value
		return d * d
	}
	if value > r.Max {
		d := value - r.Max
		return d * d
	}
	return 0
}

// point6 is a quantized climate sample in six dimensions.
type point6 struct{ t, h, c, e, d, w int64 }

func (q point6) dist(p ParameterPoint) int64 {
	return distanceToRange(q.t, p.Temperature) +
		distanceToRange(q.h, p.Humidity) +
		distanceToRange(q.c, p.Continentalness) +
		distanceToRange(q.e, p.Erosion) +
		distanceToRange(q.d, p.Depth) +
		distanceToRange(q.w, p.Weirdness) +
		p.Offset
}

type result struct {
	dist  int64
	biome string
}

// bvhNode is a bounding-volume-hierarchy node for fast nearest-box search.
type bvhNode struct {
	bounds ParameterPoint
	left   *bvhNode
	right  *bvhNode
	leaf   []int // indices into original slice; non-nil only for leaves
}

func (m *MultiNoiseData) buildBVHs() error {
	m.bvls = make(map[string]*bvhNode, len(m.Presets))
	for preset, points := range m.Presets {
		indices := make([]int, len(points))
		for i := range indices {
			indices[i] = i
		}
		m.bvls[preset] = buildBVH(points, indices, 0)
	}
	return nil
}

func buildBVH(points []ParameterPoint, indices []int, depth int) *bvhNode {
	node := &bvhNode{}
	node.bounds = boundsOf(points, indices)

	// Small leaves avoid BVH overhead for tiny lists.
	if len(indices) <= 16 {
		node.leaf = make([]int, len(indices))
		copy(node.leaf, indices)
		return node
	}

	// Pick axis by largest variance.
	axis := depth % 6
	if axis < 0 {
		axis = 0
	}

	// Partition by median of the chosen axis' center.
	mid := len(indices) / 2
	quickSelectByAxis(points, indices, axis, mid)

	leftIdx := indices[:mid]
	rightIdx := indices[mid:]

	// Degenerate split: fall back to leaf to avoid infinite recursion.
	if len(leftIdx) == 0 || len(rightIdx) == 0 {
		node.leaf = make([]int, len(indices))
		copy(node.leaf, indices)
		return node
	}

	node.left = buildBVH(points, leftIdx, depth+1)
	node.right = buildBVH(points, rightIdx, depth+1)
	return node
}

func boundsOf(points []ParameterPoint, indices []int) ParameterPoint {
	if len(indices) == 0 {
		return ParameterPoint{}
	}
	b := points[indices[0]]
	for _, i := range indices[1:] {
		p := points[i]
		if p.Temperature.Min < b.Temperature.Min {
			b.Temperature.Min = p.Temperature.Min
		}
		if p.Temperature.Max > b.Temperature.Max {
			b.Temperature.Max = p.Temperature.Max
		}
		if p.Humidity.Min < b.Humidity.Min {
			b.Humidity.Min = p.Humidity.Min
		}
		if p.Humidity.Max > b.Humidity.Max {
			b.Humidity.Max = p.Humidity.Max
		}
		if p.Continentalness.Min < b.Continentalness.Min {
			b.Continentalness.Min = p.Continentalness.Min
		}
		if p.Continentalness.Max > b.Continentalness.Max {
			b.Continentalness.Max = p.Continentalness.Max
		}
		if p.Erosion.Min < b.Erosion.Min {
			b.Erosion.Min = p.Erosion.Min
		}
		if p.Erosion.Max > b.Erosion.Max {
			b.Erosion.Max = p.Erosion.Max
		}
		if p.Depth.Min < b.Depth.Min {
			b.Depth.Min = p.Depth.Min
		}
		if p.Depth.Max > b.Depth.Max {
			b.Depth.Max = p.Depth.Max
		}
		if p.Weirdness.Min < b.Weirdness.Min {
			b.Weirdness.Min = p.Weirdness.Min
		}
		if p.Weirdness.Max > b.Weirdness.Max {
			b.Weirdness.Max = p.Weirdness.Max
		}
		if p.Offset < b.Offset {
			b.Offset = p.Offset
		}
	}
	return b
}

func (n *bvhNode) lowerBound(q point6) int64 {
	return distanceToRange(q.t, n.bounds.Temperature) +
		distanceToRange(q.h, n.bounds.Humidity) +
		distanceToRange(q.c, n.bounds.Continentalness) +
		distanceToRange(q.e, n.bounds.Erosion) +
		distanceToRange(q.d, n.bounds.Depth) +
		distanceToRange(q.w, n.bounds.Weirdness) +
		n.bounds.Offset
}

func (n *bvhNode) lookup(q point6, points []ParameterPoint, best *result) {
	if n.lowerBound(q) >= best.dist {
		return
	}

	if n.leaf != nil {
		for _, idx := range n.leaf {
			p := points[idx]
			d := q.dist(p)
			if d < best.dist {
				best.dist = d
				best.biome = p.Biome
			}
		}
		return
	}

	// Visit the child with the smaller lower bound first for better pruning.
	lbL := n.left.lowerBound(q)
	lbR := n.right.lowerBound(q)
	if lbL < lbR {
		n.left.lookup(q, points, best)
		n.right.lookup(q, points, best)
	} else {
		n.right.lookup(q, points, best)
		n.left.lookup(q, points, best)
	}
}

func quickSelectByAxis(points []ParameterPoint, indices []int, axis, k int) {
	if len(indices) <= 1 {
		return
	}
	for {
		pivot := points[indices[len(indices)/2]]
		pivotVal := axisCenter(pivot, axis)
		left, right := 0, len(indices)-1
		for left <= right {
			for axisCenter(points[indices[left]], axis) < pivotVal {
				left++
			}
			for axisCenter(points[indices[right]], axis) > pivotVal {
				right--
			}
			if left <= right {
				indices[left], indices[right] = indices[right], indices[left]
				left++
				right--
			}
		}
		if k <= right {
			indices = indices[:right+1]
		} else if k >= left {
			indices = indices[left:]
			k -= left
		} else {
			return
		}
	}
}

func axisCenter(p ParameterPoint, axis int) int64 {
	switch axis {
	case 0:
		return (p.Temperature.Min + p.Temperature.Max) / 2
	case 1:
		return (p.Humidity.Min + p.Humidity.Max) / 2
	case 2:
		return (p.Continentalness.Min + p.Continentalness.Max) / 2
	case 3:
		return (p.Erosion.Min + p.Erosion.Max) / 2
	case 4:
		return (p.Depth.Min + p.Depth.Max) / 2
	case 5:
		return (p.Weirdness.Min + p.Weirdness.Max) / 2
	}
	return 0
}
