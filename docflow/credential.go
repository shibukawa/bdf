package docflow

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// The secrets the exits and the notify nodes use (an SMTP server's
// password, an S3 access key, a webhook's URL or a bearer token) are
// named by a node (Node.Credential), never kept with it: the value is
// the environment variable DOCFLOW_CRED_<name>, the name's - as _, or
// the file <name> of Options.CredentialDir, so that no secret one can
// take out and use is in the store (requirement:print-manager-auth).
// What a value holds is the kind's: user:password for SMTP,
// access:secret for S3, a token or a URL.

// credentialRE is the shape of a credential's name.
var credentialRE = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{0,63}$`)

// ErrCredential is the error of a credential named that is not there.
var ErrCredential = errors.New("docflow: no such credential")

// CheckCredential checks the shape of a credential's name: "" is none.
func CheckCredential(name string) error {
	if name != "" && !credentialRE.MatchString(name) {
		return fmt.Errorf("%w: credential %q: letters, digits, _ and -", ErrInvalid, name)
	}
	return nil
}

// CredentialEnv is the environment variable of a credential's value.
func CredentialEnv(name string) string {
	return "DOCFLOW_CRED_" + strings.ReplaceAll(name, "-", "_")
}

// Credential reads the value of a credential: its environment variable,
// else its file in Options.CredentialDir (without a trailing line
// break). The error names the credential, never a value.
func (m *Manager) Credential(name string) (string, error) {
	if !credentialRE.MatchString(name) {
		return "", fmt.Errorf("%w: %q", ErrCredential, name)
	}
	if v, ok := os.LookupEnv(CredentialEnv(name)); ok && v != "" {
		return v, nil
	}
	if m.o.CredentialDir != "" {
		b, err := os.ReadFile(filepath.Join(m.o.CredentialDir, name))
		if err == nil {
			return strings.TrimRight(string(b), "\r\n"), nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("%w: %s: %v", ErrCredential, name, err)
		}
	}
	return "", fmt.Errorf("%w: %s (%s, or a file of that name in the credential directory)", ErrCredential, name, CredentialEnv(name))
}
