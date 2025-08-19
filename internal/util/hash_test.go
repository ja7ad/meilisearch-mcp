package util

import (
	"encoding/hex"
	"math/rand"
	"testing"
	"time"

	"golang.org/x/crypto/blake2b"
)

func TestHashVectors(t *testing.T) {
	cases := []string{
		"",
		"a",
		"abc",
		"The quick brown fox jumps over the lazy dog",
		string(make([]byte, 0)),
		func() string {
			b := make([]byte, 1024)
			for i := range b {
				b[i] = 'a'
			}
			return string(b)
		}(),
	}

	for _, in := range cases {
		wantBytes := blake2b.Sum256([]byte(in))
		want := hex.EncodeToString(wantBytes[:])
		got := Hash(in)
		if got != want {
			// Include first 16 chars for brevity on failure
			t.Fatalf("hash mismatch for input %q\nwant: %s\n got: %s", in, want, got)
		}
		if len(got) != 64 {
			t.Fatalf("expected hex length 64, got %d (input %q)", len(got), in)
		}
	}
}

func TestHashDeterministic(t *testing.T) {
	const in = "deterministic-input"
	first := Hash(in)
	for i := 0; i < 10; i++ {
		if got := Hash(in); got != first {
			t.Fatalf("hash not deterministic: iteration %d got %s first %s", i, got, first)
		}
	}
}

func TestHashLargeInput(t *testing.T) {
	// 1 MiB random data
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 1<<20)
	if _, err := r.Read(b); err != nil {
		t.Fatalf("rand read: %v", err)
	}
	_ = Hash(string(b)) // Just ensure it doesn't panic / allocate excessively.
}
