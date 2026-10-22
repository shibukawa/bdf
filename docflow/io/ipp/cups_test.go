package ipp

import "testing"

// cupsDest: the server, and the node under CUPS's and IPP's paths.
func TestCupsDest(t *testing.T) {
	for path, want := range map[string]string{"/": "", "/admin": "", "/ipp/print": "cups1", "/printers/cups1": "cups1", "/classes/cups1": "cups1", "/ipp/print/cups1/": "cups1"} {
		if got, ok := cupsDest(path, "cups1"); !ok || got != want {
			t.Errorf("%s: %q %v", path, got, ok)
		}
	}
	for _, path := range []string{"/jobs/1", "/printers/wh-a", "/classes/tokyo", "/ipp/print/x", "/printers/cups1x", "/printers/"} {
		if _, ok := cupsDest(path, "cups1"); ok {
			t.Errorf("%s is a destination", path)
		}
	}
}
