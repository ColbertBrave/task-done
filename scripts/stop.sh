#!usr/bin/env bash
# 只要出错就退出
set -o errexit
# 只要有未定义的变量就退出
set -o nounset
# 只要管道中任何一个命令失败就退出
set -o pipefail
# 开启shell调试模式，打印命令和输出
set -x
# 导入公用脚本
source service.sh

# 停止进程，并检查是否退出
for i in {1..10}
do
    if ! ps aux | grep -v grep | grep "$pkg_name" > /dev/null; then
        log "Process $pkg_name is not running, so exit"
        exit 0
    fi

    log "Process $pkg_name is running. Killing the process..."    
    # 获取进程ID并杀死进程
    process_id=$(ps aux | grep -v grep | grep "$pkg_name" | awk '{print $2}')
    kill "$process_id"

    if ! ps aux | grep -v grep | grep "$pkg_name" > /dev/null; then
        log "Process $pkg_name has been killed successfully and the service exited"
        exit 0
    fi
    sleep 1
done

# 仍然未能杀死进程，异常退出
exit 1