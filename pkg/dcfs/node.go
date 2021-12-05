package dcfs

import (
	"io/fs"
)

type NodeType int

const (
	NodeTypeUndefined NodeType = iota
	NodeTypeDirectory
	NodeTypeFile
	NodeTypeArchive
)

type Node interface {
	Inode() uint64
	Type() NodeType
	Open() (fs.File, error)
}

type NodeContent interface {
	Type() NodeType
}

type NodeEntry struct {
	Inode   uint64 `boltholdKey:"Inode"`
	Content NodeContent
}

//
type NodeDirectoryContent struct {
	Children []NodeChild
}

func (_ NodeDirectoryContent) Type() NodeType { return NodeTypeDirectory }

type NodeChild struct {
	Inode uint64 `boltholdIndex:"Inode"`
	Name  string
}

//
type Blake3Hash []byte

type NodeFileContent struct {
	Versions []Blake3Hash
}

func (_ NodeFileContent) Type() NodeType { return NodeTypeFile }

//
type NodeArchiveContent struct {
	Entries []NodeArchiveEntry
}

func (_ NodeArchiveContent) Type() NodeType { return NodeTypeArchive }

type NodeArchiveEntry struct {
	Hash Blake3Hash
	Name string
}
