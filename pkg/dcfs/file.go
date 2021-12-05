package dcfs

var (
	_ Node = (*File)(nil)
)

type File struct {
	entry *NodeEntry
}

// Node
func (file *File) Inode() uint64  { return file.entry.Inode }
func (file *File) Type() NodeType { return file.entry.Content.Type() }
