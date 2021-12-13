package dcfs

import (
	"context"
)

func (fsys *Filesystem) ForEachRecord(f func(context.Context, *NodeRecord) error) {
	if f != nil {
		fsys.db.ForEach(nil, func(node *NodeRecord) error {
			return f(fsys.ctx, node)
		})
	}
}
