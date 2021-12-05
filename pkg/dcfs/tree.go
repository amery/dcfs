package dcfs

import (
	"syscall"

	"github.com/armon/go-radix"

	"go.sancus.dev/fs"
)

var (
	_ Node = (*Directory)(nil)
)

type Directory struct {
	entry   *NodeEntry
	content *NodeDirectoryContent
	tree    *radix.Tree
}

// Node
func (dir *Directory) Inode() uint64  { return dir.entry.Inode }
func (dir *Directory) Type() NodeType { return dir.entry.Content.Type() }

func (fsys *Filesystem) Open(name string) (fs.File, error) {
	return nil, syscall.ENOENT
}
