package bucket

import (
	"log"
	"sort"
	"strings"

	"go.sancus.dev/core/errors"
)

type Files []File

func (m Files) Len() int      { return len(m) }
func (m Files) Sort()         { sort.Sort(m) }
func (m Files) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

func (m Files) Less(i, j int) bool {
	// empty ones go at the end for packing
	if a := m[i].Name; a == "" {
		return false
	} else if b := m[j].Name; b == "" {
		return true
	} else {
		return a < b
	}
}

func (m Files) Reset(i int) {
	if i < len(m) {
		m[i] = File{}
	}
}

func (m Files) Repack() {
	m.Sort()

	n := len(m)
	for n > 0 {
		if m[n-1].Name == "" {
			n--
		}
	}

	m = m[:n]
}

func (m *Files) Append(f File) {
	if f.Name != "" {
		*m = append(*m, f)
	}
}

func (m Files) ForEach(prefix string, fn func(int, string)) {

	switch {
	case fn == nil:
		// done
	case prefix == "." || prefix == "":
		// all
		for i := range m {
			if s := m[i].Name; s != "" {
				fn(i, s)
			}
		}
	default:
		// within path
		prefix += "/"

		for i := range m {
			s := m[i].Name
			s1 := strings.TrimPrefix(s, prefix)
			if s1 != "" && s != s1 {
				fn(i, s1)
			}
		}
	}
}

func (m *Bucket) Move(dir string, dest *Bucket, prefix string) {
	if prefix == "." {
		prefix = ""
	}

	m.Files.ForEach(dir, func(i int, name string) {

		f := m.Files[i]

		if prefix != "" {
			f.Name = prefix + "/" + name
		} else {
			f.Name = name
		}

		log.Printf("%+n: %s:%q %s:%v %s:%q -> %q",
			errors.Here(),
			"dir", dir,
			"i", i,
			"name", name, f.Name)

		dest.Files.Append(f)
		m.Files.Reset(i)
	})

	m.Files.Repack()
	dest.Files.Repack()
}
