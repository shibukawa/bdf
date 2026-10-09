// Command secure-reader is a sample architecture (docs/examples/secure-reader):
// a site where readers log in and read the PDF books they own, which are
// never sent as files. The server converts each book to bdf once and keeps
// it; the browser asks for ten pages at a time, and each answer is a
// segment sealed for a key pair the browser made for that request alone
// (package segment, @bdfkit/core's SegmentLoader). Once a segment is open,
// neither side keeps the key that opens it: what a recording of the
// traffic holds stays sealed even if the server's TLS key, a reader's
// password or a session cookie is taken later.
//
// A reader who does not own a book may read its first segment (a sample);
// the server refuses the pages after it, and a reader who asks for pages
// faster than a person reads them.
//
//	cd examples/secure-reader && node web/build.mjs && go run .
//	open http://127.0.0.1:8084/   (alice / alice-pass owns both books, bob / bob-pass one)
package main

import (
	"bytes"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
	_ "github.com/shibukawa/bdf/converter/pdf"
	"github.com/shibukawa/bdf/segment"
	"github.com/shibukawa/bdf/thumbnail"
)

func main() {
	booksDir := flag.String("books", "./books", "directory of the PDF books")
	cacheDir := flag.String("cache", "./cache", "directory the converted books and their covers are kept in (not served as files)")
	webDir := flag.String("web", "./web", "directory of the built front end (see web/build.mjs)")
	addr := flag.String("listen", "127.0.0.1:8084", "address to listen on")
	pages := flag.Int("pages", segment.Pages, "pages of a segment: each is sealed for its own request")
	rate := flag.Int("rate", 20, "segments a reader may ask for in a minute")
	flag.Parse()

	if err := os.MkdirAll(*cacheDir, 0o700); err != nil {
		log.Fatal(err)
	}
	books, err := loadBooks(*booksDir, *cacheDir)
	if err != nil {
		log.Fatal(err)
	}
	s := &server{
		books: books,
		users: map[string]*user{
			"alice": newUser("alice", "alice-pass", "reading-room.pdf", "field-notes.pdf"),
			"bob":   newUser("bob", "bob-pass", "reading-room.pdf"),
		},
		sessions: map[[32]byte]*session{},
		pages:    *pages,
		rate:     *rate,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.library)
	mux.HandleFunc("GET /login", s.loginPage)
	mux.HandleFunc("POST /login", s.login)
	mux.HandleFunc("POST /logout", s.logout)
	mux.HandleFunc("GET /covers/{name}", s.cover)
	// the book itself: ten pages at a time, each segment sealed for its request
	mux.Handle("POST /segments/{name}", &segment.Handler{
		Open:  s.openBook,
		Allow: s.allow,
		Pages: *pages,
		Log:   func(r *http.Request, err error) { log.Printf("%q: %v", r.URL.Path, err) }, // quoted: a path is the client's
	})
	// the viewer: a page and its scripts, which hold nothing of any book
	mux.Handle("GET /read/", http.StripPrefix("/read/", http.FileServer(http.Dir(*webDir))))

	log.Printf("secure-reader: %d book(s), listening on http://%s/ (alice / alice-pass, bob / bob-pass)", len(books), *addr)
	// forms and segment requests of other sites are refused (with the
	// SameSite cookie, a second line)
	// a connection that sends its request slowly, or none, is not kept for ever
	srv := &http.Server{Addr: *addr, Handler: headers(http.NewCrossOriginProtection().Handler(mux)), ReadHeaderTimeout: 10 * time.Second, ReadTimeout: time.Minute, IdleTimeout: 2 * time.Minute}
	log.Fatal(srv.ListenAndServe())
}

// book is a PDF converted to bdf once, kept open to cut segments from.
type book struct {
	Name  string // the PDF's file name
	Title string
	Pages int
	doc   *bdf.Reader
	cover string // path of its cover picture
}

// loadBooks converts each PDF that has no converted copy yet, and opens them.
func loadBooks(dir, cacheDir string) (map[string]*book, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	books := map[string]*book{}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.EqualFold(filepath.Ext(name), ".pdf") {
			continue
		}
		bdfPath := filepath.Join(cacheDir, name+".bdf")
		cover := filepath.Join(cacheDir, name+".png")
		if _, err := os.Stat(bdfPath); err != nil {
			if err := convert(filepath.Join(dir, name), bdfPath, cover); err != nil {
				log.Printf("secure-reader: %s: not converted: %v", name, err)
				continue
			}
		}
		r, err := bdf.OpenSingleFile(bdfPath)
		if err != nil {
			return nil, err
		}
		b := &book{Name: name, Title: r.Manifest.Meta.DC.Title.First(), doc: r, cover: cover}
		if b.Title == "" {
			b.Title = strings.TrimSuffix(name, filepath.Ext(name))
		}
		if len(r.Manifest.Views) > 0 {
			b.Pages = len(r.Manifest.Views[0].Pages)
		}
		books[name] = b
	}
	return books, nil
}

