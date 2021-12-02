package dcfs

import (
	"context"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/armon/go-radix"
	"github.com/timshannon/bolthold"
)

const (
	DbFilename        = "dcfs.db"
	DbFilePermissions = 0644
)

type Filesystem struct {
	datadir string
	tree    *radix.Tree
	db      *bolthold.Store
	mu      sync.RWMutex
	wg      sync.WaitGroup

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

	// database
	if datadir == "" {
		datadir = "."
	}
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
		tree:    radix.New(),
		db:      db,

		ctx:    ctx,
		cancel: cancel,
	}

	return m, nil
}
