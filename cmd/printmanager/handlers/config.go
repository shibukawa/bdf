package handlers

import (
	"time"

	"github.com/shibukawa/popcornweb/pwconfig"
)

// ManagerConfig is the [manager] table of config.{APP_ENV}.toml, or the
// MANAGER_* variables: how the manager runs. What it shows where
// (locations, printers, devices) is data, changed through the API.
type ManagerConfig struct {
	Records       string        `default:"bbolt:printmanager.db" help:"store of the records: bbolt:path, dynamodb://table?region=…, firestore://project/database, s3://bucket?region=…, r2://bucket?endpoint=…, gcs://bucket"`
	Documents     string        `default:"" help:"store of the slips' documents (bbolt, s3, r2, gcs); empty for the records' store"`
	Secret        string        `default:"" help:"signs the devices' sessions: 32 bytes or more, the same for every instance (required outside dev)"`
	AdminToken    string        `default:"" help:"bearer token of the API (/api), for scripts; people sign in with OIDC"`
	TickToken     string        `default:"" help:"bearer token a scheduler calls POST /api/tick with"`
	SubmitTokens  []string      `help:"name:secret pairs that may send jobs; none lets anyone send"`
	TickInterval  time.Duration `default:"30s" help:"how often the process ticks; 0 leaves it to a scheduler (serverless)"`
	Serverless    bool          `default:"false" help:"refuse a store only one process can use, and tick only when called"`
	FontDirs      []string      `help:"directories of the fonts of documents"`
	NoSystemFonts bool          `default:"false" help:"use only font_dirs"`
	MaxDocumentMB int           `default:"64" help:"largest document taken, in MiB"`
	SecureCookies bool          `default:"true" help:"Secure on the devices' cookie (false only for plain http in development)"`
}

// OIDCConfig is the [oidc] table, or the OIDC_* variables: who of the
// people the IdP signs in (popcornweb's [auth.oidc]: issuer, client_id,
// client_secret, redirect_url) gets in, their groups, and the admins.
type OIDCConfig struct {
	GroupsClaim  string   `default:"groups" help:"claim of the groups: a name, or a dotted path (realm_access.roles); a string or an array of them"`
	GroupsPrefix string   `default:"" help:"take only the groups of this prefix"`
	Admins       []string `help:"admins: issuer#subject, or a verified email; at least one admin rule outside dev"`
	AdminGroups  []string `help:"groups whose members are admins"`
	Allowed      AllowedConfig
}

// AllowedConfig narrows who gets in; empty lets in anyone the IdP signs
// in. A person passes each list that is set (one of its values); admins
// always get in.
type AllowedConfig struct {
	Issuers []string `help:"issuers whose people get in"`
	Domains []string `help:"email domains (or Google's hd) whose people get in"`
	Groups  []string `help:"groups whose members get in"`
}

// RegisterConfig binds ManagerConfig to the "manager" prefix and
// OIDCConfig to "oidc"; main calls it before the configuration is parsed.
func RegisterConfig() {
	pwconfig.Register[ManagerConfig]("manager")
	pwconfig.Register[OIDCConfig]("oidc")
}
