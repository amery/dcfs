package dcfs

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"syscall"

	"github.com/ancientlore/go-avltree"
	"github.com/timshannon/bolthold"

	"go.sancus.dev/core/errors"
)

type NodeType int

const (
	NodeTypeUndefined NodeType = iota
	NodeTypeDirectory
	NodeTypeFile
	NodeTypeArchive
)

func (t NodeType) String() string {
	const mnemonic = "?DFA"

	if t >= 0 && int(t) < len(mnemonic) {
		return string(mnemonic[t])
	} else {
		return fmt.Sprintf("%v", int(t))
	}
}

type Node interface {
	Inode() uint64
	Type() NodeType
	Open() (fs.File, error)
}

func (record *NodeRecord) NewNode() (Node, error) {
	log.Printf("%+n: %s:%v", errors.Here(), "inode", record.Inode)

	switch record.Type {
	case NodeTypeDirectory:
		node := &DirectoryNode{
			record: record,
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
	log.Printf("%+n: %s:%v %s:%v", errors.Here(), "inode", inode, "type", typ)

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

func (fsys *Filesystem) newDirectory(inode uint64) (*DirectoryNode, error) {
	if v, err := fsys.newNode(inode, NodeTypeDirectory); err != nil {
		return nil, err
	} else {
		return v.(*DirectoryNode), nil
	}
}

func (fsys *Filesystem) newFile(inode uint64) (*FileNode, error) {
	if v, err := fsys.newNode(inode, NodeTypeFile); err != nil {
		return nil, err
	} else {
		return v.(*FileNode), nil
	}
}

func (fsys *Filesystem) newArchive(inode uint64) (*ArchiveNode, error) {
	if v, err := fsys.newNode(inode, NodeTypeArchive); err != nil {
		return nil, err
	} else {
		return v.(*ArchiveNode), nil
	}
}

func (fsys *Filesystem) init() error {
	log.Printf("%+n", errors.Here())
	fsys.nodes = avltree.New(compareNode, 0)

	root, err := fsys.getNode(1)
	if err == bolthold.ErrNotFound {
		// create root
		root, err = fsys.newDirectory(1)
	}

	if err == nil {
		fsys.root = root.(*DirectoryNode)
		fsys.root.populate(fsys, false)
		fsys.spawn(func(_ context.Context, _ *Filesystem) {
			fsys.root.populate(fsys, true)
		})
	}

	return err
}

func (fsys *Filesystem) getNode(inode uint64) (Node, error) {
	log.Printf("%+n: %s:%v", errors.Here(), "inode", inode)

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
