package gen

import (
	"crypto/md5"
	"encoding/binary"
)

func rotl64(x uint64, b uint8) uint64 {
	return (x << b) | (x >> (64 - b))
}

type Xoroshiro struct {
	L uint64
	H uint64
}

func (x *Xoroshiro) Next() uint64 {
	l := x.L
	h := x.H
	n := rotl64(l+h, 17) + l
	h ^= l
	x.L = rotl64(l, 49) ^ h ^ (h << 21)
	x.H = rotl64(h, 28)
	return n
}

func (x *Xoroshiro) NextBound(n uint32) int32 {
	r := (x.Next() & 0xFFFFFFFF) * uint64(n)
	if uint32(r) < n {
		for uint32(r) < (^n+1)%n {
			r = (x.Next() & 0xFFFFFFFF) * uint64(n)
		}
	}
	return int32(r >> 32)
}

func (x *Xoroshiro) NextFloat() float32 {
	return float32(x.Next()>>(64-24)) * 5.9604645e-8
}

func (x *Xoroshiro) NextFloat64() float64 {
	return float64(x.Next()>>(64-53)) * 1.1102230246251565e-16
}

// XoroshiroFromString creates an Xoroshiro from an MD5 hash of the given string.
func XoroshiroFromString(s string) (x Xoroshiro) {
	sum := md5.Sum([]byte(s))
	x.L = binary.BigEndian.Uint64(sum[:8])
	x.H = binary.BigEndian.Uint64(sum[8:])
	return
}

// XoroshiroFromSeed creates an Xoroshiro from an int64 seed.
func XoroshiroFromSeed(seed int64) (x Xoroshiro) {
	x.L = uint64(seed) ^ 0x6a09e667f3bcc909
	x.H = x.L + 0x9e3779b97f4a7c15
	return
}

func mixStafford13(s uint64) uint64 {
	s = (s ^ (s >> 30)) * 0xbf58476d1ce4e5b9
	s = (s ^ (s >> 27)) * 0x94d049bb133111eb
	return s ^ (s >> 31)
}

func (x *Xoroshiro) Mix() {
	x.L = mixStafford13(x.L)
	x.H = mixStafford13(x.H)
}

func (x *Xoroshiro) Xor(x2 Xoroshiro) {
	x.L ^= x2.L
	x.H ^= x2.H
}

func (x *Xoroshiro) XorString(s string) {
	x.Xor(XoroshiroFromString(s))
}

func (x *Xoroshiro) Split() (x2 Xoroshiro) {
	x2.L = x.Next()
	x2.H = x.Next()
	return
}
