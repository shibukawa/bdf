package fakeprinter

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shibukawa/bdf/print/ipp"
)

// Listener is a raw transport of a printer on a port: the raw port 9100
// or LPD.
type Listener struct {
	// URL is the printer's URL: socket://host:port or lpd://host:port/queue.
	URL  string
	p    *Printer
	kind string // "socket" or "lpd"
	addr string
	wg   sync.WaitGroup

	mu sync.Mutex
	ln net.Listener
}

// ListenSocket serves the printer on the raw port (9100) at addr
// ("127.0.0.1:0" for any free port of the loopback).
func (p *Printer) ListenSocket(addr string) (*Listener, error) {
	return p.listen("socket", addr, "socket://%s")
}

// ListenLPD serves the printer as an LPD queue (RFC 1179) at addr.
func (p *Printer) ListenLPD(addr, queue string) (*Listener, error) {
	return p.listen("lpd", addr, "lpd://%s/"+queue)
}

func (p *Printer) listen(kind, addr, url string) (*Listener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	l := &Listener{URL: fmt.Sprintf(url, ln.Addr()), p: p, kind: kind, addr: ln.Addr().String()}
	l.serve(ln)
	return l, nil
}

func (l *Listener) serve(ln net.Listener) {
	l.mu.Lock()
	l.ln = ln
	l.mu.Unlock()
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			l.wg.Add(1)
			go func() {
				defer l.wg.Done()
				defer c.Close()
				if l.kind == "lpd" {
					l.p.ServeLPD(c)
				} else {
					l.p.ServeSocket(c)
				}
			}()
		}
	}()
}

// Addr is the address the printer listens at.
func (l *Listener) Addr() string { return l.addr }

// Close stops the listener: connections to it are refused.
func (l *Listener) Close() {
	l.mu.Lock()
	ln := l.ln
	l.mu.Unlock()
	ln.Close()
	l.wg.Wait()
}

// Restart listens again at the address of a listener that was closed.
func (l *Listener) Restart() error {
	ln, err := net.Listen("tcp", l.addr)
	if err != nil {
		return err
	}
	l.serve(ln)
	return nil
}

// ServeSocket serves a connection of the raw port: everything sent until
// the sender closes its side is a job. A connection that sends nothing
// is a check, not a job.
func (p *Printer) ServeSocket(c net.Conn) {
	p.lock()
	latency := p.Latency
	p.mu.Unlock()
	if latency > 0 {
		time.Sleep(latency)
	}
	r := bufio.NewReader(p.throttle(c))
	if _, err := r.Peek(1); err != nil {
		return
	}
	p.lock()
	cut, max := p.cutAfter, p.MaxDocument
	p.cutAfter = -1
	p.mu.Unlock()
	if max <= 0 {
		max = 256 << 20
	}
	if cut >= 0 {
		io.CopyN(io.Discard, r, cut)
		// dropped: the connection is reset
		if tc, ok := c.(*net.TCPConn); ok {
			tc.SetLinger(0)
		}
		return
	}
	data, err := p.readRaw(r, c, max)
	if err != nil {
		return
	}
	p.lock()
	j := &Job{Transport: "socket", Documents: [][]byte{data}, Format: p.rawFormat(data)}
	p.add(j)
	p.done(j)
	back := p.back
	p.mu.Unlock()
	p.save(j)
	if back != "" {
		io.WriteString(c, back)
	}
}

// readRaw reads a raw job to its end, answering the status queries of
// the printer's language as they come.
func (p *Printer) readRaw(r io.Reader, w io.Writer, max int64) ([]byte, error) {
	var data bytes.Buffer
	buf := make([]byte, 32<<10)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			before := data.Len()
			data.Write(buf[:n])
			if int64(data.Len()) > max {
				return nil, fmt.Errorf("a job of more than %d bytes", max)
			}
			p.answer(data.Bytes()[max64(0, before-8):], w)
		}
		if err == io.EOF {
			return data.Bytes(), nil
		}
		if err != nil {
			return nil, err
		}
	}
}

