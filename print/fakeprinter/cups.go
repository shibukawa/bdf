package fakeprinter

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/shibukawa/bdf/print/ipp"
)

// CUPS is a CUPS server of fake printers, as CUPS clients see one: each
// printer at /printers/<name>, each class at /classes/<name> (a job sent
// to a class goes to a member that is idle and accepting, in turn), and
// CUPS's own operations: CUPS-Get-Printers, CUPS-Get-Classes,
// CUPS-Get-Default, CUPS-Move-Job, CUPS-Accept-Jobs and CUPS-Reject-Jobs.
// The printers share one count of job ids, as CUPS's queues do.
type CUPS struct {
	mu       sync.Mutex
	printers []*Printer
	classes  []*class
	// def is the name of the default destination.
	def string
	ids *idSource
}

// class is a class of printers.
type class struct {
	name      string
	members   []string
	next      int
	accepting bool
	uuid      string
}

// NewCUPS returns a CUPS server without printers.
func NewCUPS() *CUPS { return &CUPS{ids: &idSource{}} }

// Add adds printers. Their job ids are the server's from now.
func (c *CUPS) Add(ps ...*Printer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, p := range ps {
		p.lock()
		p.ids = c.ids
		p.mu.Unlock()
		c.printers = append(c.printers, p)
		if c.def == "" {
			c.def = p.Name
		}
	}
}

// AddClass adds a class of the printers named.
func (c *CUPS) AddClass(name string, members ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.classes = append(c.classes, &class{name: name, members: members, accepting: true, uuid: NewUUID()})
}

// SetDefault makes a printer or a class the default destination.
func (c *CUPS) SetDefault(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.def = name
}

// Printer returns a printer by name, or nil.
func (c *CUPS) Printer(name string) *Printer {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.printer(name)
}

func (c *CUPS) printer(name string) *Printer {
	for _, p := range c.printers {
		if p.Name == name {
			return p
		}
	}
	return nil
}

func (c *CUPS) class(name string) *class {
	for _, k := range c.classes {
		if k.name == name {
			return k
		}
	}
	return nil
}

// Printers returns the printers of the server.
func (c *CUPS) Printers() []*Printer {
	c.mu.Lock()
	defer c.mu.Unlock()
	return slices.Clone(c.printers)
}

// ServeHTTP answers an IPP request to the server, a printer or a class.
func (c *CUPS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "IPP is POST", http.StatusMethodNotAllowed)
		return
	}
	kind, name, _ := strings.Cut(strings.Trim(r.URL.Path, "/"), "/")
	if i := strings.IndexByte(name, '/'); i >= 0 {
		name = name[:i] // a job's URI: /printers/<name>/<id>
	}
	c.mu.Lock()
	var target *Printer
	var cls *class
	switch kind {
	case "printers":
		target = c.printer(name)
	case "classes":
		cls = c.class(name)
	case "", "jobs", "admin":
	default:
		c.mu.Unlock()
		http.NotFound(w, r)
		return
	}
	c.mu.Unlock()
	if (kind == "printers" && target == nil) || (kind == "classes" && cls == nil) {
		http.Error(w, "no such printer", http.StatusNotFound)
		return
	}
	if target != nil && !target.admit(w, r) {
		return
	}
	path := "/" + kind + "/" + name
	if kind == "" || kind == "jobs" || kind == "admin" {
		path = "/"
	}
	body := r.Body
	q, ok := decode(w, r, body, nil, path)
	if !ok {
		return
	}
	defer q.write(w)
	switch q.req.Op() {
	case ipp.OpCupsGetPrinters, ipp.OpCupsGetClasses, ipp.OpCupsGetDefault, ipp.OpCupsMoveJob, ipp.OpCupsAcceptJobs, ipp.OpCupsRejectJobs, ipp.OpCupsSetDefault:
		if !check(q, false, false) {
			return
		}
		c.serveCUPS(q)
		return
	}
	switch {
	case target != nil:
		target.handle(q)
	case cls != nil:
		c.serveClass(q, cls)
	default:
		c.serveServer(q)
	}
}

