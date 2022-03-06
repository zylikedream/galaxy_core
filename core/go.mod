module github.com/zylikedream/galaxy/core

go 1.16

require (
	github.com/fsnotify/fsnotify v1.5.1
	github.com/gookit/goutil v0.4.4
	github.com/mitchellh/mapstructure v1.4.3
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.9.1
	github.com/rpcxio/rpcx-etcd v0.0.0-20210907081219-a9e31da236e8
	github.com/smallnest/rpcx v1.6.11
	github.com/spf13/viper v1.9.0
	go.mongodb.org/mongo-driver v1.8.0
	go.uber.org/zap v1.17.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)

replace github.com/smallnest/rpcx v1.6.11 => ../../zyrpcx
