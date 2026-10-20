package spool

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/shibukawa/bdf/print"
	"github.com/shibukawa/bdf/print/ipp"
)

// Step is the part of a job's processing it is at: each is short and can
// be started again after a crash.
type Step string

// Steps.
const (
	// StepConvert turns a document that is not a .bdf into one.
	StepConvert Step = "convert"
	// StepSend makes the job at a printer and sends the document.
	StepSend Step = "send"
	// StepWatch follows the job at the printer until it is done.
	StepWatch Step = "watch"
)

// Job is a job of the spool, as it is kept. Its state and reasons are
// IPP's (job-state, job-state-reasons), with "bdf-" for the reasons IPP
// has no keyword for.
type Job struct {
	// ID is a positive integer, as IPP's job-id; UUID is its job-uuid.
	ID   int32  `json:"id"`
	UUID string `json:"uuid"`
	// Queue is the queue the job is in; Name its name (the document's
	// file name, usually); User who submitted it.
	Queue string `json:"queue"`
	Name  string `json:"name"`
	User  string `json:"user"`

	Document Document     `json:"document"`
	Ticket   print.Ticket `json:"ticket"`
	Options  JobOptions   `json:"options"`

	State   ipp.JobState `json:"state"`
	Reasons []string     `json:"reasons,omitempty"`
	Message string       `json:"message,omitempty"`

	// Priority is job-priority, 1 to 100 (50 by default): the higher
	// the sooner.
	Priority int `json:"priority"`
	// HoldUntil is job-hold-until ("indefinite", "night" …) and
	// HoldUntilTime job-hold-until-time.
	HoldUntil     string    `json:"hold-until,omitempty"`
	HoldUntilTime time.Time `json:"hold-until-time,omitzero"`

	// Step is the step to do next, and Due when (a retry, the end of a
	// hold). A job that is done has no step.
	Step Step      `json:"step,omitempty"`
	Due  time.Time `json:"due,omitzero"`
	// Attempts are the times the job was sent; RetrySince is when the
	// printer was first not reached for the current attempt.
	Attempts   int       `json:"attempts"`
	RetrySince time.Time `json:"retry-since,omitzero"`
	// Retries counts the tries since RetrySince, for the delay of the next.
	Retries int `json:"retries,omitempty"`
	// Printer is the printer of the current attempt.
	Printer *PrinterJob `json:"printer,omitempty"`

	Submitted time.Time `json:"submitted"`
	Started   time.Time `json:"started,omitzero"`
	Finished  time.Time `json:"finished,omitzero"`

	// What was sent: sheets (pages with the copies the client made),
	// impressions the printer reported, bytes, the format and media.
	Sheets      int      `json:"sheets,omitempty"`
	Impressions int      `json:"impressions,omitempty"`
	Bytes       int64    `json:"bytes,omitempty"`
	Format      string   `json:"format,omitempty"`
	Media       string   `json:"media,omitempty"`
	Unsupported []string `json:"unsupported,omitempty"`
	Warnings    []string `json:"warnings,omitempty"`

	// DocumentsRemoved is set on a finished job once its documents are
	// removed, after their retention.
	DocumentsRemoved bool `json:"documents-removed,omitempty"`
}

// Document is the document of a job as it was received.
type Document struct {
	// Name is the file name it came with.
	Name string `json:"name"`
	// Format is "bdf", or the name of the converter's format ("pdf",
	// "xlsx" …) of a document that is converted.
	Format    string `json:"format"`
	Size      int64  `json:"size"`
	SHA256    string `json:"sha256"`
	Encrypted bool   `json:"encrypted,omitempty"`
	// BDF is the document converted to .bdf, once it is.
	BDF *Converted `json:"bdf,omitempty"`
}

// Converted is the .bdf made of a document that was not one.
type Converted struct {
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
}

// JobOptions are how a job is printed, besides its ticket.
type JobOptions struct {
	// Format is "pdf", "pwg-raster" or "" for the queue's.
	Format string `json:"format,omitempty"`
	// View is the id of the view to print; "" the first view of pages.
	View string `json:"view,omitempty"`
	// SkipRoles are the roles of the layers left out.
	SkipRoles []string `json:"skip,omitempty"`
	// Strict makes what the printer does not support an error.
	Strict bool `json:"strict,omitempty"`
	// RasterIfLossy overrides the queue's when set.
	RasterIfLossy *bool `json:"raster-if-lossy,omitempty"`
}

