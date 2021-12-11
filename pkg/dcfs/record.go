package dcfs

import (
	"fmt"

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

	if node.Inode == 0 {
		key = bh.NextSequence()
	} else {
		key = node.Inode
	}

	if err := fsys.db.Insert(key, node); err != nil {
		return 0, err
	} else {
		return node.Inode, nil
	}
}

func (fsys *Filesystem) updateRecord(node *NodeRecord) error {
	return fsys.db.Update(node.Inode, node)
}

func (fsys *Filesystem) deleteRecord(node *NodeRecord) error {
	return fsys.db.Delete(node.Inode, node)
}
