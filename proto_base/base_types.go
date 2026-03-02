package proto_base

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	UnimplementedErr = errors.New("unimplemented")
)

type Direction int

const (
	ToServer Direction = iota
	ToClient
)

func (d Direction) String() string {
	switch d {
	case ToServer:
		return "ToServer"
	case ToClient:
		return "ToClient"
	default:
		return "Unknown"
	}
}

func (d Direction) Opposite() Direction {
	return d ^ 1
}

type State int32

const (
	Handshaking State = iota
	Status
	Login
	Configuration
	Play
)

type Actor int

const (
	Servee Actor = iota // Servee is the one being served
	Client
)

func (a Actor) String() string {
	switch a {
	case Servee:
		return "Servee"
	case Client:
		return "Client"
	default:
		return "Unknown"
	}
}

func (a Actor) SendDirection() Direction {
	switch a {
	case Servee:
		return ToClient
	case Client:
		return ToServer
	default:
		panic("invalid actor")
	}
}

func (a Actor) ReceiveDirection() Direction {
	switch a {
	case Servee:
		return ToServer
	case Client:
		return ToClient
	default:
		panic("invalid actor")
	}
}

type (
	ToDo       struct{}
	RestBuffer []byte
)

var ToDoError = errors.New("to do")
var BadTypeError = errors.New("bad type")

func (b *RestBuffer) Encode(w io.Writer) (err error) {
	n, err := w.Write(*b)
	if n != len(*b) {
		panic("fucking write not writing all")
	}
	return
}

func (b *RestBuffer) Decode(r io.Reader) (err error) {
	*b, err = io.ReadAll(r)
	return
}

func EncodeString(w io.Writer, s string) (err error) {
	err = EncodeVarInt(w, int32(len(s)))
	if err != nil {
		return
	}
	_, err = w.Write([]byte(s))
	return
}

func DecodeString(r io.Reader) (ret string, err error) {
	l, err := DecodeVarInt(r)
	if err != nil {
		return
	}
	rawString, err := io.ReadAll(io.LimitReader(r, int64(l)))
	if err != nil {
		return
	}
	if len(rawString) != int(l) {
		err = errors.New("bad string length")
		return
	}
	return string(rawString), nil
}

func (t ToDo) Encode(_ io.Writer) (err error) {
	err = ToDoError
	return
}

func (t ToDo) Decode(_ io.Reader) (ret ToDo, err error) {
	err = ToDoError
	return
}

func EncodeVarInt(w io.Writer, v int32) (err error) {
	uv := uint32(v)
	for {
		b := uint8(uv & 0x7F)
		uv >>= 7
		if uv != 0 {
			b |= 0x80
		}
		err = binary.Write(w, binary.BigEndian, b)
		if err != nil {
			return err
		}
		if uv == 0 {
			break
		}
	}
	return nil
}

func DecodeVarInt(r io.Reader) (ret int32, err error) {
	for i := range 5 { // 32/7
		var b uint8 = 0
		err = binary.Read(r, binary.BigEndian, &b)
		if err != nil {
			return
		}
		ret |= int32(uint32(b&0x7F) << (7 * uint32(i)))
		if b&0x80 == 0 {
			return
		}
	}
	err = errors.New("VarInt too long")
	return
}

func EncodeVarLong(w io.Writer, v int64) (err error) {
	uv := uint64(v)
	for {
		b := uint8(uv & 0x7F)
		uv >>= 7
		if uv != 0 {
			b |= 0x80
		}
		err = binary.Write(w, binary.BigEndian, b)
		if err != nil {
			return err
		}
		if uv == 0 {
			break
		}
	}
	return nil
}

func DecodeVarLong(r io.Reader) (ret int64, err error) {
	for i := range 10 { // ceil(64/7)
		var b uint8
		err = binary.Read(r, binary.BigEndian, &b)
		if err != nil {
			return
		}
		ret |= int64(uint64(b&0x7F) << (7 * uint64(i)))
		if b&0x80 == 0 {
			return
		}
	}
	err = errors.New("VarInt too long")
	return
}

func ErroringIndex[K comparable, V any, M map[K]V](m M, i K) (v V, err error) {
	var ok bool
	v, ok = m[i]
	if !ok {
		err = errors.New("index not found")
	}
	return
}

type EncodeDecodeAble interface {
	Encode(w io.Writer) (err error)
	Decode(r io.Reader) (err error)
}

type LpVec3d struct {
	X, Y, Z float64
}

func (v LpVec3d) Encode(w io.Writer) (err error) {
	err = UnimplementedErr
	return
}

func (v LpVec3d) Decode(r io.Reader) (err error) {
	err = UnimplementedErr
	return
}

type Void struct{}

func (v Void) Encode(w io.Writer) (err error) {
	return
}
func (v Void) Decode(r io.Reader) (err error) {
	return
}

type PacketInfo struct {
	Type            EncodeDecodeAble
	Name            string
	Direction       Direction
	State           State
	PacketId        int32
	ProtocolVersion int32
}

func Bool2uint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func Bool2uint16(b bool) uint16 {
	if b {
		return 1
	}
	return 0
}

func Bool2uint32(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}

func Bool2uint64(b bool) uint64 {
	if b {
		return 1
	}
	return 0
}