// serveServer answers a request to the server itself ("/"): the jobs of
// all its printers, and a job's operations by its id.
func (c *CUPS) serveServer(q *request) {
	if !check(q, false, true) {
		return
	}
	switch q.req.Op() {
	case ipp.OpGetJobs:
		c.getJobs(q, c.Printers())
	case ipp.OpGetPrinterAttributes:
		q.fail(ipp.StatusClientErrorNotFound, "the server is not a printer: ask /printers/<name>")
	default:
		if p := c.holder(q); p != nil {
			p.handle(q)
			return
		}
		if q.res.Code == 0 {
			q.fail(ipp.StatusClientErrorNotFound, "no such job")
		}
	}
}

// holder is the printer that has the job a request names.
func (c *CUPS) holder(q *request) *Printer {
	id, ok := q.op.Get("job-id").Int()
	if !ok {
		if u := q.op.Get("job-uri").Text(); u != "" {
			n, _ := strconv.Atoi(u[strings.LastIndexByte(u, '/')+1:])
			id, ok = int32(n), n > 0
		}
	}
	if !ok {
		return nil
	}
	for _, p := range c.Printers() {
		if p.Job(id) != nil {
			return p
		}
	}
	return nil
}

// serveClass answers a request to a class: a new job goes to a member,
// a job's operations to the member that has it.
func (c *CUPS) serveClass(q *request, k *class) {
	if !check(q, false, true) {
		return
	}
	c.mu.Lock()
	var members []*Printer
	for _, n := range k.members {
		if p := c.printer(n); p != nil {
			members = append(members, p)
		}
	}
	c.mu.Unlock()
	if len(members) == 0 {
		q.fail(ipp.StatusServerErrorNotAcceptingJobs, "the class %s has no printers", k.name)
		return
	}
	switch q.req.Op() {
	case ipp.OpPrintJob, ipp.OpCreateJob, ipp.OpValidateJob:
		if !k.accepting {
			q.fail(ipp.StatusServerErrorNotAcceptingJobs, "the class %s is not accepting jobs", k.name)
			return
		}
		c.pick(k, members).handle(q)
	case ipp.OpGetPrinterAttributes:
		c.classAttributes(q, k, members)
	case ipp.OpGetJobs:
		c.getJobs(q, members)
	default:
		for _, p := range members {
			if j := c.holderOf(q, p); j {
				p.handle(q)
				return
			}
		}
		q.fail(ipp.StatusClientErrorNotFound, "no such job in %s", k.name)
	}
}

func (c *CUPS) holderOf(q *request, p *Printer) bool {
	id, ok := q.op.Get("job-id").Int()
	return ok && p.Job(id) != nil
}

// pick chooses the member a class's new job goes to: the next in turn
// that is idle and accepting, else the next that accepts, else the next.
func (c *CUPS) pick(k *class, members []*Printer) *Printer {
	c.mu.Lock()
	start := k.next
	k.next = (k.next + 1) % len(members)
	c.mu.Unlock()
	ready := func(p *Printer, idle bool) bool {
		p.lock()
		defer p.mu.Unlock()
		return p.Accepting && p.State != ipp.PrinterStopped && (!idle || p.queued() == 0)
	}
	for _, idle := range []bool{true, false} {
		for i := range members {
			p := members[(start+i)%len(members)]
			if ready(p, idle) {
				c.mu.Lock()
				k.next = (start + i + 1) % len(members)
				c.mu.Unlock()
				return p
			}
		}
	}
	return members[start%len(members)]
}

// classAttributes are a class's attributes: its first member's, named as
// the class, with its members.
func (c *CUPS) classAttributes(q *request, k *class, members []*Printer) {
	first := members[0]
	first.lock()
	attrs := first.attributes(q)
	first.mu.Unlock()
	var names, uris []ipp.Value
	state, accepting := ipp.PrinterStopped, false
	for _, p := range members {
		p.lock()
		names = append(names, ipp.Name(p.Name))
		uris = append(uris, ipp.URI(strings.TrimSuffix(q.printerURI(), q.path)+"/printers/"+p.Name))
		if p.State != ipp.PrinterStopped {
			state = min(state, p.State)
		}
		accepting = accepting || p.Accepting
		p.mu.Unlock()
	}
	for i, a := range attrs {
		switch a.Name {
		case "printer-name":
			attrs[i] = ipp.Attr("printer-name", ipp.Name(k.name))
		case "printer-info":
			attrs[i] = ipp.Attr("printer-info", ipp.Text(k.name))
		case "printer-uuid":
			attrs[i] = ipp.Attr("printer-uuid", ipp.URI(k.uuid))
		case "printer-state":
			attrs[i] = ipp.Attr("printer-state", ipp.Enum(int(state)))
		case "printer-is-accepting-jobs":
			attrs[i] = ipp.Attr("printer-is-accepting-jobs", ipp.Boolean(accepting && k.accepting))
		}
	}
	attrs = append(attrs, ipp.Attr("member-names", names...), ipp.Attr("member-uris", uris...),
		ipp.Attr("printer-type", ipp.Integer(cupsPrinterType(first)|1)))
	req := requested(q.op, nil)
	if !slices.Contains(req, "media-col-database") {
		attrs = slices.DeleteFunc(attrs, func(a ipp.Attribute) bool { return a.Name == "media-col-database" })
	}
	q.res.Add(ipp.TagPrinter, filter(attrs, req, printerGroups(attrs))...)
}

