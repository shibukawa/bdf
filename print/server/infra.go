package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/spool"
	"github.com/shibukawa/bdf/vector/pdfbdf"
)

// The IPP Shared Infrastructure Extensions (PWG 5100.18), as far as a
// proxy of a branch needs them: a queue whose printer is one of proxies is
// an infrastructure printer, whose jobs the proxies (bdf printd with
// upstreams, or another proxy of INFRA) fetch and report. Proxies are
// tokens of the role proxy, with Basic authentication, whatever ipp.auth
// says: a proxy reads others' documents.

// infraOps are the operations of proxies.
var infraOps = map[ipp.Op]bool{
	ipp.OpUpdateOutputDeviceAttributes: true, ipp.OpDeregisterOutputDevice: true, ipp.OpFetchJob: true, ipp.OpAcknowledgeJob: true,
	ipp.OpFetchDocument: true, ipp.OpAcknowledgeDocument: true, ipp.OpUpdateJobStatus: true, ipp.OpUpdateActiveJobs: true,
	ipp.OpGetOutputDeviceAttributes: true, ipp.OpUpdateDocumentStatus: true,
}

// bdfFormat is the media type of .bdf documents.
const bdfFormat = "application/x-bdf"

// isInfra tells a request of a proxy: an operation of INFRA, or Get-Jobs
// of the fetchable jobs.
func isInfra(q *ippRequest) bool {
	op := q.req.Op()
	return infraOps[op] || op == ipp.OpGetJobs && q.op.Get("which-jobs").Text() == "fetchable"
}

// authorizeProxy finds the proxy of a request: a token of the role
// proxy that may use the queue. It writes the 401 itself.
func (x *ippServer) authorizeProxy(w http.ResponseWriter, q *ippRequest) bool {
	name, secret, ok := q.r.BasicAuth()
	if ok {
		q.caller = x.s.auth.basic(name, secret)
	}
	if q.caller == nil || !q.caller.can(roleProxy) || !q.caller.mayQueue(q.queue.name) {
		w.Header().Set("WWW-Authenticate", `Basic realm="bdf"`)
		http.Error(w, "a proxy's token is required", http.StatusUnauthorized)
		return false
	}
	return true
}

// proxyPrinter is the printer of proxies of a request's queue, and the
// output-device-uuid of the proxy.
func (x *ippServer) proxyPrinter(q *ippRequest) (printer, device string, ok bool) {
	qc := x.s.cfg.Queues[q.queue.name]
	if len(qc.Printers) != 1 || !x.s.sp.IsProxy(qc.Printers[0]) {
		q.fail(ipp.StatusServerErrorOperationNotSupported, "%s is not a printer of proxies", q.queue.name)
		return "", "", false
	}
	device = q.op.Get("output-device-uuid").Text()
	if !strings.HasPrefix(device, "urn:uuid:") || len(device) > 64 {
		q.fail(ipp.StatusClientErrorBadRequest, "output-device-uuid is a urn:uuid")
		return "", "", false
	}
	return qc.Printers[0], device, true
}

