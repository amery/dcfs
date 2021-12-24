package bucket

import (
	"go.sancus.dev/fs"
)

func IsRoot(fsys fs.FS, dir string) bool {
	s := fs.Join(dir, BucketFileName)
	fi, err := fs.Stat(fsys, s)

	if err == nil && !fi.IsDir() {
		return true
	}

	return false
}
