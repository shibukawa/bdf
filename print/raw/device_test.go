package raw_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/shibukawa/bdf/print/fakeprinter"
	"github.com/shibukawa/bdf/print/raw"
)

func TestParseDevices(t *testing.T) {
	for _, c := range []struct {
		url  string
		want raw.Printer
	}{
		{"serial:///dev/ttyUSB0?lang=escpos&baud=19200&flow=rtscts", raw.Printer{URL: "serial:///dev/ttyUSB0", Scheme: "serial", Device: "/dev/ttyUSB0", Format: "escpos",
			Serial: raw.Serial{Baud: 19200, Bits: 8, Parity: "none", Stop: 1, Flow: "rtscts"}}},
		{"serial://COM3?format=zpl&profile=zd421-300&parity=even&bits=7&stop=2", raw.Printer{URL: "serial://COM3", Scheme: "serial", Device: "COM3", Format: "zpl", Profile: "zd421-300",
			Serial: raw.Serial{Baud: 9600, Bits: 7, Parity: "even", Stop: 2, Flow: "none"}}},
		{"usb:///dev/usb/lp0?lang=ESC/P&width=136*18", raw.Printer{URL: "usb:///dev/usb/lp0", Scheme: "usb", Device: "/dev/usb/lp0", Format: "escp", Width: "136*18"}},
		{"winspool:///Zebra%20ZD421?lang=zpl&dpi=300", raw.Printer{URL: "winspool:///Zebra%20ZD421", Scheme: "winspool", Device: "Zebra ZD421", Format: "zpl", DPI: "300"}},
		{"winspool://server/Labels", raw.Printer{URL: "winspool://server/Labels", Scheme: "winspool", Device: `\\server\Labels`, Format: raw.FormatPostScript}},
		{"winspool://Receipts?lang=escpos", raw.Printer{URL: "winspool://Receipts", Scheme: "winspool", Device: "Receipts", Format: "escpos"}},
		{"socket://printer:9100?lang=sbpl", raw.Printer{URL: "socket://printer:9100", Scheme: "socket", Addr: "printer:9100", Format: "sbpl"}},
		{"socket://printer?profile=tm-t88", raw.Printer{URL: "socket://printer", Scheme: "socket", Addr: "printer:9100", Profile: "tm-t88"}},
	} {
		p, err := raw.Parse(c.url)
		if err != nil || *p != c.want {
			t.Errorf("%s: %+v %v, want %+v", c.url, p, err, c.want)
		}
	}
	for _, bad := range []string{"serial://", "serial:///dev/x?baud=1234", "serial:///dev/x?bits=6", "serial:///dev/x?parity=mark", "serial:///dev/x?flow=dtr",
		"usb:///dev/x?baud=9600", "winspool://", "socket://h?format=ps&profile=zd421-203", "socket://h?dpi=300", "socket://h?format=pdf&lang=zpl", "socket://h?lang=pcl"} {
		if p, err := raw.Parse(bad); err == nil {
			t.Errorf("%s: %+v", bad, p)
		}
	}
	if !raw.IsURL("serial:///dev/x") || !raw.IsURL("USB:///dev/x") || !raw.IsURL("winspool://q") {
		t.Error("IsURL")
	}
	p, _ := raw.Parse("socket://h?lang=zpl&pjl=true")
	if p.UsesPJL("zpl") {
		t.Error("PJL wraps PostScript and PDF only")
	}
}

// TestUSB sends a job to a device file: a file stands for it.
func TestUSB(t *testing.T) {
	path := filepath.Join(t.TempDir(), "lp0")
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	p, err := raw.Parse("usb://" + filepath.ToSlash(path) + "?lang=zpl")
	if err != nil {
		t.Fatal(err)
	}
	data := bytes.Repeat([]byte("^XA^FDlabel^FS^XZ\n"), 5000)
	n, err := p.Send(context.Background(), &raw.Job{Name: "labels", Format: "zpl", Data: data})
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		if !errors.Is(err, raw.ErrNotHere) || p.Check(context.Background()) == nil {
			t.Fatalf("%s has no device files of USB printers: %v", runtime.GOOS, err)
		}
		return
	}
	if err != nil || n != int64(len(data)) {
		t.Fatal(n, err)
	}
	got, _ := os.ReadFile(path)
	if !bytes.Equal(got, data) {
		t.Fatalf("the device has %d bytes of %d", len(got), len(data))
	}
	if err := p.Check(context.Background()); err != nil {
		t.Fatal(err)
	}
}

// TestSerial sends a job to a serial port: a fake printer reading the
// far end of a pseudo-terminal stands for the printer.
func TestSerial(t *testing.T) {
	fp := fakeprinter.New()
	name := fakeprinter.StartSerial(t, fp)
	p, err := raw.Parse("serial://" + name + "?lang=escpos&baud=115200")
	if err != nil {
		t.Fatal(err)
	}
	data := make([]byte, 100000)
	for i := range data {
		data[i] = byte(i * 7) // every byte, the control ones too: the port is raw
	}
	begun := false
	n, err := p.Send(context.Background(), &raw.Job{Name: "receipt", Format: "escpos", Data: data, Begin: func() error { begun = true; return nil }})
	if err != nil || n != int64(len(data)) || !begun {
		t.Fatal(n, err, begun)
	}
	deadline := time.Now().Add(10 * time.Second)
	for len(fp.Jobs()) == 0 && time.Now().Before(deadline) {
		time.Sleep(50 * time.Millisecond)
	}
	jobs := fp.Jobs()
	if len(jobs) != 1 || !bytes.Equal(jobs[0].Data(), data) || jobs[0].Transport != "serial" {
		t.Fatalf("the printer has %d jobs, not the one sent", len(jobs))
	}
}

// TestSerialStalled cancels a job whose printer takes no more (its
// flow control holds it): the write ends with the context.
func TestSerialStalled(t *testing.T) {
	// a printer that reads nothing
	far, name, err := fakeprinter.OpenPTY()
	if err != nil {
		t.Skip("no pseudo-terminals:", err)
	}
	defer far.Close()
	p, err := raw.Parse("serial://" + name + "?lang=escpos")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	// nobody reads the far end: the terminal's buffer fills
	_, err = p.Send(ctx, &raw.Job{Format: "escpos", Data: make([]byte, 4<<20)})
	if !errors.Is(err, raw.ErrPartial) {
		t.Fatalf("%v, want a partial job", err)
	}
	_ = io.Discard
}
