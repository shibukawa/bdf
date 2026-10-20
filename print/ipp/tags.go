package ipp

import (
	"fmt"
	"strconv"
)

// Tag is the tag of a value or of a group of attributes in the encoding
// of a message (RFC 8010 §3.5). The tags below 0x10 begin groups; the
// others name the type of a value, or that a value is out of band.
type Tag uint8

// Group tags (RFC 8010 §3.5.1).
const (
	TagOperation         Tag = 0x01
	TagJob               Tag = 0x02
	TagEnd               Tag = 0x03 // end-of-attributes
	TagPrinter           Tag = 0x04
	TagUnsupportedGroup  Tag = 0x05
	TagSubscription      Tag = 0x06
	TagEventNotification Tag = 0x07
	TagResource          Tag = 0x08
	TagDocument          Tag = 0x09
	TagSystem            Tag = 0x0a
)

// Out-of-band value tags: a value that is not there, and why (RFC 8010
// §3.5.2, RFC 8011 §5.1.1).
const (
	TagUnsupported     Tag = 0x10
	TagDefault         Tag = 0x11
	TagUnknown         Tag = 0x12
	TagNoValue         Tag = 0x13
	TagNotSettable     Tag = 0x15
	TagDeleteAttribute Tag = 0x16
	TagAdminDefine     Tag = 0x17
)

// Value tags of the types of values (RFC 8010 §3.5.2).
const (
	TagInteger Tag = 0x21
	TagBoolean Tag = 0x22
	TagEnum    Tag = 0x23

	TagOctetString     Tag = 0x30
	TagDateTime        Tag = 0x31
	TagResolution      Tag = 0x32
	TagRange           Tag = 0x33 // rangeOfInteger
	TagBeginCollection Tag = 0x34
	TagTextLang        Tag = 0x35 // textWithLanguage
	TagNameLang        Tag = 0x36 // nameWithLanguage
	TagEndCollection   Tag = 0x37

	TagText           Tag = 0x41 // textWithoutLanguage
	TagName           Tag = 0x42 // nameWithoutLanguage
	TagReservedString Tag = 0x43
	TagKeyword        Tag = 0x44
	TagURI            Tag = 0x45
	TagURIScheme      Tag = 0x46
	TagCharset        Tag = 0x47
	TagLanguage       Tag = 0x48 // naturalLanguage
	TagMimeType       Tag = 0x49 // mimeMediaType
	TagMemberName     Tag = 0x4a // memberAttrName

	TagExtension Tag = 0x7f
)

var tagNames = map[Tag]string{
	TagOperation: "operation-attributes", TagJob: "job-attributes", TagEnd: "end-of-attributes",
	TagPrinter: "printer-attributes", TagUnsupportedGroup: "unsupported-attributes",
	TagSubscription: "subscription-attributes", TagEventNotification: "event-notification-attributes",
	TagResource: "resource-attributes", TagDocument: "document-attributes", TagSystem: "system-attributes",
	TagUnsupported: "unsupported", TagDefault: "default", TagUnknown: "unknown", TagNoValue: "no-value",
	TagNotSettable: "not-settable", TagDeleteAttribute: "delete-attribute", TagAdminDefine: "admin-define",
	TagInteger: "integer", TagBoolean: "boolean", TagEnum: "enum",
	TagOctetString: "octetString", TagDateTime: "dateTime", TagResolution: "resolution", TagRange: "rangeOfInteger",
	TagBeginCollection: "collection", TagTextLang: "textWithLanguage", TagNameLang: "nameWithLanguage",
	TagEndCollection: "endCollection",
	TagText:          "text", TagName: "name", TagReservedString: "reservedString", TagKeyword: "keyword", TagURI: "uri",
	TagURIScheme: "uriScheme", TagCharset: "charset", TagLanguage: "naturalLanguage", TagMimeType: "mimeMediaType",
	TagMemberName: "memberAttrName", TagExtension: "extension",
}

// String returns the name of the tag in the registry, or its number.
func (t Tag) String() string {
	if s, ok := tagNames[t]; ok {
		return s
	}
	return "tag-0x" + strconv.FormatUint(uint64(t), 16)
}

// IsGroup reports whether the tag begins a group of attributes.
func (t Tag) IsGroup() bool { return t < 0x10 }

