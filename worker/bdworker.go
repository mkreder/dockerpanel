
package worker

import "log"
import (
	"github.com/mkreder/dockerpanel/model"
	"time"
	"os"
	"os/exec"
	"strconv"
	"io/ioutil"
	"bytes"
)

func calcularPuerto(id uint) string{
	if ( id < 10 ){
		return "300" + strconv.Itoa(int(id))
	} else if (id < 100){
		return "30" + strconv.Itoa(int(id))
	} else if (id < 1000 ){
		return "3" + strconv.Itoa(int(id))
	} else {
		log.Panic("No soportamos mas de 1000 BDs")
	}
	return ""
}

func crearDirectorioConfigBD(nombre string){
	if _ , err := os.Stat("configs/bd/Dockerfile"); os.IsNotExist(err){
		log.Printf("Creando archivos de configuración")
		srcFolder := "defaults/bd"
		os.MkdirAll("configs/bd",0755)
		destFolder := "configs/"
		cpCmd := exec.Command("cp", "-rf", srcFolder, destFolder)
		err := cpCmd.Run()
		check(err)
		os.MkdirAll("data/bd/" + nombre,0755)
	}
}

func crearConfigPMA(){
	config := "<?php \n$cfg['blowfish_secret'] = 'abcdkfj2dsfsdfdfasjlkdjlkdsjfkdsjfadslkfajdkfjkldsjfdaflkdsafdx83kx';\n$i = 0;\n"
	for _, bd := range model.Mgr.GetAllBDs(){
		config = config + "$i++; \n" +
			"$cfg['Servers'][$i]['verbose'] = '"+ bd.Nombre  +  "';\n" +
			"$cfg['Servers'][$i]['host'] = 'dp-bd-" + bd.Nombre + "';\n" +
		    "$cfg['Servers'][$i]['port'] = '3306';\n" +
		    "$cfg['Servers'][$i]['connect_type'] = 'tcp';\n" +
		    "$cfg['Servers'][$i]['extension'] = 'mysqli';\n" +
			"$cfg['Servers'][$i]['auth_type'] = 'cookie';\n" +
			"$cfg['Servers'][$i]['AllowNoPassword'] = false;\n"
	}

	err := ioutil.WriteFile("configs/bd/phpmyadmin/config.inc.php", []byte(config), 0644)
	check(err)
}

func crearSQL(bd model.BD){
	sql := ""
	var usuario model.UsuarioBD
	abds := model.Mgr.GetAllAsociacionBDs()
	for _ , abd := range abds {
		if abd.BDID == bd.ID && ( abd.Estado == 1 || abd.Estado == 3) {
			usuario = model.Mgr.GetUsuarioBD(strconv.Itoa(int(abd.UsuarioBDID)))
			for _, ip := range bd.IPs {
				if usuario.Estado == 1 && ip.Estado == 1 && abd.Estado == 1 {
					sql = sql + "GRANT ALL ON *.* TO '" + usuario.Nombre + "'@'" + ip.Valor + "' IDENTIFIED BY '" + usuario.Password + "' ;\n"
				} else {
					sql = sql + "REVOKE ALL ON *.* FROM '" + usuario.Nombre + "'@'" + ip.Valor + "' ;\n"
				}
			}
			ip := "172.17.0.0/255.255.255.0"
			if usuario.Estado == 1 && abd.Estado == 1 {
				sql = sql + "GRANT ALL ON *.* TO '" + usuario.Nombre + "'@'" + ip + "' IDENTIFIED BY '" + usuario.Password + "' ;\n"
			} else {
				sql = sql + "REVOKE ALL ON *.* FROM '" + usuario.Nombre + "'@'" + ip + "' ;\n"
			}

			if abd.Estado == 1 {
				abd.Estado = 2
				model.Mgr.UpdateAsociacionBD(&abd)
			}


			if abd.Estado == 3 {
				model.Mgr.RemoveAsociacionBD(abd)
			}

		}

	}

	for _, ubd := range model.Mgr.GetUsuariosDeBD(strconv.Itoa(int(bd.ID))){
		if ubd.Estado == 2 {
			model.Mgr.RemoveUsuarioBD(strconv.Itoa(int(ubd.ID)))
		}
	}

	for _, ip := range bd.IPs {
		if ip.Estado == 2 {
			model.Mgr.RemoveAssociationIP(&bd,&ip)
		}
	}

	sql = sql + "FLUSH PRIVILEGES; \n"
	os.MkdirAll("configs/bd/" + bd.Nombre + "/conf",0755)
	log.Printf(sql)
	err := ioutil.WriteFile("configs/bd/" + bd.Nombre + "/conf/userdata.sql", []byte(sql), 0644)
	check(err)
}


func construirImagenPMA(){
	crearConfigPMA()
	log.Printf("Construyendo la imagen")
	cmdString := "cd configs/bd/phpmyadmin" + "; docker stop dp-bd-phpmyadmin ; docker rm dp-bd-phpmyadmin;  docker rmi dp-img-bd-phpmyadmin; docker build -t dp-img-bd-phpmyadmin" + " ."
	buildCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var out2 bytes.Buffer
	var stderr2 bytes.Buffer
	buildCmd.Stdout = &out2
	buildCmd.Stderr = &stderr2
	err := buildCmd.Run()
	log.Printf(out2.String())
	checkCmd(err,stderr2)
	log.Printf("Imagen creada")
}

