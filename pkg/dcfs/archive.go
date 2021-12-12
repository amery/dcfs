package dcfs

import (
	"io/fs"
	"log"
	"syscall"

	"go.sancus.dev/core/errors"
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
	log.Printf("%#v: %v", errors.Here(), node.Inode())
	return nil, syscall.ENOSYS
}
