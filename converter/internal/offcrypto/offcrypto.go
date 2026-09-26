// Package offcrypto decrypts password-protected Office Open XML documents
// ([MS-OFFCRYPTO]). An encrypted .pptx, .docx or .xlsx is a compound file
// with two streams: EncryptionInfo, which says how the key is derived from
// the password, and EncryptedPackage, the encrypted ZIP package. Agile
// encryption (Office 2010 and later) and Standard encryption (Office 2007)
// with AES are supported.
package offcrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"hash"
	"io"
	"unicode/utf16"

	"github.com/shibukawa/bdf/converter/internal/cfb"
)

// ErrWrongPassword is returned when the password does not open the document.
var ErrWrongPassword = errors.New("offcrypto: wrong password")

// ErrIntegrity is returned (with the decrypted package) when the package
// fails the HMAC of Agile encryption.
var ErrIntegrity = errors.New("offcrypto: the encrypted package fails its integrity check")

const (
	infoStream    = "EncryptionInfo"
	packageStream = "EncryptedPackage"
)

// IsEncrypted reports whether r is an encrypted Office Open XML document: a
// compound file with an EncryptedPackage stream.
func IsEncrypted(r io.ReaderAt, size int64) bool {
	head := make([]byte, len(cfb.Signature))
	if _, err := r.ReadAt(head, 0); err != nil || !bytes.Equal(head, cfb.Signature) {
		return false
	}
	f, err := cfb.Open(r, size)
	return err == nil && f.Has(packageStream)
}

// Verify checks password against the document without decrypting the
// package: nil, ErrWrongPassword or an error about the file.
func Verify(r io.ReaderAt, size int64, password string) error {
	_, info, err := open(r, size)
	if err != nil {
		return err
	}
	_, err = info.secretKey(password)
	return err
}

// Decrypt returns the Office Open XML package (a ZIP file) of an encrypted
// document. When the package fails its integrity check it is returned with
// ErrIntegrity.
func Decrypt(r io.ReaderAt, size int64, password string) ([]byte, error) {
	f, info, err := open(r, size)
	if err != nil {
		return nil, err
	}
	key, err := info.secretKey(password)
	if err != nil {
		return nil, err
	}
	pkg, err := f.Stream(packageStream)
	if err != nil {
		return nil, fmt.Errorf("offcrypto: %s: %w", packageStream, err)
	}
	if len(pkg) < 8 {
		return nil, errors.New("offcrypto: encrypted package too short")
	}
	n := binary.LittleEndian.Uint64(pkg)
	if n > uint64(len(pkg)-8) {
		n &= 0xffffffff // some writers of Standard encryption leave the high half undefined
	}
	if n > uint64(len(pkg)-8) {
		return nil, errors.New("offcrypto: encrypted package shorter than its size")
	}
	integrity := info.verifyIntegrity(key, pkg) // over the stream as stored
	data := pkg[8:]
	data = data[:len(data)/aes.BlockSize*aes.BlockSize]
	if uint64(len(data)) < n {
		return nil, errors.New("offcrypto: encrypted package shorter than its size")
	}
	if err := info.decrypt(key, data); err != nil {
		return nil, err
	}
	if integrity != nil {
		return data[:n], integrity
	}
	return data[:n], nil
}

func open(r io.ReaderAt, size int64) (*cfb.File, encryption, error) {
	f, err := cfb.Open(r, size)
	if err != nil {
		return nil, nil, err
	}
	b, err := f.Stream(infoStream)
	if errors.Is(err, cfb.ErrNotExist) && f.Has(packageStream) {
		return nil, nil, errors.New("offcrypto: documents protected with rights management (IRM) are not supported")
	}
	if err != nil {
		return nil, nil, fmt.Errorf("offcrypto: %s: %w", infoStream, err)
	}
	info, err := parseInfo(b)
	if err != nil {
		return nil, nil, err
	}
	return f, info, nil
}

// encryption is one of the encryption methods.
type encryption interface {
	// secretKey derives the key the package is encrypted with, checking
	// the password against the verifier.
	secretKey(password string) ([]byte, error)
	// decrypt decrypts the package data (after its size) in place.
	decrypt(key, data []byte) error
	// verifyIntegrity checks the package stream against its HMAC, if any.
	verifyIntegrity(key, stream []byte) error
}