// getJobs answers Get-Jobs of printers: their jobs, the newest first.
func (c *CUPS) getJobs(q *request, ps []*Printer) {
	type found struct {
		id    int32
		group ipp.Group
	}
	var all []found
	for _, p := range ps {
		sub := &request{r: q.r, req: q.req, op: q.op, res: ipp.NewResponse(ipp.StatusOK, 1, "utf-8", "en"), path: "/printers/" + p.Name, secure: q.secure}
		p.lock()
		limit := q.op.Get("limit")
		p.getJobs(sub)
		p.mu.Unlock()
		_ = limit
		if sub.res.Code != 0 {
			q.res.Code = sub.res.Code
			return
		}
		for _, g := range sub.res.Groups {
			if g.Tag == ipp.TagJob {
				id, _ := g.Get("job-id").Int()
				all = append(all, found{id, g})
			}
		}
	}
	slices.SortFunc(all, func(a, b found) int { return int(b.id - a.id) })
	limit, _ := q.op.Get("limit").Int()
	for i, f := range all {
		if limit > 0 && i >= int(limit) {
			break
		}
		q.res.Groups = append(q.res.Groups, f.group)
	}
}

// serveCUPS answers CUPS's own operations.
func (c *CUPS) serveCUPS(q *request) {
	switch q.req.Op() {
	case ipp.OpCupsGetPrinters:
		// the printers, then the classes (as cupsd lists them: a class is
		// a destination too, told by the bit of its printer-type)
		for _, p := range c.Printers() {
			p.lock()
			attrs := p.cupsAttrs(q, c.uri(q, "printers", p.Name))
			p.mu.Unlock()
			q.res.Groups = append(q.res.Groups, ipp.Group{Tag: ipp.TagPrinter, Attrs: filter(attrs, requested(q.op, nil), nil)})
		}
		c.listClasses(q)
	case ipp.OpCupsGetClasses:
		c.listClasses(q)
	case ipp.OpCupsGetDefault:
		c.mu.Lock()
		def := c.def
		p, k := c.printer(def), c.class(def)
		c.mu.Unlock()
		switch {
		case p != nil:
			p.lock()
			attrs := p.cupsAttrs(q, c.uri(q, "printers", p.Name))
			p.mu.Unlock()
			q.res.Add(ipp.TagPrinter, filter(attrs, requested(q.op, nil), nil)...)
		case k != nil:
			q.res.Add(ipp.TagPrinter, ipp.Attr("printer-name", ipp.Name(k.name)), ipp.Attr("printer-uri-supported", ipp.URI(c.uri(q, "classes", k.name))),
				ipp.Attr("printer-type", ipp.Integer(1)))
		default:
			q.fail(ipp.StatusClientErrorNotFound, "no default destination")
		}
	case ipp.OpCupsSetDefault:
		_, name := c.destination(q)
		if name == "" {
			q.fail(ipp.StatusClientErrorNotFound, "no such destination")
			return
		}
		c.SetDefault(name)
	case ipp.OpCupsAcceptJobs, ipp.OpCupsRejectJobs:
		ok := q.req.Op() == ipp.OpCupsAcceptJobs
		kind, name := c.destination(q)
		c.mu.Lock()
		defer c.mu.Unlock()
		switch kind {
		case "printers":
			if p := c.printer(name); p != nil {
				p.SetAccepting(ok)
				return
			}
		case "classes":
			if k := c.class(name); k != nil {
				k.accepting = ok
				return
			}
		}
		q.fail(ipp.StatusClientErrorNotFound, "no such destination")
	case ipp.OpCupsMoveJob:
		c.moveJob(q)
	}
}