// IsOutOfBand reports whether the tag stands for a value that is not
// there.
func (t Tag) IsOutOfBand() bool { return t >= 0x10 && t < 0x20 }

// IsString reports whether the value of the tag is a character string
// (RFC 8010's character-string types, 0x40 to 0x5f).
func (t Tag) IsString() bool { return t >= 0x40 && t < 0x60 && t != TagMemberName }

// Op is the operation of a request (RFC 8011 §5.4.15 and the IANA IPP
// registry).
type Op uint16

// Operations.
const (
	OpPrintJob                   Op = 0x0002
	OpPrintURI                   Op = 0x0003
	OpValidateJob                Op = 0x0004
	OpCreateJob                  Op = 0x0005
	OpSendDocument               Op = 0x0006
	OpSendURI                    Op = 0x0007
	OpCancelJob                  Op = 0x0008
	OpGetJobAttributes           Op = 0x0009
	OpGetJobs                    Op = 0x000a
	OpGetPrinterAttributes       Op = 0x000b
	OpHoldJob                    Op = 0x000c
	OpReleaseJob                 Op = 0x000d
	OpRestartJob                 Op = 0x000e
	OpPausePrinter               Op = 0x0010
	OpResumePrinter              Op = 0x0011
	OpPurgeJobs                  Op = 0x0012
	OpSetPrinterAttributes       Op = 0x0013
	OpSetJobAttributes           Op = 0x0014
	OpGetPrinterSupportedValues  Op = 0x0015
	OpCreatePrinterSubscriptions Op = 0x0016
	OpCreateJobSubscriptions     Op = 0x0017
	OpGetSubscriptionAttributes  Op = 0x0018
	OpGetSubscriptions           Op = 0x0019
	OpRenewSubscription          Op = 0x001a
	OpCancelSubscription         Op = 0x001b
	OpGetNotifications           Op = 0x001c
	OpHoldNewJobs                Op = 0x0022
	OpReleaseHeldNewJobs         Op = 0x0023
	OpCancelCurrentJob           Op = 0x0027
	OpCancelJobs                 Op = 0x0038
	OpCancelMyJobs               Op = 0x0039
	OpResubmitJob                Op = 0x003a
	OpCloseJob                   Op = 0x003b
	OpIdentifyPrinter            Op = 0x003c
	OpValidateDocument           Op = 0x003d

	// Operations of the IPP Shared Infrastructure Extensions (PWG
	// 5100.18): a proxy fetches jobs from an infrastructure printer and
	// reports them.
	OpAcknowledgeDocument          Op = 0x003f
	OpAcknowledgeIdentifyPrinter   Op = 0x0040
	OpAcknowledgeJob               Op = 0x0041
	OpFetchDocument                Op = 0x0042
	OpFetchJob                     Op = 0x0043
	OpGetOutputDeviceAttributes    Op = 0x0044
	OpUpdateActiveJobs             Op = 0x0045
	OpDeregisterOutputDevice       Op = 0x0046
	OpUpdateDocumentStatus         Op = 0x0047
	OpUpdateJobStatus              Op = 0x0048
	OpUpdateOutputDeviceAttributes Op = 0x0049

	// Private operations of CUPS (its cups-ipp registrations).
	OpCupsGetDefault       Op = 0x4001
	OpCupsGetPrinters      Op = 0x4002
	OpCupsAddModifyPrinter Op = 0x4003
	OpCupsDeletePrinter    Op = 0x4004
	OpCupsGetClasses       Op = 0x4005
)

