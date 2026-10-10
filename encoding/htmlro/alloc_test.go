package htmlro

import (
	"bytes"
	"runtime"
	"testing"
)

// allocatedBytes returns the bytes allocated by runs calls of fn. It counts
// on both compilers: TinyGo's testing.AllocsPerRun is a constant zero, and
// its MemStats.Mallocs too, but TotalAlloc moves there as well.
func allocatedBytes(runs int, fn func()) uint64 {
	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)
	for i := 0; i < runs; i++ {
		fn()
	}
	runtime.ReadMemStats(&after)
	return after.TotalAlloc - before.TotalAlloc
}

// The helper itself must see an allocation, or every test built on it is
// vacuous.
func TestAllocatedBytesSeesAllocations(t *testing.T) {
	var keep [][]byte
	if b := allocatedBytes(10, func() { keep = append(keep, make([]byte, 1000)) }); b < 10*1000 {
		t.Errorf("allocatedBytes saw %d bytes of 10000 allocated; the tests built on it would be vacuous", b)
	}
	if len(keep) != 10 {
		t.Error(len(keep))
	}
}

// The reader allocates nothing in steady state: the token loop, the names,
// the attributes and the decoded values all come out of its own buffers,
// whichever way the input arrives.
func TestReaderAllocatesNothing(t *testing.T) {
	doc := []byte(`<!DOCTYPE html>
<html lang=en><head><title>T &amp; x</title><meta charset="utf-8"></head>
<body class="a b" id=main data-x="1 &lt; 2">
<p>Hello <b>world</b> &copy; 2026<br/><IMG SRC="x.png" ALT="y">
<script>var a = "<p>"; // &amp;</script><!-- c --><textarea>&lt;</textarea>
<svg viewBox="0 0 1 1"><![CDATA[x]]></svg>
</p></body></html>`)
	r := NewBytesReader(doc, Options{})
	src := bytes.NewReader(doc)
	var buf []byte
	var n int
	loop := func() {
		for {
			k, err := r.Next()
			if err != nil || k == EOF {
				return
			}
			switch k {
			case StartTag, SelfClosingTag:
				n += len(r.Name())
				if v, ok := r.Attr("class"); ok && v.Equal("a b") {
					n++
				}
				for {
					name, v, ok := r.NextAttr()
					if !ok {
						break
					}
					n += len(name)
					buf = v.AppendTo(buf[:0])
				}
			case EndTag:
				if r.NameIs("p") {
					n++
				}
			case Text, Comment, Doctype:
				buf = r.Text().AppendTo(buf[:0])
			}
		}
	}
	cases := []struct {
		name string
		fn   func()
	}{
		{"bytes", func() {
			r.ResetBytes(doc)
			loop()
		}},
		{"stream", func() {
			src.Reset(doc)
			r.Reset(src)
			loop()
		}},
	}
	for _, c := range cases {
		c.fn() // the first run sizes buf and the reader's buffers
		if b := allocatedBytes(50, c.fn); b != 0 {
			t.Errorf("%s: %d bytes allocated in 50 runs", c.name, b)
		}
	}
	if n == 0 {
		t.Error("nothing read")
	}
}
