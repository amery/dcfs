package dcfs

import (
	"syscall"

	"go.sancus.dev/fs"
)

// Tree
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