// infra serves an operation of proxies.
func (x *ippServer) infra(q *ippRequest) {
	printer, device, ok := x.proxyPrinter(q)
	if !ok {
		return
	}
	ctx := q.r.Context()
	switch q.req.Op() {
	case ipp.OpUpdateOutputDeviceAttributes:
		var caps *ipp.Printer
		if g := q.req.Group(ipp.TagPrinter); g != nil {
			caps = ipp.ParsePrinter(&ipp.Message{Groups: []ipp.Group{*g}})
		}
		if err := x.s.sp.ProxySeen(printer, device, caps); err != nil {
			x.failSpool(q, err)
		}
	case ipp.OpDeregisterOutputDevice:
		x.s.sp.ProxyGone(printer, device)
	case ipp.OpGetOutputDeviceAttributes:
		q.fail(ipp.StatusServerErrorOperationNotSupported, "")
	case ipp.OpGetJobs:
		x.fetchableJobs(q, printer, device)
	case ipp.OpFetchJob:
		id, ok := x.jobID(q)
		if !ok {
			return
		}
		i := slices.IndexFunc(x.s.sp.Fetchable(printer, device), func(j spool.Job) bool { return j.ID == id })
		if i < 0 {
			q.fail(ipp.StatusClientErrorNotFetchable, "job %d is not fetchable", id)
			return
		}
		j, err := x.s.sp.Get(ctx, id)
		if err != nil {
			x.failSpool(q, err)
			return
		}
		attrs := x.jobAttrs(q, &j, true)
		// the queue's defaults go with the ticket: the proxy's queue has
		// its own
		defaults, _ := x.s.cfg.Queues[j.Queue].ticket()
		t, _ := json.Marshal(j.Ticket.Merge(defaults))
		o, _ := json.Marshal(j.Options)
		attrs = append(attrs, ipp.Attr("bdf-ticket", ipp.Text(string(t))), ipp.Attr("bdf-job-options", ipp.Text(string(o))),
			ipp.Attr("job-priority", ipp.Integer(j.Priority)))
		q.res.Add(ipp.TagJob, dedupe(attrs)...)
	case ipp.OpAcknowledgeJob:
		id, ok := x.jobID(q)
		if !ok {
			return
		}
		if _, err := x.s.sp.AcknowledgeJob(ctx, id, printer, device); err != nil {
			if errors.Is(err, spool.ErrNotPossible) {
				q.fail(ipp.StatusClientErrorNotFetchable, "%v", err)
				return
			}
			x.failSpool(q, err)
			return
		}
		x.s.log.Info("ipp: a proxy took a job", "job", id, "printer", printer, "device", device, "user", q.caller.user)
	case ipp.OpFetchDocument:
		x.fetchDocument(q, printer, device)
	case ipp.OpAcknowledgeDocument:
		id, ok := x.jobID(q)
		if !ok {
			return
		}
		if _, err := x.s.sp.AcknowledgeDocument(ctx, id, printer, device); err != nil {
			x.failSpool(q, err)
		}
	case ipp.OpUpdateDocumentStatus:
		// one document a job: its job's status says it all
	case ipp.OpUpdateJobStatus:
		id, ok := x.jobID(q)
		if !ok {
			return
		}
		g := q.req.Group(ipp.TagJob)
		if g == nil {
			g = &ipp.Group{}
		}
		n, isEnum := g.Get("output-device-job-state").Int()
		if !isEnum || n < int32(ipp.JobPending) || n > int32(ipp.JobCompleted) {
			q.fail(ipp.StatusClientErrorBadRequest, "output-device-job-state is a job-state")
			return
		}
		impressions, _ := g.Get("job-impressions-completed").Int()
		j, err := x.s.sp.ProxyReport(ctx, id, printer, device, ipp.JobState(n), g.Get("output-device-job-state-reasons").Strings(),
			int(impressions), g.Get("output-device-job-state-message").Text())
		if err != nil && !(errors.Is(err, spool.ErrNotPossible) && j.State.Done()) {
			x.failSpool(q, err)
			return
		}
		// the job's state here: the proxy cancels a job canceled here
		q.res.Add(ipp.TagJob, ipp.Attr("job-id", ipp.Integer(int(j.ID))), ipp.Attr("job-state", ipp.Enum(int(j.State))),
			ipp.Attr("job-state-reasons", keywordValues(publicReasons(orNone(j.Reasons)))...))
	case ipp.OpUpdateActiveJobs:
		x.updateActiveJobs(q, printer, device)
	}
}

func orNone(r []string) []string {
	if len(r) == 0 {
		return []string{"none"}
	}
	return r
}

// dedupe keeps the last of attributes of one name.
func dedupe(attrs []ipp.Attribute) []ipp.Attribute {
	var out []ipp.Attribute
	for i, a := range attrs {
		if slices.IndexFunc(attrs[i+1:], func(b ipp.Attribute) bool { return b.Name == a.Name }) < 0 {
			out = append(out, a)
		}
	}
	return out
}

// jobID is the job-id of a request of a proxy.
func (x *ippServer) jobID(q *ippRequest) (int32, bool) {
	n, ok := q.op.Get("job-id").Int()
	if !ok || n <= 0 {
		q.fail(ipp.StatusClientErrorBadRequest, "job-id is a positive integer")
		return 0, false
	}
	return n, true
}

// fetchableJobs lists the jobs a proxy may fetch.
func (x *ippServer) fetchableJobs(q *ippRequest, printer, device string) {
	limit := 0
	if n, ok := q.op.Get("limit").Int(); ok && n > 0 {
		limit = int(n)
	}
	req := requested(q.op, []string{"job-id", "job-uri"})
	for i, j := range x.s.sp.Fetchable(printer, device) {
		if limit > 0 && i >= limit {
			break
		}
		attrs := x.jobAttrs(q, &j, true)
		for k := range attrs {
			if attrs[k].Name == "job-state-reasons" {
				attrs[k] = ipp.Attr("job-state-reasons", ipp.Keyword("job-fetchable"))
			}
		}
		q.res.Groups = append(q.res.Groups, ipp.Group{Tag: ipp.TagJob, Attrs: filterAttrs(attrs, req, jobGroups)})
	}
}

