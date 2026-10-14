package tag

import (
	"encoding/base64"
	"fmt"
	"io/fs"
	"strings"
)

// MaxPicture bounds the bytes of a picture of the data.
const MaxPicture = 32 << 20

// PictureData reads a picture that the data names for a picture tag: the
// path of a file among files, or a data: URL in base64.
func PictureData(files fs.FS, name string) ([]byte, error) {
	if strings.HasPrefix(name, "data:") {
		meta, body, ok := strings.Cut(name[len("data:"):], ",")
		if !ok || !strings.HasSuffix(meta, ";base64") {
			return nil, fmt.Errorf("a data: URL that is not base64")
		}
		if base64.StdEncoding.DecodedLen(len(body)) > MaxPicture {
			return nil, fmt.Errorf("it is larger than %d MiB", MaxPicture>>20)
		}
		data, err := base64.StdEncoding.DecodeString(strings.TrimSpace(body))
		if err != nil {
			return nil, fmt.Errorf("a data: URL that is not base64")
		}
		return data, nil
	}
	if files == nil {
		return nil, fmt.Errorf("no files were given to read %q from", name)
	}
	name = strings.TrimPrefix(strings.ReplaceAll(name, "\\", "/"), "./")
	if !fs.ValidPath(name) {
		return nil, fmt.Errorf("%q is not the path of a file beside the template", name)
	}
	st, err := fs.Stat(files, name)
	if err != nil {
		return nil, fmt.Errorf("%q is not among the files", name)
	}
	if st.Size() > MaxPicture {
		return nil, fmt.Errorf("%q is larger than %d MiB", name, MaxPicture>>20)
	}
	return fs.ReadFile(files, name)
}
