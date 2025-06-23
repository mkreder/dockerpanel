package worker

import (
	"bytes"
	"log"
	"os/exec"
)

func restartContainer(name string){
	cmdString := "docker restart -t 1 " + name
	tarCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	tarCmd.Stderr = &stderr
	err := tarCmd.Run()
	checkCmd(err,stderr)
	log.Printf("Contenedor reiniciado " + name)
}

func isRunning(name string) bool{
	cmdString := "docker ps -q -f  name=" + name + "$"
	psCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	psCmd.Stderr = &stderr
	psCmd.Stdout = &out
	err := psCmd.Run()
	checkCmd(err,stderr)
	if len(out.String()) != 0 {
		return true
	} else {
		return false
	}
}