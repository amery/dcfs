package dcfs

import (
	"io/fs"
	"syscall"

	"github.com/armon/go-radix"
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

func (record *NodeRecord) NewNode() (Node, error) {
	switch record.Type {
	case NodeTypeDirectory:
		node := &DirectoryNode{
			record: record,
			tree:   radix.New(),
		}
		return node, nil
	case NodeTypeFile:
		node := &FileNode{
			record: record,
		}
		return node, nil
	case NodeTypeArchive:
		node := &ArchiveNode{
			record: record,
		}
		return node, nil
	default:
		return nil, syscall.EINVAL
	}
}

func (fsys *Filesystem) newNode(inode uint64, typ NodeType) (Node, error) {
	record := &NodeRecord{
		Inode: inode,
		Type:  typ,
	}

	if _, err := fsys.putRecord(record); err != nil {
		return nil, err
	}

	if node, err := record.NewNode(); err != nil {
		fsys.deleteRecord(record)
		return nil, err
	} else {
		return node, nil
	}
}

func (fsys *Filesystem) newDirectory(inode uint64) (Node, error) {
	return fsys.newNode(inode, NodeTypeDirectory)
}

func (fsys *Filesystem) newFile(inode uint64) (Node, error) {
	return fsys.newNode(inode, NodeTypeFile)
}

func (fsys *Filesystem) newArchive(inode uint64) (Node, error) {
	return fsys.newNode(inode, NodeTypeArchive)
}
