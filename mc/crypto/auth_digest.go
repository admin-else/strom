package crypto

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

func twosComplement(p []byte) []byte {
	carry := true
	for i := len(p) - 1; i >= 0; i-- {
		p[i] = byte(^p[i])
		if carry {
			carry = p[i] == 0xff
			p[i]++
		}
	}
	return p
}

// AuthDigest stolen from https://gist.github.com/toqueteos/5372776
func AuthDigest(elems ...[]byte) string {
	h := sha1.New()
	for _, elem := range elems {
		h.Write(elem)
	}
	hash := h.Sum(nil)

	negative := (hash[0] & 0x80) == 0x80
	if negative {
		hash = twosComplement(hash)
	}

	res := strings.TrimLeft(hex.EncodeToString(hash), "0")
	if negative {
		res = "-" + res
	}

	return res
}
