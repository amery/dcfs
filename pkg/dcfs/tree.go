package dcfs

import (
	"syscall"

	"github.com/armon/go-radix"

	"go.sancus.dev/fs"
)

type Directory struct {
	entry   *NodeEntry
	content *NodeDirectoryContent
	tree    *radix.Tree
}

func (fsys *Filesystem) Open(name string) (fs.File, error) {
	return nil, syscall.ENOENT
}
