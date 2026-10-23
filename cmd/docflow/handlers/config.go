package handlers

import (
	"time"

	"github.com/shibukawa/popcornweb/pwconfig"
)

// DocflowConfig is the [docflow] table of config.{APP_ENV}.toml, or the
// DOCFLOW_* variables: how the engine runs. What it shows where
// (locations, printers, devices) is data, changed through the API.
type DocflowConfig struct {
	Records         string        `default:"bbolt:docflow.db" help:"store of the records: bbolt:path, dynamodb://table?region=…, firestore://project/database, s3://bucket?region=…, r2://bucket?endpoint=…, gcs://bucket"`
	Documents       string        `default:"" help:"store of the slips' documents (bbolt, s3, r2, gcs); empty for the records' store"`
	Secret          string        `default:"" help:"signs the devices' sessions, the apps' tokens and the pages' sessions: 32 bytes or more, the same for every instance (required outside dev)"`
	PreviousSecrets []string      `help:"retired secrets that still open what was signed with them, while a rotation goes on"`
	AdminToken      string        `default:"" help:"bearer token of the API (/api), for scripts; people sign in with OIDC"`
	TickToken       string        `default:"" help:"bearer token a scheduler calls POST /api/tick with"`
	SubmitTokens    []string      `help:"name:secret pairs that may send jobs; none lets anyone send"`
	TickInterval    time.Duration `default:"30s" help:"how often the process ticks; 0 leaves it to a scheduler (serverless)"`
	JobKeep         time.Duration `default:"336h" help:"how long a job that ended is kept (its record, and its document once nothing shows it)"`
	Serverless      bool          `default:"false" help:"refuse a store only one process can use, and tick only when called"`
	FontDirs        []string      `help:"directories of the fonts of documents"`
	NoSystemFonts   bool          `default:"false" help:"use only font_dirs"`
	MaxDocumentMB   int           `default:"64" help:"largest document taken, in MiB"`
	OnError         []string      `help:"nodes a job goes to when a send or a step fails for good at a node with no on_error edges (a notify node)"`
	CredentialDir   string        `default:"" help:"directory of the credentials nodes name, a file each (DOCFLOW_CRED_<name> comes first)"`
	SecureCookies   bool          `default:"true" help:"Secure on the devices' cookie (false only for plain http in development)"`
	DeviceURL       string        `default:"" help:"origin the devices of displays open (an enrollment's URL and QR code), e.g. https://docflow.example.com; empty: this server as the admin reached it, at this computer's address on the network for localhost, on https_port when there is one"`
	HTTPSPort       int           `default:"0" help:"port of an HTTPS listener serving what the HTTP one does, for devices on the network (a device's key needs HTTPS); 0 for none"`
	TLSCertFile     string        `default:"" help:"file of the certificate (PEM) of https_port, followed by its chain; empty in development: one signed by itself, kept in the user's cache"`
	TLSKeyFile      string        `default:"" help:"file of the private key (PEM) of tls_cert_file"`
	SignKeyFile     string        `default:"" help:"file of the private key (PEM) that sign nodes sign documents with, as a witness"`
	SignCertFile    string        `default:"" help:"file of the certificate of the signing key (PEM), followed by those of its chain"`
	SignTSA         string        `default:"" help:"URL of a time-stamp authority (RFC 3161) that stamps the signatures of sign nodes"`
}

// OIDCConfig is the [oidc] table, or the OIDC_* variables: who of the
// people the IdP signs in (popcornweb's [auth.oidc]: issuer, client_id,
// client_secret, redirect_url) gets in, their groups, and the admins.
type OIDCConfig struct {
	GroupsClaim  string   `default:"groups" help:"claim of the groups: a name, or a dotted path (realm_access.roles); a string or an array of them"`
	GroupsPrefix string   `default:"" help:"take only the groups of this prefix"`
	Admins       []string `help:"admins: issuer#subject, or a verified email; at least one admin rule outside dev"`
	AdminGroups  []string `help:"groups whose members are admins"`
	// The roles short of an admin: topology changes the topology and what
	// is at it (the admin pages); audit reads the audit log and every job.
	TopologyAdmins      []string `help:"people with the role topology: issuer#subject, or a verified email"`
	TopologyAdminGroups []string `help:"groups whose members have the role topology"`
	Auditors            []string `help:"people with the role audit: issuer#subject, or a verified email"`
	AuditorGroups       []string `help:"groups whose members have the role audit"`
	Allowed             AllowedConfig
}

// AllowedConfig narrows who gets in; empty lets in anyone the IdP signs
// in. A person passes each list that is set (one of its values); admins
// always get in.
type AllowedConfig struct {
	Issuers []string `help:"issuers whose people get in"`
	Domains []string `help:"email domains (or Google's hd) whose people get in"`
	Groups  []string `help:"groups whose members get in"`
}

// RegisterConfig binds DocflowConfig to the "docflow" prefix and
// OIDCConfig to "oidc"; main calls it before the configuration is parsed.
func RegisterConfig() {
	pwconfig.Register[DocflowConfig]("docflow")
	pwconfig.Register[OIDCConfig]("oidc")
}
