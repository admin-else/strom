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

// sectionFormatVersion26_4 is the first version where the section wire format
// added the fluid count short and grew the biome container to 16^3 entries.
const sectionFormatVersion26_4 = "26.4-snapshot-2"

// uses26_4SectionFormat reports whether the version's section wire format is the
// 26.4 form (nonEmpty+fluid shorts, 16^3 block and biome containers).
func uses26_4SectionFormat(version string) bool {
	return version == sectionFormatVersion26_4
}

// biomeEntries returns the number of biome entries per section for the version.
func biomeEntries(version string) int32 {
	if uses26_4SectionFormat(version) {
		return BlocksPerChunkSection
	}
	return BiomesPerChunkSection
}

// makeBiomeFormatForVersion returns the StorageFormat for biome data at the given version.
func makeBiomeFormatForVersion(version string) StorageFormat {
	directBpe := uint8(math.Ceil(math.Log2(float64(len(data2.BiomesForVersion(version))))))
	if uses26_4SectionFormat(version) {
		bpes := []uint8{0, 1, 2, 3, 4, 5, 6, 7, 8}
		if !slices.Contains(bpes, directBpe) {
			bpes = append(bpes, directBpe)
			slices.Sort(bpes)
		}
		return StorageFormat{
			AvailableBpes: bpes,
			BiggestDirect: true,
			Len:           BlocksPerChunkSection,
		}
	}
	return StorageFormat{
		AvailableBpes: []uint8{0, 1, 2, 3, directBpe},
		BiggestDirect: true,
		Len:           BiomesPerChunkSection,
	}
}

type Section struct {
	BlockCount int16
	FluidCount int16
	Blocks, Biomes *Storage
}

// MakeBiomeFormat returns the StorageFormat for biome data at the given version.
func MakeBiomeFormat(version string) StorageFormat {
	return makeBiomeFormatForVersion(version)
}

// MakeBlockFormat returns the StorageFormat for block data at the given version.
func MakeBlockFormat(version string) StorageFormat {
	directBpe := uint8(math.Ceil(math.Log2(float64(len(data2.BlocksForVersion(version))))))
	return StorageFormat{
		RedirectingBpes: []uint8{1, 2, 3},
		AvailableBpes:   []uint8{0, 4, 5, 6, 7, 8, directBpe},
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
	if uses26_4SectionFormat(version) {
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
	if uses26_4SectionFormat(version) {
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
