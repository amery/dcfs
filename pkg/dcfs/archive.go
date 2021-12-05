package dcfs

import (
	"io/fs"
	"syscall"
)

var (
	_ Node = (*Archive)(nil)
)

type Archive struct {
	entry *NodeEntry
}

// Node
func (file *Archive) Inode() uint64  { return file.entry.Inode }
func (file *Archive) Type() NodeType { return file.entry.Content.Type() }

func (file *Archive) Open() (fs.File, error) {
	return nil, syscall.ENOSYS
}
