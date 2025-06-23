package worker

import (
	"github.com/mkreder/dockerpanel/model"
	"log"
	"os"
	"os/exec"
	"time"

	"strings"
	"strconv"
	"bytes"
	"regexp"
	"io/ioutil"
)

func check(e error){
	if e != nil {
		panic(e)
	}
}

func checkCmd(e error, stderr bytes.Buffer){
	if e != nil {
		log.Printf("Error: " + stderr.String())
	}
}

func createDirectoryConfigDNS(zone model.Zone){
	if _ , err := os.Stat("configs/dns/Dockerfile"); os.IsNotExist(err){
		log.Printf("Creating configuration files")
		// No hay un Dockerfile creado
		srcFolder := "defaults/dns/Dockerfile"
		os.MkdirAll("configs/dns",0755)
		destFolder := "configs/dns/"
		cpCmd := exec.Command("cp", "-rf", srcFolder, destFolder)
		err := cpCmd.Run()
		check(err)
		srcFolder = "defaults/dns/named.conf"
		os.MkdirAll("configs/dns/bind",0755)
		destFolder = "configs/dns/bind/"
		cpCmd2 := exec.Command("cp", "-rf", srcFolder, destFolder)
		err = cpCmd2.Run()
		check(err)
	}

	if _ , err := os.Stat("configs/dns/bind/conf.d/" + zone.Domain + ".conf"); os.IsNotExist(err) {
		os.MkdirAll("configs/dns/bind/conf.d/", 0755)
		f, err := os.Create("configs/dns/bind/conf.d/" + zone.Domain + ".conf")
		check(err)

		_, _ = f.WriteString("zone " + zone.Domain + " IN { \n" +
			"   type master; \n" +
			"    file \"/etc/bind/zones/" + zone.Domain + ".conf\"; \n" +
			"};\n")
		f.Sync()
		f.Close()


		namedf, err := os.OpenFile("configs/dns/bind/named.conf", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}

		namedf.WriteString("include \"/etc/bind/conf.d/" + zone.Domain  + ".conf\"; \n")
		namedf.Sync()
		namedf.Close()

	}
}

func processRecords(zone model.Zone){
	log.Printf("Processing records for " + zone.Domain)
	nsname := ""
	for _ , record := range zone.Records {
		if ( record.Type == "NS" ){
			nsname = record.Value
			break
		}
	}
	serial := strconv.Itoa(int(time.Now().Unix()))
	newZone := "@	IN	SOA " + nsname + "\t"+ strings.Replace(zone.Email,"@",".",-1)+ " ( \n"+
		"\t\t\t\t\t\t" + serial + "\n" +
		"\t\t\t\t\t\t 8H \n" +
		"\t\t\t\t\t\t 2H \n" +
		"\t\t\t\t\t\t 4W \n" +
		"\t\t\t\t\t\t 1D ) \n"

	for _ , record := range zone.Records {
		if ( record.Type == "NS" ){
			newZone = newZone + "\t\t\t\t\t\t\t NS \t" + record.Value + "\n"
		}
	}

	for _ , record := range zone.Records {
		if ( record.Type == "MX" ){
			newZone = newZone + "\t\t\t\t\t\t\t MX \t" + record.Priority + " " + record.Value + "\n"
		}
	}

	for _ , record := range zone.Records {
		if ( record.Type != "MX" ) && ( record.Type != "NS" ) {
			newZone = newZone + record.Name + "\t\t" + record.Type + "\t\t" + record.Value + "\n"
		}
	}

	os.Mkdir("configs/dns/bind/zones/",0755)
	f2, err := os.Create("configs/dns/bind/zones/" + zone.Domain + ".conf" )
	check(err)
	_, _ = f2.WriteString(newZone)
	f2.Sync()
	f2.Close()
}


func buildDNSContainer(){
	cmdString := "docker images -q dp-img-dns"
	imgCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	imgCmd.Stderr = &stderr
	imgCmd.Stdout = &out
	err := imgCmd.Run()
	checkCmd(err,stderr)
	if len(out.String()) == 0 {
		log.Printf("Building the image")
		cmdString := "cd configs/dns/; docker build -t dp-img-dns ."
		buildCmd := exec.Command("/bin/sh" , "-c", cmdString)
		var out2 bytes.Buffer
		var stderr2 bytes.Buffer
		buildCmd.Stdout = &out2
		buildCmd.Stderr = &stderr2
		err := buildCmd.Run()
		log.Printf(out2.String())
		checkCmd(err,stderr2)
		log.Printf("Image created")
	} else {
		log.Printf("Image already exists, skipping step.")
	}
}

func runDNSContainer(name string){
	wd, _ := os.Getwd()
	cmdString := "cd configs/dns/; docker stop dp-dns; docker rm dp-dns; docker run -d -p 53:53/udp -v " + wd + "/configs/dns/bind:/etc/bind:ro" + " --name " + name + " dp-img-dns"
	tarCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	tarCmd.Stderr = &stderr
	tarCmd.Stdout = &out
	err := tarCmd.Run()
	checkCmd(err,stderr)
	log.Printf("Container started " + name + " " + out.String() )
}



func removeZone(zone model.Zone){
	os.Remove("configs/dns/bind/conf.d/" + zone.Domain + ".conf")
	os.Remove("configs/dns/bind/zones/" + zone.Domain + ".conf" )

	namedBuf, err := ioutil.ReadFile("configs/dns/bind/named.conf")
	check(err)

	re := regexp.MustCompile("(?m)^.*" + zone.Domain + ".*$[\r\n]+")
	res := re.ReplaceAllString(string(namedBuf), "")

	err = ioutil.WriteFile("configs/dns/bind/named.conf", []byte(res), 0644)
	if err != nil {
		panic(err)
	}

	zoneID := strconv.Itoa(int(zone.ID))
	log.Printf("Deleting zone " + zone.Domain + " " + zoneID)
	model.Mgr.RemoveZone(zoneID)
}

func RunDNSWorker (){
	log.Printf("Starting DNS worker")
	// Loop forever
	for {
		for _ , zone := range model.Mgr.GetAllZones() {
			if zone.Status == 1 {
				log.Printf("Working on zone " + zone.Domain )

				createDirectoryConfigDNS(zone)
				processRecords(zone)
				buildDNSContainer()

				zone.Status = 2
				model.Mgr.UpdateZone(&zone)
			} else if zone.Status == 2 || zone.Status == 4 {
				if ! isRunning("dp-dns"){
					runDNSContainer("dp-dns")
				} else {
					restartContainer("dp-dns")
				}
				zone.Status = 3
				model.Mgr.UpdateZone(&zone)
			} else if zone.Status == 3 {
				if ! isRunning("dp-dns"){
					zone.Status = 4
					model.Mgr.UpdateZone(&zone)
				}
			} else if zone.Status == 5 {
				removeZone(zone)
				restartContainer("dp-dns")
			}

		}
		time.Sleep(2 * time.Second)
	}
}