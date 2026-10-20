//go:build linux || darwin || freebsd || netbsd || openbsd || dragonfly

package raw

import (
	"os"

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
	return os.NewFile(uintptr(fd), path), nil
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
