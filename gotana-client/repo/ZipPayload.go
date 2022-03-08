package repo

import "bytes"

// ZipPayload is produced after successfully zipping json. Zip is the attribute
// that stores the compressed json, and ID is the unique identifier for the
// zipped contents.
type ZipPayload struct {
	ID  string
	Zip bytes.Buffer
}