// convert writes the bdf of a PDF and the picture of its first page.
func convert(src, bdfPath, cover string) error {
	res, err := converter.ConvertFile(src, "", &converter.Options{Params: map[string]string{"remote": "false"}})
	if err != nil {
		return err
	}
	var b bytes.Buffer
	if err := res.Doc.WriteSingle(&b); err != nil {
		return err
	}
	if err := os.WriteFile(bdfPath, b.Bytes(), 0o600); err != nil {
		return err
	}
	th, err := thumbnail.Make(res.Doc, &thumbnail.Options{Size: 320})
	if err != nil {
		log.Printf("secure-reader: %s: cover: %v", src, err)
		return nil
	}
	var png bytes.Buffer
	if err := thumbnail.Encode(&png, th.Image, thumbnail.PNG); err != nil {
		return err
	}
	return os.WriteFile(cover, png.Bytes(), 0o600)
}

// user is a reader's account. A real site keeps these in its database; the
// password is kept as a PBKDF2 hash here too.
type user struct {
	Name  string
	salt  []byte
	hash  []byte
	owns  []string
	mu    sync.Mutex
	asked []time.Time // when the segments of the last minute were asked for
}

const passwordIterations = 600_000

func newUser(name, password string, owns ...string) *user {
	u := &user{Name: name, salt: make([]byte, 16), owns: owns}
	rand.Read(u.salt)
	u.hash = hashPassword(password, u.salt)
	return u
}

func hashPassword(password string, salt []byte) []byte {
	h, err := pbkdf2.Key(sha256.New, password, salt, passwordIterations, 32)
	if err != nil {
		panic(err)
	}
	return h
}

func (u *user) has(name string) bool { return slices.Contains(u.owns, name) }

type session struct {
	user    *user
	expires time.Time
}

const sessionCookie = "session"
const sessionLife = 8 * time.Hour

type server struct {
	books map[string]*book
	users map[string]*user
	pages int
	rate  int

	mu       sync.Mutex
	sessions map[[32]byte]*session // by the hash of the cookie: a copy of the table logs nobody in
}

// reader returns the user a request's session cookie logs in, or nil.
func (s *server) reader(r *http.Request) *user {
	c, err := r.Cookie(sessionCookie)
	if err != nil {
		return nil
	}
	key := sha256.Sum256([]byte(c.Value))
	s.mu.Lock()
	defer s.mu.Unlock()
	ss := s.sessions[key]
	if ss == nil || time.Now().After(ss.expires) {
		delete(s.sessions, key)
		return nil
	}
	return ss.user
}

func (s *server) login(w http.ResponseWriter, r *http.Request) {
	u := s.users[r.PostFormValue("user")]
	// a user that does not exist costs the same time as a wrong password
	salt, want := make([]byte, 16), make([]byte, 32)
	if u != nil {
		salt, want = u.salt, u.hash
	}
	ok := subtle.ConstantTimeCompare(hashPassword(r.PostFormValue("password"), salt), want) == 1
	if u == nil || !ok {
		w.WriteHeader(http.StatusUnauthorized)
		s.render(w, loginTemplate, "The user name or the password is wrong.")
		return
	}
	token := make([]byte, 32)
	rand.Read(token)
	value := base64.RawURLEncoding.EncodeToString(token)
	s.mu.Lock()
	s.sessions[sha256.Sum256([]byte(value))] = &session{user: u, expires: time.Now().Add(sessionLife)}
	s.mu.Unlock()
	http.SetCookie(w, &http.Cookie{Name: sessionCookie, Value: value, Path: "/", HttpOnly: true, Secure: r.TLS != nil, SameSite: http.SameSiteStrictMode, MaxAge: int(sessionLife.Seconds())})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *server) logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(sessionCookie); err == nil {
		s.mu.Lock()
		delete(s.sessions, sha256.Sum256([]byte(c.Value)))
		s.mu.Unlock()
	}
	http.SetCookie(w, &http.Cookie{Name: sessionCookie, Path: "/", HttpOnly: true, Secure: r.TLS != nil, SameSite: http.SameSiteStrictMode, MaxAge: -1})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// openBook is segment.Handler's Open: the book, for a reader who is logged in.
func (s *server) openBook(r *http.Request) (*bdf.Reader, error) {
	if s.reader(r) == nil {
		return nil, &segment.StatusError{Status: http.StatusUnauthorized, Err: errors.New("not logged in")}
	}
	b := s.books[r.PathValue("name")]
	if b == nil {
		return nil, fs.ErrNotExist
	}
	return b.doc, nil
}

