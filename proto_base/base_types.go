package proto_base

import (
	"encoding/binary"
	"errors"
	"io"
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

type State int

const (
	Handshaking State = iota
	Status
	Login
	Configuration
	Play
)

type Actor int

const (
	Server Actor = iota
	Client
)

func (a Actor) SendDirection() Direction {
	switch a {
	case Server:
		return ToClient
	case Client:
		return ToServer
	default:
		panic("invalid actor")
	}
}

func (a Actor) ReceiveDirection() Direction {
	switch a {
	case Server:
		return ToServer
	case Client:
		return ToClient
	default:
		panic("invalid actor")
	}
}

type PacketIdentifier struct {
	Direction Direction
	State     State
	Name      string
}

type (
	ToDo       struct{}
	Void       struct{}
	RestBuffer []byte
)

func (v Void) Encode(_ io.Writer) (err error) {
	return
}

func (v Void) Decode(_ io.Reader) (ret Void, err error) {
	return
}

var ToDoError = errors.New("to do")
var BadTypeError = errors.New("bad type")

func (r RestBuffer) Encode(w io.Writer) (err error) {
	_, err = w.Write(r)
	return
}

func (RestBuffer) Decode(r io.Reader) (ret RestBuffer, err error) {
	b, err := io.ReadAll(r)
	ret = b
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
		err := binary.Write(w, binary.BigEndian, b)
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

func ErroringIndex[K comparable, V any, M map[K]V](m M, i K) (v V, err error) {
	var ok bool
	v, ok = m[i]
	if !ok {
		err = errors.New("index not found")
	}
	return
}

type EncodeDecodeAble[Self any] interface {
	Encode(w io.Writer) (err error)
	Decode(r io.Reader) (ret Self, err error)
}
