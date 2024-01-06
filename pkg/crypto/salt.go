package crypto

import (
	"crypto/rand"
	"encoding/hex"
)

const LenSalt = 8

func NewSalt() string {
	b := make([]byte, hex.DecodedLen(LenSalt))
	rand.Read(b)

	saltHex := make([]byte, LenSalt)
	hex.Encode(saltHex, b)
	return string(saltHex)
}
