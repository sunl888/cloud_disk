package hasher

import (
	"encoding/hex"
	"golang.org/x/crypto/argon2"
)

type Argon2Hasher struct {
	salt    []byte
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

// Hash the given value.
func (b *Argon2Hasher) Make(value string) string {
	return hex.EncodeToString(argon2.Key([]byte(value), b.salt, b.time, b.memory, b.threads, b.keyLen))
}

// Check the given plain value against a hash.
func (b *Argon2Hasher) Check(value string, hashedValue string) bool {
	return b.Make(value) == hashedValue
}

func NewArgon2Hasher(salt []byte, time, memory uint32, threads uint8, keyLen uint32) *Argon2Hasher {
	return &Argon2Hasher{salt, time, memory, threads, keyLen}
}
