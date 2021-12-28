package bucket

import (
	"sort"
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
