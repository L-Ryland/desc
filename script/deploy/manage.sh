#!/bin/bash

serverName=$1

# 添加 Supervisor 配置
cat <<EOF > /etc/supervisor/conf.d/golangserver_$serverName.conf
[program:golangserver_$serverName]
command=bash "./golang-start.sh"
directory=/root/webStorage_$serverName/golangserver
autostart=true
autorestart=true
delay=10
stopasgroup=true
stopsignal=INT
redirect_stderr=true
stdout_logfile=/root/webStorage_$serverName/golangserver/golangserver.log
EOF

cat <<EOF > /etc/supervisor/conf.d/reactserver_$serverName.conf
[program:reactserver_$serverName]
command=bash "./react-start.sh"
directory=/root/webStorage_$serverName/reactserver
autostart=true
autorestart=true
delay=10
stopasgroup=true
stopsignal=INT
redirect_stderr=true
stdout_logfile=/root/webStorage_$serverName/reactserver/reactserver.log
EOF

# 重新加载 Supervisor 配置
supervisorctl reread && supervisorctl update

# 添加 Logrotate 配置
cat <<EOF > /etc/logrotate.d/golangserver_$serverName
/root/webStorage_$serverName/golangserver/golangserver.log {
    rotate 7
    daily
    missingok
    notifempty
    compress
    delaycompress
    create 644 root root
}
EOF

cat <<EOF > /etc/logrotate.d/reactserver_$serverName
/root/webStorage_$serverName/reactserver/reactserver.log {
    rotate 7
    daily
    missingok
    notifempty
    compress
    delaycompress
    create 644 root root
}
EOF
