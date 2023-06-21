//ProtoBuf使用方法
1.安装下载protoc，很多种安装方法，下载地址https://github.com/google/protobuf/releases
2.安装下载proto的go插件，命令是go get github.com/golang/protobuf/protoc-gen-go，也可以自己手动下载安装（如果使用go get则会自动生成protoc-gen-go的可执行文件）
3.将protoc-gen-go可执行文件路径加到PATH环境变量中，如果是go get安装是会在GOBIN路径下生成protoc-gen-go，执行export PATH=$PATH:$GOBIN（原因在于, protoc-gen-go可执行文件需要被protoc调用）
4.安装goprotobuf库（注意，protoc-gen-go只是一个插件，goprotobuf的其他功能比如marshal、unmarshal等功能还需要由protobuf库提供）go get github.com/golang/protobuf/proto
5.写example.proto文件以及.go文件测试。由于proto生成go文件的命令是protoc --go_out=./ *（官方文档:https://developers.google.com/protocol-buffers/docs/gotutorial)

go run main.go -c config/dev/config.go -t scanner

/Users/liqing/GoPath/src/github.com/go-xorm/cmd/xorm/xorm  reverse mysql "skygo_detection:05df42d112f90122@(10.228.64.139:2147)/skygo_detection?charset=utf8" /Users/liqing/GoPath/src/github.com/go-xorm/cmd/xorm/templates/goxorm ./