func max64(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// answer answers the status queries in what was just read: ZPL's ~HS
// (the host status, three lines) and ESC/POS's DLE EOT n (a status byte).
// A printer stopped for media-empty says so.
func (p *Printer) answer(b []byte, w io.Writer) {
	p.lock()
	lang, stopped := p.Language, p.State == ipp.PrinterStopped
	reasons := strings.Join(p.Reasons, ",")
	p.mu.Unlock()
	empty := stopped && strings.Contains(reasons, "media-empty")
	open := stopped && strings.Contains(reasons, "door-open")
	switch lang {
	case "zpl":
		for range bytes.Count(b, []byte("~HS")) {
			paperOut, pause := 0, 0
			if empty {
				paperOut = 1
			}
			if stopped {
				pause = 1
			}
			head := 0
			if open {
				head = 1
			}
			fmt.Fprintf(w, "\x02030,%d,%d,1245,000,0,0,0,000,0,0,0\x03\r\n\x02000,0,%d,0,1,2,6,0,00000000,1,000\x03\r\n\x021234,0\x03\r\n", paperOut, pause, head)
		}
	case "escpos":
		for i := 0; i+2 < len(b); i++ {
			if b[i] != 0x10 || b[i+1] != 0x04 {
				continue
			}
			st := byte(0x12) // the fixed bits of every status byte
			switch b[i+2] {
			case 2: // the off-line status
				if open {
					st |= 0x04
				}
				if stopped {
					st |= 0x08
				}
			case 4: // the paper sensors
				if empty {
					st |= 0x60
				}
			}
			w.Write([]byte{st})
		}
	}
}

// rawFormat tells what a raw job is: the printer's language, or what its
// first bytes say.
func (p *Printer) rawFormat(data []byte) string {
	if f := Sniff(data); f != "" {
		return f
	}
	return p.Language
}

// save writes a raw job to the printer's Dir.
func (p *Printer) save(j *Job) {
	p.lock()
	dir := p.Dir
	p.mu.Unlock()
	if dir == "" {
		return
	}
	writeFile(dir, fmt.Sprintf("%d-1.%s", j.ID, extension(j.Format, j.Data())), j.Data())
}

// ServeLPD takes a receive-job command (RFC 1179) and its files on a
// connection.
func (p *Printer) ServeLPD(c net.Conn) {
	r := bufio.NewReader(p.throttle(c))
	line, err := r.ReadString('\n')
	if err != nil || len(line) < 2 {
		return
	}
	switch line[0] {
	case 2:
	case 3, 4:
		// the state of the queue: one line of it
		p.lock()
		n := p.queued()
		p.mu.Unlock()
		fmt.Fprintf(c, "%s is ready\n%d jobs queued\n", strings.TrimSpace(line[1:]), n)
		return
	default:
		c.Write([]byte{1})
		return
	}
	p.lock()
	cut, cutControl, refuse := p.cutAfter, p.cutControl, p.refuse || !p.Accepting
	p.cutAfter, p.cutControl = -1, false
	p.mu.Unlock()
	if refuse {
		c.Write([]byte{1})
		return
	}
	job := &Job{Transport: "lpd", Queue: strings.TrimSuffix(line[1:], "\n")}
	record := func() {
		p.lock()
		job.Format = p.rawFormat(job.Data())
		p.add(job)
		p.done(job)
		p.mu.Unlock()
		p.save(job)
	}
	c.Write([]byte{0})
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			// the sender is done: a job whose control file came first is
			// complete once its data files have
			if err == io.EOF && line == "" && job.Control != "" && len(job.Documents) > 0 {
				record()
			}
			return
		}
		if len(line) < 2 || line[0] != 2 && line[0] != 3 {
			c.Write([]byte{1})
			return
		}
		size, _, _ := strings.Cut(line[1:], " ")
		n, err := strconv.Atoi(size)
		if err != nil || n < 0 || n > 256<<20 {
			c.Write([]byte{1})
			return
		}
		if line[0] == 2 && cut >= 0 && len(job.Documents) > 0 {
			return
		}
		c.Write([]byte{0})
		if line[0] == 2 && cutControl {
			io.CopyN(io.Discard, r, int64(n)/2)
			return
		}
		buf := make([]byte, n+1)
		if _, err := io.ReadFull(r, buf); err != nil || buf[n] != 0 {
			return
		}
		if line[0] == 3 {
			job.Documents = append(job.Documents, buf[:n])
			c.Write([]byte{0})
			continue
		}
		job.Control = string(buf[:n])
		for _, l := range strings.Split(job.Control, "\n") {
			switch {
			case strings.HasPrefix(l, "P"):
				job.User = l[1:]
			case strings.HasPrefix(l, "J"):
				job.Name = l[1:]
			}
		}
		if len(job.Documents) > 0 {
			// the data came first: the control file ends the job (it is
			// kept before the sender hears it was taken)
			record()
			c.Write([]byte{0})
			return
		}
		c.Write([]byte{0})
	}
}

