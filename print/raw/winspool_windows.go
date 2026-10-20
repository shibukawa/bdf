package raw

import (
	"errors"
	"fmt"
	"os"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	winspool         = windows.NewLazySystemDLL("winspool.drv")
	procOpenPrinter  = winspool.NewProc("OpenPrinterW")
	procClosePrinter = winspool.NewProc("ClosePrinter")
	procStartDoc     = winspool.NewProc("StartDocPrinterW")
	procEndDoc       = winspool.NewProc("EndDocPrinter")
	procStartPage    = winspool.NewProc("StartPagePrinter")
	procEndPage      = winspool.NewProc("EndPagePrinter")
	procWritePrinter = winspool.NewProc("WritePrinter")
)

// docInfo1 is DOC_INFO_1W.
type docInfo1 struct {
	docName, outputFile, datatype *uint16
}

// spool is a job of RAW data on a queue of Windows: the data goes to the
// printer as it is, past the driver, which a driver of a printer of
// labels or receipts passes on (Zebra, SATO, Epson TM …). The job is
// printed once it is closed.
type spool struct {
	h    uintptr
	open bool
}

func openPrinter(name string) (uintptr, error) {
	p, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return 0, err
	}
	var h uintptr
	if r, _, err := procOpenPrinter.Call(uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(&h)), 0); r == 0 {
		return 0, &os.PathError{Op: "OpenPrinter", Path: name, Err: err}
	}
	return h, nil
}

func checkSpool(name string) error {
	h, err := openPrinter(name)
	if err != nil {
		return err
	}
	procClosePrinter.Call(h)
	return nil
}

func openSpool(name, job string) (device, error) {
	h, err := openPrinter(name)
	if err != nil {
		return nil, err
	}
	doc, _ := windows.UTF16PtrFromString(job)
	raw, _ := windows.UTF16PtrFromString("RAW")
	info := docInfo1{docName: doc, datatype: raw}
	if r, _, err := procStartDoc.Call(h, 1, uintptr(unsafe.Pointer(&info))); r == 0 {
		procClosePrinter.Call(h)
		return nil, fmt.Errorf("raw: StartDocPrinter %s: %w", name, err)
	}
	if r, _, err := procStartPage.Call(h); r == 0 {
		procEndDoc.Call(h)
		procClosePrinter.Call(h)
		return nil, fmt.Errorf("raw: StartPagePrinter %s: %w", name, err)
	}
	return &spool{h: h, open: true}, nil
}

func (s *spool) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	var n uint32
	if r, _, err := procWritePrinter.Call(s.h, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)), uintptr(unsafe.Pointer(&n))); r == 0 {
		return int(n), fmt.Errorf("raw: WritePrinter: %w", err)
	}
	if int(n) < len(b) {
		return int(n), errors.New("raw: WritePrinter took part of the data")
	}
	return int(n), nil
}

func (s *spool) Close() error {
	if !s.open {
		return nil
	}
	s.open = false
	var err error
	if r, _, e := procEndPage.Call(s.h); r == 0 {
		err = fmt.Errorf("raw: EndPagePrinter: %w", e)
	}
	if r, _, e := procEndDoc.Call(s.h); r == 0 && err == nil {
		err = fmt.Errorf("raw: EndDocPrinter: %w", e)
	}
	procClosePrinter.Call(s.h)
	return err
}

// SetWriteDeadline is not had: the spooler takes the data at once.
func (s *spool) SetWriteDeadline(time.Time) error { return os.ErrNoDeadline }
