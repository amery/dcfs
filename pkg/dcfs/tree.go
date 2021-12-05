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

func (dir *Directory) Open() (fs.File, error) {
	return nil, syscall.ENOSYS
}

func (fsys *Filesystem) locate(name string) (Node, error) {
	return nil, syscall.ENOSYS
}

func (fsys *Filesystem) Open(name string) (fs.File, error) {
	if name == "." {
		// special case
		if f, err := fsys.root.Open(); err != nil {
			return nil, &fs.PathError{"open", name, err}
		} else {
			return f, nil
		}

	} else if !fs.ValidPath(name) {
		return nil, &fs.PathError{"open", name, syscall.EINVAL}
	} else if node, err := fsys.locate(name); err != nil {
		return nil, &fs.PathError{"open", name, err}
	} else if f, err := node.Open(); err != nil {
		return nil, &fs.PathError{"open", name, err}
	} else {
		return f, nil
	}
}