// ServeSerial reads jobs from a serial line (the far end of a
// pseudo-terminal): a job ends when nothing comes for idle (a serial
// printer has no end of job but the gap after it).
func (p *Printer) ServeSerial(rw io.ReadWriter, idle time.Duration) error {
	if idle <= 0 {
		idle = 500 * time.Millisecond
	}
	type chunk struct {
		b   []byte
		err error
	}
	in := make(chan chunk)
	go func() {
		buf := make([]byte, 32<<10)
		for {
			n, err := rw.Read(buf)
			in <- chunk{bytes.Clone(buf[:n]), err}
			if err != nil {
				close(in)
				return
			}
		}
	}()
	var data bytes.Buffer
	flush := func() {
		if data.Len() == 0 {
			return
		}
		b := bytes.Clone(data.Bytes())
		data.Reset()
		p.lock()
		j := &Job{Transport: "serial", Documents: [][]byte{b}, Format: p.rawFormat(b)}
		p.add(j)
		p.done(j)
		p.mu.Unlock()
		p.save(j)
	}
	timer := time.NewTimer(idle)
	defer timer.Stop()
	for {
		select {
		case c, ok := <-in:
			if !ok {
				flush()
				return nil
			}
			if len(c.b) > 0 {
				before := data.Len()
				data.Write(c.b)
				p.answer(data.Bytes()[max64(0, before-8):], rw)
				timer.Reset(idle)
			}
			if c.err != nil {
				flush()
				if c.err == io.EOF {
					return nil
				}
				return c.err
			}
		case <-timer.C:
			flush()
			timer.Reset(idle)
		}
	}
}

// Sniff tells the format of a document from its first bytes: a MIME type
// for what IPP knows (application/pdf, application/postscript,
// image/pwg-raster, image/urf, image/jpeg, image/png), a language of
// print/pdl for labels and receipts ("zpl", "sbpl", "tpcl", "escpos",
// "escp"), "pcl", or "". A PJL header is looked through.
func Sniff(b []byte) string {
	if bytes.HasPrefix(b, []byte("\x1b%-12345X")) {
		if i := bytes.Index(b, []byte("@PJL ENTER LANGUAGE")); i >= 0 {
			if j := bytes.IndexByte(b[i:], '\n'); j >= 0 {
				return Sniff(b[i+j+1:])
			}
		}
		return ""
	}
	switch {
	case bytes.HasPrefix(b, []byte("%PDF-")):
		if bytes.Contains(b[:min(len(b), 4096)], []byte("PCLm")) {
			return "application/PCLm"
		}
		return "application/pdf"
	case bytes.HasPrefix(b, []byte("%!")), bytes.HasPrefix(b, []byte("\x04%!")):
		return "application/postscript"
	case bytes.HasPrefix(b, []byte("RaS2")):
		return "image/pwg-raster"
	case bytes.HasPrefix(b, []byte("UNIRAST")):
		return "image/urf"
	case bytes.HasPrefix(b, []byte("\xff\xd8\xff")):
		return "image/jpeg"
	case bytes.HasPrefix(b, []byte("\x89PNG")):
		return "image/png"
	case bytes.HasPrefix(b, []byte("^XA")), bytes.HasPrefix(b, []byte("~")) && bytes.Contains(b[:min(len(b), 64)], []byte("^XA")):
		return "zpl"
	case bytes.HasPrefix(b, []byte("\x02\x1bA")):
		return "sbpl"
	case bytes.HasPrefix(b, []byte("{D")):
		return "tpcl"
	case bytes.HasPrefix(b, []byte("\x1b@")):
		head := b[:min(len(b), 64)]
		if bytes.Contains(head, []byte("\x1dP")) || bytes.Contains(head, []byte("\x1dv0")) {
			return "escpos"
		}
		return "escp"
	case bytes.HasPrefix(b, []byte("\x1bE")):
		return "pcl"
	}
	return ""
}

// PJL returns the lines of the PJL header of a document, and the document
// in it; ok is false for a document without one.
func PJL(data []byte) (lines []string, doc []byte, ok bool) {
	const uel = "\x1b%-12345X"
	if !bytes.HasPrefix(data, []byte(uel)) {
		return nil, data, false
	}
	rest := data[len(uel):]
	for {
		i := bytes.IndexByte(rest, '\n')
		if i < 0 {
			return lines, nil, false
		}
		l := strings.TrimRight(string(rest[:i]), "\r")
		rest = rest[i+1:]
		lines = append(lines, l)
		if strings.HasPrefix(l, "@PJL ENTER LANGUAGE") {
			break
		}
	}
	end := bytes.LastIndex(rest, []byte(uel+"@PJL EOJ"))
	if end < 0 {
		return lines, rest, false
	}
	return lines, rest[:end], true
}
