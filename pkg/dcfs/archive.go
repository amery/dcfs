package dcfs

import (
	"io/fs"
	"syscall"
)

var (
	_ Node = (*ArchiveNode)(nil)
)

// Archive
type ArchiveNode struct {
	record *NodeRecord
}

// Node
func (node *ArchiveNode) Inode() uint64  { return node.record.Inode }
func (node *ArchiveNode) Type() NodeType { return node.record.Type }

func (node *ArchiveNode) Open() (fs.File, error) {
	return nil, syscall.ENOSYS
}
