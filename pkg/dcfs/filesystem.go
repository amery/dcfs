package dcfs

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"go.sancus.dev/fs"

	"github.com/ancientlore/go-avltree"
	"github.com/timshannon/bolthold"
)

const (
	DbFilename         = "dcfs.db"
	DbFilePermissions  = 0644
	DataDirPermissions = 0755
)

// Interfaces
var (
	_ fs.FS = (*Filesystem)(nil)
)

type Filesystem struct {
	datadir string
	db      *bolthold.Store
	mu      sync.RWMutex
	root    *DirectoryNode
	nodes   *avltree.Tree

	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	cancelled int32
}

func (m *Filesystem) spawn(fn func(ctx context.Context, fs *Filesystem)) {
	if fn != nil && m.cancelled == 0 {
		m.wg.Add(1)
		go func() {
			defer m.wg.Done()
			fn(m.ctx, m)
		}()
	}
}

func (m *Filesystem) Cancel() {
	if atomic.CompareAndSwapInt32(&m.cancelled, 0, 1) {
		// just once
		m.cancel()
	}
}

func (m *Filesystem) Close() error {
	m.Cancel()
	m.wg.Wait()
	return m.db.Close()
}

func New(ctx context.Context, datadir string) (*Filesystem, error) {

	if datadir == "" {
		datadir = "."
	}

	// make sure datadir exists
	if err := os.MkdirAll(datadir, DataDirPermissions); err != nil {
		return nil, err
	}

	// database
	filename := filepath.Join(datadir, DbFilename)
	db, err := bolthold.Open(filename, DbFilePermissions, nil)
	if err != nil {
		return nil, err
	}

	// cancelation
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithCancel(ctx)

	m := &Filesystem{
		datadir: datadir,
		db:      db,

		ctx:    ctx,
		cancel: cancel,
	}

	if err := m.init(); err != nil {
		defer m.Close()
		return nil, err
	}

	return m, nil
}
