package nbt

import (
	"math"
	"strconv"
)

// NbtAccounter mirrors net.minecraft.nbt.NbtAccounter. It bounds the heap and
// stack depth a decode may consume, protecting callers from hostile NBT.
const (
	NbtAccounterDEFAULT_NBT_QUOTA      = 2097152
	NbtAccounterUNCOMPRESSED_NBT_QUOTA = 104857600
	NbtAccounterMAX_STACK_DEPTH        = 512
)

// NbtAccounterException mirrors net.minecraft.nbt.NbtAccounterException.
type NbtAccounterException struct {
	Message string
}

func (e *NbtAccounterException) Error() string { return e.Message }

// NbtFormatException mirrors net.minecraft.nbt.NbtFormatException.
type NbtFormatException struct {
	Message string
}

func (e *NbtFormatException) Error() string { return e.Message }

// NbtAccounter mirrors net.minecraft.nbt.NbtAccounter.
type NbtAccounter struct {
	quota    int64
	usage    int64
	maxDepth int
	depth    int
}

// NewNbtAccounter mirrors the NbtAccounter(long, int) constructor.
func NewNbtAccounter(quota int64, maxDepth int) *NbtAccounter {
	return &NbtAccounter{quota: quota, maxDepth: maxDepth}
}

// NbtAccounterCreate mirrors NbtAccounter.create(long).
func NbtAccounterCreate(quota int64) *NbtAccounter {
	return NewNbtAccounter(quota, NbtAccounterMAX_STACK_DEPTH)
}

// NbtAccounterDefaultQuota mirrors NbtAccounter.defaultQuota().
func NbtAccounterDefaultQuota() *NbtAccounter {
	return NewNbtAccounter(NbtAccounterDEFAULT_NBT_QUOTA, NbtAccounterMAX_STACK_DEPTH)
}

// NbtAccounterUncompressedQuota mirrors NbtAccounter.uncompressedQuota().
func NbtAccounterUncompressedQuota() *NbtAccounter {
	return NewNbtAccounter(NbtAccounterUNCOMPRESSED_NBT_QUOTA, NbtAccounterMAX_STACK_DEPTH)
}

// NbtAccounterUnlimitedHeap mirrors NbtAccounter.unlimitedHeap().
func NbtAccounterUnlimitedHeap() *NbtAccounter {
	return NewNbtAccounter(math.MaxInt64, NbtAccounterMAX_STACK_DEPTH)
}

// AccountBytesN mirrors NbtAccounter.accountBytes(long, long).
func (a *NbtAccounter) AccountBytesN(bytesPerEntry int64, count int64) error {
	return a.AccountBytes(bytesPerEntry * count)
}

// AccountBytes mirrors NbtAccounter.accountBytes(long).
func (a *NbtAccounter) AccountBytes(size int64) (err error) {
	if size < 0 {
		return &NbtAccounterException{Message: "Tried to account NBT tag with negative size: " + strconv.FormatInt(size, 10)}
	}
	if a.usage+size > a.quota {
		return &NbtAccounterException{
			Message: "Tried to read NBT tag that was too big; tried to allocate: " + strconv.FormatInt(a.usage, 10) + " + " + strconv.FormatInt(size, 10) + " bytes where max allowed: " + strconv.FormatInt(a.quota, 10),
		}
	}
	a.usage += size
	return nil
}

// PushDepth mirrors NbtAccounter.pushDepth().
func (a *NbtAccounter) PushDepth() (err error) {
	if a.depth >= a.maxDepth {
		return &NbtAccounterException{Message: "Tried to read NBT tag with too high complexity, depth > " + strconv.Itoa(a.maxDepth)}
	}
	a.depth++
	return nil
}

// PopDepth mirrors NbtAccounter.popDepth().
func (a *NbtAccounter) PopDepth() (err error) {
	if a.depth <= 0 {
		return &NbtAccounterException{Message: "NBT-Accounter tried to pop stack-depth at top-level"}
	}
	a.depth--
	return nil
}

// Usage mirrors NbtAccounter.getUsage().
func (a *NbtAccounter) Usage() int64 { return a.usage }

// Depth mirrors NbtAccounter.getDepth().
func (a *NbtAccounter) Depth() int { return a.depth }