func parseInfo(b []byte) (encryption, error) {
	if len(b) < 8 {
		return nil, errors.New("offcrypto: EncryptionInfo too short")
	}
	major, minor := binary.LittleEndian.Uint16(b), binary.LittleEndian.Uint16(b[2:])
	switch {
	case major == 4 && minor == 4:
		return parseAgile(b[8:])
	case (major == 2 || major == 3 || major == 4) && minor == 2:
		return parseStandard(b[4:])
	case (major == 3 || major == 4) && minor == 3:
		return nil, errors.New("offcrypto: extensible encryption is not supported")
	}
	return nil, fmt.Errorf("offcrypto: unknown EncryptionInfo version %d.%d", major, minor)
}

// passwordBytes is the password as the key derivations hash it: UTF-16LE.
func passwordBytes(password string) []byte {
	u := utf16.Encode([]rune(password))
	b := make([]byte, 2*len(u))
	for i, c := range u {
		binary.LittleEndian.PutUint16(b[2*i:], c)
	}
	return b
}

// iterate hashes salt and password, then rehashes with an iterator spin
// times ([MS-OFFCRYPTO] 2.3.4.7 and 2.3.4.11).
func iterate(newHash func() hash.Hash, salt, password []byte, spin int) []byte {
	h := newHash()
	h.Write(salt)
	h.Write(password)
	sum := h.Sum(nil)
	var it [4]byte
	for i := range spin {
		binary.LittleEndian.PutUint32(it[:], uint32(i))
		h.Reset()
		h.Write(it[:])
		h.Write(sum)
		sum = h.Sum(sum[:0])
	}
	return sum
}

// fit truncates b to n bytes or pads it with pad.
func fit(b []byte, n int, pad byte) []byte {
	out := make([]byte, n)
	m := copy(out, b)
	for i := m; i < n; i++ {
		out[i] = pad
	}
	return out
}

// Agile encryption ([MS-OFFCRYPTO] 2.3.4.10–2.3.4.15).

// b64 is a base64 attribute of the Agile EncryptionInfo XML.
type b64 []byte

func (b *b64) UnmarshalText(text []byte) error {
	v, err := base64.StdEncoding.DecodeString(string(text))
	*b = v
	return err
}

type agileParams struct {
	SaltSize        int    `xml:"saltSize,attr"`
	BlockSize       int    `xml:"blockSize,attr"`
	KeyBits         int    `xml:"keyBits,attr"`
	HashSize        int    `xml:"hashSize,attr"`
	CipherAlgorithm string `xml:"cipherAlgorithm,attr"`
	CipherChaining  string `xml:"cipherChaining,attr"`
	HashAlgorithm   string `xml:"hashAlgorithm,attr"`
	SaltValue       b64    `xml:"saltValue,attr"`
}

type agile struct {
	KeyData       agileParams `xml:"keyData"`
	DataIntegrity struct {
		EncryptedHmacKey   b64 `xml:"encryptedHmacKey,attr"`
		EncryptedHmacValue b64 `xml:"encryptedHmacValue,attr"`
	} `xml:"dataIntegrity"`
	KeyEncryptors []struct {
		URI          string `xml:"uri,attr"`
		EncryptedKey struct {
			agileParams
			SpinCount                  int `xml:"spinCount,attr"`
			EncryptedVerifierHashInput b64 `xml:"encryptedVerifierHashInput,attr"`
			EncryptedVerifierHashValue b64 `xml:"encryptedVerifierHashValue,attr"`
			EncryptedKeyValue          b64 `xml:"encryptedKeyValue,attr"`
		} `xml:"encryptedKey"`
	} `xml:"keyEncryptors>keyEncryptor"`

	password int // index of the password key encryptor
}

const passwordEncryptor = "http://schemas.microsoft.com/office/2006/keyEncryptor/password"

// Block keys of the derivations.
var (
	blockVerifierInput = []byte{0xfe, 0xa7, 0xd2, 0x76, 0x3b, 0x4b, 0x9e, 0x79}
	blockVerifierValue = []byte{0xd7, 0xaa, 0x0f, 0x6d, 0x30, 0x61, 0x34, 0x4e}
	blockKeyValue      = []byte{0x14, 0x6e, 0x0b, 0xe7, 0xab, 0xac, 0xd0, 0xd6}
	blockHmacKey       = []byte{0x5f, 0xb2, 0xad, 0x01, 0x0c, 0xb9, 0xe1, 0xf6}
	blockHmacValue     = []byte{0xa0, 0x67, 0x7f, 0x02, 0xb2, 0x2c, 0x84, 0x33}
)

// maxSpin bounds the spin count a file can ask for; Office uses 100000.
const maxSpin = 10_000_000