// fetchDocument sends a proxy the document of a job it acknowledged: the
// .bdf as it is kept (with the password of a locked one) when the proxy
// takes .bdf, or a PDF of it.
func (x *ippServer) fetchDocument(q *ippRequest, printer, device string) {
	id, ok := x.jobID(q)
	if !ok {
		return
	}
	if n, ok := q.op.Get("document-number").Int(); ok && n != 1 {
		q.fail(ipp.StatusClientErrorNotFound, "a job has one document")
		return
	}
	ctx := q.r.Context()
	bound := false
	for _, j := range x.s.sp.ProxyJobs(printer, device) {
		bound = bound || j.ID == id
	}
	if !bound {
		q.fail(ipp.StatusClientErrorNotFetchable, "job %d is not this proxy's", id)
		return
	}
	accepted := q.op.Get("document-format-accepted").Strings()
	switch {
	case len(accepted) == 0 || slices.Contains(accepted, bdfFormat):
		rc, size, err := x.s.sp.OpenBDF(ctx, id)
		if err != nil {
			x.failSpool(q, err)
			return
		}
		q.res.Add(ipp.TagOperation, ipp.Attr("document-format", ipp.MimeType(bdfFormat)), ipp.Attr("compression", ipp.Keyword("none")))
		if pw := x.s.sp.Password(id); pw != "" {
			q.res.Add(ipp.TagOperation, ipp.Attr("document-password", ipp.OctetString([]byte(pw))))
		}
		q.after = func(w io.Writer) error {
			defer rc.Close()
			_, err := io.CopyN(w, rc, size)
			return err
		}
	case slices.Contains(accepted, "application/pdf"):
		doc, err := x.s.sp.Document(ctx, id)
		if err != nil {
			x.failSpool(q, err)
			return
		}
		j, _ := x.s.sp.Get(ctx, id)
		q.res.Add(ipp.TagOperation, ipp.Attr("document-format", ipp.MimeType("application/pdf")), ipp.Attr("compression", ipp.Keyword("none")))
		var views []string
		if j.Options.View != "" {
			views = []string{j.Options.View}
		}
		q.after = func(w io.Writer) error {
			_, err := pdfbdf.Write(w, doc, &pdfbdf.Options{Views: views, SkipRoles: j.Options.SkipRoles, FontFS: x.s.o.FontFS,
				FontDirs: x.s.cfg.Fonts.Dirs, NoSystemFonts: x.s.cfg.Fonts.NoSystem})
			return err
		}
	default:
		q.fail(ipp.StatusClientErrorDocumentFormatNotSupported, "the documents are application/x-bdf or application/pdf")
	}
}

// updateActiveJobs puts the jobs of a proxy right after either stopped:
// a job the proxy has that ended here is told (the proxy cancels it), a
// job the proxy reports done ends here, and a job the proxy had the
// document of and no longer has is aborted.
func (x *ippServer) updateActiveJobs(q *ippRequest, printer, device string) {
	ctx := q.r.Context()
	ids := q.op.Get("job-ids").Ints()
	states := q.op.Get("output-device-job-states").Ints()
	if len(ids) != len(states) {
		q.fail(ipp.StatusClientErrorBadRequest, "job-ids and output-device-job-states differ in number")
		return
	}
	var endedIDs, endedStates []ipp.Value
	for i, id := range ids {
		st := ipp.JobState(states[i])
		j, err := x.s.sp.Get(ctx, id)
		if err != nil {
			continue
		}
		if j.State.Done() {
			if !st.Done() {
				endedIDs, endedStates = append(endedIDs, ipp.Integer(int(id))), append(endedStates, ipp.Enum(int(j.State)))
			}
			continue
		}
		if st.Done() {
			x.s.sp.ProxyReport(ctx, id, printer, device, st, nil, 0, "")
		}
	}
	for _, j := range x.s.sp.ProxyJobs(printer, device) {
		if j.Step == spool.StepRemote && !slices.Contains(ids, j.ID) {
			x.s.log.Warn("ipp: a proxy no longer has a job", "job", j.ID, "printer", printer, "device", device)
			x.s.sp.ProxyLost(ctx, j.ID, printer, device)
		}
	}
	if len(endedIDs) > 0 {
		q.res.Add(ipp.TagOperation, ipp.Attr("job-ids", endedIDs...), ipp.Attr("output-device-job-states", endedStates...))
	}
}
