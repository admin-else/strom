package crypto

import (
	"crypto/md5"

	"github.com/google/uuid"
)

// UUIDv3 Variation of google.uuid.NewHash
// We need this because we want a 14-byte namespace and Google forces 16
func UUIDv3(space []byte, data []byte) uuid.UUID {
	h := md5.New()
	h.Reset()
	h.Write(space) //nolint:errcheck
	h.Write(data)  //nolint:errcheck
	s := h.Sum(nil)
	var u uuid.UUID
	copy(u[:], s)
	u[6] = (u[6] & 0x0f) | uint8((3&0xf)<<4) // Version 3
	u[8] = (u[8] & 0x3f) | 0x80              // RFC 4122 variant
	return u
}

func FromOfflinePlayer(displayName string) uuid.UUID {
	return UUIDv3([]byte("OfflinePlayer:"), []byte(displayName))
}
