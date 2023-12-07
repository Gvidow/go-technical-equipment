package crypto

import (
	"encoding/hex"

	"golang.org/x/crypto/argon2"
)

const (
	_time    = 3
	_memory  = 32 * 1024
	_threads = 4
	_keyLen  = 32
)

func Hash(row string, salt string) string {
	rowHash := argon2.Key([]byte(row), []byte(salt), _time, _memory, _threads, _keyLen)

	rowHashHex := make([]byte, hex.EncodedLen(_keyLen))
	hex.Encode(rowHashHex, rowHash)

	return string(rowHashHex)
}
