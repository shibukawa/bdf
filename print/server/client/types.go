package client

import (
	"encoding/json"
	"time"
)

// Job is a job of the server. Its state and reasons are IPP's
// (job-state, job-state-reasons), with "bdf-" for those IPP has no
// keyword for.
type Job struct {
	ID       int32           `json:"id"`
	UUID     string          `json:"uuid"`
	Queue    string          `json:"queue"`
	Name     string          `json:"name"`
	User     string          `json:"user"`
	Document Document        `json:"document"`
	Ticket   json.RawMessage `json:"ticket,omitempty"`
	Options  JobOptions      `json:"options"`

	// State is pending, pending-held, processing, processing-stopped,
	// canceled, aborted or completed.
	State   string   `json:"state"`
	Reasons []string `json:"reasons,omitempty"`
	Message string   `json:"message,omitempty"`

	Priority      int       `json:"priority"`
	HoldUntil     string    `json:"hold-until,omitempty"`
	HoldUntilTime time.Time `json:"hold-until-time,omitzero"`

	// Step is what the server does next (convert, send, watch), and Due
	// when (a retry, the end of a hold).
	Step       string      `json:"step,omitempty"`
	Due        time.Time   `json:"due,omitzero"`
	Attempts   int         `json:"attempts"`
	RetrySince time.Time   `json:"retry-since,omitzero"`
	Retries    int         `json:"retries,omitempty"`
	Printer    *PrinterJob `json:"printer,omitempty"`

	Submitted time.Time `json:"submitted"`
	Started   time.Time `json:"started,omitzero"`
	Finished  time.Time `json:"finished,omitzero"`

	Sheets      int      `json:"sheets,omitempty"`
	Impressions int      `json:"impressions,omitempty"`
	Bytes       int64    `json:"bytes,omitempty"`
	Format      string   `json:"format,omitempty"`
	Media       string   `json:"media,omitempty"`
	Unsupported []string `json:"unsupported,omitempty"`
	Warnings    []string `json:"warnings,omitempty"`

	DocumentsRemoved bool `json:"documents-removed,omitempty"`
}

// Done tells a job that is canceled, aborted or completed.
func (j *Job) Done() bool {
	return j.State == "canceled" || j.State == "aborted" || j.State == "completed"
}

// Document is the document of a job as the server received it.
type Document struct {
	Name      string     `json:"name"`
	Format    string     `json:"format"`
	Size      int64      `json:"size"`
	SHA256    string     `json:"sha256"`
	Encrypted bool       `json:"encrypted,omitempty"`
	BDF       *Converted `json:"bdf,omitempty"`
}

// Converted is the .bdf the server made of a document that was not one.
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

// PrinterJob is the job of an attempt at a printer.
type PrinterJob struct {
	Name           string    `json:"name"`
	URL            string    `json:"url"`
	AttemptStarted time.Time `json:"attempt-started"`
	JobID          int32     `json:"job-id,omitempty"`
	JobURI         string    `json:"job-uri,omitempty"`
	Sent           bool      `json:"sent,omitempty"`
}

// Queue is a queue of the server.
type Queue struct {
	Name     string   `json:"name"`
	Printers []string `json:"printers"`
	Policy   string   `json:"policy"`
	// Paused queues take jobs and do not send them; Rejecting ones do
	// not take jobs.
	Paused    bool `json:"paused"`
	Rejecting bool `json:"rejecting"`
	// Jobs counts the jobs not done, by state.
	Jobs map[string]int `json:"jobs"`
}

// Printer is a printer of the server and what it said last.
type Printer struct {
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Up        bool      `json:"up"`
	State     string    `json:"state,omitempty"`
	Reasons   []string  `json:"reasons,omitempty"`
	Message   string    `json:"message,omitempty"`
	Accepting bool      `json:"accepting"`
	MakeModel string    `json:"make-and-model,omitempty"`
	Checked   time.Time `json:"checked,omitzero"`
	Error     string    `json:"error,omitempty"`
	// Jobs are the jobs sent to it and not done.
	Jobs int `json:"jobs"`
	// Capabilities are what it supports (Client.Printer only).
	Capabilities *Capabilities `json:"capabilities,omitempty"`
}

// Capabilities are what a printer supports, with IPP's names.
type Capabilities struct {
	Formats          []string `json:"document-format-supported,omitempty"`
	Media            []string `json:"media-supported,omitempty"`
	MediaReady       []string `json:"media-ready,omitempty"`
	MediaDefault     string   `json:"media-default,omitempty"`
	MediaSources     []string `json:"media-source-supported,omitempty"`
	MediaTypes       []string `json:"media-type-supported,omitempty"`
	Sides            []string `json:"sides-supported,omitempty"`
	SidesDefault     string   `json:"sides-default,omitempty"`
	ColorModes       []string `json:"print-color-mode-supported,omitempty"`
	ColorModeDefault string   `json:"print-color-mode-default,omitempty"`
	Qualities        []int32  `json:"print-quality-supported,omitempty"`
	Finishings       []int32  `json:"finishings-supported,omitempty"`
	OutputBins       []string `json:"output-bin-supported,omitempty"`
	CopiesMax        int32    `json:"copies-max,omitempty"`
	Location         string   `json:"printer-location,omitempty"`
	Info             string   `json:"printer-info,omitempty"`
}
