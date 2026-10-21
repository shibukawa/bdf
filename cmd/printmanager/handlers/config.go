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
	AdminToken    string        `default:"" help:"bearer token of the API (/api)"`
	TickToken     string        `default:"" help:"bearer token a scheduler calls POST /api/tick with"`
	SubmitTokens  []string      `help:"name:secret pairs that may send jobs; none lets anyone send"`
	TickInterval  time.Duration `default:"30s" help:"how often the process ticks; 0 leaves it to a scheduler (serverless)"`
	JobKeep       time.Duration `default:"336h" help:"how long a job that ended is kept (its record, and its document once nothing shows it)"`
	Serverless    bool          `default:"false" help:"refuse a store only one process can use, and tick only when called"`
	FontDirs      []string      `help:"directories of the fonts of documents"`
	NoSystemFonts bool          `default:"false" help:"use only font_dirs"`
	MaxDocumentMB int           `default:"64" help:"largest document taken, in MiB"`
	SecureCookies bool          `default:"true" help:"Secure on the devices' cookie (false only for plain http in development)"`
}

// RegisterConfig binds ManagerConfig to the "manager" prefix; main calls
// it before the configuration is parsed.
func RegisterConfig() { pwconfig.Register[ManagerConfig]("manager") }
