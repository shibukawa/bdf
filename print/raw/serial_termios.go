//go:build linux || darwin || freebsd || netbsd || openbsd || dragonfly

package raw

import (
	"errors"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

// openSerial opens a serial port in raw mode, as the setting says, without
// making it the controlling terminal.
func openSerial(path string, s Serial) (device, error) {
	fd, err := unix.Open(path, unix.O_RDWR|unix.O_NOCTTY|unix.O_NONBLOCK|unix.O_CLOEXEC, 0)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: path, Err: err}
	}
	if err := setSerial(fd, s); err != nil {
		unix.Close(fd)
		return nil, &os.PathError{Op: "set", Path: path, Err: err}
	}
	// a non-blocking descriptor makes a file whose writes take deadlines
	return &serialPort{File: os.NewFile(uintptr(fd), path), fd: fd}, nil
}

// serialPort is an open serial port: it is closed once what was written
// is sent, which a port closed at once may drop (macOS does).
type serialPort struct {
	*os.File
	fd int
}

// drainTimeout bounds the wait for what was written to be sent, at the
// close: a printer that holds the port (its flow control) may never take
// it.
var drainTimeout = 30 * time.Second

func (s *serialPort) Close() error {
	done := make(chan error, 1)
	go func() { done <- drain(s.fd) }()
	var err error
	select {
	case err = <-done:
	case <-time.After(drainTimeout):
		err = errors.New("the port did not send what was written")
	}
	if cerr := s.File.Close(); err == nil {
		err = cerr
	}
	return err
}

// Abort closes the port, dropping what it has not sent: the job was
// canceled, or the printer stopped taking it.
func (s *serialPort) Abort() error {
	discard(s.fd)
	return s.File.Close()
}

func setSerial(fd int, s Serial) error {
	t, err := unix.IoctlGetTermios(fd, getTermios)
	if err != nil {
		return err
	}
	t.Iflag &^= unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON | unix.IXOFF | unix.IXANY
	t.Oflag &^= unix.OPOST
	t.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN
	t.Cflag &^= unix.CSIZE | unix.PARENB | unix.PARODD | unix.CSTOPB | unix.CRTSCTS
	t.Cflag |= unix.CREAD | unix.CLOCAL
	if s.Bits == 7 {
		t.Cflag |= unix.CS7
	} else {
		t.Cflag |= unix.CS8
	}
	switch s.Parity {
	case "even":
		t.Cflag |= unix.PARENB
	case "odd":
		t.Cflag |= unix.PARENB | unix.PARODD
	}
	if s.Stop == 2 {
		t.Cflag |= unix.CSTOPB
	}
	switch s.Flow {
	case "rtscts":
		t.Cflag |= unix.CRTSCTS
	case "xonxoff":
		t.Iflag |= unix.IXON | unix.IXOFF
	}
	t.Cc[unix.VMIN], t.Cc[unix.VTIME] = 0, 0
	if err := setSpeed(t, s.Baud); err != nil {
		return err
	}
	return unix.IoctlSetTermios(fd, setTermios, t)
}
