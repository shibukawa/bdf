package raw

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Serial is the setting of a serial port: 9600 bauds, 8 bits, no parity,
// 1 stop bit and no flow control unless the URL says otherwise.
type Serial struct {
	Baud int
	// Bits are the data bits: 7 or 8.
	Bits int
	// Parity is "none", "even" or "odd".
	Parity string
	// Stop are the stop bits: 1 or 2.
	Stop int
	// Flow is "none", "rtscts" (hardware) or "xonxoff" (software).
	Flow string
}

// bauds are the speeds a serial port is set to.
var bauds = []int{1200, 2400, 4800, 9600, 19200, 38400, 57600, 115200, 230400}

func parseSerial(q url.Values) (Serial, error) {
	s := Serial{Baud: 9600, Bits: 8, Parity: "none", Stop: 1, Flow: "none"}
	num := func(name string, v *int, ok func(int) bool) error {
		t := q.Get(name)
		if t == "" {
			return nil
		}
		n, err := strconv.Atoi(t)
		if err != nil || !ok(n) {
			return fmt.Errorf("%s %q", name, t)
		}
		*v = n
		return nil
	}
	if err := num("baud", &s.Baud, func(n int) bool {
		for _, b := range bauds {
			if b == n {
				return true
			}
		}
		return false
	}); err != nil {
		return s, fmt.Errorf("%w: 1200 to 230400", err)
	}
	if err := num("bits", &s.Bits, func(n int) bool { return n == 7 || n == 8 }); err != nil {
		return s, fmt.Errorf("%w: 7 or 8", err)
	}
	if err := num("stop", &s.Stop, func(n int) bool { return n == 1 || n == 2 }); err != nil {
		return s, fmt.Errorf("%w: 1 or 2", err)
	}
	switch v := strings.ToLower(q.Get("parity")); v {
	case "":
	case "none", "even", "odd":
		s.Parity = v
	default:
		return s, fmt.Errorf("parity %q: none, even or odd", v)
	}
	switch v := strings.ToLower(q.Get("flow")); v {
	case "":
	case "none", "rtscts", "xonxoff":
		s.Flow = v
	default:
		return s, fmt.Errorf("flow %q: none, rtscts or xonxoff", v)
	}
	return s, nil
}

// ErrNotHere is returned for a device that this OS does not reach that
// way (a USB device file on Windows, a queue of Windows elsewhere).
var ErrNotHere = errors.New("raw: the printer is not reached that way on this system")

// device is an open serial port, device file or queue of Windows.
type device interface {
	io.Writer
	Close() error
	// SetWriteDeadline bounds the writes after it; a device that cannot
	// returns an error, and is then written to without a bound.
	SetWriteDeadline(t time.Time) error
}

// openDevice opens the printer's device for writing.
func (p *Printer) openDevice(name string) (device, error) {
	switch p.Scheme {
	case "serial":
		return openSerial(p.Device, p.Serial)
	case "usb":
		return openUSB(p.Device)
	case "winspool":
		return openSpool(p.Device, name)
	}
	return nil, fmt.Errorf("raw: scheme %q", p.Scheme)
}

// checkDevice tells whether the device is there: the serial port or the
// device file exists, the queue of Windows opens.
func (p *Printer) checkDevice() error {
	switch p.Scheme {
	case "winspool":
		return checkSpool(p.Device)
	case "usb":
		if err := usbHere(); err != nil {
			return err
		}
	}
	_, err := os.Stat(p.Device)
	return err
}

// openUSB opens the device file of a USB printer (usblp's /dev/usb/lp0 on
// Linux, ulpt on the BSDs), which takes what is written as the printer's
// own data.
func openUSB(path string) (device, error) {
	if err := usbHere(); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|nonblock, 0)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// sendDevice sends a document to a serial port, a device file or a queue
// of Windows: all of it, in blocks, each within the timeout where the
// device can bound its writes. The printer may print it from its first
// byte, so Begin is called before it.
func (p *Printer) sendDevice(ctx context.Context, j *Job, data []byte) (int64, error) {
	d, err := p.openDevice(firstOr(j.Name, "document"))
	if err != nil {
		return 0, err
	}
	closed := false
	done := make(chan struct{})
	defer func() {
		close(done)
		if !closed {
			d.Close()
		}
	}()
	go func() {
		select {
		case <-ctx.Done():
			d.SetWriteDeadline(time.Now())
		case <-done:
		}
	}()
	if j.Begin != nil {
		if err := j.Begin(); err != nil {
			return 0, err
		}
	}
	deadlines := d.SetWriteDeadline(time.Now().Add(p.timeout())) == nil
	var n int64
	for len(data) > 0 {
		if err := ctx.Err(); err != nil {
			return n, fmt.Errorf("%w: %w", ErrPartial, err)
		}
		k := min(len(data), 16<<10)
		if deadlines {
			d.SetWriteDeadline(time.Now().Add(p.timeout()))
		}
		w, err := d.Write(data[:k])
		n += int64(w)
		if err != nil {
			return n, fmt.Errorf("%w: %w", ErrPartial, ctxErr(ctx, err))
		}
		data = data[k:]
	}
	closed = true
	if err := d.Close(); err != nil {
		// a queue of Windows takes the job when it is closed
		return n, fmt.Errorf("%w: %w", ErrPartial, err)
	}
	return n, nil
}

func firstOr(s, or string) string {
	if s == "" {
		return or
	}
	return s
}
