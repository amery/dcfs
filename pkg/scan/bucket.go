package scan

import (
	"context"
	"sync"
	"syscall"

	"go.sancus.dev/fs"

	"github.com/amery/dcfs/pkg/bucket"
)

type BucketData struct {
	mu      sync.Mutex
	scanner *Scanner
	fsys    fs.FS
	count   int
	bucket  *bucket.Bucket
}

func (m *BucketData) Stat(name string) (fs.FileInfo, error) {
	return fs.Stat(m.fsys, name)
}

func (m *BucketData) Commit() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.bucket.Commit()
}

func (m *BucketData) put() *BucketData {
	m.count++
	return m
}

func (m *BucketData) pop() bool {
	m.count--
	if m.count <= 0 {
		m.count = 0
		return true
	}

	return false
}

func (m *BucketData) Pop() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.pop()
}

func (m *BucketData) Close() error {
	if m.Pop() {
		return m.Commit()
	}

	return nil
}

func (m *BucketData) Add(ctx context.Context, name string) error {
	return syscall.ENOSYS
}

func (m *BucketData) Split(dir string) (*BucketData, error) {
	return nil, syscall.ENOSYS
}

func (m *Scanner) Bucket(fsys fs.FS, dir string) (*BucketData, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// for each volume fs we keep a list of path/bucket pairs
	vol, ok := m.data[fsys]
	if !ok {
		vol = make(map[string]*BucketData, 1)
		m.data[fsys] = vol
	}

	if data, ok := vol[dir]; ok {
		// hit
		return data.put(), nil
	}

	// new bucket
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		return nil, err
	}

	bkt, err := bucket.New(sub)
	if err != nil {
		return nil, err
	}

	data := &BucketData{
		scanner: m,
		fsys:    sub,
		bucket:  bkt,
	}

	vol[dir] = data
	return data.put(), nil
}
