#!/bin/bash

# 如果没有提供参数，则使用默认值
if [ $# -eq 0 ]; then
    serverID="0"
    serverName="www"
elif [ $# -eq 1 ]; then
    # 只有一个参数时，构建 serverName
    serverID=$1
    serverName="test${serverID}"
else
    # 有两个参数时，使用提供的参数
    serverID=$1
    serverName=$2
fi

serverpath="webStorage_"
serverpath+=$serverName

script_directory=$(dirname "$(readlink -f "$0")")

# 本地脚本的路径
local_script_path="$script_directory/deploy/add-nginx.sh"
# 远程服务器上保存脚本的路径
remote_script_path="~/add-nginx.sh"
# 将本地脚本上传到远程服务器
scp $local_script_path webStorage:$remote_script_path
# 执行远程脚本并传递参数
ssh webStorage "bash $remote_script_path" "$serverID" "$serverName" "$serverpath"

local_script_path="$script_directory/deploy/golang-start.sh"
remote_script_path="~/$serverpath/golangserver"
scp $local_script_path webStorage:$remote_script_path

local_script_path="$script_directory/deploy/react-start.sh"
remote_script_path="~/$serverpath/reactserver"
scp $local_script_path webStorage:$remote_script_path

local_script_path="$script_directory/deploy/manage.sh"
remote_script_path="~/$serverpath/manage.sh"
scp $local_script_path webStorage:$remote_script_path
ssh webStorage "bash $remote_script_path" "$serverName"
