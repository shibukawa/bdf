package bdf

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/binary"
	"errors"
	"fmt"

	"golang.org/x/text/unicode/norm"
)

// Encryption (docs/spec.md §3.5): every part is sealed on its own with
// AES-256-GCM under a random content key, so parts can still be fetched one
// by one. The content key is stored wrapped (AES-KW) with a key derived from
// a password. The manifest the document is read with is itself a sealed part;
// what is stored in the clear is an outer manifest that names the sealed parts
// by the hash of their stored bytes, which says nothing about their content.

// Encryption values.
const (
	// CipherA256GCM seals parts with AES-256-GCM; the stored bytes are the
	// 12-byte nonce followed by the ciphertext and the 16-byte tag.
	CipherA256GCM = "A256GCM"
	// KeyPassword is a key slot that a password unlocks.
	KeyPassword = "password"
	// KDFPBKDF2SHA256 derives key-encryption keys with PBKDF2-HMAC-SHA-256
	// from the password normalized to NFC and encoded as UTF-8.
	KDFPBKDF2SHA256 = "PBKDF2-SHA256"
	// DefaultIterations is the PBKDF2 iteration count of new password slots.
	DefaultIterations = 600_000
	// MaxIterations bounds the iteration count readers accept, so a file
	// cannot keep a reader busy for long.
	MaxIterations = 10_000_000

	// FlagEncrypted is the single-file header flag of encrypted documents.
	FlagEncrypted = 1

	nonceSize = 12
	keySize   = 32
	saltSize  = 16
)

// manifestAAD is the additional data the manifest is sealed with; parts use
// their hash.
var manifestAAD = []byte("manifest")

// ErrWrongPassword is returned by Reader.Unlock when no key slot opens with
// the password.
var ErrWrongPassword = errors.New("bdf: wrong password")

// ErrLocked is returned when the parts of an encrypted document are read
// before Reader.Unlock.
var ErrLocked = errors.New("bdf: the document is encrypted; unlock it with its password")

// Encryption is the "encryption" member of an encrypted document's outer
// manifest.
type Encryption struct {
	Cipher string    `json:"cipher"`
	Keys   []KeySlot `json:"keys"`
	// Manifest is the sealed part that holds the document's manifest.
	Manifest SealedManifest `json:"manifest"`
}

// KeySlot holds the content key wrapped with a key derived from a password.
type KeySlot struct {
	Type string `json:"type"`
	KDF  string `json:"kdf"`
	Iter int    `json:"iter"`
	Salt []byte `json:"salt"`
	// Key is the content key wrapped with AES-KW (RFC 3394).
	Key []byte `json:"key"`
}

// SealedManifest locates the sealed manifest.
type SealedManifest struct {
	Part Hash `json:"part"`
	// Enc is the encoding of the manifest JSON inside the seal.
	Enc string `json:"enc"`
}

// A Lock encrypts documents: a random content key and the key slots that
// unlock it. Set it as Document.Lock before writing.
type Lock struct {
	key   []byte
	slots []KeySlot
}

// NewPasswordLock makes a lock with a fresh content key that password
// unlocks. iter is the PBKDF2 iteration count (0: DefaultIterations).
func NewPasswordLock(password string, iter int) (*Lock, error) {
	if password == "" {
		return nil, errors.New("bdf: empty password")
	}
	if iter == 0 {
		iter = DefaultIterations
	}
	if iter < 1 || iter > MaxIterations {
		return nil, fmt.Errorf("bdf: iteration count %d out of range", iter)
	}
	key := make([]byte, keySize)
	salt := make([]byte, saltSize)
	rand.Read(key)
	rand.Read(salt)
	kek, err := passwordKEK(password, salt, iter)
	if err != nil {
		return nil, err
	}
	wrapped, err := wrapKey(kek, key)
	if err != nil {
		return nil, err
	}
	return &Lock{key: key, slots: []KeySlot{{Type: KeyPassword, KDF: KDFPBKDF2SHA256, Iter: iter, Salt: salt, Key: wrapped}}}, nil
}

func (l *Lock) aead() (cipher.AEAD, error) { return newAEAD(l.key) }

