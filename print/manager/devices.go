package manager

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/yuruidb"
)

// A device enrolls with a one-time code an administrator made for a
// location: its page makes an ECDSA P-256 key that cannot be exported
// (WebCrypto, non-extractable, kept in IndexedDB) and sends the public
// key. Afterwards it signs a challenge for a session that lasts minutes,
// and signs again when the session ends. A URL or a cookie copied to
// another device does not let it in: the key stays on the device.
//
// The challenges and the sessions are signed with the manager's secret,
// not stored: any instance checks them.

// enrollAlphabet is Crockford's base 32: no I, L, O or U to misread.
const enrollAlphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

// EnrollCodeLength is the length of an enrollment code: 50 bits.
const EnrollCodeLength = 10

// Enrollment is a code that adds one device to a location.
type Enrollment struct {
	Code     string    `json:"code"`
	Location string    `json:"location"`
	Expires  time.Time `json:"expires"`
}

// NewEnrollment makes a one-time code for a location, valid for ttl
// (default 15 minutes).
func (m *Manager) NewEnrollment(ctx context.Context, loc string, ttl time.Duration) (*Enrollment, error) {
	if _, err := m.Location(ctx, loc); err != nil {
		return nil, err
	}
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}
	b := make([]byte, EnrollCodeLength)
	rand.Read(b)
	code := make([]byte, EnrollCodeLength)
	for i, c := range b {
		code[i] = enrollAlphabet[int(c)%len(enrollAlphabet)]
	}
	e := Enrollment{Code: string(code), Location: loc, Expires: m.now().Add(ttl)}
	if _, err := m.st.db.Insert(ctx, e.Code, enrollRecord{Location: loc, Expires: e.Expires}); err != nil {
		return nil, mapErr(err)
	}
	return &e, nil
}

// normalizeCode reads a code as typed: lower case, spaces and hyphens,
// and the letters Crockford's base 32 reads as digits.
func normalizeCode(s string) string {
	s = strings.ToUpper(strings.NewReplacer(" ", "", "-", "").Replace(s))
	return strings.NewReplacer("O", "0", "I", "1", "L", "1").Replace(s)
}

// Enroll adds a device with a code: the code is used up, and the device
// is of the code's location. key is the device's public key, SPKI DER in
// base64 (WebCrypto's exportKey("spki")).
func (m *Manager) Enroll(ctx context.Context, code, name, key string) (*Device, error) {
	pub, err := parseKey(key)
	if err != nil {
		return nil, err
	}
	code = normalizeCode(code)
	if len(code) != EnrollCodeLength {
		return nil, fmt.Errorf("%w: the code", ErrUnauthorized)
	}
	rec, err := m.st.db.Get[enrollRecord](ctx, code)
	if errors.Is(err, yuruidb.ErrNotFound) {
		return nil, fmt.Errorf("%w: the code", ErrUnauthorized)
	}
	if err != nil {
		return nil, mapErr(err)
	}
	// the code is used up first: of two devices with it, one gets in
	if err := m.st.db.Delete[enrollRecord](ctx, code, rec.Revision()); err != nil {
		if errors.Is(err, yuruidb.ErrConflict) || errors.Is(err, yuruidb.ErrNotFound) {
			return nil, fmt.Errorf("%w: the code", ErrUnauthorized)
		}
		return nil, mapErr(err)
	}
	if m.now().After(rec.Data.Expires) {
		return nil, fmt.Errorf("%w: the code expired", ErrUnauthorized)
	}
	id, err := yuruidb.NewID()
	if err != nil {
		return nil, err
	}
	name = clip(strings.TrimSpace(name), 80)
	if name == "" {
		name = "device " + id[len(id)-4:]
	}
	der, _ := x509.MarshalPKIXPublicKey(pub)
	d := Device{ID: id, Location: rec.Data.Location, Name: name, Key: base64.StdEncoding.EncodeToString(der), Enrolled: m.now()}
	if _, err := m.st.db.Insert(ctx, id, deviceRecord(d)); err != nil {
		return nil, mapErr(err)
	}
	if err := m.touch(ctx, d.Location, func(st *locationState) { st.Devices = append(st.Devices, id) }); err != nil {
		return nil, err
	}
	m.o.Log.Info("device enrolled", "location", d.Location, "device", id, "name", name)
	return &d, nil
}

func parseKey(key string) (*ecdsa.PublicKey, error) {
	der, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("%w: the key is not base64", ErrInvalid)
	}
	k, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, fmt.Errorf("%w: the key: %v", ErrInvalid, err)
	}
	pub, ok := k.(*ecdsa.PublicKey)
	if !ok || pub.Curve != elliptic.P256() {
		return nil, fmt.Errorf("%w: the key is not ECDSA P-256", ErrInvalid)
	}
	return pub, nil
}

// Device returns a device.
func (m *Manager) Device(ctx context.Context, id string) (*Device, error) {
	rec, err := m.st.db.Get[deviceRecord](ctx, id)
	if err != nil {
		return nil, mapErr(err)
	}
	d := Device(rec.Data)
	return &d, nil
}

// Devices returns the devices of a location, by when they enrolled.
func (m *Manager) Devices(ctx context.Context, loc string) ([]Device, error) {
	rec, err := m.st.db.Get[locationState](ctx, loc)
	if errors.Is(err, yuruidb.ErrNotFound) {
		return []Device{}, nil
	}
	if err != nil {
		return nil, mapErr(err)
	}
	out := make([]Device, 0, len(rec.Data.Devices))
	for _, id := range rec.Data.Devices {
		d, err := m.Device(ctx, id)
		if errors.Is(err, ErrNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}
		out = append(out, *d)
	}
	return out, nil
}

