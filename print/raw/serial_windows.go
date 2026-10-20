package raw

import (
	"os"

	"golang.org/x/sys/windows"
)

// The flags of a DCB (fBinary, fParity, fOutxCtsFlow, fDtrControl,
// fOutX, fInX, fRtsControl).
const (
	dcbBinary       = 0x1
	dcbParity       = 0x2
	dcbOutxCtsFlow  = 0x4
	dcbDtrEnable    = 0x1 << 4
	dcbOutX         = 0x100
	dcbInX          = 0x200
	dcbRtsEnable    = 0x1 << 12
	dcbRtsHandshake = 0x2 << 12
)

// openSerial opens a COM port (COM3, or \\.\COM12) as the setting says.
func openSerial(name string, s Serial) (device, error) {
	path := name
	if len(path) < 4 || path[:4] != `\\.\` {
		path = `\\.\` + name
	}
	p, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	h, err := windows.CreateFile(p, windows.GENERIC_READ|windows.GENERIC_WRITE, 0, nil, windows.OPEN_EXISTING, 0, 0)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: name, Err: err}
	}
	if err := setSerial(h, s); err != nil {
		windows.CloseHandle(h)
		return nil, &os.PathError{Op: "set", Path: name, Err: err}
	}
	return os.NewFile(uintptr(h), name), nil
}

func setSerial(h windows.Handle, s Serial) error {
	var d windows.DCB
	d.DCBlength = uint32(64)
	if err := windows.GetCommState(h, &d); err != nil {
		return err
	}
	d.BaudRate = uint32(s.Baud)
	d.ByteSize = uint8(s.Bits)
	d.Flags = dcbBinary | dcbDtrEnable | dcbRtsEnable
	switch s.Parity {
	case "even":
		d.Parity, d.Flags = 2, d.Flags|dcbParity
	case "odd":
		d.Parity, d.Flags = 1, d.Flags|dcbParity
	default:
		d.Parity = 0
	}
	d.StopBits = 0
	if s.Stop == 2 {
		d.StopBits = 2
	}
	switch s.Flow {
	case "rtscts":
		d.Flags = d.Flags&^dcbRtsEnable | dcbOutxCtsFlow | dcbRtsHandshake
	case "xonxoff":
		d.Flags |= dcbOutX | dcbInX
		d.XonChar, d.XoffChar = 0x11, 0x13
	}
	if err := windows.SetCommState(h, &d); err != nil {
		return err
	}
	// a write that the printer does not take in 30 seconds fails
	return windows.SetCommTimeouts(h, &windows.CommTimeouts{WriteTotalTimeoutConstant: 30000})
}
