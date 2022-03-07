package repo

import "bytes"

type ZipPayload struct {
	ID  string
	Zip bytes.Buffer
}
