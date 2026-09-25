package bdf

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Hash identifies a part by content: the first 16 bytes of SHA-256.
type Hash [16]byte

// HashOf computes the part hash of data.
func HashOf(data []byte) Hash {
	sum := sha256.Sum256(data)
	var h Hash
	copy(h[:], sum[:16])
	return h
}

// String returns the 32-character lowercase hex form used in manifests and file names.
func (h Hash) String() string { return hex.EncodeToString(h[:]) }

// ParseHash parses the hex form.
func ParseHash(s string) (Hash, error) {
	var h Hash
	if len(s) != 32 {
		return h, fmt.Errorf("bdf: bad hash %q", s)
	}
	if _, err := hex.Decode(h[:], []byte(s)); err != nil {
		return h, fmt.Errorf("bdf: bad hash %q: %w", s, err)
	}
	return h, nil
}

// MarshalText implements encoding.TextMarshaler.
func (h Hash) MarshalText() ([]byte, error) { return []byte(h.String()), nil }

// UnmarshalText implements encoding.TextUnmarshaler.
func (h *Hash) UnmarshalText(b []byte) error {
	v, err := ParseHash(string(b))
	if err != nil {
		return err
	}
	*h = v
	return nil
}
