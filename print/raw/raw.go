// Package raw sends documents to printers that do not speak IPP: over the
// raw socket of port 9100 (AppSocket, JetDirect), the line printer
// daemon protocol (LPD, RFC 1179), a serial port (and a Bluetooth printer
// of the serial port profile, which the OS shows as one), the device file
// of a USB printer (Linux and the BSDs), and a queue of Windows whose
// driver passes RAW data on (winspool). The document is a file of the
// printer's language, which package print makes: PostScript or PDF, or
// the language of a printer of labels, receipts or forms (ZPL, SBPL,
// TPCL, ESC/POS, ESC/P; package print/pdl). The printer is told nothing
// else of the job but what the file says and an optional PJL header (the
// copies, sides and paper of PostScript are in the file, those of PDF in
// PJL).
//
// None of them tells the state of a job: Send ends when the printer has
// the whole document. None has a way to ask the printer what it
// supports, so the format sent is told by the URL, as in
// socket://host:9100?format=pdf or serial:///dev/ttyUSB0?lang=escpos.
//
// A job that was sent in part may have been printed in part: from the
// first byte of a 9100 job, and from the control file of an LPD job (an
// LPD printer prints a job once it has its control file, which is sent
// after the data). Send calls Begin then, and an error after it is
// ErrPartial, which a spool does not send again on its own.
package raw

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"slices"
	"strings"
	"sync/atomic"
	"time"
)

// Formats of a raw printer: what package print writes for it.
const (
	FormatPostScript = "ps"
	FormatPDF        = "pdf"
)

// languages are the languages of printers of labels, receipts and forms
// (print/pdl's), which a raw printer may take as its format.
var languages = map[string]string{
	"zpl": "zpl", "zplii": "zpl", "sbpl": "sbpl", "tpcl": "tpcl",
	"escpos": "escpos", "esc/pos": "escpos", "escp": "escp", "esc/p": "escp", "escpj": "escp", "esc/p-j": "escp", "escp2": "escp",
}

// IsLanguage reports whether a format is the language of a printer of
// labels, receipts or forms.
func IsLanguage(format string) bool {
	_, ok := languages[format]
	return ok && languages[format] == format
}

// Printer is a printer reached over 9100, LPD, a serial port, a USB
// device file or a queue of Windows, as its URL says:
//
//	socket://host[:9100][?format=…][&pjl=true|false]
//	lpd://host[:515]/queue[?format=…][&pjl=true|false][&user=name]
//	serial:///dev/ttyUSB0?format=…[&baud=9600][&bits=8][&parity=none][&stop=1][&flow=none]
//	serial://COM3?format=…                       (Windows)
//	usb:///dev/usb/lp0?format=…                 (Linux, the BSDs)
//	winspool:///Printer%20Name?format=…         (Windows)
//
// The format is ps (the default) or pdf, or a language of print/pdl
// (zpl, sbpl, tpcl, escpos, escp; lang is the same parameter), which
// takes a profile, dpi and width too.
type Printer struct {
	// URL is the printer's URL, without its query.
	URL string
	// Scheme is "socket", "lpd", "serial", "usb" or "winspool".
	Scheme string
	// Addr is the host and port.
	Addr string
	// Queue is the queue of an LPD printer.
	Queue string
	// Device is the serial port or device file, or the name of the
	// queue of Windows.
	Device string
	// Serial is the setting of a serial port.
	Serial Serial
	// Format is the format the printer takes: FormatPostScript (the
	// default), FormatPDF, or a language of print/pdl; "" for the
	// language of the Profile, which the URL names without one.
	Format string
	// Profile, DPI and Width describe a printer of a language of
	// print/pdl: a profile's name, and its dots and width as the URL
	// overrides them ("" for the profile's).
	Profile, DPI, Width string
	// PJL wraps the document in a PJL header: "true", "false", or "" for
	// PDF only (which many printers do not tell from other files without
	// it).
	PJL string
	// User is the user an LPD job is of ("bdf" when "").
	User string
	// Timeout bounds the connection and each write and read: 30 seconds
	// when 0.
	Timeout time.Duration
}

