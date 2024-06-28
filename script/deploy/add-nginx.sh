#!/bin/bash

serverID=$1
serverName=$2
serverpath=$3

DBpath="webStorage_"
DBpath+=$serverName
golangPort=8080
reactPort=3000
if [ -n "$serverID" ]; then
    # 如果 serverID 不为空，进行算术运算
    golangPort=$((golangPort + serverID))
    reactPort=$((reactPort + serverID))
fi

# 输出配置
echo "当前配置："
echo "serverpath: $serverpath"
echo "serverName: $serverName"
echo "golangPort: $golangPort"
echo "reactPort: $reactPort"

# set path
cd ~
mkdir -p $serverpath
cd $serverpath
mkdir -p golangserver
mkdir -p reactserver

# set config
echo "
serverID=\"$serverID\"
serverName=\"$serverName\"
serverpath=\"$serverpath\"
DBpath=\"$DBpath\"
golangPort=\"$golangPort\"
reactPort=\"$reactPort\"
" >~/$serverpath/config.conf

echo "setting nginx"

# set Nginx
server_block="server {
    listen 80;
    server_name $serverName.dytx2tyxt.com; # 域名地址

    location / {
        proxy_pass http://localhost:$reactPort; # 转发到React开发服务器
    }

    location /v1 {
        proxy_pass http://localhost:$golangPort; # 转发到Go应用
    }
}"
nginx_conf="/etc/nginx/nginx.conf"

# 检查是否存在相同的 server_name 或端口
if grep -qF "$serverName.dytx2tyxt.com" "$nginx_conf" ||
    grep -qF "http://localhost:$golangPort;" "$nginx_conf" ||
    grep -qF "http://localhost:$reactPort;" "$nginx_conf"; then
    echo "Nginx configuration already exists for server_name or port. Aborting."
    exit 1
else
    # 如果不存在，添加配置
    awk -v config="$server_block" '/http {/ {print; print config; next} 1' "$nginx_conf" | sudo tee "$nginx_conf" >/dev/null
    sudo nginx -s reload
fi
