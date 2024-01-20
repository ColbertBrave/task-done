#!usr/bin/env bash
# 只要出错就退出
set -o errexit
# 只要有未定义的变量就退出
set -o nounset
# 只要管道中任何一个命令失败就退出
set -o pipefail
# 导入公用脚本
source service.sh


# 设置环境变量
export GO111MODULE=on
export GOPROXY=https://goproxy.cn

# 进入项目目录
cd "${workspace}/cmd"

# 编译
go build -gcflags=-trimpath="$GOPATH" -ldflags "-w -s" -buildmode=pie -o ${pkg_name}
strip ${pkg_name}
