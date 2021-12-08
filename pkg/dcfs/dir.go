package dcfs

import (
	"syscall"

	"github.com/armon/go-radix"

	"go.sancus.dev/fs"
)

var (
	_ Node = (*DirectoryNode)(nil)
)

type DirectoryEntry struct {
	Inode uint64 `boltholdIndex:"Inode"`
	Name  string
}

// DirectoryNode
type DirectoryNode struct {
	record *NodeRecord
	tree   *radix.Tree
}

// Node
func (node *DirectoryNode) Inode() uint64  { return node.record.Inode }
func (node *DirectoryNode) Type() NodeType { return node.record.Type }

func (node *DirectoryNode) Open() (fs.File, error) {
	return nil, syscall.ENOSYS
}
