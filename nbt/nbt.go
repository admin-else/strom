package nbt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
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

type Tag struct {
	Name  string
	Value any
}

type Anon struct {
	Value any
}

func (s Anon) Decode(r io.Reader) (ret Anon, err error) {
	id := int8(0)
	err = binary.Read(r, Order, &id)
	if err != nil {
		return
	}
	s.Value, err = readPayload(id, r)
	if err != nil {
		return
	}
	ret = s
	return
}

func (s Anon) Encode(w io.Writer) (err error) {
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

func (s *Tag) Read(r io.Reader) (err error) {
	id := int8(0)
	err = binary.Read(r, Order, &id)
	if err != nil {
		return
	}
	var name any
	name, err = readPayload(String, r)
	if err != nil {
		return
	}
	s.Value, err = readPayload(id, r)
	if err != nil {
		return
	}
	s.Name = name.(string)
	return
}

func (s *Tag) Write(w io.Writer) (err error) {
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

func readPayload(id int8, r io.Reader) (ret any, err error) {
	switch id {
	case End:
		ret = struct{}{}
	case Byte:
		rett := int8(0)
		err = binary.Read(r, Order, &rett)
		ret = rett
	case Short:
		rett := int16(0)
		err = binary.Read(r, Order, &rett)
		ret = rett
	case Int:
		rett := int32(0)
		err = binary.Read(r, Order, &rett)
		ret = rett
	case Long:
		rett := int64(0)
		err = binary.Read(r, Order, &rett)
		ret = rett
	case Float:
		rett := float32(0)
		err = binary.Read(r, Order, &rett)
		ret = rett
	case Double:
		rett := float64(0)
		err = binary.Read(r, Order, &rett)
		ret = rett
	case ByteArray:
		var l int32
		err = binary.Read(r, Order, &l)
		if err != nil {
			return
		}
		if l < 0 {
			err = errors.New("bad array length")
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
		var l uint16
		err = binary.Read(r, Order, &l)
		if err != nil {
			return
		}
		var s []byte
		s, err = io.ReadAll(io.LimitReader(r, int64(l)))
		if err != nil {
			return
		}
		if len(s) != int(l) {
			err = errors.New("bad string length")
			return
		}
		ret = string(s)
	case List:
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
		if l <= 0 || id == End {
			ret = []any{}
			return
		}
		var rett []any
		for range l {
			var result any
			result, err = readPayload(id, r)
			if err != nil {
				return
			}
			rett = append(rett, result)
		}
		ret = rett
	case Compound:
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
			var namea any
			namea, err = readPayload(String, r)
			if err != nil {
				return
			}
			name := namea.(string)
			rett[name], err = readPayload(id, r)
			if err != nil {
				return
			}
		}
		ret = rett
	case IntArray:
		var l int32
		err = binary.Read(r, Order, &l)
		if err != nil {
			return
		}
		if l < 0 {
			err = errors.New("bad array length")
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
		var l int32
		err = binary.Read(r, Order, &l)
		if err != nil {
			return
		}
		if l < 0 {
			err = errors.New("bad array length")
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
		err = errors.New("unkown nbt type id")
		return
	}
	return
}

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
		err = errors.New("unkown nbt type")
	}
	return
}
