package util

import (
	"encoding/hex"

	"golang.org/x/crypto/blake2b"
)

// Hash returns the BLAKE2b-256 digest (hex encoded) of the provided string.
// It uses the optimized one-shot Sum256 function for speed and zero allocations.
func Hash(v string) string {
	sum := blake2b.Sum256([]byte(v))
	return hex.EncodeToString(sum[:])
}
