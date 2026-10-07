package level

import (
	"encoding/binary"
	"errors"
	"io"
	"log/slog"
	"math"
	"slices"

	data2 "github.com/admin-else/strom/mc/data"
	"github.com/admin-else/strom/mc/proto_base"
	"github.com/admin-else/strom/mc/proto_generated"
	"github.com/admin-else/strom/mc/util"
)

type BlockState struct {
	Id         string
	Properties map[string]string
}

const (
	ChunkWidth            = 16
	ChunkColumns          = ChunkWidth * ChunkWidth
	BlocksPerChunkSection = ChunkWidth * ChunkWidth * ChunkWidth
	ChunkBiomesWidth      = 4
	BiomesPerChunkSection = ChunkBiomesWidth * ChunkBiomesWidth * ChunkBiomesWidth
)

// sectionFluidCountFirstVersion is the first version whose chunk section wire
// format stores the fluid count as a second short (ViaVersion: ChunkSection
// fluid count is available for 26.1+).
const sectionFluidCountFirstVersion = "26.1"

// sectionBiome4096FirstVersion is the first version whose chunk section biome
// container holds 16^3 = 4096 entries; earlier versions (including 26.2) still
// use 4^3 = 64.
const sectionBiome4096FirstVersion = "26.4-snapshot-2"

// versionAtLeast reports whether version is at or after boundary in release order.
func versionAtLeast(version, boundary string) (ret bool) {
	i := slices.Index(proto_generated.SupportedVersions, version)
	j := slices.Index(proto_generated.SupportedVersions, boundary)
	if i == -1 || j == -1 {
		return version == boundary
	}
	return i >= j
}

// hasSectionFluidCount reports whether the version's section wire format
// includes the fluid count short.
func hasSectionFluidCount(version string) (ret bool) {
	return versionAtLeast(version, sectionFluidCountFirstVersion)
}

// sectionBiomeEntries returns the number of biome entries per section for the version.
func sectionBiomeEntries(version string) (ret int32) {
	if versionAtLeast(version, sectionBiome4096FirstVersion) {
		return BlocksPerChunkSection
	}
	return BiomesPerChunkSection
}

// makeBiomeFormatForVersion returns the StorageFormat for biome data at the given version.
func makeBiomeFormatForVersion(version string) StorageFormat {
	directBpe := uint8(math.Ceil(math.Log2(float64(len(data2.BiomesForVersion(version))))))
	return StorageFormat{
		AvailableBpes: makeBpeRange(0, directBpe),
		BiggestDirect: true,
		Len:           sectionBiomeEntries(version),
	}
}

// makeBpeRange returns every bits-per-entry value a PalettedContainer can use:
// 0 for a single-value storage, every value from minBpe up to the direct
// (global palette) size, which is always the last element so it is read as a
// direct palette with no inline palette.
func makeBpeRange(minBpe, directBpe uint8) (ret []uint8) {
	ret = []uint8{0}
	start := minBpe
	if start < 1 {
		start = 1
	}
	for v := start; v <= directBpe; v++ {
		ret = append(ret, v)
	}
	if !slices.Contains(ret, directBpe) {
		ret = append(ret, directBpe)
	}
	slices.Sort(ret)
	return slices.Compact(ret)
}

type Section struct {
	BlockCount     int16
	FluidCount     int16
	Blocks, Biomes *Storage
}

// MakeBiomeFormat returns the StorageFormat for biome data at the given version.
func MakeBiomeFormat(version string) StorageFormat {
	return makeBiomeFormatForVersion(version)
}

// MakeBlockFormat returns the StorageFormat for block data at the given version.
func MakeBlockFormat(version string) StorageFormat {
	directBpe := uint8(math.Ceil(math.Log2(float64(data2.BlockStateCount(version)))))
	return StorageFormat{
		RedirectingBpes: []uint8{1, 2, 3},
		AvailableBpes:   makeBpeRange(4, directBpe),
		BiggestDirect:   true,
		Len:             BlocksPerChunkSection,
	}
}

