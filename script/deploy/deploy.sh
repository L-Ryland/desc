#!/bin/bash
# Ubuntu Server 22.04 LTS 64bit

# set network
echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDEnVEHV9vPJFe2VU6N27u1KEvOXF+pzGL9yaZkdZDTx4lTtwh1bRO3fFQPn7qd/kBrLI4ohweMSnKLlHOPAVPZEAx3vLHzy1vCYH7/pEEh3ReiFanCYrZpWmrra4BISWcvriWxe8ZVzk0KTzF171hwjQgCbLSIlNQnQAIqZxSZ8VRU8BEOY9GFerulcHkjj9ow7cMj1JJdRrptnt1zUaskISCDbYW8nYUW0akd5Ozf4PR1NToVT9mXpGWIiJj1Zl2JHaz4JVtUpIzI7X/Q/2kcxfDMRgaNI2YBAFCoDBU+z2qqRm0nJ9Bc1kxp0zagUmDPu9vEzR99BlHl6u+VtjJn skey-ql12zev3" >>./.ssh/authorized_keys


sudo apt-get update
sudo apt-get install -y gnupg curl

# set nginx
sudo apt-get install -y nginx

# set nodejs
sudo apt-get install -y ca-certificates curl gnupg
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
NODE_MAJOR=18
echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | sudo tee /etc/apt/sources.list.d/nodesource.list
sudo apt-get update
sudo apt-get install -y nodejs
sudo npm install -g serve

# set mongodb
curl -fsSL https://pgp.mongodb.com/server-6.0.asc |
    sudo gpg -o /usr/share/keyrings/mongodb-server-6.0.gpg \
        --dearmor
echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-6.0.gpg ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/6.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-6.0.list
sudo apt-get update
sudo apt-get install -y mongodb-org
sudo systemctl start mongod
sudo systemctl enable mongod

sudo apt-get install supervisor
