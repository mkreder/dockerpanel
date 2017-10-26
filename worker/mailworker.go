package worker

import (
	"github.com/mkreder/dockerpanel/model"
	"os"
	"os/exec"
	"time"

	"bytes"
	"io/ioutil"
	"log"
)

func crearDirectorioConfigMail(){
	if _ , err := os.Stat("configs/mail/postfix/Dockerfile"); os.IsNotExist(err){
		log.Printf("Creando archivos de configuración")
		os.MkdirAll("configs/mail",0755)
		srcFolder := "defaults/mail/postfix"
		destFolder := "configs/mail/"
		cpCmd2 := exec.Command("cp", "-rf", srcFolder, destFolder)
		err = cpCmd2.Run()
		check(err)
		srcFolder = "defaults/mail/dovecot"
		cpCmd3:= exec.Command("cp", "-rf", srcFolder, destFolder)
		err = cpCmd3.Run()
		check(err)
		srcFolder = "defaults/mail/mailman"
		cpCmd4:= exec.Command("cp", "-rf", srcFolder, destFolder)
		err = cpCmd4.Run()
		check(err)

		os.MkdirAll("data/mail/spool",0755)
		os.MkdirAll("data/mail/mailboxes",0755)
	}
}

func generarConfig(){
	vdomains := ""
	vmailboxs := ""

	// Generar virtual_host y vmailbox
	for _ , dominio := range model.Mgr.GetAllDominios(){
		vdomains = vdomains + dominio.Nombre + "\n"

		for _ , cuenta := range dominio.Cuentas {
			vmailboxs = vmailboxs + cuenta.Nombre + "@" + dominio.Nombre + " " + dominio.Nombre + "/" + cuenta.Nombre + "/ \n"
		}
	}

	err := ioutil.WriteFile("configs/mail/postfix/conf/virtual_domains", []byte(vdomains), 0644)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("configs/mail/postfix/conf/vmailbox", []byte(vmailboxs), 0644)
	if err != nil {
		panic(err)
	}

	// Generar passwd y shadow
	for _ , dominio := range model.Mgr.GetAllDominios() {
		os.MkdirAll("configs/mail/dovecot/conf/auth/"+dominio.Nombre, 0755)
		passwd := ""
		shadow := ""
		for _, cuenta := range dominio.Cuentas {
			passwd = passwd + cuenta.Nombre + "@" + dominio.Nombre + "::101:104:/var/mail/vhosts/" + dominio.Nombre + "/"+  cuenta.Nombre + "\n"
			shadow = shadow + cuenta.Nombre + "@" + dominio.Nombre + ":{DIGEST-MD5}" + cuenta.Password + "\n"
		}
		err := ioutil.WriteFile("configs/mail/dovecot/conf/auth/"+dominio.Nombre+"/passwd", []byte(passwd), 0644)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile("configs/mail/dovecot/conf/auth/"+dominio.Nombre+"/shadow", []byte(shadow), 0644)
		if err != nil {
			panic(err)
		}

	}
}