var opNames = map[Op]string{
	OpPrintJob: "Print-Job", OpPrintURI: "Print-URI", OpValidateJob: "Validate-Job", OpCreateJob: "Create-Job",
	OpSendDocument: "Send-Document", OpSendURI: "Send-URI", OpCancelJob: "Cancel-Job", OpGetJobAttributes: "Get-Job-Attributes",
	OpGetJobs: "Get-Jobs", OpGetPrinterAttributes: "Get-Printer-Attributes", OpHoldJob: "Hold-Job", OpReleaseJob: "Release-Job",
	OpRestartJob: "Restart-Job", OpPausePrinter: "Pause-Printer", OpResumePrinter: "Resume-Printer", OpPurgeJobs: "Purge-Jobs",
	OpSetPrinterAttributes: "Set-Printer-Attributes", OpSetJobAttributes: "Set-Job-Attributes",
	OpGetPrinterSupportedValues: "Get-Printer-Supported-Values", OpCreatePrinterSubscriptions: "Create-Printer-Subscriptions",
	OpCreateJobSubscriptions: "Create-Job-Subscriptions", OpGetSubscriptionAttributes: "Get-Subscription-Attributes",
	OpGetSubscriptions: "Get-Subscriptions", OpRenewSubscription: "Renew-Subscription", OpCancelSubscription: "Cancel-Subscription",
	OpGetNotifications: "Get-Notifications", OpHoldNewJobs: "Hold-New-Jobs", OpReleaseHeldNewJobs: "Release-Held-New-Jobs",
	OpCancelCurrentJob: "Cancel-Current-Job", OpCancelJobs: "Cancel-Jobs", OpCancelMyJobs: "Cancel-My-Jobs",
	OpResubmitJob: "Resubmit-Job", OpCloseJob: "Close-Job", OpIdentifyPrinter: "Identify-Printer", OpValidateDocument: "Validate-Document",
	OpAcknowledgeDocument: "Acknowledge-Document", OpAcknowledgeIdentifyPrinter: "Acknowledge-Identify-Printer",
	OpAcknowledgeJob: "Acknowledge-Job", OpFetchDocument: "Fetch-Document", OpFetchJob: "Fetch-Job",
	OpGetOutputDeviceAttributes: "Get-Output-Device-Attributes", OpUpdateActiveJobs: "Update-Active-Jobs",
	OpDeregisterOutputDevice: "Deregister-Output-Device", OpUpdateDocumentStatus: "Update-Document-Status",
	OpUpdateJobStatus: "Update-Job-Status", OpUpdateOutputDeviceAttributes: "Update-Output-Device-Attributes",
	OpCupsGetDefault: "CUPS-Get-Default", OpCupsGetPrinters: "CUPS-Get-Printers", OpCupsAddModifyPrinter: "CUPS-Add-Modify-Printer",
	OpCupsDeletePrinter: "CUPS-Delete-Printer", OpCupsGetClasses: "CUPS-Get-Classes",
}

// String returns the name of the operation, or its number.
func (o Op) String() string {
	if s, ok := opNames[o]; ok {
		return s
	}
	return "operation-0x" + strconv.FormatUint(uint64(o), 16)
}

// Status is the status code of a response (RFC 8011 §5.4.15 and
// appendix B).
type Status uint16