// unlockKey returns the content key a password unwraps from one of the slots.
func unlockKey(e *Encryption, password string) ([]byte, error) {
	if e.Cipher != CipherA256GCM {
		return nil, &FormatError{Msg: fmt.Sprintf("unknown cipher %q", e.Cipher)}
	}
	for _, s := range e.Keys {
		if s.Type != KeyPassword {
			continue
		}
		if s.KDF != KDFPBKDF2SHA256 {
			return nil, &FormatError{Msg: fmt.Sprintf("unknown key derivation %q", s.KDF)}
		}
		if s.Iter < 1 || s.Iter > MaxIterations {
			return nil, &FormatError{Msg: fmt.Sprintf("iteration count %d out of range", s.Iter)}
		}
		kek, err := passwordKEK(password, s.Salt, s.Iter)
		if err != nil {
			return nil, err
		}
		if key, err := unwrapKey(kek, s.Key); err == nil && len(key) == keySize {
			return key, nil
		}
	}
	return nil, ErrWrongPassword
}

func passwordKEK(password string, salt []byte, iter int) ([]byte, error) {
	return pbkdf2.Key(sha256.New, norm.NFC.String(password), salt, iter, keySize)
}

func newAEAD(key []byte) (cipher.AEAD, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(b)
}

// seal encrypts plain under a fresh random nonce: nonce, ciphertext, tag.
func seal(a cipher.AEAD, plain, aad []byte) []byte {
	out := make([]byte, nonceSize, nonceSize+len(plain)+a.Overhead())
	rand.Read(out)
	return a.Seal(out, out, plain, aad)
}

// open decrypts the output of seal.
func open(a cipher.AEAD, sealed, aad []byte) ([]byte, error) {
	if len(sealed) < nonceSize+a.Overhead() {
		return nil, &FormatError{Msg: "sealed part too short"}
	}
	plain, err := a.Open(nil, sealed[:nonceSize], sealed[nonceSize:], aad)
	if err != nil {
		return nil, &FormatError{Msg: "sealed part fails authentication"}
	}
	return plain, nil
}

// kwIV is the default initial value of RFC 3394 key wrap.
var kwIV = []byte{0xa6, 0xa6, 0xa6, 0xa6, 0xa6, 0xa6, 0xa6, 0xa6}

var errUnwrap = errors.New("bdf: key unwrap failed")

// wrapKey wraps key with kek (AES key wrap, RFC 3394 §2.2.1).
func wrapKey(kek, key []byte) ([]byte, error) {
	if len(key)%8 != 0 || len(key) < 16 {
		return nil, errors.New("bdf: bad key length to wrap")
	}
	b, err := aes.NewCipher(kek)
	if err != nil {
		return nil, err
	}
	n := len(key) / 8
	out := make([]byte, 8+len(key))
	copy(out, kwIV)
	copy(out[8:], key)
	var blk [16]byte
	for j := range 6 {
		for i := 1; i <= n; i++ {
			copy(blk[:8], out[:8])
			copy(blk[8:], out[8*i:8*i+8])
			b.Encrypt(blk[:], blk[:])
			t := uint64(n*j + i)
			binary.BigEndian.PutUint64(out[:8], binary.BigEndian.Uint64(blk[:8])^t)
			copy(out[8*i:], blk[8:])
		}
	}
	return out, nil
}

// unwrapKey reverses wrapKey, checking the integrity value.
func unwrapKey(kek, wrapped []byte) ([]byte, error) {
	if len(wrapped)%8 != 0 || len(wrapped) < 24 {
		return nil, errUnwrap
	}
	b, err := aes.NewCipher(kek)
	if err != nil {
		return nil, err
	}
	n := len(wrapped)/8 - 1
	a := make([]byte, 8)
	copy(a, wrapped[:8])
	r := make([]byte, 8*n)
	copy(r, wrapped[8:])
	var blk [16]byte
	for j := 5; j >= 0; j-- {
		for i := n; i >= 1; i-- {
			t := uint64(n*j + i)
			binary.BigEndian.PutUint64(blk[:8], binary.BigEndian.Uint64(a)^t)
			copy(blk[8:], r[8*(i-1):8*i])
			b.Decrypt(blk[:], blk[:])
			copy(a, blk[:8])
			copy(r[8*(i-1):], blk[8:])
		}
	}
	if subtle.ConstantTimeCompare(a, kwIV) != 1 {
		return nil, errUnwrap
	}
	return r, nil
}
