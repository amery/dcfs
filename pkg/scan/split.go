package scan

import (
	"log"
	"path/filepath"
	"strings"

	"go.sancus.dev/core/errors"
	"go.sancus.dev/fs"
	"go.sancus.dev/fs/dirfs"

	"github.com/amery/dcfs/pkg/bucket"
)

func (m *Scanner) splitVolume(vol, path string) (fs.FS, string, error) {
	root := vol + string(filepath.Separator)

	m.mu.Lock()
	defer m.mu.Unlock()

	// find or create fs.FS
	fsys, ok := m.vol[vol]
	if !ok {
		var err error

		fsys, err = dirfs.New(root)
		if err != nil {
			return nil, "", err
		}
		m.addVolume(fsys, vol)
	}

	// and trim path accordingly
	if s := strings.TrimPrefix(path, root); s == "" {
		return fsys, ".", nil
	} else {
		s = filepath.ToSlash(s)
		return fsys, s, nil
	}
}

// Split turns a OS path into volume's root fs and fs.FS friendly path.
func (m *Scanner) SplitVolume(path string) (fs.FS, string, error) {

	// absolute path
	if s, err := filepath.Abs(path); err == nil {
		path = s
	} else {
		return nil, "", err
	}

	// volume name
	vol := filepath.VolumeName(path)

	// split
	return m.splitVolume(vol, path)
}

// Split turns a path into volume's root fs, path to the bucket, and bucket relative path.
func (m *Scanner) Split(path string) (fs.FS, string, string, error) {
	log.Printf("%+n: %s:%q", errors.Here(), "path", path)

	// volume and absolute path
	fsys, path, err := m.SplitVolume(path)
	if err != nil {
		return nil, path, "", err
	}

	// find existing bucket
	dir, base, name := path, "", ""
	for {
		log.Printf("%+n: %s %s:%q %s:%q", errors.Here(),
			fsys, "dir", dir, "name", name)

		if bucket.IsRoot(fsys, dir) {
			// hit
			return fsys, dir, name, nil
		} else if dir == "." {
			// give up
			break
		}

		// up one level
		dir, base = fs.Split(dir)

		if name == "" {
			name = base
		} else {
			name = fs.Join(base, name)
		}
	}

	// find dir to create bucket
	fi, err := fs.Stat(fsys, path)
	if err != nil {
		dir, name = ".", path
	} else if !fi.IsDir() {
		dir, name = fs.Split(path)
	} else {
		dir, name = path, ""
	}

	return fsys, dir, name, err
}
