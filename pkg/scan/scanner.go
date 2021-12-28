package scan

import (
	"context"
	"sync"
	"sync/atomic"

	"go.sancus.dev/fs"
)

type Scanner struct {
	mu        sync.Mutex
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	cancelled uint32
	err       error

	vol  map[string]fs.FS // Volumes
	data map[fs.FS]*Node  // Nodes
}

func NewScanner(ctx context.Context) (*Scanner, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithCancel(ctx)

	m := &Scanner{
		ctx:    ctx,
		cancel: cancel,
		vol:    make(map[string]fs.FS),
		data:   make(map[fs.FS]*Node),
	}

	return m, nil
}

func (m *Scanner) Cancel(err error) {
	if atomic.CompareAndSwapUint32(&m.cancelled, 0, 1) {
		// once
		m.cancel()
		m.err = err
	}
}

func (m *Scanner) Spawn(f func(context.Context)) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		f(m.ctx)
	}()
}

func (m *Scanner) Wait() error {
	m.wg.Wait()
	return m.err
}

func (m *Scanner) Scan(ctx context.Context, fsys fs.FS, dir, name string) error {
	bkt, err := m.Bucket(fsys, dir)
	if err != nil {
		return err
	}
	defer bkt.Close()

	if name == "" || name == "." {
		// recursively add everything
		name = ""
	} else if fi, err := bkt.Stat(name); err != nil {
		// invalid name
		return err
	} else if !fi.IsDir() {
		// file to add
	} else if bkt, err = bkt.Split(name); err != nil {
		// subdir split failed
		return err
	} else {
		// subdir split
		name = ""
	}

	return bkt.Add(m.ctx, name)
}

func (m *Scanner) Close() error {
	m.Cancel(nil)

	return m.Wait()
}
