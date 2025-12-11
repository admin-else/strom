package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/admin-else/queser"
)

func UnpackIndirectPalette(r io.Reader, bitsPerEntry uint8, numberOfEntries int) (err error) {
	var palletLen queser.VarInt
	palletLen, err = palletLen.Decode(r)
	if err != nil {
		return
	}
	var pallet []int32 // remember we can't use make because this is user-controlled data
	for range palletLen {
		var entry queser.VarInt
		entry, err = queser.VarInt(0).Decode(r)
		if err != nil {
			return
		}
		pallet = append(pallet, int32(entry))
	}
	entriesPerLong := int(64 / bitsPerEntry)
	numberOfLongs := int(math.Ceil(float64(numberOfEntries) / float64(entriesPerLong)))
	var data []uint64
	for range numberOfLongs {
		var entry uint64
		err = binary.Read(r, binary.BigEndian, &entry)
		if err != nil {
			return
		}
		data = append(data, entry)
	}
	fmt.Println(data)
	return
}

func GetBit(data []uint64, index, bpe int32) uint64 {
	// net.minecraft.util.collection.PackedIntegerArray
	return 1
}

func UnpackPalette(r io.Reader, numberOfEntries int) (err error) {
	var bitsPerEntry uint8
	err = binary.Read(r, binary.BigEndian, &bitsPerEntry)
	if err != nil {
		return
	}
	err = UnpackIndirectPalette(r, bitsPerEntry, numberOfEntries)
	return
}

func UnpackSection(r io.Reader) (err error) {
	var blockCount int16
	err = binary.Read(r, binary.BigEndian, &blockCount)
	if err != nil {
		return
	}
	err = UnpackPalette(r, 4096) // 4096 is the number of blocks in a chunk
	return
}

func mainE() (err error) {
	fmt.Println(packet.X*16, packet.Z*16)
	b := bytes.NewBuffer(packet.ChunkData.Val)
	err = UnpackSection(b)
	return
}

func main() {
	err := mainE()
	if err != nil {
		panic(err)
	}
}
