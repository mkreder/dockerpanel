
# Install instructions for Centos7 (Should be similar for RHEL7)
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