// Status codes.
const (
	StatusOK                     Status = 0x0000
	StatusOKIgnoredOrSubstituted Status = 0x0001
	StatusOKConflicting          Status = 0x0002
	StatusOKIgnoredSubscriptions Status = 0x0003
	StatusOKTooManyEvents        Status = 0x0005
	StatusOKEventsComplete       Status = 0x0007

	StatusClientErrorBadRequest                     Status = 0x0400
	StatusClientErrorForbidden                      Status = 0x0401
	StatusClientErrorNotAuthenticated               Status = 0x0402
	StatusClientErrorNotAuthorized                  Status = 0x0403
	StatusClientErrorNotPossible                    Status = 0x0404
	StatusClientErrorTimeout                        Status = 0x0405
	StatusClientErrorNotFound                       Status = 0x0406
	StatusClientErrorGone                           Status = 0x0407
	StatusClientErrorRequestEntityTooLarge          Status = 0x0408
	StatusClientErrorRequestValueTooLong            Status = 0x0409
	StatusClientErrorDocumentFormatNotSupported     Status = 0x040a
	StatusClientErrorAttributesOrValuesNotSupported Status = 0x040b
	StatusClientErrorURISchemeNotSupported          Status = 0x040c
	StatusClientErrorCharsetNotSupported            Status = 0x040d
	StatusClientErrorConflictingAttributes          Status = 0x040e
	StatusClientErrorCompressionNotSupported        Status = 0x040f
	StatusClientErrorCompressionError               Status = 0x0410
	StatusClientErrorDocumentFormatError            Status = 0x0411
	StatusClientErrorDocumentAccessError            Status = 0x0412
	StatusClientErrorAttributesNotSettable          Status = 0x0413
	StatusClientErrorIgnoreAllSubscriptions         Status = 0x0414
	StatusClientErrorTooManySubscriptions           Status = 0x0415
	StatusClientErrorDocumentPasswordError          Status = 0x0418
	StatusClientErrorDocumentPermissionError        Status = 0x0419
	StatusClientErrorDocumentSecurityError          Status = 0x041a
	StatusClientErrorDocumentUnprintableError       Status = 0x041b
	StatusClientErrorAccountInfoNeeded              Status = 0x041c
	StatusClientErrorAccountClosed                  Status = 0x041d
	StatusClientErrorAccountLimitReached            Status = 0x041e
	StatusClientErrorAccountAuthorizationFailed     Status = 0x041f
	StatusClientErrorNotFetchable                   Status = 0x0420

	StatusServerErrorInternal                         Status = 0x0500
	StatusServerErrorOperationNotSupported            Status = 0x0501
	StatusServerErrorServiceUnavailable               Status = 0x0502
	StatusServerErrorVersionNotSupported              Status = 0x0503
	StatusServerErrorDeviceError                      Status = 0x0504
	StatusServerErrorTemporaryError                   Status = 0x0505
	StatusServerErrorNotAcceptingJobs                 Status = 0x0506
	StatusServerErrorBusy                             Status = 0x0507
	StatusServerErrorJobCanceled                      Status = 0x0508
	StatusServerErrorMultipleDocumentJobsNotSupported Status = 0x0509
	StatusServerErrorPrinterIsDeactivated             Status = 0x050a
	StatusServerErrorTooManyJobs                      Status = 0x050b
	StatusServerErrorTooManyDocuments                 Status = 0x050c
)

var statusNames = map[Status]string{
	StatusOK: "successful-ok", StatusOKIgnoredOrSubstituted: "successful-ok-ignored-or-substituted-attributes",
	StatusOKConflicting: "successful-ok-conflicting-attributes", StatusOKIgnoredSubscriptions: "successful-ok-ignored-subscriptions",
	StatusOKTooManyEvents: "successful-ok-too-many-events", StatusOKEventsComplete: "successful-ok-events-complete",
	StatusClientErrorBadRequest: "client-error-bad-request", StatusClientErrorForbidden: "client-error-forbidden",
	StatusClientErrorNotAuthenticated: "client-error-not-authenticated", StatusClientErrorNotAuthorized: "client-error-not-authorized",
	StatusClientErrorNotPossible: "client-error-not-possible", StatusClientErrorTimeout: "client-error-timeout",
	StatusClientErrorNotFound: "client-error-not-found", StatusClientErrorGone: "client-error-gone",
	StatusClientErrorRequestEntityTooLarge:          "client-error-request-entity-too-large",
	StatusClientErrorRequestValueTooLong:            "client-error-request-value-too-long",
	StatusClientErrorDocumentFormatNotSupported:     "client-error-document-format-not-supported",
	StatusClientErrorAttributesOrValuesNotSupported: "client-error-attributes-or-values-not-supported",
	StatusClientErrorURISchemeNotSupported:          "client-error-uri-scheme-not-supported",
	StatusClientErrorCharsetNotSupported:            "client-error-charset-not-supported",
	StatusClientErrorConflictingAttributes:          "client-error-conflicting-attributes",
	StatusClientErrorCompressionNotSupported:        "client-error-compression-not-supported",
	StatusClientErrorCompressionError:               "client-error-compression-error",
	StatusClientErrorDocumentFormatError:            "client-error-document-format-error",
	StatusClientErrorDocumentAccessError:            "client-error-document-access-error",
	StatusClientErrorAttributesNotSettable:          "client-error-attributes-not-settable",
	StatusClientErrorIgnoreAllSubscriptions:         "client-error-ignore-all-subscriptions",
	StatusClientErrorTooManySubscriptions:           "client-error-too-many-subscriptions",
	StatusClientErrorDocumentPasswordError:          "client-error-document-password-error",
	StatusClientErrorDocumentPermissionError:        "client-error-document-permission-error",
	StatusClientErrorDocumentSecurityError:          "client-error-document-security-error",
	StatusClientErrorDocumentUnprintableError:       "client-error-document-unprintable-error",
	StatusClientErrorAccountInfoNeeded:              "client-error-account-info-needed", StatusClientErrorAccountClosed: "client-error-account-closed",
	StatusClientErrorAccountLimitReached:        "client-error-account-limit-reached",
	StatusClientErrorAccountAuthorizationFailed: "client-error-account-authorization-failed",
	StatusClientErrorNotFetchable:               "client-error-not-fetchable",
	StatusServerErrorInternal:                   "server-error-internal-error", StatusServerErrorOperationNotSupported: "server-error-operation-not-supported",
	StatusServerErrorServiceUnavailable: "server-error-service-unavailable", StatusServerErrorVersionNotSupported: "server-error-version-not-supported",
	StatusServerErrorDeviceError: "server-error-device-error", StatusServerErrorTemporaryError: "server-error-temporary-error",
	StatusServerErrorNotAcceptingJobs: "server-error-not-accepting-jobs", StatusServerErrorBusy: "server-error-busy",
	StatusServerErrorJobCanceled:                      "server-error-job-canceled",
	StatusServerErrorMultipleDocumentJobsNotSupported: "server-error-multiple-document-jobs-not-supported",
	StatusServerErrorPrinterIsDeactivated:             "server-error-printer-is-deactivated", StatusServerErrorTooManyJobs: "server-error-too-many-jobs",
	StatusServerErrorTooManyDocuments: "server-error-too-many-documents",
}

