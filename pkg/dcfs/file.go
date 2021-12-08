package dcfs

import (
	"io/fs"
	"syscall"
)

var (
	_ Node = (*FileNode)(nil)
)

// FileNode
type FileNode struct {
	record *NodeRecord
}

// Node
func (node *FileNode) Inode() uint64  { return node.record.Inode }
func (node *FileNode) Type() NodeType { return node.record.Type }

func (node *FileNode) Open() (fs.File, error) {
	return nil, syscall.ENOSYS
}
