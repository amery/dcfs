package scan

import (
	"context"
	"log"
	"syscall"

	"go.sancus.dev/core/errors"
	"go.sancus.dev/fs"

	"github.com/amery/dcfs/pkg/bucket"
)

func (m *Node) getBucketNode(dir string) (*Node, error) {
	var bkt *bucket.Bucket

	log.Printf("%+n: %s %s:%q", errors.Here(), m.fsys, "dir", dir)

	// find or create subnode
	child, err := m.getNode(dir)
	if err != nil {
		// failed to create subnode
		child = nil
	} else if child.bucket != nil {
		// already contains a bucket
	} else if bkt, err = bucket.New(child.fsys); err != nil {
		// failed to create bucket
		child = nil
	} else {
		// associate new bucket
		child.bucket = bkt
	}

	return child, err
}

func (m *Node) Commit() error {
	log.Printf("%+n", errors.Here())

	m.mu.Lock()
	defer m.mu.Unlock()

	return m.commit()
}

func (m *Node) commit() error {
	log.Printf("%+n", errors.Here())

	return m.bucket.Commit(m.fsys)
}

func (m *Node) put() *Node {
	m.count++
	log.Printf("%+n: %s:%v", errors.Here(), "count", m.count)
	return m
}

func (m *Node) pop() bool {
	log.Printf("%+n: %s:%v", errors.Here(), "count", m.count)

	m.count--
	if m.count <= 0 {
		m.count = 0
		return true
	}

	return false
}

func (m *Node) Pop() bool {
	log.Printf("%+n", errors.Here())

	m.mu.Lock()
	defer m.mu.Unlock()

	return m.pop()
}

func (m *Node) close() error {
	log.Printf("%+n", errors.Here())

	if m.pop() {
		return m.commit()
	}

	return nil
}

func (m *Node) Close() error {
	log.Printf("%+n", errors.Here())

	if m.Pop() {
		return m.Commit()
	}

	return nil
}

func (m *Node) Add(ctx context.Context, name string) error {
	log.Printf("%+n: %s:%q", errors.Here(), "name", name)

	return syscall.ENOSYS
}

func (m *Scanner) Bucket(fsys fs.FS, dir string) (*Node, error) {
	log.Printf("%+n: %s %s:%q", errors.Here(), fsys, "dir", dir)

	m.mu.Lock()
	defer m.mu.Unlock()

	// get root node
	vol, ok := m.data[fsys]
	if !ok {
		return nil, syscall.EINVAL
	}

	// and the bucket
	return vol.getBucketNode(dir)
}
