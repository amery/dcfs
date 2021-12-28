package scan

import (
	"strings"
	"sync"
	"syscall"

	"github.com/armon/go-radix"

	"go.sancus.dev/fs"

	"github.com/amery/dcfs/pkg/bucket"
)

type Node struct {
	mu      sync.Mutex
	scanner *Scanner
	fsys    fs.FS
	tree    *radix.Tree
	count   int
	bucket  *bucket.Bucket
}

func (m *Scanner) addVolume(fsys fs.FS, vol string) {
	n := &Node{
		scanner: m,
		fsys:    fsys,
		tree:    radix.New(),
	}

	m.vol[vol] = fsys
	m.data[fsys] = n
}

func (m *Node) getNode(dir string) (*Node, error) {

	m.mu.Lock()
	for {

		best, v, ok := m.tree.LongestPrefix(dir)
		if !ok {
			break
		}

		n, ok := v.(*Node)
		if !ok {
			// can't happen
			m.mu.Unlock()
			return nil, syscall.EINVAL
		}

		extra := strings.TrimPrefix(best, dir)

		if extra == "" {
			// match
			m.mu.Unlock()
			return n, nil
		}

		// loop
		n.mu.Lock()
		m.mu.Unlock()

		m = n
		dir = extra
	}

	defer m.mu.Unlock()

	return m.split(dir)
}

func (m *Node) split(dir string) (*Node, error) {
	fsys, err := fs.Sub(m.fsys, dir)
	if err != nil {
		return nil, err
	}

	n := &Node{
		scanner: m.scanner,
		fsys:    fsys,
		tree:    radix.New(),
	}

	// when splitting a bucket, get a bucket
	if m.bucket != nil {
		bkt, err := bucket.New(fsys)
		if err != nil {
			return nil, err
		}
		n.bucket = bkt
	}

	m.tree.Insert(dir, n)

	// move child nodes in
	m.tree.WalkPrefix(dir, func(path string, v interface{}) bool {
		extra := strings.TrimPrefix(path, dir)
		if extra != "" && extra[0] == '/' {
			extra = extra[1:] // remove leading '/'

			m.tree.Delete(path)
			n.tree.Insert(extra, v)
		}

		return false // continue
	})

	// and move bucket's content
	if m.bucket != nil {
		defer m.close()
		defer n.close()

		m.put().bucket.Move(dir, n.put().bucket, ".")
	}

	return n, nil
}

func (m *Node) Split(dir string) (*Node, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.split(dir)
}

func (m *Node) Stat(name string) (fs.FileInfo, error) {
	return fs.Stat(m.fsys, name)
}
