package dcfs

import (
	"io/fs"
	"log"
	"syscall"

	"github.com/ancientlore/go-avltree"
	"github.com/armon/go-radix"
	"github.com/timshannon/bolthold"
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

func (fsys *Filesystem) init() error {
	fsys.nodes = avltree.New(compareNode, 0)

	root, err := fsys.getNode(1)
	if err == bolthold.ErrNotFound {
		// create root
		root, err = fsys.newDirectory(1)
	}

	if err == nil {
		fsys.root = root.(*DirectoryNode)
	}

	return err
}

func (fsys *Filesystem) getNode(inode uint64) (Node, error) {
	fsys.mu.Lock()
	defer fsys.mu.Unlock()

	// check on the tree
	if v := fsys.nodes.Find(inode); v != nil {
		return v.(Node), nil
	}

	// check on the database
	record, err := fsys.getRecord(inode)
	if err != nil {
		return nil, err
	}

	// convert db entry to Node
	node, err := record.NewNode()
	if err != nil {
		return nil, err
	}

	// insert node onto the tree
	if _, dupe := fsys.nodes.Add(node); dupe {
		// unexpected error, we previously confirmed
		// the node wasn't present on the tree
		log.Panicf("Node %v duplicated", node.Inode())
	}

	// inserted
	return node, nil
}

// compares two Nodes in the AVL tree
func compareNode(a interface{}, b interface{}) int {
	na := a.(Node).Inode()
	nb := b.(Node).Inode()

	if na < nb {
		return -1
	} else if na > nb {
		return 1
	} else {
		return 0
	}
}