// Parse reads the URL of a raw printer.
func Parse(s string) (*Printer, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("raw: %w", err)
	}
	p := &Printer{Scheme: strings.ToLower(u.Scheme), Format: FormatPostScript}
	port := "9100"
	options := []string{"format", "lang", "pjl", "profile", "dpi", "width"}
	switch p.Scheme {
	case "socket":
		if strings.Trim(u.Path, "/") != "" {
			return nil, fmt.Errorf("raw: %s: a socket printer has no path", s)
		}
	case "lpd":
		port = "515"
		p.Queue = strings.Trim(u.Path, "/")
		if p.Queue == "" || strings.ContainsAny(p.Queue, "/ \t\n\x00") || len(p.Queue) > 255 {
			return nil, fmt.Errorf("raw: %s: an LPD printer is lpd://host/queue", s)
		}
		options = append(options, "user")
	case "serial", "usb":
		// serial:///dev/ttyUSB0, serial://COM3
		p.Device = u.Path
		if u.Host != "" {
			p.Device = u.Host + u.Path
		}
		if p.Device == "" || strings.ContainsRune(p.Device, 0) {
			return nil, fmt.Errorf("raw: %s: no device (%s:///dev/… or %s://COM3)", s, p.Scheme, p.Scheme)
		}
		if p.Scheme == "serial" {
			options = append(options, "baud", "bits", "parity", "stop", "flow")
		}
	case "winspool":
		// winspool:///Queue%20Name, winspool://Queue (a name a host may
		// be), and winspool://server/queue for \\server\queue
		switch {
		case u.Host == "":
			p.Device = strings.TrimPrefix(u.Path, "/")
		case u.Path == "" || u.Path == "/":
			p.Device = u.Host
		default:
			p.Device = `\\` + u.Host + strings.ReplaceAll(u.Path, "/", `\`)
		}
		if p.Device == "" || strings.ContainsRune(p.Device, 0) {
			return nil, fmt.Errorf("raw: %s: no queue (winspool:///Queue%%20Name)", s)
		}
	default:
		return nil, fmt.Errorf("raw: %s: not a socket://, lpd://, serial://, usb:// or winspool:// URL", s)
	}
	if p.Scheme == "socket" || p.Scheme == "lpd" {
		if u.Hostname() == "" {
			return nil, fmt.Errorf("raw: %s: no host", s)
		}
		if u.Port() != "" {
			port = u.Port()
		}
		p.Addr = net.JoinHostPort(u.Hostname(), port)
	}
	q := u.Query()
	for k := range q {
		if !slices.Contains(options, k) {
			return nil, fmt.Errorf("raw: %s: %q is not an option (%s)", s, k, strings.Join(options, ", "))
		}
	}
	format := strings.ToLower(q.Get("format"))
	if l := strings.ToLower(q.Get("lang")); l != "" {
		if format != "" && languages[format] != languages[l] {
			return nil, fmt.Errorf("raw: %s: format %q and lang %q", s, format, l)
		}
		format = l
	}
	switch format {
	case "", "ps", "postscript":
	case "pdf":
		p.Format = FormatPDF
	default:
		l, ok := languages[format]
		if !ok {
			return nil, fmt.Errorf("raw: %s: format %q: ps, pdf, zpl, sbpl, tpcl, escpos or escp", s, format)
		}
		p.Format = l
	}
	p.Profile, p.DPI, p.Width = q.Get("profile"), q.Get("dpi"), q.Get("width")
	if (p.Profile != "" || p.DPI != "" || p.Width != "") && !IsLanguage(p.Format) {
		if format != "" || p.Profile == "" {
			return nil, fmt.Errorf("raw: %s: profile, dpi and width are those of a printer of labels, receipts or forms (lang=zpl …)", s)
		}
		p.Format = "" // the profile's language
	}
	switch v := strings.ToLower(q.Get("pjl")); v {
	case "", "true", "false":
		p.PJL = v
	default:
		return nil, fmt.Errorf("raw: %s: pjl %q: true or false", s, v)
	}
	p.User = q.Get("user")
	if p.Scheme == "serial" {
		if p.Serial, err = parseSerial(q); err != nil {
			return nil, fmt.Errorf("raw: %s: %w", s, err)
		}
	}
	u.RawQuery, u.Fragment = "", ""
	p.URL = u.String()
	return p, nil
}

// IsURL tells a URL of a raw printer (socket://, lpd://, serial://, usb://
// or winspool://).
func IsURL(s string) bool {
	l := strings.ToLower(s)
	for _, p := range []string{"socket://", "lpd://", "serial://", "usb://", "winspool://"} {
		if strings.HasPrefix(l, p) {
			return true
		}
	}
	return false
}

// UsesPJL tells whether a document in a format is sent with a PJL header.
func (p *Printer) UsesPJL(format string) bool {
	if format != FormatPDF && format != FormatPostScript {
		return false // PJL is that of printers of pages
	}
	switch p.PJL {
	case "true":
		return true
	case "false":
		return false
	}
	return format == FormatPDF
}

func (p *Printer) timeout() time.Duration {
	if p.Timeout > 0 {
		return p.Timeout
	}
	return 30 * time.Second
}

// Errors of Send.
var (
	// ErrPartial wraps an error after the printer may have begun to print
	// the job: it is not known whether, or how much of it, it printed.
	ErrPartial = errors.New("raw: the printer may have printed a part of the job")
	// ErrRefused is returned when an LPD printer refuses the job (or its
	// queue).
	ErrRefused = errors.New("raw: the printer refused the job")
)

// Job is a document to send.
type Job struct {
	// Name is the job's name (PJL and LPD), and User its user (LPD).
	Name string
	// Format is the document's: FormatPostScript, FormatPDF or a
	// language of print/pdl.
	Format string
	// Data is the document.
	Data []byte
	// PJL are the lines of the PJL header that set the job (@PJL SET
	// DUPLEX=ON …), without the lines that begin and end it, which Send
	// writes when the printer uses PJL for the format.
	PJL []string
	// Begin is called before what the printer may print of the job is
	// sent; an error stops the job there.
	Begin func() error
}

// Send sends a document to the printer, and returns the bytes sent.
func (p *Printer) Send(ctx context.Context, j *Job) (int64, error) {
	data := j.Data
	if p.UsesPJL(j.Format) {
		data = withPJL(j, data)
	}
	switch p.Scheme {
	case "socket":
		return p.sendSocket(ctx, j, data)
	case "lpd":
		return p.sendLPD(ctx, j, data)
	case "serial", "usb", "winspool":
		return p.sendDevice(ctx, j, data)
	}
	return 0, fmt.Errorf("raw: scheme %q", p.Scheme)
}

// uel is PJL's Universal Exit Language: it ends what the printer was
// reading, and begins PJL.
const uel = "\x1b%-12345X"

// withPJL wraps a document in a PJL job: its name, the settings, and the
// language it is in.
func withPJL(j *Job, data []byte) []byte {
	name := pjlString(j.Name)
	lang := "POSTSCRIPT"
	if j.Format == FormatPDF {
		lang = "PDF"
	}
	var b strings.Builder
	b.WriteString(uel + "@PJL\r\n")
	fmt.Fprintf(&b, "@PJL JOB NAME=\"%s\"\r\n", name)
	for _, l := range j.PJL {
		b.WriteString(l + "\r\n")
	}
	fmt.Fprintf(&b, "@PJL ENTER LANGUAGE=%s\r\n", lang)
	out := make([]byte, 0, b.Len()+len(data)+64)
	out = append(out, b.String()...)
	out = append(out, data...)
	out = append(out, fmt.Sprintf("%s@PJL EOJ NAME=\"%s\"\r\n%s", uel, name, uel)...)
	return out
}

// pjlString makes a string a PJL string holds: printable ASCII, without
// quotes, of 80 characters at most.
func pjlString(s string) string {
	var b strings.Builder
	for _, r := range s {
		if b.Len() >= 80 {
			break
		}
		if r < 0x20 || r > 0x7e || r == '"' {
			r = '_'
		}
		b.WriteRune(r)
	}
	return b.String()
}

// conn dials the printer, and closes the connection when the context is
// done.
func (p *Printer) conn(ctx context.Context) (net.Conn, func(), error) {
	d := net.Dialer{Timeout: p.timeout()}
	c, err := d.DialContext(ctx, "tcp", p.Addr)
	if err != nil {
		return nil, nil, err
	}
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			c.SetDeadline(time.Now())
		case <-done:
		}
	}()
	return c, func() { close(done); c.Close() }, nil
}

// Check tells whether the printer takes connections: it connects and
// closes the connection, sending nothing.
func (p *Printer) Check(ctx context.Context) error {
	switch p.Scheme {
	case "serial", "usb", "winspool":
		return p.checkDevice()
	}
	_, closer, err := p.conn(ctx)
	if err != nil {
		return err
	}
	closer()
	return nil
}

// sendSocket sends a document over the raw socket: all of it, then the end
// of the stream, then it waits a while for the printer to close the
// connection (what it says back is read and dropped).
func (p *Printer) sendSocket(ctx context.Context, j *Job, data []byte) (int64, error) {
	c, closer, err := p.conn(ctx)
	if err != nil {
		return 0, err
	}
	defer closer()
	// what the printer says back (PJL status, errors of PostScript) is
	// read, so that it never stops for want of a reader
	read := make(chan struct{})
	var readErr error
	go func() {
		_, readErr = io.Copy(io.Discard, c)
		close(read)
	}()
	if j.Begin != nil {
		if err := j.Begin(); err != nil {
			return 0, err
		}
	}
	n, err := p.write(ctx, c, data)
	if err != nil {
		return n, fmt.Errorf("%w: %w", ErrPartial, ctxErr(ctx, err))
	}
	if tc, ok := c.(interface{ CloseWrite() error }); ok {
		tc.CloseWrite()
		// a printer closes the connection once it has the job, or once it
		// printed it; a connection closed before that could lose its end
		select {
		case <-read:
			if readErr != nil && ctx.Err() == nil {
				// a reset: the printer dropped what it had not read
				return n, fmt.Errorf("%w: the printer closed the connection before it had the whole job: %w", ErrPartial, readErr)
			}
		case <-time.After(socketLinger):
		case <-ctx.Done():
		}
	}
	return n, nil
}

// socketLinger is how long the end of a 9100 job waits for the printer to
// close the connection.
var socketLinger = 10 * time.Second

// write writes data in blocks, each within the timeout.
func (p *Printer) write(ctx context.Context, c net.Conn, data []byte) (int64, error) {
	var n int64
	for len(data) > 0 {
		if err := ctx.Err(); err != nil {
			return n, err
		}
		k := min(len(data), 64<<10)
		c.SetWriteDeadline(time.Now().Add(p.timeout()))
		w, err := c.Write(data[:k])
		n += int64(w)
		if err != nil {
			return n, err
		}
		data = data[k:]
	}
	return n, nil
}

// ctxErr is the error of a context that ended a write, or the error.
func ctxErr(ctx context.Context, err error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return err
}

// jobNumbers make the numbers of LPD jobs, 0 to 999.
var jobNumbers atomic.Uint32

func init() { jobNumbers.Store(uint32(os.Getpid()) % 1000) }

// sendLPD sends a document as an LPD job: the receive-job command, the
// data file, then the control file that has the printer print it, each
// acknowledged.
func (p *Printer) sendLPD(ctx context.Context, j *Job, data []byte) (int64, error) {
	c, closer, err := p.conn(ctx)
	if err != nil {
		return 0, err
	}
	defer closer()
	r := bufio.NewReader(c)
	ack := func(what string) error {
		c.SetReadDeadline(time.Now().Add(p.timeout()))
		b, err := r.ReadByte()
		if err != nil {
			return ctxErr(ctx, err)
		}
		if b != 0 {
			return fmt.Errorf("%w: %s (%d)", ErrRefused, what, b)
		}
		return nil
	}
	send := func(b []byte) error {
		_, err := p.write(ctx, c, b)
		return ctxErr(ctx, err)
	}
	host, _ := os.Hostname()
	host = lpdName(host, 31)
	if host == "" {
		host = "bdf"
	}
	num := fmt.Sprintf("%03d", jobNumbers.Add(1)%1000)
	dfName, cfName := "dfA"+num+host, "cfA"+num+host
	user := lpdName(p.User, 31)
	if user == "" {
		user = "bdf"
	}
	name := lpdName(j.Name, 99)
	control := fmt.Sprintf("H%s\nP%s\nJ%s\nN%s\nl%s\nU%s\n", host, user, name, name, dfName, dfName)

	if err := send([]byte("\x02" + p.Queue + "\n")); err != nil {
		return 0, err
	}
	if err := ack("the queue " + p.Queue); err != nil {
		return 0, err
	}
	if err := send([]byte(fmt.Sprintf("\x03%d %s\n", len(data), dfName))); err != nil {
		return 0, err
	}
	if err := ack("the data file"); err != nil {
		return 0, err
	}
	if err := send(append(data[:len(data):len(data)], 0)); err != nil {
		return 0, err
	}
	if err := ack("the data"); err != nil {
		return int64(len(data)), err
	}
	n := int64(len(data))
	if err := send([]byte(fmt.Sprintf("\x02%d %s\n", len(control), cfName))); err != nil {
		return n, err
	}
	if err := ack("the control file"); err != nil {
		return n, err
	}
	// from here the printer may print the job
	if j.Begin != nil {
		if err := j.Begin(); err != nil {
			return n, err
		}
	}
	if err := send(append([]byte(control), 0)); err != nil {
		return n, fmt.Errorf("%w: %w", ErrPartial, err)
	}
	if err := ack("the control"); err != nil {
		if errors.Is(err, ErrRefused) {
			return n, err
		}
		return n, fmt.Errorf("%w: %w", ErrPartial, err)
	}
	return n, nil
}

// lpdName makes a name an LPD control file holds: printable ASCII without
// spaces, of n characters at most.
func lpdName(s string, n int) string {
	var b strings.Builder
	for _, r := range s {
		if b.Len() >= n {
			break
		}
		if r <= 0x20 || r > 0x7e {
			r = '_'
		}
		b.WriteRune(r)
	}
	return b.String()
}

// PJLPaper is the PJL name of a paper of a PWG media name, "" for papers
// PJL does not name.
func PJLPaper(media string) string {
	for _, p := range []struct{ prefix, name string }{
		{"iso_a3_", "A3"}, {"iso_a4_", "A4"}, {"iso_a5_", "A5"}, {"iso_b5_", "B5"}, {"jis_b4_", "JISB4"}, {"jis_b5_", "JISB5"},
		{"na_letter_", "LETTER"}, {"na_legal_", "LEGAL"}, {"na_executive_", "EXECUTIVE"}, {"na_ledger_", "LEDGER"},
		{"jpn_hagaki_", "HAGAKI"},
	} {
		if strings.HasPrefix(media, p.prefix) {
			return p.name
		}
	}
	return ""
}
