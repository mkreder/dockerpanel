setenforce 0
sed "/SELINUX=/s/enforcing/disabled/" /etc/selinux/config -i
systemctl stop postfix
systemctl disable postfix
systemctl  start docker
systemctl  enable docker
docker pull debian
docker pull ubuntu
docker pull alpine