package bdf

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/hkdf"
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
// a password, or agreed with ECDH with a reader's public key. The manifest
// the document is read with is itself a sealed part; what is stored in the
// clear is an outer manifest that names the sealed parts by the hash of their
// stored bytes, which says nothing about their content.

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
	// cannot keep a reader busy for long: that of a key slot, and the sum of
	// those of the slots a reader tries.
	MaxIterations = 10_000_000

	// KeyECDH is a key slot that one reader's private key unlocks: the
	// content key is wrapped with a key agreed with ECDH between the
	// reader's public key and a key the writer makes for the document and
	// forgets (docs/spec.md §3.5).
	KeyECDH = "ecdh"
	// CurveP256 is the curve of ecdh slots (NIST P-256).
	CurveP256 = "P-256"
	// KDFHKDFSHA256 derives the key-encryption key of an ecdh slot from the
	// shared secret with HKDF-SHA-256.
	KDFHKDFSHA256 = "HKDF-SHA256"
	// MaxECDHSlots bounds the ecdh slots a reader tries: each costs a key
	// agreement.
	MaxECDHSlots = 16

	// FlagEncrypted is the single-file header flag of encrypted documents.
	FlagEncrypted = 1

	nonceSize = 12
	keySize   = 32
	saltSize  = 16
)

// manifestAAD is the additional data the manifest is sealed with; parts use
// their hash.
var manifestAAD = []byte("manifest")

// ecdhInfo is the HKDF info of ecdh slots.
const ecdhInfo = "bdf ecdh v1"

// ErrWrongPassword is returned by Reader.Unlock when no key slot opens with
// the password.
var ErrWrongPassword = errors.New("bdf: wrong password")

// ErrWrongKey is returned by Reader.UnlockECDH when no key slot opens with
// the private key.
var ErrWrongKey = errors.New("bdf: the document is not sealed for this key")

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

// KeySlot holds the content key wrapped with a key derived from a password
// (Iter and Salt) or agreed with ECDH (Crv and EPK).
type KeySlot struct {
	Type string `json:"type"`
	KDF  string `json:"kdf"`
	Iter int    `json:"iter,omitempty"`
	Salt []byte `json:"salt,omitempty"`
	Crv  string `json:"crv,omitempty"`
	// EPK is the writer's public key of an ecdh slot, an uncompressed point.
	EPK []byte `json:"epk,omitempty"`
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

// NewECDHLock makes a lock with a fresh content key that only the private
// key of reader unlocks (a P-256 key). The key the content key is wrapped
// with is agreed with ECDH between reader and a key made here, which is
// forgotten once the lock is made: when the reader forgets its private key
// too, as a key pair made for one request does, nobody can open what the
// lock sealed any more.
func NewECDHLock(reader *ecdh.PublicKey) (*Lock, error) {
	if reader.Curve() != ecdh.P256() {
		return nil, errors.New("bdf: the reader's key is not a P-256 key")
	}
	eph, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	z, err := eph.ECDH(reader)
	if err != nil {
		return nil, err
	}
	epk := eph.PublicKey().Bytes()
	kek, err := ecdhKEK(z, epk, reader.Bytes())
	if err != nil {
		return nil, err
	}
	key := make([]byte, keySize)
	rand.Read(key)
	wrapped, err := wrapKey(kek, key)
	if err != nil {
		return nil, err
	}
	return &Lock{key: key, slots: []KeySlot{{Type: KeyECDH, KDF: KDFHKDFSHA256, Crv: CurveP256, EPK: epk, Key: wrapped}}}, nil
}

// ecdhKEK derives the key-encryption key of an ecdh slot from the shared
// secret, binding both public keys into it.
func ecdhKEK(z, epk, rpk []byte) ([]byte, error) {
	salt := sha256.Sum256(append(append([]byte{}, epk...), rpk...))
	return hkdf.Key(sha256.New, z, salt[:], ecdhInfo, keySize)
}

func (l *Lock) aead() (cipher.AEAD, error) { return newAEAD(l.key) }

// unlockKey returns the content key a password unwraps from one of the slots.
func unlockKey(e *Encryption, password string) ([]byte, error) {
	if e.Cipher != CipherA256GCM {
		return nil, &FormatError{Msg: fmt.Sprintf("unknown cipher %q", e.Cipher)}
	}
	spent := 0
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
		// each slot is within the range, but many of them are not
		if spent += s.Iter; spent > MaxIterations {
			return nil, &FormatError{Msg: fmt.Sprintf("the key slots take more than %d iterations", MaxIterations)}
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

// unlockKeyECDH returns the content key a private key unwraps from one of
// the ecdh slots.
func unlockKeyECDH(e *Encryption, priv *ecdh.PrivateKey) ([]byte, error) {
	if e.Cipher != CipherA256GCM {
		return nil, &FormatError{Msg: fmt.Sprintf("unknown cipher %q", e.Cipher)}
	}
	if priv.Curve() != ecdh.P256() {
		return nil, errors.New("bdf: the key is not a P-256 key")
	}
	rpk := priv.PublicKey().Bytes()
	tried := 0
	for _, s := range e.Keys {
		if s.Type != KeyECDH {
			continue
		}
		if tried++; tried > MaxECDHSlots {
			return nil, &FormatError{Msg: fmt.Sprintf("more than %d ecdh key slots", MaxECDHSlots)}
		}
		if s.Crv != CurveP256 {
			return nil, &FormatError{Msg: fmt.Sprintf("unknown curve %q", s.Crv)}
		}
		if s.KDF != KDFHKDFSHA256 {
			return nil, &FormatError{Msg: fmt.Sprintf("unknown key derivation %q", s.KDF)}
		}
		epk, err := ecdh.P256().NewPublicKey(s.EPK)
		if err != nil {
			return nil, &FormatError{Msg: "bad public key in an ecdh key slot"}
		}
		z, err := priv.ECDH(epk)
		if err != nil {
			return nil, err
		}
		kek, err := ecdhKEK(z, s.EPK, rpk)
		if err != nil {
			return nil, err
		}
		if key, err := unwrapKey(kek, s.Key); err == nil && len(key) == keySize {
			return key, nil
		}
	}
	return nil, ErrWrongKey
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