func construirImagenBD(){
	cmdString := "docker images -q dp-img-bd"
	imgCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	imgCmd.Stderr = &stderr
	imgCmd.Stdout = &out
	err := imgCmd.Run()
	checkCmd(err,stderr)
	if len(out.String()) == 0 {
		log.Printf("Construyendo la imagen")
		cmdString := "cd configs/bd/" + "; docker build -t dp-img-bd" + " ."
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

func correrContenedorPMA(){
	link := ""
	for _ , bd := range model.Mgr.GetAllBDs(){
		link = link + " --link dp-bd-" + bd.Nombre + ":dp-bd-" + bd.Nombre
	}
	cmdString := "cd configs/bd/phpmyadmin; docker stop dp-bd-phpmyadmin; docker rm dp-bd-phpmyadmin; docker run -d -p 58080:8080 " + link + " --name dp-bd-phpmyadmin dp-img-bd-phpmyadmin"
	runCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	runCmd.Stdout = &out
	runCmd.Stderr = &stderr
	err := runCmd.Run()
	checkCmd(err,stderr)
	log.Printf("Contenedor iniciado dp-bd-phpmyadmin "  + out.String() )
}

func correrContenedorBD(bd model.BD){
	wd, _ := os.Getwd()
	cmdString := "cd configs/bd/" + bd.Nombre + "; docker stop dp-bd-" + bd.Nombre + "; docker rm dp-bd-"+ bd.Nombre +"; docker run -d -p " + calcularPuerto(bd.ID) + ":3306"   +   " -v " + wd + "/configs/bd/" + bd.Nombre + "/conf:/conf:ro"  + " -v " + wd + "/data/bd/" + bd.Nombre + ":/var/lib/mysql" + " --name dp-bd-" + bd.Nombre + " dp-img-bd"
	runCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	runCmd.Stdout = &out
	runCmd.Stderr = &stderr
	err := runCmd.Run()
	checkCmd(err,stderr)
	log.Printf("Contenedor iniciado dp-bd-" + bd.Nombre + " " + out.String() )
}

func ejecutarSQL(bd model.BD){
	cmdString := "docker exec dp-bd-" + bd.Nombre + " /scripts/loadconf.sh"
	execCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	execCmd.Stdout = &out
	execCmd.Stderr = &stderr
	err := execCmd.Run()
	checkCmd(err,stderr)
	log.Printf("SQL ejecutado dp-bd-" + bd.Nombre + " " + out.String() )
}

func borrarBD(bd model.BD){
	cmdString := "docker stop dp-bd-" + bd.Nombre + " ; docker rm dp-bd-" + bd.Nombre
	execCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	execCmd.Stdout = &out
	execCmd.Stderr = &stderr
	err := execCmd.Run()
	checkCmd(err,stderr)
	model.Mgr.RemoveBD(strconv.Itoa(int(bd.ID)))
	log.Printf("BD dp-bd-" + bd.Nombre + " borrada ")
}

func RunBDWorker() {
	log.Printf("Iniciando BD worker")
	// Loop para siempre
	for {
		for _ , bd := range model.Mgr.GetAllBDs() {
			if bd.Estado == 1 {
				log.Printf("Trabajando en la BD " + bd.Nombre )

				if ! isRunning("dp-bd-" + bd.Nombre){
					crearDirectorioConfigBD(bd.Nombre)
				}
				crearSQL(bd)
				if ! isRunning("dp-bd-" + bd.Nombre) {
					construirImagenBD()
					construirImagenPMA()
				} else {
					ejecutarSQL(bd)
				}
				bd.Estado = 2
				model.Mgr.UpdateBD(&bd)

			} else if bd.Estado == 2 || bd.Estado == 4 {
				if ! isRunning("dp-bd-" + bd.Nombre){
					correrContenedorBD(bd)
				}
				if ! isRunning("dp-bd-phpmyadmin"){
					correrContenedorPMA()
				}
				bd.Estado = 3
				model.Mgr.UpdateBD(&bd)
				time.Sleep(10*time.Second)
			} else if bd.Estado == 3 {
				if ! isRunning("dp-bd-" + bd.Nombre){
					bd.Estado = 4
					model.Mgr.UpdateBD(&bd)
				}
			} else if bd.Estado == 5 {
				borrarBD(bd)
			}
		}
		for _ , user := range model.Mgr.GetAllUsuarioBDs() {
			if user.Estado == 2 {
				count := 0
				for _, abd := range model.Mgr.GetAllAsociacionBDs(){
					if ( abd.UsuarioBDID == user.ID ){
						count = count + 1;
					}
				}
				if count == 0 {
					model.Mgr.RemoveUsuarioBD(strconv.Itoa(int(user.ID)))
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
}