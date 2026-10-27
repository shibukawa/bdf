package transfer

import (
	"testing"
)

// A node's URL: the schemes and their ports, the user, the directory
// (absolute, under the login's own, or the login's own), no password.
func TestParseURL(t *testing.T) {
	ports := map[string]int{"sftp": 22}
	for _, c := range []struct {
		url             string
		host            string
		port            int
		user, base, dir string
		path, fileURL   string
		err             string
	}{
		{url: "sftp://u@h/in", host: "h", port: 22, user: "u", base: "/in", dir: "/in", path: "/in/x/y", fileURL: "sftp://h/in/x/y"},
		{url: "sftp://u@h:2222/", host: "h", port: 2222, user: "u", base: "/", dir: "/", path: "/x/y", fileURL: "sftp://h:2222/x/y"},
		{url: "sftp://u@h", host: "h", port: 22, user: "u", base: "", dir: ".", path: "x/y", fileURL: "sftp://h/~/x/y"},
		{url: "sftp://u@h/~/in/", host: "h", port: 22, user: "u", base: "in", dir: "in", path: "in/x/y", fileURL: "sftp://h/~/in/x/y"},
		{url: "sftp://h/a/../b/./c", host: "h", port: 22, user: "", base: "/a/b/c", dir: "/a/b/c", path: "/a/b/c/x/y", fileURL: "sftp://h/a/b/c/x/y"}, // .. and . are dropped, not resolved (CleanPath)
		{url: "sftp://u:pw@h/", err: "password"},
		{url: "ftp://u@h/", err: "sftp://user@host:port/dir"},
		{url: "sftp://u@h:0/", err: "port"},
		{url: "sftp://u@h/in?x=1", err: "sftp://"},
		{url: "sftp:///in", err: "sftp://"},
	} {
		s, err := ParseURL(c.url, ports)
		if c.err != "" {
			if err == nil || !contains(err.Error(), c.err) {
				t.Errorf("%s: %v, want %q", c.url, err, c.err)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: %v", c.url, err)
			continue
		}
		if s.Host != c.host || s.Port != c.port || s.User != c.user || s.Base != c.base || s.Dir() != c.dir {
			t.Errorf("%s: %+v (dir %s)", c.url, s, s.Dir())
		}
		if got := s.Path("x/y"); got != c.path {
			t.Errorf("%s: path %s, want %s", c.url, got, c.path)
		}
		if got := s.URL(s.Path("x/y")); got != c.fileURL {
			t.Errorf("%s: url %s, want %s", c.url, got, c.fileURL)
		}
	}
}

// The settings of a poll: the pattern, what is done after.
func TestCheckPoll(t *testing.T) {
	for _, c := range []struct {
		s   PollSettings
		err string
	}{
		{PollSettings{}, ""},
		{PollSettings{Pattern: "*.pdf", After: AfterDelete}, ""},
		{PollSettings{After: AfterMove, DoneDir: "done"}, ""},
		{PollSettings{Pattern: "[a-"}, "pattern"},
		{PollSettings{After: AfterMove}, "done_dir"},
		{PollSettings{After: "keep"}, "after"},
		{PollSettings{MaxPerTick: -1}, "max_per_tick"},
	} {
		err := CheckPoll(c.s)
		if (c.err == "") != (err == nil) || err != nil && !contains(err.Error(), c.err) {
			t.Errorf("%+v: %v, want %q", c.s, err, c.err)
		}
	}
}

func contains(s, sub string) bool {
	return len(sub) == 0 || len(s) >= len(sub) && (s == sub || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