func parseAgile(b []byte) (encryption, error) {
	a := &agile{password: -1}
	if err := xml.Unmarshal(b, a); err != nil {
		return nil, fmt.Errorf("offcrypto: EncryptionInfo: %w", err)
	}
	for i, k := range a.KeyEncryptors {
		if k.URI == passwordEncryptor {
			a.password = i
		}
	}
	if a.password < 0 {
		return nil, errors.New("offcrypto: no password key encryptor (certificate-protected documents are not supported)")
	}
	k := &a.KeyEncryptors[a.password].EncryptedKey
	for _, p := range []*agileParams{&a.KeyData, &k.agileParams} {
		if err := p.check(); err != nil {
			return nil, err
		}
	}
	if k.SpinCount < 0 || k.SpinCount > maxSpin {
		return nil, fmt.Errorf("offcrypto: spin count %d out of range", k.SpinCount)
	}
	return a, nil
}

func (p *agileParams) check() error {
	if p.CipherAlgorithm != "AES" {
		return fmt.Errorf("offcrypto: cipher %q is not supported", p.CipherAlgorithm)
	}
	if p.CipherChaining != "ChainingModeCBC" {
		return fmt.Errorf("offcrypto: cipher chaining %q is not supported", p.CipherChaining)
	}
	if p.KeyBits != 128 && p.KeyBits != 192 && p.KeyBits != 256 {
		return fmt.Errorf("offcrypto: key size %d is not supported", p.KeyBits)
	}
	if p.BlockSize != aes.BlockSize {
		return fmt.Errorf("offcrypto: block size %d is not supported", p.BlockSize)
	}
	if p.newHash() == nil {
		return fmt.Errorf("offcrypto: hash %q is not supported", p.HashAlgorithm)
	}
	if len(p.SaltValue) == 0 {
		return errors.New("offcrypto: missing salt")
	}
	return nil
}

func (p *agileParams) newHash() func() hash.Hash {
	switch p.HashAlgorithm {
	case "SHA1", "SHA-1":
		return sha1.New
	case "SHA256":
		return sha256.New
	case "SHA384":
		return sha512.New384
	case "SHA512":
		return sha512.New
	}
	return nil
}

// cbc decrypts data (a whole number of blocks) with key and iv.
func cbc(key, iv, data []byte) ([]byte, error) {
	if len(data)%aes.BlockSize != 0 {
		return nil, errors.New("offcrypto: encrypted value is not a whole number of blocks")
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(data))
	cipher.NewCBCDecrypter(c, iv).CryptBlocks(out, data)
	return out, nil
}

func (a *agile) secretKey(password string) ([]byte, error) {
	k := &a.KeyEncryptors[a.password].EncryptedKey
	newHash := k.newHash()
	h := iterate(newHash, k.SaltValue, passwordBytes(password), k.SpinCount)
	derive := func(block []byte) []byte {
		d := newHash()
		d.Write(h)
		d.Write(block)
		return fit(d.Sum(nil), k.KeyBits/8, 0x36)
	}
	iv := fit(k.SaltValue, k.BlockSize, 0x36)
	input, err := cbc(derive(blockVerifierInput), iv, k.EncryptedVerifierHashInput)
	if err != nil {
		return nil, err
	}
	value, err := cbc(derive(blockVerifierValue), iv, k.EncryptedVerifierHashValue)
	if err != nil {
		return nil, err
	}
	d := newHash()
	d.Write(input[:min(k.SaltSize, len(input))])
	want := d.Sum(nil)
	if len(value) < len(want) || subtle.ConstantTimeCompare(want, value[:len(want)]) != 1 {
		return nil, ErrWrongPassword
	}
	key, err := cbc(derive(blockKeyValue), iv, k.EncryptedKeyValue)
	if err != nil {
		return nil, err
	}
	if len(key) < a.KeyData.KeyBits/8 {
		return nil, errors.New("offcrypto: encrypted key too short")
	}
	return key[:a.KeyData.KeyBits/8], nil
}

// ivFor is the IV of a package segment or an integrity value: the hash of
// the key data salt and a block key, fitted to the block size.
func (a *agile) ivFor(block []byte) []byte {
	d := a.KeyData.newHash()()
	d.Write(a.KeyData.SaltValue)
	d.Write(block)
	return fit(d.Sum(nil), a.KeyData.BlockSize, 0x36)
}

// segmentSize is the length of the independently chained package segments.
const segmentSize = 4096

func (a *agile) decrypt(key, data []byte) error {
	c, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	var idx [4]byte
	for i, off := 0, 0; off < len(data); i, off = i+1, off+segmentSize {
		seg := data[off:min(off+segmentSize, len(data))]
		binary.LittleEndian.PutUint32(idx[:], uint32(i))
		cipher.NewCBCDecrypter(c, a.ivFor(idx[:])).CryptBlocks(seg, seg)
	}
	return nil
}

