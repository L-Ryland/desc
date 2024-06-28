#!/bin/bash

# 读取配置文件
if [ -f "../config.conf" ]; then
    source "../config.conf"
fi

# 输出配置
echo "当前配置："
echo "serverpath: $serverpath"
echo "reactPort: $reactPort"

cd ~/$serverpath/reactserver
serve -s build -l $reactPort
