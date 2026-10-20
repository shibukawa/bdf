package main

import (
	"context"
	"errors"
	"flag"
	"net"
	"net/url"
	"os"
	"os/user"
	"strings"

	"github.com/shibukawa/bdf/print/server/client"
)

// The environment variables of a print server's URL and token, the
// defaults of -server and -token-file.
const (
	serverEnv = "BDF_SERVER"
	tokenEnv  = "BDF_TOKEN"
)

// serverFlags are the flags of talking to a print server (bdf printd).
type serverFlags struct {
	url       *string
	tokenFile *string
}

// addServerFlags adds -server and -token-file; env makes $BDF_SERVER the
// default of -server.
func addServerFlags(fs *flag.FlagSet, env bool) *serverFlags {
	usage := "the print server (bdf printd): http://127.0.0.1:8640, https://print.example.local:8640"
	def := ""
	if env {
		usage += " (default: $" + serverEnv + ")"
		def = os.Getenv(serverEnv)
	}
	return &serverFlags{
		url:       fs.String("server", def, usage),
		tokenFile: fs.String("token-file", "", "file holding the server's token (- for the standard input; default: $"+tokenEnv+")"),
	}
}

// client makes the client of the server; user names the user of the
// jobs ("" the user of the process), which a server takes from a token
// that may act for others or when it has no tokens.
func (f *serverFlags) client(userName string) *client.Client {
	u := *f.url
	if !strings.Contains(u, "://") {
		u = "http://" + u
	}
	c := &client.Client{URL: u, Token: os.Getenv(tokenEnv), User: userName}
	if *f.tokenFile != "" {
		tok, err := readPassword(*f.tokenFile)
		check(err)
		c.Token = strings.TrimSpace(tok)
	}
	if c.User == "" {
		c.User = processUser()
	}
	return c
}

// processUser is the name of the user of the process, without a Windows
// domain.
func processUser() string {
	name := os.Getenv("USER")
	if u, err := user.Current(); err == nil && u.Username != "" {
		name = u.Username
	}
	if i := strings.LastIndexByte(name, '\\'); i >= 0 {
		name = name[i+1:]
	}
	return name
}

// plaintext tells a server reached without TLS on another machine.
func plaintext(server string) bool {
	u, err := url.Parse(server)
	if err != nil || u.Scheme == "https" {
		return false
	}
	host := u.Hostname()
	if host == "localhost" {
		return false
	}
	ip := net.ParseIP(host)
	return ip == nil || !ip.IsLoopback()
}

// serverExit is the exit status of an error of a print server, with the
// meanings of bdf print's.
func serverExit(err error) int {
	var se *client.Error
	var ue *url.Error
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return exitTimeout
	case errors.As(err, &se):
		switch {
		case se.Code == "not-accepting", se.Code == "closed", se.Code == "busy":
			return exitStopped
		case se.Code == "unauthorized":
			return exitConnection
		case se.Status >= 400 && se.Status < 500, se.Status == 507:
			return exitUnsupported
		}
	case errors.As(err, &ue):
		return exitConnection
	}
	return exitFailed
}