// String returns the name of the status code, or its number.
func (s Status) String() string {
	if n, ok := statusNames[s]; ok {
		return n
	}
	return "status-0x" + strconv.FormatUint(uint64(s), 16)
}

// OK reports whether the status is a success (0x0000 to 0x00ff).
func (s Status) OK() bool { return s < 0x0100 }

// ClientError reports whether the status is a client error (0x0400 to
// 0x04ff): the request was wrong.
func (s Status) ClientError() bool { return s >= 0x0400 && s < 0x0500 }

// ServerError reports whether the status is a server error (0x0500 to
// 0x05ff): the printer could not do it.
func (s Status) ServerError() bool { return s >= 0x0500 && s < 0x0600 }

// JobState is the state of a job (RFC 8011 §5.3.7, "job-state").
type JobState int32

// Job states.
const (
	JobPending           JobState = 3
	JobPendingHeld       JobState = 4
	JobProcessing        JobState = 5
	JobProcessingStopped JobState = 6
	JobCanceled          JobState = 7
	JobAborted           JobState = 8
	JobCompleted         JobState = 9
)

var jobStateNames = map[JobState]string{
	JobPending: "pending", JobPendingHeld: "pending-held", JobProcessing: "processing", JobProcessingStopped: "processing-stopped",
	JobCanceled: "canceled", JobAborted: "aborted", JobCompleted: "completed",
}

// String returns the keyword of the state.
func (s JobState) String() string {
	if n, ok := jobStateNames[s]; ok {
		return n
	}
	return "job-state-" + strconv.Itoa(int(s))
}

// Done reports whether the job is over: canceled, aborted or completed.
func (s JobState) Done() bool { return s >= JobCanceled }

// MarshalText writes the keyword of the state ("pending" …).
func (s JobState) MarshalText() ([]byte, error) {
	if _, ok := jobStateNames[s]; !ok {
		return nil, fmt.Errorf("ipp: job-state %d", int32(s))
	}
	return []byte(s.String()), nil
}

// UnmarshalText reads the keyword of a state.
func (s *JobState) UnmarshalText(b []byte) error {
	for k, n := range jobStateNames {
		if n == string(b) {
			*s = k
			return nil
		}
	}
	return fmt.Errorf("ipp: job-state %q", b)
}

// PrinterState is the state of a printer (RFC 8011 §5.4.11,
// "printer-state").
type PrinterState int32

// Printer states.
const (
	PrinterIdle       PrinterState = 3
	PrinterProcessing PrinterState = 4
	PrinterStopped    PrinterState = 5
)

// String returns the keyword of the state.
func (s PrinterState) String() string {
	switch s {
	case PrinterIdle:
		return "idle"
	case PrinterProcessing:
		return "processing"
	case PrinterStopped:
		return "stopped"
	}
	return "printer-state-" + strconv.Itoa(int(s))
}
