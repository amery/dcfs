package dcfs

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/armon/go-radix"

	"go.sancus.dev/fs"
)

var (
	_ fs.FileInfo = dirinfo{}
	_ fs.File     = (*Directory)(nil)
	_ Node        = (*DirectoryNode)(nil)
)

// fs.FileInfo
type dirinfo struct {
	basename string
	node     *DirectoryNode
}

func (fi dirinfo) Name() string {
	return fi.basename
}

func (fi dirinfo) Size() int64 {
	return 0
}

func (fi dirinfo) Mode() fs.FileMode {
	return fs.ModeDir | 0755
}

func (fi dirinfo) ModTime() time.Time {
	return time.Time{}
}

func (fi dirinfo) IsDir() bool {
	return true
}

func (fi dirinfo) Sys() interface{} {
	return fi.node
}

// Directory
type Directory struct {
	basename string
	node     *DirectoryNode
}

func (dir *Directory) String() string {
	return fmt.Sprintf("%s name=%q ptr=%p", dir.node, dir.basename, dir)
}

func (dir *Directory) Close() error {
	return nil
}

func (dir *Directory) Read(b []byte) (int, error) {
	return 0, io.EOF
}

func (dir *Directory) Stat() (fs.FileInfo, error) {
	return dirinfo{dir.basename, dir.node}, nil
}

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

func (node *DirectoryNode) String() string {
	return node.record.String()
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
	dir := &Directory{"", node}
	return dir, nil
}