func (a *agile) verifyIntegrity(key, stream []byte) error {
	di := &a.DataIntegrity
	if len(di.EncryptedHmacKey) == 0 || len(di.EncryptedHmacValue) == 0 {
		return nil
	}
	hk, err := cbc(key, a.ivFor(blockHmacKey), di.EncryptedHmacKey)
	if err != nil {
		return err
	}
	hv, err := cbc(key, a.ivFor(blockHmacValue), di.EncryptedHmacValue)
	if err != nil {
		return err
	}
	newHash := a.KeyData.newHash()
	size := newHash().Size()
	if len(hk) < size || len(hv) < size {
		return ErrIntegrity
	}
	m := hmac.New(newHash, hk[:size])
	m.Write(stream)
	if !hmac.Equal(m.Sum(nil), hv[:size]) {
		return ErrIntegrity
	}
	return nil
}

// Standard encryption ([MS-OFFCRYPTO] 2.3.4.5–2.3.4.9).

type standard struct {
	keyBits               int
	salt                  []byte
	encryptedVerifier     []byte
	encryptedVerifierHash []byte
}

const (
	flagCryptoAPI = 0x04
	flagAES       = 0x20
	standardSpin  = 50000
)

func parseStandard(b []byte) (encryption, error) {
	le := binary.LittleEndian
	if len(b) < 8 {
		return nil, errors.New("offcrypto: EncryptionInfo too short")
	}
	flags, hsize := le.Uint32(b), int(le.Uint32(b[4:]))
	b = b[8:]
	if hsize < 32 || hsize > len(b) {
		return nil, errors.New("offcrypto: bad encryption header size")
	}
	hdr, v := b[:hsize], b[hsize:]
	if flags&flagCryptoAPI == 0 || flags&flagAES == 0 {
		return nil, errors.New("offcrypto: standard encryption without AES is not supported")
	}
	algID, algIDHash, keyBits := le.Uint32(hdr[8:]), le.Uint32(hdr[12:]), int(le.Uint32(hdr[16:]))
	switch algID {
	case 0, 0x660e, 0x660f, 0x6610:
	default:
		return nil, fmt.Errorf("offcrypto: algorithm 0x%x is not supported", algID)
	}
	if algIDHash != 0 && algIDHash != 0x8004 {
		return nil, fmt.Errorf("offcrypto: hash algorithm 0x%x is not supported", algIDHash)
	}
	if keyBits != 128 && keyBits != 192 && keyBits != 256 {
		return nil, fmt.Errorf("offcrypto: key size %d is not supported", keyBits)
	}
	if len(v) < 4+16+16+4+32 || le.Uint32(v) != 16 {
		return nil, errors.New("offcrypto: bad encryption verifier")
	}
	return &standard{keyBits: keyBits, salt: v[4:20], encryptedVerifier: v[20:36], encryptedVerifierHash: v[40:72]}, nil
}

func (s *standard) secretKey(password string) ([]byte, error) {
	h := iterate(sha1.New, s.salt, passwordBytes(password), standardSpin)
	d := sha1.New()
	d.Write(h)
	d.Write([]byte{0, 0, 0, 0}) // block 0
	final := d.Sum(nil)
	x := func(pad byte) []byte {
		buf := bytes.Repeat([]byte{pad}, 64)
		for i, c := range final {
			buf[i] ^= c
		}
		sum := sha1.Sum(buf)
		return sum[:]
	}
	key := append(x(0x36), x(0x5c)...)[:s.keyBits/8]

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	verifier := make([]byte, 16)
	ecb(c, verifier, s.encryptedVerifier)
	vh := make([]byte, 32)
	ecb(c, vh, s.encryptedVerifierHash)
	want := sha1.Sum(verifier)
	if subtle.ConstantTimeCompare(want[:], vh[:sha1.Size]) != 1 {
		return nil, ErrWrongPassword
	}
	return key, nil
}

func ecb(c cipher.Block, dst, src []byte) {
	for i := 0; i+aes.BlockSize <= len(src); i += aes.BlockSize {
		c.Decrypt(dst[i:i+aes.BlockSize], src[i:i+aes.BlockSize])
	}
}

func (s *standard) decrypt(key, data []byte) error {
	c, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	ecb(c, data, data)
	return nil
}

func (s *standard) verifyIntegrity(key, stream []byte) error { return nil }
