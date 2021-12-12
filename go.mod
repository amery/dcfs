module github.com/amery/dcfs

go 1.17

require (
	github.com/ancientlore/go-avltree v1.0.2
	github.com/armon/go-radix v1.0.0
	github.com/spf13/cobra v1.3.0
	github.com/timshannon/bolthold v0.0.0-20210913165410-232392fc8a6a
	github.com/ulikunitz/xz v0.5.10
	github.com/zeebo/blake3 v0.2.1
	go.sancus.dev/core v0.16.0
	go.sancus.dev/fs v0.0.1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.etcd.io/bbolt v1.3.6 // indirect
	golang.org/x/sys v0.0.0-20211205182925-97ca703d548d // indirect
)

replace (
	go.sancus.dev/core => ../../../go.sancus.dev/core
	go.sancus.dev/fs => ../../../go.sancus.dev/fs
)
