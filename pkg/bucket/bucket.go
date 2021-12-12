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
	Metadata interface{} `yaml:",omitempty"`
	Files    Files       `yaml:",omitempty"`
}

func (m *Bucket) Load(fsys fs.FS) error {
	log.Printf("%+n", errors.Here())

	buf, err := fs.ReadFile(fsys, BucketFileName)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(buf, m)
}

func (m *Bucket) Commit(fsys fs.FS) error {
	log.Printf("%+n", errors.Here())

	buf, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	return fs.WriteFile(fsys, BucketFileName, buf, 0644)
}

func New(fsys fs.FS) (*Bucket, error) {
	log.Printf("%+n: %s", errors.Here(), fsys)

	m := &Bucket{}

	err := m.Load(fsys)
	if fs.IsNotExist(err) {
		// attempt to create one
		err = m.Commit(fsys)
	}

	if err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
