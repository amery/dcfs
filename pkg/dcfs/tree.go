package dcfs

import (
	"log"
	"strings"
	"syscall"

	"go.sancus.dev/core/errors"
	"go.sancus.dev/fs"
)

// Tree
func (fsys *Filesystem) locateBest(name string) (Node, string, string, error) {
	log.Printf("%+n: %s:%q", errors.Here(), "name", name)

	node := fsys.root
	p0, p1 := "", name

	for {
		log.Printf("%+n: <%s> %s:%q %s:%q", errors.Here(), node,
			"p0", p0, "p1", p1)

		next, s0, s1, err := node.locate(fsys, p1)

		log.Printf("%+n: -> <%s> %s:%q %s:%q: %s", errors.Here(), next,
			"s0", s0, "s1", s1, err)

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
	log.Printf("%+n: %s:%q", errors.Here(), "name", name)

	node, _, extra, err := fsys.locateBest(name)

	if extra == "" && err == nil {
		return node, nil
	}

	err = fs.AsPathError("locate", name, syscall.ENOENT)
	return nil, err
}

func (fsys *Filesystem) Open(name string) (fs.File, error) {
	log.Printf("%+n: %s:%q", errors.Here(), "name", name)

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

func (fsys *Filesystem) Mkdir(name string, perm fs.FileMode) error {
	log.Printf("%+n: %s:%q", errors.Here(), "name", name)

	if !fs.ValidPath(name) {
		return fs.AsPathError("mkdir", name, syscall.EINVAL)
	} else if n, _, p1, err := fsys.locateBest(name); p1 == "" {
		return fs.AsPathError("mkdir", name, syscall.EEXIST)
	} else if err != syscall.ENOENT {
		return fs.AsPathError("mkdir", name, err)
	} else if i := strings.IndexRune(p1, '/'); i >= 0 {
		return fs.AsPathError("mkdir", name, syscall.ENOENT)
	} else if dir, ok := n.(*DirectoryNode); !ok {
		return fs.AsPathError("mkdir", name, syscall.ENOTDIR)
	} else if _, err := dir.mkdir(fsys, p1); err != nil {
		return fs.AsPathError("mkdir", name, err)
	} else {
		return nil
	}
}
