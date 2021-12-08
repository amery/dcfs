package fuse

import (
	"log"

	"go.sancus.dev/core/errors"
	"go.sancus.dev/fs/fuse"

	"github.com/amery/dcfs/pkg/dcfs"
)

const (
	FSName = "dcfs"
)

type Filesystem struct {
	*fuse.Filesystem

	store *dcfs.Filesystem
}

func (fsys *Filesystem) Abort() error {
	log.Printf("%+n", errors.Here())
	return fsys.Close()
}

func NewWithBackend(store *dcfs.Filesystem, mountpoint string, options ...fuse.MountOption) (*Filesystem, error) {
	// prepend ours
	opts := make([]fuse.MountOption, 2, len(options)+2)
	opts[0] = fuse.FSName(FSName)
	opts[1] = fuse.Subtype(FSName)
	opts = append(opts, options...)

	// fuse frontend
	ffs, err := fuse.New(store, mountpoint, opts...)
	if err == nil {
		// wrap them up
		fsys := &Filesystem{
			Filesystem: ffs,
			store:      store,
		}
		return fsys, nil
	}

	return nil, err
}

func New(datadir, mountpoint string, options ...fuse.MountOption) (*Filesystem, error) {

	// backend
	store, err := dcfs.New(nil, datadir)
	if err != nil {
		return nil, err
	}

	// fuse frontend
	fsys, err := NewWithBackend(store, mountpoint, options...)
	if err != nil {
		store.Close()
		return nil, err
	}

	return fsys, nil
}
