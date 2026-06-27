package proto_base

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
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

func (s State) String() string {
	switch s {
	case Handshaking:
		return "Handshaking"
	case Status:
		return "Status"
	case Login:
		return "Login"
	case Configuration:
		return "Configuration"
	case Play:
		return "Play"
	default:
		return "Unknown"
	}
}

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
	Decode(r io.ReadSeeker) (err error)
}

type Void struct{}

func (v Void) Encode(w io.Writer) (err error) {
	return
}
func (v Void) Decode(r io.ReadSeeker) (err error) {
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

type LpVec3d struct {
	X, Y, Z float64
}

func (v *LpVec3d) Encode(w io.Writer) error {
	return writeVelocity(w, *v)
}
func (v *LpVec3d) Decode(r io.ReadSeeker) error {
	v3, err := readVelocity(r)
	if err != nil {
		return err
	}
	v.X = v3.X
	v.Y = v3.Y
	v.Z = v3.Z
	return nil
}
func readVelocity(r io.Reader) (LpVec3d, error) {
	var marker, second uint8
	if err := binary.Read(r, binary.BigEndian, &marker); err != nil {
		return LpVec3d{}, err
	}
	if marker == 0 {
		return LpVec3d{}, nil
	}
	if err := binary.Read(r, binary.BigEndian, &second); err != nil {
		return LpVec3d{}, err
	}
	var packed uint32
	if err := binary.Read(r, binary.BigEndian, &packed); err != nil {
		return LpVec3d{}, err
	}
	m := (uint64(packed) << 16) | (uint64(second) << 8) | uint64(marker)
	n := int64(marker & 3)
	if (marker & 4) != 0 {
		// fast marker bit set, read varint for extended scale
		scale, err := DecodeVarInt(r)
		if err != nil {
			return LpVec3d{}, err
		}
		n |= (int64(scale) & 0xFFFFFFFF) << 2
	}
	x := fromLong(int64(m>>3), n)
	y := fromLong(int64(m>>18), n)
	z := fromLong(int64(m>>33), n)
	return LpVec3d{X: x, Y: y, Z: z}, nil
}
func writeVelocity(w io.Writer, v LpVec3d) error {
	const maxVal = 1.7179869183e10
	const scaleMax = 32766.0
	x := clamp(v.X, maxVal)
	y := clamp(v.Y, maxVal)
	z := clamp(v.Z, maxVal)
	g := absMax(absMax(x, y), z)
	if g < 3.051944088384301e-5 {
		_, err := w.Write([]byte{0})
		return err
	}
	l := int64(math.Ceil(g))
	bl := (l & 3) != l
	m := l
	if bl {
		m = (l & 3) | 4
	}
	lf := float64(l)
	n := toLong(x/lf) << 3
	o := toLong(y/lf) << 18
	p := toLong(z/lf) << 33
	q := m | n | o | p
	var buf [7]byte
	buf[0] = byte(q)
	buf[1] = byte(q >> 8)
	binary.BigEndian.PutUint32(buf[2:6], uint32(q>>16))
	if _, err := w.Write(buf[:6]); err != nil {
		return err
	}
	if bl {
		if err := EncodeVarInt(w, int32(l>>2)); err != nil {
			return err
		}
	}
	return nil
}
func clamp(v, max float64) float64 {
	if math.IsNaN(v) {
		return 0
	}
	if v > max {
		return max
	}
	if v < -max {
		return -max
	}
	return v
}
func absMax(a, b float64) float64 {
	if math.Abs(a) > math.Abs(b) {
		return a
	}
	return b
}
func toLong(v float64) int64 {
	return int64(math.Round((v*0.5 + 0.5) * 32766.0))
}
func fromLong(bits int64, scale int64) float64 {
	v := bits & 0x7FFF
	if v > 32766 {
		v = 32766
	}
	return (float64(v)*2.0/32766.0 - 1.0) * float64(scale)
}
