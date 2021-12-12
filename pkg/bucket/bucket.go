package bucket

import (
	"log"

	"gopkg.in/yaml.v3"

	"go.sancus.dev/core/errors"
	"go.sancus.dev/fs"
)

const (
	BucketFileName = "dcfs.yaml"
)

type Bucket struct {
	fsys     fs.FS       `yaml:"-"`
	Metadata interface{} `yaml:",omitempty"`
	Files    Files       `yaml:",omitempty"`
}

func (m *Bucket) Load() error {
	log.Printf("%+n", errors.Here())

	buf, err := fs.ReadFile(m.fsys, BucketFileName)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(buf, m)
}

func (m *Bucket) Commit() error {
	log.Printf("%+n", errors.Here())

	buf, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	return fs.WriteFile(m.fsys, BucketFileName, buf, 0644)
}

func New(fsys fs.FS) (*Bucket, error) {
	log.Printf("%+n: %s", errors.Here(), fsys)

	m := &Bucket{
		fsys: fsys,
	}

	err := m.Load()
	if fs.IsNotExist(err) {
		// attempt to create one
		err = m.Commit()
	}

	if err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
