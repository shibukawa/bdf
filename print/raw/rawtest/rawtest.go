// Package rawtest is a printer of 9100 or LPD for tests: it keeps the
// jobs it is sent, and can fail as printers fail.
package rawtest

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
)

// Job is a job the printer was sent.
type Job struct {
	// Data is the document (the data file of an LPD job).
	Data []byte
	// Queue and Control are the queue and the control file of an LPD job.
	Queue   string
	Control string
}

// Server is a printer of 9100 or LPD on a port of the loopback.
type Server struct {
	// URL is the printer's URL: socket://127.0.0.1:port, or
	// lpd://127.0.0.1:port/queue.
	URL  string
	lpd  bool
	addr string
	ln   net.Listener
	wg   sync.WaitGroup

	mu         sync.Mutex
	jobs       []Job
	refuse     bool
	cutAfter   int64
	cutControl bool
	back       string
}

// Refuse makes an LPD printer refuse its queue, or take it again.
func (s *Server) Refuse(refuse bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.refuse = refuse
}

// CutAfter closes the connection of the next job after n bytes of its
// document (9100), or once its data file is sent, before the control file
// (LPD, any n).
func (s *Server) CutAfter(n int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cutAfter = n
}

// CutControl closes the connection of the next LPD job in its control
// file.
func (s *Server) CutControl() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cutControl = true
}

// Back sets what a 9100 printer says back once it has a job.
func (s *Server) Back(back string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.back = back
}

// NewSocket starts a 9100 printer.
func NewSocket() *Server { return start("socket://%s", false) }

// NewLPD starts an LPD printer with a queue.
func NewLPD(queue string) *Server { return start("lpd://%s/"+queue, true) }

func start(url string, lpd bool) *Server {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	s := &Server{URL: fmt.Sprintf(url, ln.Addr()), lpd: lpd, addr: ln.Addr().String()}
	s.serve(ln)
	return s
}

// serve takes the connections of a listener.
func (s *Server) serve(ln net.Listener) {
	s.mu.Lock()
	s.ln = ln
	s.mu.Unlock()
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			s.wg.Add(1)
			go func() {
				defer s.wg.Done()
				defer c.Close()
				if s.lpd {
					s.serveLPD(c)
				} else {
					s.serveSocket(c)
				}
			}()
		}
	}()
}

// Close stops the printer: connections to it are refused.
func (s *Server) Close() {
	s.mu.Lock()
	ln := s.ln
	s.mu.Unlock()
	ln.Close()
	s.wg.Wait()
}

// Restart starts a printer that was closed again, at its address.
func (s *Server) Restart() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.serve(ln)
	return nil
}

// Jobs returns the jobs the printer has.
func (s *Server) Jobs() []Job {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]Job(nil), s.jobs...)
}

// cut takes the failure of the next job.
func (s *Server) cut() (after int64, control bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	after, control = s.cutAfter, s.cutControl
	s.cutAfter, s.cutControl = 0, false
	return after, control
}

func (s *Server) serveSocket(c net.Conn) {
	r := bufio.NewReader(c)
	if _, err := r.Peek(1); err != nil {
		// a connection that sends nothing is a check, not a job
		return
	}
	after, _ := s.cut()
	var in io.Reader = r
	if after > 0 {
		in = io.LimitReader(r, after)
	}
	data, err := io.ReadAll(in)
	if after > 0 {
		// dropped: the connection is reset
		if tc, ok := c.(*net.TCPConn); ok {
			tc.SetLinger(0)
		}
		return
	}
	if err != nil {
		return
	}
	s.mu.Lock()
	s.jobs = append(s.jobs, Job{Data: data})
	back := s.back
	s.mu.Unlock()
	if back != "" {
		io.WriteString(c, back)
	}
}

// serveLPD takes a receive-job command (RFC 1179) and its files.
func (s *Server) serveLPD(c net.Conn) {
	r := bufio.NewReader(c)
	line, err := r.ReadString('\n')
	if err != nil || len(line) < 2 {
		return
	}
	if line[0] != 2 {
		c.Write([]byte{1})
		return
	}
	after, cutControl := s.cut()
	s.mu.Lock()
	refuse := s.refuse
	s.mu.Unlock()
	if refuse {
		c.Write([]byte{1})
		return
	}
	job := Job{Queue: strings.TrimSuffix(line[1:], "\n")}
	c.Write([]byte{0})
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return
		}
		if len(line) < 2 || line[0] != 2 && line[0] != 3 {
			c.Write([]byte{1})
			return
		}
		size, _, _ := strings.Cut(line[1:], " ")
		n, err := strconv.Atoi(size)
		if err != nil || n < 0 {
			c.Write([]byte{1})
			return
		}
		if line[0] == 2 && after > 0 {
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
			job.Data = buf[:n]
			c.Write([]byte{0})
			continue
		}
		job.Control = string(buf[:n])
		s.mu.Lock()
		s.jobs = append(s.jobs, job)
		s.mu.Unlock()
		c.Write([]byte{0})
		return
	}
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
