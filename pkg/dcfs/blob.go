package dcfs

import (
	"time"
)

type Blake3Hash []byte

type Blob struct {
	Blake3   Blake3Hash `boltholdUnique:"Blake3"`
	Size     uint64
	Created  time.Time
	Mimetype string
}
