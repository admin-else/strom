package level

import (
	"encoding/binary"
	"errors"
	"io"
	"log/slog"
	"math"

	data2 "git.anygate.cloud/anygatecloud/strom/mc/data"
	"git.anygate.cloud/anygatecloud/strom/mc/proto_base"
	"git.anygate.cloud/anygatecloud/strom/mc/util"
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

type Section struct {
	BlockCount     int16
	Blocks, Biomes *Storage
}

// MakeBiomeFormat returns the StorageFormat for biome data at the given version.
func MakeBiomeFormat(version string) StorageFormat {
	directBpe := uint8(math.Ceil(math.Log2(float64(len(data2.BiomesForVersion(version))))))
	return StorageFormat{
		AvailableBpes: []uint8{0, 1, 2, 3, directBpe},
		BiggestDirect: true,
		Len:           BiomesPerChunkSection,
	}
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
	s.Blocks, err = ReadSectionStorage(r, MakeBlockFormat(version))
	if err != nil {
		return
	}
	s.Biomes, err = ReadSectionStorage(r, MakeBiomeFormat(version))
	return
}

// SectionEncodePacket encodes a chunk section to the Minecraft packet wire format.
func SectionEncodePacket(w io.Writer, s Section) (err error) {
	err = binary.Write(w, binary.BigEndian, s.BlockCount)
	if err != nil {
		return
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
	if a.BlockCount != b.BlockCount {
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
