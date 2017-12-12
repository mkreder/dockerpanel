setenforce 0
systemctl stop postfix
systemctl disable postfix
systemctl  start docker
systemctl  enable docker