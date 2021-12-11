package dcfs

import (
	"fmt"
	"log"

	bh "github.com/timshannon/bolthold"
)

type NodeRecord struct {
	Inode   uint64 `boltholdKey:"Inode"`
	Type    NodeType
	Content []interface{}
}

func (record *NodeRecord) String() string {
	s := fmt.Sprintf("node:%v/%s", record.Inode, record.Type)

	if l := len(record.Content); l > 0 {
		s += fmt.Sprintf(" [%v]%T", l, record.Content[0])
	}

	return s
}

func (fsys *Filesystem) resetRecordSequence() {
	var max uint64

	// find max
	fsys.db.ForEach(nil, func(node *NodeRecord) error {
		if node.Inode > max {
			max = node.Inode
		}
		return nil
	})

	// set NodeRecord's sequence to the maximum value found
	tx, err := fsys.db.Bolt().Begin(true)
	if err != nil {
		log.Fatal(err)
	}

	bkt := tx.Bucket([]byte("NodeRecord"))
	if err := bkt.SetSequence(max); err != nil {
		log.Fatal(err)
	} else {
		tx.Commit()
	}
}

func (fsys *Filesystem) getRecord(inode uint64) (*NodeRecord, error) {
	node := &NodeRecord{}
	if err := fsys.db.FindOne(node, bh.Where("Inode").Eq(inode)); err != nil {
		return nil, err
	} else {
		return node, nil
	}
}

func (fsys *Filesystem) putRecord(node *NodeRecord) (uint64, error) {
	var key interface{}
	var auto bool

	if node.Inode == 0 {
		auto = true
		key = bh.NextSequence()
	} else {
		key = node.Inode
	}

	for {
		if err := fsys.db.Insert(key, node); err == nil {
			return node.Inode, nil
		} else if err == bh.ErrKeyExists && auto {
			// again
			fsys.resetRecordSequence()
		} else {
			return 0, err
		}
	}
}

func (fsys *Filesystem) appendRecordContent(parent *NodeRecord, child interface{}) error {
	parent.Content = append(parent.Content, child)
	if err := fsys.updateRecord(parent); err != nil {
		parent.Content = parent.Content[:len(parent.Content)-1]
		return err
	}
	return nil
}

func (fsys *Filesystem) updateRecord(node *NodeRecord) error {
	return fsys.db.Update(node.Inode, node)
}

func (fsys *Filesystem) deleteRecord(node *NodeRecord) error {
	return fsys.db.Delete(node.Inode, node)
}