// PrinterJob is the job of an attempt at a printer. JobID is set as soon
// as the printer answers Create-Job, before the document is sent, and
// Sent once the printer has the whole document: after a crash they tell
// whether the printer has the job already.
type PrinterJob struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	// AttemptStarted is when the attempt began, before Create-Job.
	AttemptStarted time.Time `json:"attempt-started"`
	JobID          int32     `json:"job-id,omitempty"`
	JobURI         string    `json:"job-uri,omitempty"`
	// Begun says that a raw printer may have begun to print the job: what
	// it may print was being sent.
	Begun bool `json:"begun,omitempty"`
	Sent  bool `json:"sent,omitempty"`
}

// Reasons of bdf's own (IPP's job-state-reasons have no keyword for them).
const (
	// ReasonRetryWait: the printer was not reached; the job is sent again
	// at Due.
	ReasonRetryWait = "bdf-retry-wait"
	// ReasonRetryExhausted: the printer was not reached for as long as the
	// spool retries.
	ReasonRetryExhausted = "bdf-retry-exhausted"
	// ReasonUnsupported: the printer does not support what the ticket
	// asks, and the job is strict.
	ReasonUnsupported = "bdf-unsupported"
	// ReasonDuplicateCheck: after a crash, the printer has more than one
	// job that may be this one; a person decides.
	ReasonDuplicateCheck = "bdf-duplicate-check"
)

// key is the job's id in the store: ten digits, so that the store's order
// is the ids' order.
func (j *Job) key() string { return jobKey(j.ID) }

func jobKey(id int32) string { return fmt.Sprintf("%010d", id) }

// parseKey returns the id of a key, or 0.
func parseKey(k string) int32 {
	n, err := strconv.ParseInt(k, 10, 32)
	if err != nil || n <= 0 {
		return 0
	}
	return int32(n)
}

// originalKey and bdfKey are the store's ids of a job's documents.
func originalKey(id int32) string { return jobKey(id) + ".original" }
func bdfKey(id int32) string      { return jobKey(id) + ".bdf" }

// printerJobName is the job-name at the printer: the attempt first, so
// that a printer that cuts long names keeps what tells the job apart
// after a crash.
func printerJobName(j *Job) string { return attemptPrefix(j) + j.Name }

// attemptPrefix is what the printer's job-name of an attempt begins with.
func attemptPrefix(j *Job) string { return fmt.Sprintf("#%d.%d ", j.ID, j.Attempts) }

// transitions are the states a job may go to from each state.
var transitions = map[ipp.JobState][]ipp.JobState{
	ipp.JobPending:           {ipp.JobPending, ipp.JobPendingHeld, ipp.JobProcessing, ipp.JobCanceled, ipp.JobAborted},
	ipp.JobPendingHeld:       {ipp.JobPendingHeld, ipp.JobPending, ipp.JobCanceled, ipp.JobAborted},
	ipp.JobProcessing:        {ipp.JobProcessing, ipp.JobPending, ipp.JobPendingHeld, ipp.JobProcessingStopped, ipp.JobCanceled, ipp.JobAborted, ipp.JobCompleted},
	ipp.JobProcessingStopped: {ipp.JobProcessingStopped, ipp.JobProcessing, ipp.JobCanceled, ipp.JobAborted, ipp.JobCompleted},
	// a job that is done goes back only by Retry
	ipp.JobCanceled: {ipp.JobPending},
	ipp.JobAborted:  {ipp.JobPending},
}

// to moves a job to a state with its reasons, refusing a move the state
// machine has not, and sets Finished when the job is done.
func (j *Job) to(now time.Time, st ipp.JobState, reasons ...string) error {
	if !slices.Contains(transitions[j.State], st) {
		return fmt.Errorf("spool: job %d cannot go from %s to %s", j.ID, j.State, st)
	}
	j.State = st
	j.Reasons = reasons
	if st.Done() {
		j.Finished = now
		j.Step = ""
		j.Due = time.Time{}
	} else {
		j.Finished = time.Time{}
	}
	return nil
}

// clone copies a job deeply enough that the copy can be changed without
// changing the job.
func (j *Job) clone() Job {
	c := *j
	c.Reasons = slices.Clone(j.Reasons)
	c.Unsupported = slices.Clone(j.Unsupported)
	c.Warnings = slices.Clone(j.Warnings)
	if j.Printer != nil {
		p := *j.Printer
		c.Printer = &p
	}
	if j.Document.BDF != nil {
		b := *j.Document.BDF
		c.Document.BDF = &b
	}
	return c
}
