package nbt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"unicode/utf16"
)

const (
	End int8 = iota
	Byte
	Short
	Int
	Long
	Float
	Double
	ByteArray
	String
	List
	Compound
	IntArray
	LongArray
)

var Order binary.ByteOrder = binary.BigEndian

var BadStringLengthError = errors.New("bad string length")

type Tag struct {
	Name  string
	Value any
}

type Anon struct {
	Value any
}

func (s *Anon) Decode(r io.Reader) (err error) {
	return s.DecodeAccounter(r, NbtAccounterUnlimitedHeap())
}

// DecodeAccounter mirrors Anon.Decode with a caller-supplied NbtAccounter, which
// bounds the heap and nesting depth like Java's NbtIo read path.
func (s *Anon) DecodeAccounter(r io.Reader, accounter *NbtAccounter) (err error) {
	id := int8(0)
	err = binary.Read(r, Order, &id)
	if err != nil {
		return
	}
	s.Value, err = readPayload(id, r, accounter)
	if err != nil {
		return
	}
	return
}

func (s *Anon) Encode(w io.Writer) (err error) {
	id, err := getId(s.Value)
	if err != nil {
		return
	}
	err = writePayload(id, w)
	if err != nil {
		return
	}
	err = writePayload(s.Value, w)
	return
}

func (s *Tag) Decode(r io.Reader) (err error) {
	return s.DecodeAccounter(r, NbtAccounterUnlimitedHeap())
}

// DecodeAccounter mirrors Tag.Decode with a caller-supplied NbtAccounter. The
// root tag's name is read but not accounted, mirroring NbtIo.readUnnamedTag's
// StringTag.skipString.
func (s *Tag) DecodeAccounter(r io.Reader, accounter *NbtAccounter) (err error) {
	id := int8(0)
	err = binary.Read(r, Order, &id)
	if err != nil {
		return
	}
	if id == End {
		return
	}
	name := ""
	name, err = readNbtString(r)
	if err != nil {
		return
	}
	s.Value, err = readPayload(id, r, accounter)
	if err != nil {
		return
	}
	s.Name = name
	return
}

func (s *Tag) Encode(w io.Writer) (err error) {
	id, err := getId(s.Value)
	if err != nil {
		return
	}
	err = writePayload(id, w)
	if err != nil {
		return
	}
	err = writePayload(s.Name, w)
	if err != nil {
		return
	}
	err = writePayload(s.Value, w)
	return
}

func getId(v any) (id int8, err error) {
	switch v.(type) {
	case struct{}:
		id = End
	case int8:
		id = Byte
	case int16:
		id = Short
	case int32:
		id = Int
	case int64:
		id = Long
	case float32:
		id = Float
	case float64:
		id = Double
	case []int8:
		id = ByteArray
	case string:
		id = String
	case []any:
		id = List
	case map[string]any:
		id = Compound
	case []int32:
		id = IntArray
	case []int64:
		id = LongArray
	default:
		err = fmt.Errorf("unsupported type %T", v)
	}
	return
}

