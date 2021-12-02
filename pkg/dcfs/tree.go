package dcfs

import (
	"syscall"

	"go.sancus.dev/fs"
)

func (fsys *Filesystem) Open(name string) (fs.File, error) {
	return nil, syscall.ENOENT
}