// allow is segment.Handler's Allow: the first segment of any book (a
// sample), the others to the book's owners, and not too many a minute.
func (s *server) allow(r *http.Request, sg bdf.Segment) error {
	u := s.reader(r)
	if u == nil {
		return &segment.StatusError{Status: http.StatusUnauthorized, Err: errors.New("not logged in")}
	}
	if sg.From >= s.pages && !u.has(r.PathValue("name")) {
		return fs.ErrPermission
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	now := time.Now()
	u.asked = slices.DeleteFunc(u.asked, func(t time.Time) bool { return now.Sub(t) > time.Minute })
	if len(u.asked) >= s.rate {
		return &segment.StatusError{Status: http.StatusTooManyRequests, Err: fmt.Errorf("%s reads too fast", u.Name)}
	}
	u.asked = append(u.asked, now)
	return nil
}

// cover sends the picture of a book's first page, to a reader who is logged in.
func (s *server) cover(w http.ResponseWriter, r *http.Request) {
	b := s.books[r.PathValue("name")]
	if s.reader(r) == nil || b == nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Cache-Control", "private, no-cache")
	http.ServeFile(w, r, b.cover)
}

// shelfBook is a book as the library page shows it.
type shelfBook struct {
	*book
	Owned  bool
	Sample int // pages a reader who does not own it may read
}

func (s *server) library(w http.ResponseWriter, r *http.Request) {
	u := s.reader(r)
	if u == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	var shelf []shelfBook
	for _, b := range s.books {
		shelf = append(shelf, shelfBook{book: b, Owned: u.has(b.Name), Sample: min(s.pages, b.Pages)})
	}
	slices.SortFunc(shelf, func(a, b shelfBook) int { return strings.Compare(a.Title, b.Title) })
	s.render(w, libraryTemplate, struct {
		User  string
		Books []shelfBook
	}{u.Name, shelf})
}

func (s *server) loginPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, loginTemplate, "")
}

func (s *server) render(w http.ResponseWriter, t *template.Template, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	if err := t.Execute(w, data); err != nil {
		log.Print(err)
	}
}

// headers adds what every answer carries: no scripts but the site's own, no
// framing by other sites, no address sent to them.
func headers(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; img-src 'self' blob: data:; style-src 'self' 'unsafe-inline'; object-src 'none'; base-uri 'none'; frame-ancestors 'none'; form-action 'self'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "same-origin")
		h.ServeHTTP(w, r)
	})
}

const style = `<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
  :root { color-scheme: light dark; font: 16px system-ui, sans-serif; }
  body { max-width: 60rem; margin: 2rem auto; padding: 0 1rem; }
  h1 { font-size: 1.375rem; }
  header { display: flex; justify-content: space-between; align-items: baseline; }
  p.lead, small { color: #666; }
  form.login { display: grid; gap: .75rem; max-width: 20rem; }
  input, button { font: inherit; padding: .4rem .6rem; }
  .error { color: #b00020; }
  .books { display: grid; grid-template-columns: repeat(auto-fill, minmax(11rem, 1fr)); gap: 1rem; list-style: none; margin: 1.5rem 0; padding: 0; }
  .books li { border: 1px solid #ccc; border-radius: 8px; overflow: hidden; }
  .books a { display: block; color: inherit; text-decoration: none; }
  .books img { display: block; width: 100%; aspect-ratio: 3 / 4; object-fit: contain; background: #eee; }
  .books span, .books small { display: block; padding: .25rem .75rem; }
  .books span { font-weight: 600; }
</style>`

var loginTemplate = template.Must(template.New("login").Parse(`<!doctype html>
<html lang="en"><meta charset="utf-8"><title>Log in – secure-reader</title>` + style + `
<h1>secure-reader</h1>
<p class="lead">A sample architecture (<a href="https://github.com/shibukawa/bdf/blob/main/docs/examples/secure-reader.md">docs/examples/secure-reader</a>): books are read ten pages at a time, each sealed for the request that asked for it.</p>
<form class="login" method="post" action="/login">
  <label>User <input name="user" autocomplete="username" required></label>
  <label>Password <input name="password" type="password" autocomplete="current-password" required></label>
  <button>Log in</button>
  {{if .}}<p class="error" role="alert">{{.}}</p>{{end}}
</form>
<p><small>Demo accounts: alice / alice-pass owns both books; bob / bob-pass owns The Reading Room.</small></p>
`))

var libraryTemplate = template.Must(template.New("library").Parse(`<!doctype html>
<html lang="en"><meta charset="utf-8"><title>Your books – secure-reader</title>` + style + `
<header><h1>{{.User}}'s books</h1><form method="post" action="/logout"><button>Log out</button></form></header>
<ul class="books">
{{range .Books}}<li><a href="/read/?book={{.Name}}"><img src="/covers/{{.Name}}" alt=""><span>{{.Title}}</span>
<small>{{.Pages}} pages · {{if .Owned}}yours{{else}}sample: the first {{.Sample}} pages{{end}}</small></a></li>
{{end}}</ul>
`))
