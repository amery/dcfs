package dcfs

import (
	"syscall"

	"go.sancus.dev/fs"
)

// Tree
func (fsys *Filesystem) locateBest(name string) (Node, string, string, error) {

	node := fsys.root
	p0, p1 := "", name

	for {
		next, _, s1, err := node.locate(fsys, p1)

		if err == syscall.EAGAIN {
			// not populated yet. retry
			node.populate(fsys, false)
		} else if err != nil {
			// error. stop.
			return node, p0, p1, err
		} else if l := len(s1); l == 0 {
			// match
			return next, name, "", nil
		} else {
			p0 = name[:len(name)-l-1]
			p1 = s1

			if dir, ok := next.(*DirectoryNode); !ok {
				// not a directory. done.
				return next, p0, p1, nil
			} else {
				// next directory
				node = dir
			}
		}
	}
}

func (fsys *Filesystem) locate(name string) (Node, error) {

	node, _, extra, err := fsys.locateBest(name)

	if extra == "" && err == nil {
		return node, nil
	}

	err = fs.AsPathError("locate", name, syscall.ENOENT)
	return nil, err
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
