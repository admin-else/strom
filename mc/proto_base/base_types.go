package proto_base

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

// UnimplementedErr is returned when a code path is not yet implemented.
var (
	UnimplementedErr = errors.New("unimplemented")
)

// Direction indicates whether a packet is sent to the server or to the client.
type Direction int

const (
	ToServer Direction = iota
	ToClient
)

// String returns a human-readable representation of the Direction.
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

// Opposite returns the reverse Direction.
func (d Direction) Opposite() Direction {
	return d ^ 1
}

// State represents the Minecraft connection state.
type State int32

const (
	Handshaking State = iota
	Status
	Login
	Configuration
	Play
)

// String returns a human-readable representation of the State.
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

// Actor represents a participant in a Minecraft connection.
type Actor int

const (
	// Servee is the one being served (e.g. the backend server in a proxy).
	Servee Actor = iota
	Client
)

// String returns a human-readable representation of the Actor.
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

// SendDirection returns the Direction in which this Actor sends packets.
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

// ReceiveDirection returns the Direction from which this Actor receives packets.
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
	// ToDo is a placeholder type for unimplemented packet fields.
	ToDo struct{}
	// RestBuffer represents the remaining unread bytes in a packet buffer.
	RestBuffer []byte
)

// ToDoError is returned when a ToDo placeholder is used.
var ToDoError = errors.New("to do")

// BadTypeError is returned when an unexpected type is encountered.
var BadTypeError = errors.New("bad type")

// Encode writes the RestBuffer to w.
func (b *RestBuffer) Encode(w io.Writer) (err error) {
	n, err := w.Write(*b)
	if n != len(*b) {
		panic("fucking write not writing all")
	}
	return
}

// Decode reads all remaining bytes from r into the RestBuffer.
func (b *RestBuffer) Decode(r io.Reader) (err error) {
	*b, err = io.ReadAll(r)
	return
}

// EncodeString encodes a length-prefixed UTF-8 string in Minecraft varint format.
func EncodeString(w io.Writer, s string) (err error) {
	err = EncodeVarInt(w, int32(len(s)))
	if err != nil {
		return
	}
	_, err = w.Write([]byte(s))
	return
}

// DecodeString decodes a length-prefixed UTF-8 string in Minecraft varint format.
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

// Encode is a no-op stub that always returns ToDoError.
func (t ToDo) Encode(_ io.Writer) (err error) {
	err = ToDoError
	return
}

// Decode is a no-op stub that always returns ToDoError.
func (t ToDo) Decode(_ io.Reader) (ret ToDo, err error) {
	err = ToDoError
	return
}

// EncodeVarInt encodes a signed 32-bit integer in Minecraft varint format.
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

// DecodeVarInt decodes a signed 32-bit integer in Minecraft varint format.
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

// EncodeVarLong encodes a signed 64-bit integer in Minecraft varlong format.
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

// DecodeVarLong decodes a signed 64-bit integer in Minecraft varlong format.
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

// ErroringIndex looks up key k in map m and returns an error if the key is not found.
func ErroringIndex[K comparable, V any, M map[K]V](m M, i K) (v V, err error) {
	var ok bool
	v, ok = m[i]
	if !ok {
		err = errors.New("index not found")
	}
	return
}

// EncodeDecodeAble is implemented by types that can be encoded and decoded for
// Minecraft protocol packets.
type EncodeDecodeAble interface {
	Encode(w io.Writer) (err error)
	Decode(r io.ReadSeeker) (err error)
}

// PacketInfo holds metadata about a Minecraft protocol packet.
type PacketInfo struct {
	Type            EncodeDecodeAble
	Name            string
	Direction       Direction
	State           State
	PacketId        int32
	ProtocolVersion int32
}

// Bool2uint8 converts a bool to a uint8 (1 for true, 0 for false).
func Bool2uint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

// Bool2uint16 converts a bool to a uint16 (1 for true, 0 for false).
func Bool2uint16(b bool) uint16 {
	if b {
		return 1
	}
	return 0
}

// Bool2uint32 converts a bool to a uint32 (1 for true, 0 for false).
func Bool2uint32(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}

// Bool2uint64 converts a bool to a uint64 (1 for true, 0 for false).
func Bool2uint64(b bool) uint64 {
	if b {
		return 1
	}
	return 0
}

// LpVec3d is a compressed 3D velocity vector used in Minecraft packets.
type LpVec3d struct {
	X, Y, Z float64
}

// Encode encodes the LpVec3d as a compressed velocity triple into w.
func (v *LpVec3d) Encode(w io.Writer) error {
	return writeVelocity(w, *v)
}

// Decode decodes a compressed velocity triple from r into the LpVec3d.
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
