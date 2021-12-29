package dcfs

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"syscall"
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
	mu     sync.RWMutex
	record *NodeRecord
	tree   *radix.Tree
}

func (node *DirectoryNode) String() string {
	return node.record.String()
}

func (node *DirectoryNode) locate(fsys *Filesystem, name string) (Node, string, string, error) {
	node.mu.RLock()
	defer node.mu.RUnlock()

	if node.tree == nil {
		// directory not yet populated, try again later
		return nil, "", name, syscall.EAGAIN
	} else if name == "" || name == "." {
		// exact dir match
		return node, name, "", nil
	} else if p0, v, ok := node.tree.LongestPrefix(name); !ok {
		// no match
		return node, "", name, syscall.ENOENT
	} else if next, ok := v.(Node); !ok {
		// can't happen
		return node, "", name, syscall.EINVAL
	} else {
		var p1 string

		if l := len(p0); len(name) > l {
			// partial match
			p1 = name[l+1:]
		}

		return next, p0, p1, nil
	}
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
