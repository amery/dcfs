package dcfs

import (
	"context"
	"log"
	"sync"
	"syscall"

	"github.com/armon/go-radix"

	"go.sancus.dev/fs"
)

var (
	_ Node = (*DirectoryNode)(nil)
)

// DirectoryNode
type DirectoryEntry struct {
	Inode uint64 `boltholdIndex:"Inode"`
	Name  string
}

type DirectoryNode struct {
	mu     sync.Mutex
	record *NodeRecord
	tree   *radix.Tree
}

func (node *DirectoryNode) populate(fsys *Filesystem, recursive bool) {
	node.mu.Lock()
	defer node.mu.Unlock()

	if node.tree == nil {
		tree := radix.New()

		for _, v := range node.record.Content {
			child := v.(*DirectoryEntry)

			if n, err := fsys.getNode(child.Inode); err != nil {
				log.Printf("%s: child %v not found", node, child.Inode)
			} else {
				tree.Insert(child.Name, n)
			}
		}

		node.tree = tree
	}

	if recursive {
		node.tree.Walk(func(name string, v interface{}) bool {
			if dir, ok := v.(*DirectoryNode); ok {
				fsys.spawn(func(ctx context.Context, _ *Filesystem) {
					select {
					case <-ctx.Done():
						return
					default:
						dir.populate(fsys, recursive)
					}
				})
			}
			return false // continue
		})
	}
}

// Node
func (node *DirectoryNode) Inode() uint64  { return node.record.Inode }
func (node *DirectoryNode) Type() NodeType { return node.record.Type }

func (node *DirectoryNode) Open() (fs.File, error) {
	return nil, syscall.ENOSYS
}