// RemoveDevice removes a device: its session ends at its next request.
func (m *Manager) RemoveDevice(ctx context.Context, id string) error {
	rec, err := m.st.db.Get[deviceRecord](ctx, id)
	if err != nil {
		return mapErr(err)
	}
	if err := m.st.db.Delete[deviceRecord](ctx, id, rec.Revision()); err != nil && !errors.Is(err, yuruidb.ErrNotFound) {
		return mapErr(err)
	}
	m.o.Log.Info("device removed", "location", rec.Data.Location, "device", id)
	return m.touch(ctx, rec.Data.Location, func(st *locationState) {
		st.Devices = slices.DeleteFunc(st.Devices, func(d string) bool { return d == id })
	})
}

// ChallengeFor is how long a challenge may be signed.
const ChallengeFor = 2 * time.Minute

// sign is the manager's MAC of a kind of token and its payload.
func (m *Manager) sign(kind string, payload []byte) []byte {
	h := hmac.New(sha256.New, m.o.Secret)
	h.Write([]byte("bdf-print-manager\x00" + kind + "\x00"))
	h.Write(payload)
	return h.Sum(nil)
}

var b64 = base64.RawURLEncoding

// token is kind.payload.mac, in base64url.
func (m *Manager) token(kind string, payload []byte) string {
	return b64.EncodeToString(payload) + "." + b64.EncodeToString(m.sign(kind, payload))
}

// open checks a token's MAC and returns its payload.
func (m *Manager) open(kind, tok string) ([]byte, bool) {
	p, mac, ok := strings.Cut(tok, ".")
	if !ok {
		return nil, false
	}
	payload, err1 := b64.DecodeString(p)
	sum, err2 := b64.DecodeString(mac)
	if err1 != nil || err2 != nil || !hmac.Equal(sum, m.sign(kind, payload)) {
		return nil, false
	}
	return payload, true
}

// Challenge returns a challenge for a device to sign: it expires after
// ChallengeFor.
func (m *Manager) Challenge() string {
	payload := make([]byte, 8+16)
	binary.BigEndian.PutUint64(payload, uint64(m.now().Add(ChallengeFor).Unix()))
	rand.Read(payload[8:])
	return m.token("challenge", payload)
}

// Authenticate checks a device's signature of a challenge (ECDSA P-256
// with SHA-256, as WebCrypto makes it: r and s of 32 bytes each, or
// ASN.1) and returns a session for it, and when it ends.
func (m *Manager) Authenticate(ctx context.Context, device, challenge, signature string) (string, time.Time, error) {
	payload, ok := m.open("challenge", challenge)
	if !ok || len(payload) != 24 || m.now().Unix() > int64(binary.BigEndian.Uint64(payload)) {
		return "", time.Time{}, fmt.Errorf("%w: the challenge", ErrUnauthorized)
	}
	d, err := m.Device(ctx, device)
	if errors.Is(err, ErrNotFound) {
		return "", time.Time{}, fmt.Errorf("%w: no such device", ErrUnauthorized)
	}
	if err != nil {
		return "", time.Time{}, err
	}
	pub, err := parseKey(d.Key)
	if err != nil {
		return "", time.Time{}, err
	}
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		sig, err = b64.DecodeString(signature)
	}
	if err != nil || !verify(pub, []byte(challenge), sig) {
		return "", time.Time{}, fmt.Errorf("%w: the signature", ErrUnauthorized)
	}
	until := m.now().Add(m.o.SessionFor)
	session := m.token("session", []byte(d.ID+"|"+d.Location+"|"+strconv.FormatInt(until.Unix(), 10)))
	m.seen(ctx, d)
	return session, until, nil
}

func verify(pub *ecdsa.PublicKey, data, sig []byte) bool {
	sum := sha256.Sum256(data)
	if len(sig) == 64 {
		r, s := new(big.Int).SetBytes(sig[:32]), new(big.Int).SetBytes(sig[32:])
		return ecdsa.Verify(pub, sum[:], r, s)
	}
	return ecdsa.VerifyASN1(pub, sum[:], sig)
}

// Session returns the device of a session: one not expired, of a device
// still enrolled at the same location.
func (m *Manager) Session(ctx context.Context, session string) (*Device, error) {
	payload, ok := m.open("session", session)
	if !ok {
		return nil, ErrUnauthorized
	}
	parts := strings.Split(string(payload), "|")
	if len(parts) != 3 {
		return nil, ErrUnauthorized
	}
	until, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil || m.now().Unix() > until {
		return nil, ErrUnauthorized
	}
	d, err := m.Device(ctx, parts[0])
	if errors.Is(err, ErrNotFound) {
		return nil, ErrUnauthorized
	}
	if err != nil {
		return nil, err
	}
	if d.Location != parts[1] {
		return nil, ErrUnauthorized
	}
	m.seen(ctx, d)
	return d, nil
}

// seenEvery bounds the writes of LastSeen: a device polls every few
// seconds.
const seenEvery = 30 * time.Second

// seen records that a device is connected.
func (m *Manager) seen(ctx context.Context, d *Device) {
	now := m.now()
	if now.Sub(d.LastSeen) < seenEvery {
		return
	}
	_, err := change(ctx, m.st.db, d.ID, false, func(r *deviceRecord) error {
		if now.Sub(r.LastSeen) < seenEvery {
			return errSkip
		}
		r.LastSeen = now
		return nil
	})
	if err != nil {
		m.o.Log.Warn("device last seen", "device", d.ID, "err", err)
		return
	}
	d.LastSeen = now
}
