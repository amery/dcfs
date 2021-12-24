package scan

import (
	"path/filepath"
	"strings"

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
		m.vol[vol] = fsys
	}

	// and trim path accordingly
	if s := strings.TrimPrefix(path, root); s == "" {
		return fsys, ".", nil
	} else {
		return fsys, s, nil
	}
}

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

func (m *Scanner) Split(path string) (fs.FS, string, string, error) {
	// volume and absolute path
	fsys, path, err := m.SplitVolume(path)
	if err != nil {
		return nil, path, "", err
	}

	// find existing bucket
	dir, base, name := path, "", ""
	for {
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
