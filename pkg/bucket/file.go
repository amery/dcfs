package bucket

import (
	_ "github.com/ulikunitz/xz"
	_ "github.com/zeebo/blake3"
)

type File struct {
	Hash string
	Name string
}
