#!usr/bin/env bash
# 项目根目录
workspace=$(cd $(dirname "$0")/.. || exit; pwd)
# 二进制包名称
pkg_name=cloud_disk
# 指定启动脚本的用户
user=ranxuening
# 脚本日志输出目录
log_path=$workspace/logs/system.log

# 定义日志函数
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S')|SYSTEM|$1" >> $log_path
}