// listClasses adds the classes to an answer.
func (c *CUPS) listClasses(q *request) {
	c.mu.Lock()
	classes := slices.Clone(c.classes)
	c.mu.Unlock()
	for _, k := range classes {
		var names, uris []ipp.Value
		for _, n := range k.members {
			names = append(names, ipp.Name(n))
			uris = append(uris, ipp.URI(c.uri(q, "printers", n)))
		}
		attrs := []ipp.Attribute{
			ipp.Attr("printer-name", ipp.Name(k.name)),
			ipp.Attr("printer-uri-supported", ipp.URI(c.uri(q, "classes", k.name))),
			ipp.Attr("printer-info", ipp.Text(k.name)),
			ipp.Attr("printer-make-and-model", ipp.Text("Class of "+strings.Join(k.members, ", "))),
			ipp.Attr("printer-state", ipp.Enum(int(ipp.PrinterIdle))),
			ipp.Attr("printer-state-reasons", ipp.Keyword("none")),
			ipp.Attr("printer-is-accepting-jobs", ipp.Boolean(k.accepting)),
			ipp.Attr("printer-type", ipp.Integer(1)), // CUPS_PRINTER_CLASS
			ipp.Attr("member-names", names...),
			ipp.Attr("member-uris", uris...),
		}
		q.res.Groups = append(q.res.Groups, ipp.Group{Tag: ipp.TagPrinter, Attrs: filter(attrs, requested(q.op, nil), nil)})
	}
}

// destination is the printer or class printer-uri names.
func (c *CUPS) destination(q *request) (kind, name string) {
	u := q.op.Get("printer-uri").Text()
	for _, k := range []string{"printers", "classes"} {
		if i := strings.Index(u, "/"+k+"/"); i >= 0 {
			return k, strings.Trim(u[i+len(k)+2:], "/")
		}
	}
	return "", ""
}

// moveJob moves a job that is not done (job-id, or job-uri) to the
// printer job-printer-uri names; without a job, every job of the printer
// printer-uri names.
func (c *CUPS) moveJob(q *request) {
	g := q.req.Group(ipp.TagJob)
	if g == nil || g.Get("job-printer-uri") == nil {
		q.fail(ipp.StatusClientErrorBadRequest, "no job-printer-uri")
		return
	}
	to := g.Get("job-printer-uri").Text()
	i := strings.Index(to, "/printers/")
	if i < 0 {
		q.fail(ipp.StatusClientErrorNotPossible, "jobs move to printers")
		return
	}
	dst := c.Printer(strings.Trim(to[i+len("/printers/"):], "/"))
	if dst == nil {
		q.fail(ipp.StatusClientErrorNotFound, "no printer %s", to)
		return
	}
	var ids []int32
	var src *Printer
	if id, ok := q.op.Get("job-id").Int(); ok {
		ids = []int32{id}
		src = c.holder(q)
	} else if kind, name := c.destination(q); kind == "printers" {
		src = c.Printer(name)
		if src != nil {
			for _, j := range src.Jobs() {
				ids = append(ids, j.ID)
			}
		}
	}
	if src == nil {
		q.fail(ipp.StatusClientErrorNotFound, "no such job")
		return
	}
	if src == dst {
		return
	}
	src.lock()
	var moved []*Job
	for _, id := range ids {
		j := src.job(id)
		if j == nil || src.state(j).Done() {
			if len(ids) == 1 {
				src.mu.Unlock()
				q.fail(ipp.StatusClientErrorNotPossible, "the job is done")
				return
			}
			continue
		}
		moved = append(moved, j)
		src.jobs = slices.DeleteFunc(src.jobs, func(x *Job) bool { return x == j })
	}
	src.notify(Event{Kind: "jobs"})
	src.mu.Unlock()
	dst.lock()
	for _, j := range moved {
		j.Printer = dst.Name
		dst.jobs = append(dst.jobs, j)
		dst.notify(Event{Kind: "job", Job: j.ID})
	}
	dst.mu.Unlock()
}

// uri is the URI of a printer or class of the server as the client
// reached it.
func (c *CUPS) uri(q *request, kind, name string) string {
	scheme := "ipp"
	if q.secure {
		scheme = "ipps"
	}
	return fmt.Sprintf("%s://%s/%s/%s", scheme, q.r.Host, kind, name)
}
