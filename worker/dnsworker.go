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
		panic(e)
	}
}

func crearDirectorioConfig(zona model.Zona){
	if _ , err := os.Stat("configs/dns/Dockerfile"); os.IsNotExist(err){
		log.Printf("Creando archivos de configuración")
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

	if _ , err := os.Stat("configs/dns/bind/conf.d/" + zona.Dominio + ".conf"); os.IsNotExist(err) {
		os.MkdirAll("configs/dns/bind/conf.d/", 0755)
		f, err := os.Create("configs/dns/bind/conf.d/" + zona.Dominio + ".conf")
		check(err)

		_, _ = f.WriteString("zone " + zona.Dominio + " IN { \n" +
			"   type master; \n" +
			"    file \"/etc/bind/zones/" + zona.Dominio + ".conf\"; \n" +
			"};\n")
		f.Sync()
		f.Close()


		namedf, err := os.OpenFile("configs/dns/bind/named.conf", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}

		namedf.WriteString("include \"/etc/bind/conf.d/" + zona.Dominio  + ".conf\"; \n")
		namedf.Sync()
		namedf.Close()

	}
}

func procesarRegistros(zona model.Zona){
	log.Printf("Procesando registros de " + zona.Dominio)
	nsname := ""
	for _ , registro := range zona.Registros {
		if ( registro.Tipo == "NS" ){
			nsname = registro.Valor
			break
		}
	}
	serial := strconv.Itoa(int(time.Now().Unix()))
	nuevaZona := "@	IN	SOA " + nsname + "\t"+ strings.Replace(zona.Email,"@",".",-1)+ " ( \n"+
		"\t\t\t\t\t\t" + serial + "\n" +
		"\t\t\t\t\t\t 8H \n" +
		"\t\t\t\t\t\t 2H \n" +
		"\t\t\t\t\t\t 4W \n" +
		"\t\t\t\t\t\t 1D ) \n"

	for _ , registro := range zona.Registros {
		if ( registro.Tipo == "NS" ){
			nuevaZona = nuevaZona + "\t\t\t\t\t\t\t NS \t" + registro.Valor + "\n"
		}
	}

	for _ , registro := range zona.Registros {
		if ( registro.Tipo == "MX" ){
			nuevaZona = nuevaZona + "\t\t\t\t\t\t\t MX \t" + registro.Prioridad + " " + registro.Valor + "\n"
		}
	}

	for _ , registro := range zona.Registros {
		if ( registro.Tipo != "MX" ) && ( registro.Tipo != "NS" ) {
			nuevaZona = nuevaZona + registro.Nombre + "\t\t" + registro.Tipo + "\t\t" + registro.Valor + "\n"
		}
	}

	os.Mkdir("configs/dns/bind/zones/",0755)
	f2, err := os.Create("configs/dns/bind/zones/" + zona.Dominio + ".conf" )
	check(err)
	_, _ = f2.WriteString(nuevaZona)
	f2.Sync()
	f2.Close()
}


func buildearContenedor(){
	cmdString := "docker images -q dp-dnsimage"
	imgCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	imgCmd.Stderr = &stderr
	imgCmd.Stdout = &out
	err := imgCmd.Run()
	checkCmd(err,stderr)
	if len(out.String()) == 0 {
		log.Printf("Construyendo la imagen")
		cmdString := "cd configs/dns/; docker build -t dp-dnsimage ."
		buildCmd := exec.Command("/bin/sh" , "-c", cmdString)
		var out2 bytes.Buffer
		var stderr2 bytes.Buffer
		buildCmd.Stdout = &out2
		buildCmd.Stderr = &stderr2
		err := buildCmd.Run()
		log.Printf(out2.String())
		checkCmd(err,stderr2)
		log.Printf("Imagen creada")
	} else {
		log.Printf("La imagen ya existe, salteando paso.")
	}
}

func correrContenedor(nombre string){
	wd, _ := os.Getwd()
	cmdString := "cd configs/dns/; docker stop dp-dns; docker rm dp-dns; docker run -d -p 53:53/udp -v " + wd + "/configs/dns/bind:/etc/bind:ro" + " --name " + nombre + " dp-dnsimage"
	tarCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	tarCmd.Stderr = &stderr
	err := tarCmd.Run()
	checkCmd(err,stderr)
	log.Printf("Corriendo contenedor " + nombre)
}

func reiniciarContenedor(nombre string){
	cmdString := "docker restart -t 1 " + nombre
	tarCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	tarCmd.Stderr = &stderr
	err := tarCmd.Run()
	checkCmd(err,stderr)
	log.Printf("Contenedor reiniciado " + nombre)
}

func isRunning(nombre string) bool{
	cmdString := "docker ps -q -f  name=" + nombre
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

func removerZona(zona model.Zona){
	os.Remove("configs/dns/bind/conf.d/" + zona.Dominio + ".conf")
	os.Remove("configs/dns/bind/zones/" + zona.Dominio + ".conf" )

	namedBuf, err := ioutil.ReadFile("configs/dns/bind/named.conf")
	check(err)

	re := regexp.MustCompile("(?m)^.*" + zona.Dominio + ".*$[\r\n]+")
	res := re.ReplaceAllString(string(namedBuf), "")

	err = ioutil.WriteFile("configs/dns/bind/named.conf", []byte(res), 0644)
	if err != nil {
		panic(err)
	}

	zonaID := strconv.Itoa(int(zona.ID))
	log.Printf("Eliminando zona " + zona.Dominio + " " + zonaID)
	model.Mgr.RemoveZona(zonaID)
}

func RunDNSWorker (){
	log.Printf("Iniciando DNS worker")
	// Loop para siempre
	for {
		for _ , zona := range model.Mgr.GetAllZonas() {
			if zona.Estado == 1 {
				log.Printf("Trabajando en la zona " + zona.Dominio )

				crearDirectorioConfig(zona)
				procesarRegistros(zona)
				buildearContenedor()

				zona.Estado = 2
				model.Mgr.UpdateZona(&zona)
			} else if zona.Estado == 2 || zona.Estado == 4 {
				if ! isRunning("dp-dns"){
					correrContenedor("dp-dns")
				} else {
					reiniciarContenedor("dp-dns")
				}
				zona.Estado = 3
				model.Mgr.UpdateZona(&zona)
			} else if zona.Estado == 3 {
				if ! isRunning("dp-dns"){
					zona.Estado = 4
					model.Mgr.UpdateZona(&zona)
				}
			} else if zona.Estado == 5 {
				removerZona(zona)
				reiniciarContenedor("dp-dns")
			}

		}
		time.Sleep(2 * time.Second)
	}
}