func readPayload(id int8, r io.Reader, accounter *NbtAccounter) (ret any, err error) {
	switch id {
	case End:
		if err = accounter.AccountBytes(8); err != nil {
			return
		}
		ret = struct{}{}
	case Byte:
		if err = accounter.AccountBytes(9); err != nil {
			return
		}
		rett := int8(0)
		err = binary.Read(r, Order, &rett)
		ret = rett
	case Short:
		if err = accounter.AccountBytes(10); err != nil {
			return
		}
		rett := int16(0)
		err = binary.Read(r, Order, &rett)
		ret = rett
	case Int:
		if err = accounter.AccountBytes(12); err != nil {
			return
		}
		rett := int32(0)
		err = binary.Read(r, Order, &rett)
		ret = rett
	case Long:
		if err = accounter.AccountBytes(16); err != nil {
			return
		}
		rett := int64(0)
		err = binary.Read(r, Order, &rett)
		ret = rett
	case Float:
		if err = accounter.AccountBytes(12); err != nil {
			return
		}
		rett := float32(0)
		err = binary.Read(r, Order, &rett)
		ret = rett
	case Double:
		if err = accounter.AccountBytes(16); err != nil {
			return
		}
		rett := float64(0)
		err = binary.Read(r, Order, &rett)
		ret = rett
	case ByteArray:
		if err = accounter.AccountBytes(24); err != nil {
			return
		}
		var l int32
		err = binary.Read(r, Order, &l)
		if err != nil {
			return
		}
		if l < 0 {
			err = errors.New("bad array length")
			return
		}
		if err = accounter.AccountBytesN(1, int64(l)); err != nil {
			return
		}
		var rett []int8
		for range l {
			var element int8
			err = binary.Read(r, Order, &element)
			if err != nil {
				return
			}
			rett = append(rett, element)
		}
		ret = rett
	case String:
		if err = accounter.AccountBytes(36); err != nil {
			return
		}
		s := ""
		s, err = readNbtString(r)
		if err != nil {
			return
		}
		if err = accounter.AccountBytesN(2, nbtStringCodeUnits(s)); err != nil {
			return
		}
		ret = s
	case List:
		if err = accounter.PushDepth(); err != nil {
			return
		}
		defer accounter.PopDepth()
		if err = accounter.AccountBytes(36); err != nil {
			return
		}
		var id int8
		err = binary.Read(r, Order, &id)
		if err != nil {
			return
		}
		var l int32
		err = binary.Read(r, Order, &l)
		if err != nil {
			return
		}
		if l < 0 {
			err = &NbtFormatException{Message: "ListTag length cannot be negative: " + fmt.Sprint(l)}
			return
		}
		if id == End && l > 0 {
			err = &NbtFormatException{Message: "Missing type on ListTag"}
			return
		}
		if err = accounter.AccountBytesN(4, int64(l)); err != nil {
			return
		}
		if l == 0 {
			ret = []any{}
			return
		}
		var rett []any
		for range l {
			var result any
			result, err = readPayload(id, r, accounter)
			if err != nil {
				return
			}
			rett = append(rett, result)
		}
		ret = rett
	case Compound:
		if err = accounter.PushDepth(); err != nil {
			return
		}
		defer accounter.PopDepth()
		if err = accounter.AccountBytes(48); err != nil {
			return
		}
		rett := map[string]any{}
		for {
			id := int8(0)
			err = binary.Read(r, Order, &id)
			if err != nil {
				return
			}
			if id == End {
				break
			}
			if err = accounter.AccountBytes(28); err != nil {
				return
			}
			name := ""
			name, err = readNbtString(r)
			if err != nil {
				return
			}
			if err = accounter.AccountBytesN(2, nbtStringCodeUnits(name)); err != nil {
				return
			}
			if _, exists := rett[name]; !exists {
				if err = accounter.AccountBytes(36); err != nil {
					return
				}
			}
			rett[name], err = readPayload(id, r, accounter)
			if err != nil {
				return
			}
		}
		ret = rett
	case IntArray:
		if err = accounter.AccountBytes(24); err != nil {
			return
		}
		var l int32
		err = binary.Read(r, Order, &l)
		if err != nil {
			return
		}
		if l < 0 {
			err = errors.New("bad array length")
			return
		}
		if err = accounter.AccountBytesN(4, int64(l)); err != nil {
			return
		}
		var rett []int32
		for range l {
			var element int32
			err = binary.Read(r, Order, &element)
			if err != nil {
				return
			}
			rett = append(rett, element)
		}
		ret = rett
	case LongArray:
		if err = accounter.AccountBytes(24); err != nil {
			return
		}
		var l int32
		err = binary.Read(r, Order, &l)
		if err != nil {
			return
		}
		if l < 0 {
			err = errors.New("bad array length")
			return
		}
		if err = accounter.AccountBytesN(8, int64(l)); err != nil {
			return
		}
		var rett []int64
		for range l {
			var element int64
			err = binary.Read(r, Order, &element)
			if err != nil {
				return
			}
			rett = append(rett, element)
		}
		ret = rett
	default:
		err = errors.New("unknown nbt type id")
		return
	}
	return
}

// readNbtString reads a length-prefixed NBT string.
//
// Java's DataInput.writeUTF/readUTF use modified UTF-8 (CESU-8): U+0000 is the
// two-byte C0 80 and supplementary code points are encoded as UTF-16 surrogate
// pairs (six bytes) rather than Go's four-byte UTF-8. This package stores the
// raw bytes as a Go string, so ASCII round-trips exactly but other encodings
// must be checked before interoperating with Java NBT.
func readNbtString(r io.Reader) (s string, err error) {
	var l uint16
	if err = binary.Read(r, Order, &l); err != nil {
		return
	}
	var data []byte
	data, err = io.ReadAll(io.LimitReader(r, int64(l)))
	if err != nil {
		return
	}
	if len(data) != int(l) {
		err = BadStringLengthError
		return
	}
	return string(data), nil
}

// nbtStringCodeUnits is the UTF-16 code-unit count Java uses for NbtAccounter
// string accounting.
func nbtStringCodeUnits(s string) int64 {
	return int64(len(utf16.Encode([]rune(s))))
}

var UnknownTagTypeErr = errors.New("unknown nbt type")

func writePayload(v any, w io.Writer) (err error) {
	switch v := v.(type) {
	case struct{}:
	case int8, int16, int32, int64, float32, float64:
		err = binary.Write(w, Order, v)
	case []int8:
		err = binary.Write(w, Order, int32(len(v)))
		if err != nil {
			return
		}
		err = binary.Write(w, Order, v)
	case []int32:
		err = binary.Write(w, Order, int32(len(v)))
		if err != nil {
			return
		}
		err = binary.Write(w, Order, v)
	case []int64:
		err = binary.Write(w, Order, int32(len(v)))
		if err != nil {
			return
		}
		err = binary.Write(w, Order, v)
	case string:
		// Raw Go string bytes; see readNbtString for the Java modified UTF-8
		// (CESU-8) interop caveat.
		err = binary.Write(w, Order, uint16(len(v)))
		if err != nil {
			return
		}
		err = binary.Write(w, Order, []byte(v))
	case []any:
		id := End
		if len(v) != 0 {
			id, err = getId(v[0])
			if err != nil {
				return
			}
		}
		err = binary.Write(w, Order, id)
		if err != nil {
			return
		}
		err = binary.Write(w, Order, int32(len(v)))
		for _, v := range v {
			err = writePayload(v, w)
			if err != nil {
				return
			}
		}
	case map[string]any:
		for k, v := range v {
			var id int8
			id, err = getId(v)
			if err != nil {
				return
			}
			err = binary.Write(w, Order, id)
			if err != nil {
				return
			}
			err = writePayload(k, w)
			if err != nil {
				return
			}
			err = writePayload(v, w)
			if err != nil {
				return
			}
		}
		err = binary.Write(w, Order, End)
	default:
		err = UnknownTagTypeErr
	}
	return
}
