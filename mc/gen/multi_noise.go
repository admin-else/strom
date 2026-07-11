package gen

import (
	"fmt"
	"math"

	"github.com/admin-else/strom/mc/data"
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
	if err := data.LoadVersionedJson(version, "multi_noise_biome_source_parameter_list", &raw); err != nil {
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

	bestBiome := ""
	bestDist := int64(math.MaxInt64)

	for _, p := range points {
		dist := distanceToRange(t, p.Temperature) +
			distanceToRange(h, p.Humidity) +
			distanceToRange(c, p.Continentalness) +
			distanceToRange(e, p.Erosion) +
			distanceToRange(d, p.Depth) +
			distanceToRange(w, p.Weirdness) +
			p.Offset

		if dist < bestDist {
			bestDist = dist
			bestBiome = p.Biome
		}
	}

	return bestBiome, nil
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