var BadPaletteLenErr = errors.New("bad palette length")

// ReadSectionStorage reads a packed section storage from a reader.
func ReadSectionStorage(r io.Reader, format StorageFormat) (s *Storage, err error) {
	var bpe uint8
	err = binary.Read(r, binary.BigEndian, &bpe)
	if err != nil {
		return
	}
	var pallette []int32
	if bpe == 0 {
		var v int32
		v, err = proto_base.DecodeVarInt(r)
		if err != nil {
			return
		}
		pallette = []int32{v}
	} else if !util.IsLastElement(format.AvailableBpes, bpe) {
		var paletteLen int32
		paletteLen, err = proto_base.DecodeVarInt(r)
		if err != nil {
			return
		}
		if paletteLen < 0 {
			err = BadPaletteLenErr
			return
		}
		for range paletteLen {
			var v int32
			v, err = proto_base.DecodeVarInt(r)
			if err != nil {
				return
			}
			if slog.Default().Enabled(nil, slog.LevelDebug) {
				b, _ := data2.LookupBlockByStateId("1.21.11", v)
				if b == nil {
					slog.Debug("unknown palette block in chunk packet", "id", v)
				} else {
					slog.Debug("palette block in chunk packet", "id", v, "name", b.Name)
				}

			}

			pallette = append(pallette, v)
		}
	}
	return format.ImportFromReader(r, bpe, pallette)
}

// WriteSectionStorage writes a packed section storage to a writer.
func WriteSectionStorage(w io.Writer, s *Storage) (err error) {
	err = binary.Write(w, binary.BigEndian, s.bpe)
	if err != nil {
		return
	}
	if s.bpe == 0 {
		err = proto_base.EncodeVarInt(w, s.palette[0])
		if err != nil {
			return
		}
	} else if !util.IsLastElement(s.format.AvailableBpes, s.bpe) {
		err = proto_base.EncodeVarInt(w, int32(len(s.palette)))
		if err != nil {
			return
		}
		for _, v := range s.palette {
			err = proto_base.EncodeVarInt(w, v)
			if err != nil {
				return
			}
		}
	}
	return binary.Write(w, binary.BigEndian, s.data)
}

// SectionDecodePacket decodes a chunk section from the Minecraft packet wire format.
func SectionDecodePacket(r io.Reader, version string) (s Section, err error) {
	s = Section{}
	err = binary.Read(r, binary.BigEndian, &s.BlockCount)
	if err != nil {
		return
	}
	if hasSectionFluidCount(version) {
		err = binary.Read(r, binary.BigEndian, &s.FluidCount)
		if err != nil {
			return
		}
	}
	s.Blocks, err = ReadSectionStorage(r, MakeBlockFormat(version))
	if err != nil {
		return
	}
	s.Biomes, err = ReadSectionStorage(r, makeBiomeFormatForVersion(version))
	return
}

// SectionEncodePacket encodes a chunk section to the Minecraft packet wire format.
func SectionEncodePacket(w io.Writer, s Section) (err error) {
	return SectionEncodePacketVersion(w, s, "")
}

// SectionEncodePacketVersion encodes a chunk section to the wire format for the given version.
func SectionEncodePacketVersion(w io.Writer, s Section, version string) (err error) {
	err = binary.Write(w, binary.BigEndian, s.BlockCount)
	if err != nil {
		return
	}
	if hasSectionFluidCount(version) {
		err = binary.Write(w, binary.BigEndian, s.FluidCount)
		if err != nil {
			return
		}
	}
	err = WriteSectionStorage(w, s.Blocks)
	if err != nil {
		return
	}
	err = WriteSectionStorage(w, s.Biomes)
	return
}

// SectionEquals returns whether two sections are equal.
func SectionEquals(a, b Section) bool {
	if a.BlockCount != b.BlockCount || a.FluidCount != b.FluidCount {
		return false
	}
	if !a.Blocks.Equals(b.Blocks) {
		return false
	}
	if !a.Biomes.Equals(b.Biomes) {
		return false
	}
	return true
}
