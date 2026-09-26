package cjkcmap

import (
	"bytes"
	"encoding/hex"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestLookup(t *testing.T) {
	if loadErr != nil {
		t.Fatal(loadErr)
	}
	for _, tc := range []struct {
		cmap string
		code uint32
		cid  uint32
	}{
		{"90ms-RKSJ-H", 0x82a0, 843},     // あ
		{"90ms-RKSJ-H", 0x41, 264},       // half-width A
		{"90ms-RKSJ-V", 0x8141, 7887},    // 、 vertical form (usecmap 90ms-RKSJ-H)
		{"90ms-RKSJ-V", 0x82a0, 843},     // from the parent
		{"UniJIS-UCS2-H", 0x3042, 843},   // stored over UniJIS-UTF32-H
		{"UniJIS-UTF16-H", 0x3042, 843},  //
		{"UniJIS-UTF8-H", 0xe38182, 843}, // あ in UTF-8
		{"UniJIS-UTF32-H", 0x3042, 843},
		{"UniJIS-UTF16-V", 0x3001, 7887},
		{"UniGB-UCS2-H", 0x4e2d, 4559}, // 中
		{"UniKS-UCS2-H", 0xac00, 1086}, // 가
		{"UniCNS-UCS2-H", 0x4e2d, 661}, // 中
		{"EUC-H", 0xa4a2, 843},         // あ in EUC-JP
	} {
		m := Lookup(tc.cmap)
		if m == nil {
			t.Errorf("%s missing", tc.cmap)
			continue
		}
		if cid, ok := m.CID(tc.code); !ok || cid != tc.cid {
			t.Errorf("%s %#x = %d %v, want %d", tc.cmap, tc.code, cid, ok, tc.cid)
		}
	}
	if m := Lookup("90ms-RKSJ-V"); !m.Vertical || m.Collection != "Adobe-Japan1" || len(m.Codespace) != 4 {
		t.Errorf("90ms-RKSJ-V = %+v", m)
	}
	if r, ok := Lookup("UniJIS-UTF16-H").CodePoint(0xd840dc0b); !ok || r != 0x2000b {
		t.Errorf("UTF-16 code point %#x %v", r, ok)
	}
	if _, ok := Lookup("90ms-RKSJ-H").CodePoint(0x82a0); ok {
		t.Error("RKSJ is not a Unicode CMap")
	}
	if Lookup("Identity-H") != nil || Lookup("nope") != nil {
		t.Error("unexpected CMap")
	}
}

func TestUnicode(t *testing.T) {
	for _, tc := range []struct {
		coll     string
		cid      uint32
		r        rune
		vertical bool
	}{
		{"Adobe-Japan1", 843, 'あ', false},
		{"Adobe-Japan1", 1, ' ', false},
		{"Adobe-Japan1", 1125, '亜', false},
		{"Adobe-Japan1", 7887, '、', true}, // the vertical 、
		{"Adobe-Japan1", 7891, 'ー', true}, // only a vertical CMap reaches it
		{"Adobe-Japan1", 264, 'A', false}, // half-width
		{"Adobe-Japan1", 231, ' ', false}, // not U+2002
		{"Adobe-Japan1", 323, '|', false}, // not U+FFE8
		{"Adobe-GB1", 4559, '中', false},
		{"Adobe-Korea1", 1086, '가', false},
		{"Adobe-CNS1", 661, '中', false},
	} {
		r, v, ok := Unicode(tc.coll, tc.cid)
		if !ok || r != tc.r || v != tc.vertical {
			t.Errorf("%s CID %d = %q %v %v, want %q %v", tc.coll, tc.cid, r, v, ok, tc.r, tc.vertical)
		}
	}
	if r, ok := VerticalForm("Adobe-Japan1", 7887); !ok || r != '︑' {
		t.Errorf("vertical form of CID 7887 = %q %v", r, ok)
	}
	if _, ok := VerticalForm("Adobe-Japan1", 843); ok {
		t.Error("あ has no vertical form")
	}
	if _, _, ok := Unicode("Adobe-Identity", 5); ok {
		t.Error("Identity has no table")
	}
}

// TestAgainstSource compares every CMap with the Adobe file it was made
// from (Ghostscript's copy, or BDF_CMAP_DIR).
func TestAgainstSource(t *testing.T) {
	dir := os.Getenv("BDF_CMAP_DIR")
	if dir == "" {
		dir = "/opt/homebrew/share/ghostscript/Resource/CMap"
	}
	if _, err := os.Stat(filepath.Join(dir, "90ms-RKSJ-H")); err != nil {
		t.Skip("no CMap resources; set BDF_CMAP_DIR")
	}
	loadOnce.Do(load)
	var resolve func(name string) map[uint32]uint32
	resolve = func(name string) map[uint32]uint32 {
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		m := map[uint32]uint32{}
		if u := useRe.FindSubmatch(b); u != nil {
			m = resolve(string(u[1]))
		}
		own := map[uint32]uint32{}
		for _, sec := range rangeRe.FindAllSubmatch(b, -1) {
			for _, e := range entryRe.FindAllSubmatch(sec[1], -1) {
				lo, hi, cid := hexCode(e[1]), hexCode(e[2]), atoi(e[3])
				for c := lo; c <= hi; c++ {
					if _, dup := own[c]; !dup {
						own[c] = cid + c - lo
					}
				}
			}
		}
		for _, sec := range charRe.FindAllSubmatch(b, -1) {
			for _, e := range charEntryRe.FindAllSubmatch(sec[1], -1) {
				if c := hexCode(e[1]); own[c] == 0 {
					own[c] = atoi(e[2])
				}
			}
		}
		for c, cid := range own {
			m[c] = cid
		}
		return m
	}
	n := 0
	for name, m := range cmaps {
		want := resolve(name)
		for code, cid := range want {
			if got, ok := m.CID(code); !ok || got != cid {
				t.Errorf("%s %#x = %d %v, want %d", name, code, got, ok, cid)
				break
			}
			n++
		}
	}
	if len(cmaps) < 150 {
		t.Errorf("%d CMaps", len(cmaps))
	}
	t.Logf("%d CMaps, %d codes checked", len(cmaps), n)
}

var (
	useRe       = regexp.MustCompile(`/(\S+)\s+usecmap`)
	rangeRe     = regexp.MustCompile(`(?s)begincidrange(.*?)endcidrange`)
	entryRe     = regexp.MustCompile(`<([0-9a-fA-F]+)>\s*<([0-9a-fA-F]+)>\s*(\d+)`)
	charRe      = regexp.MustCompile(`(?s)begincidchar(.*?)endcidchar`)
	charEntryRe = regexp.MustCompile(`<([0-9a-fA-F]+)>\s*(\d+)`)
)

func hexCode(b []byte) uint32 {
	d, _ := hex.DecodeString(string(bytes.TrimSpace(b)))
	var c uint32
	for _, x := range d {
		c = c<<8 | uint32(x)
	}
	return c
}

func atoi(b []byte) uint32 {
	v, _ := strconv.Atoi(strings.TrimSpace(string(b)))
	return uint32(v)
}
