# Docker Panel

This repo contains my engineering thesis, a web control panel to manage several services (web, email, dns, ftp, databases). 
Each service is deployed in a separate container. This improves security, as each service is isolated on its own container. An attacker gaining access to a container, for instance, a web server, will not be able to get access to the databases that this server will host as they will be running in different containers.  

PRs are welcome. 

## Install instructions for CentOS7 (Should be similar for RHEL7)
1. Install git and docker

```
yum install -y git docker
```

2. Install Go

```
sudo rpm --import https://mirror.go-repo.io/centos/RPM-GPG-KEY-GO-REPO
curl -s https://mirror.go-repo.io/centos/go-repo.repo | sudo tee /etc/yum.repos.d/go-repo.repo
sudo yum -y install golang
```

3. Run preparation scripts:

```
./prepare_centos.sh
./install_deps.sh
```

4. Run main.go:

```
go run main.go
```

Login to server IP on port 9090 with username "admin@admin.com" and password "admin"

## TODO: 

* English translation
* Improve user interface (specially user management)