func prepararImagenPostfix(){
	cmdString := "docker images -q dp-img-mail-postfix"
	imgCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	imgCmd.Stderr = &stderr
	imgCmd.Stdout = &out
	err := imgCmd.Run()
	checkCmd(err,stderr)
	if len(out.String()) == 0 {
		log.Printf("Construyendo la imagen")
		cmdString := "cd configs/mail/postfix/; docker build -t dp-img-mail-postfix ."
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

func prepararImagenDovecot(){
	cmdString := "docker images -q dp-img-mail-dovecot"
	imgCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	imgCmd.Stderr = &stderr
	imgCmd.Stdout = &out
	err := imgCmd.Run()
	checkCmd(err,stderr)
	if len(out.String()) == 0 {
		log.Printf("Construyendo la imagen")
		cmdString := "cd configs/mail/dovecot/; docker build -t dp-img-mail-dovecot ."
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

func correrContenedorDovecot(){
	wd, _ := os.Getwd()
	cmdString := "cd configs/mail/dovecot; docker stop dp-mail-dovecot; docker rm dp-mail-dovecot; docker run -d -p 110:110 -p 143:143 -v " + wd + "/configs/mail/dovecot/conf:/etc/dovecot:ro -v " + wd + "/data/mail/mailboxes:/var/mail" + " --name dp-mail-dovecot dp-img-mail-dovecot"
	tarCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	tarCmd.Stderr = &stderr
	tarCmd.Stdout = &out
	err := tarCmd.Run()
	checkCmd(err,stderr)
	log.Printf("Contenedor iniciado dp-mail-dovecot iniciado" + out.String() )
	time.Sleep(2*time.Second)
}

func correrContenedorPostfix(){
	wd, _ := os.Getwd()
	cmdString := "cd configs/mail/postfix; docker stop dp-mail-postfix; docker rm dp-mail-postfix; docker run -d -p 25:25 -v " + wd + "/configs/mail/postfix/conf:/etc/postfix -v " + wd + "/data/mail/mailboxes:/var/mail -v " + wd + "/data/mail/spool:/var/spool/postfix" + " --name dp-mail-postfix  --link dp-mail-dovecot:dp-mail-dovecot  dp-img-mail-postfix"
	tarCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	tarCmd.Stderr = &stderr
	tarCmd.Stdout = &out
	err := tarCmd.Run()
	checkCmd(err,stderr)
	log.Printf("Contenedor iniciado dp-mail-postfix " + out.String())
}

//
//func removerDominio(dominio model.Dominio){
//	os.Remove("configs/mail/bind/conf.d/" + dominio.Nombre + ".conf")
//	os.Remove("configs/mail/bind/zones/" + dominio.Nombre + ".conf" )
//
//	namedBuf, err := ioutil.ReadFile("configs/mail/bind/named.conf")
//	check(err)
//
//	re := regexp.MustCompile("(?m)^.*" + dominio.Nombre + ".*$[\r\n]+")
//	res := re.ReplaceAllString(string(namedBuf), "")
//
//	err = ioutil.WriteFile("configs/mail/bind/named.conf", []byte(res), 0644)
//	if err != nil {
//		panic(err)
//	}
//
//	dominioID := strconv.Itoa(int(dominio.ID))
//	log.Printf("Eliminando dominio " + dominio.Nombre + " " + dominioID)
//	model.Mgr.RemoveDominio(dominioID)
//}

func RunMailWorker () {
	log.Printf("Iniciando Mail worker")
	// Loop para siempre
	for {
		for _, dominio := range model.Mgr.GetAllDominios() {
			if dominio.Estado == 1 {
				log.Printf("Trabajando en la dominio " + dominio.Nombre)

				crearDirectorioConfigMail()
				generarConfig()
				prepararImagenDovecot()
				prepararImagenPostfix()

				dominio.Estado = 2
				model.Mgr.UpdateDominio(&dominio)
			} else if dominio.Estado == 2 || dominio.Estado == 4 {
				if ! isRunning("dp-mail-dovecot") {
					correrContenedorDovecot()
				} else {
					reiniciarContenedor("dp-mail-dovecot")
				}
				if ! isRunning("dp-mail-postfix") {
					correrContenedorPostfix()
				} else {
					reiniciarContenedor("dp-mail-postfix")
				}
				dominio.Estado = 3
				model.Mgr.UpdateDominio(&dominio)
			} else if dominio.Estado == 3 {
				if ! isRunning("dp-mail-dovecot") || ! isRunning("dp-mail-postfix") {
					dominio.Estado = 4
					model.Mgr.UpdateDominio(&dominio)
				}
				//} else if dominio.Estado == 5 {
				//	removerDominio(dominio)
				//	reiniciarContenedor("dp-mail")
				//}

			}
			time.Sleep(2 * time.Second)
		}
	}
}

