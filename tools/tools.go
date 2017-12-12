package tools

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"github.com/mkreder/dockerpanel/model"
	"net"
	"log"
	"io"
	"os/exec"
	"crypto/md5"
	"encoding/hex"
	"strings"
)


func GetIPAddress() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	//return localAddr.IP.String()
	return "35.153.46.76"
}


func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func GetRunningContainers() []model.Container{
	var containers []model.Container

	cmdString := "docker ps --format '{{.ID}};{{.Names}};{{.Image}};{{.Command}};{{.Status}};{{.Ports}}'"
	psCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	psCmd.Stderr = &stderr
	psCmd.Stdout = &out
	psCmd.Run()
	reader := csv.NewReader(bufio.NewReader(strings.NewReader(out.String())))
	reader.Comma = ';'

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		containers = append(containers, model.Container{
			ID: line[0],
			Nombre: line[1],
			Imagen: line[2],
			Comando: line[3],
			Estado: line[4],
			Puertos: line[5],
		})
	}

	return containers
}