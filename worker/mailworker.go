package worker

import (
	"github.com/mkreder/dockerpanel/model"
	"os"
	"os/exec"
	"time"

	"bytes"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func crearDirectorioConfigMail() {
	if _, err := os.Stat("configs/mail/postfix/Dockerfile"); os.IsNotExist(err) {
		log.Printf("Creando archivos de configuración")
		os.MkdirAll("configs/mail", 0755)
		srcFolder := "defaults/mail/postfix"
		destFolder := "configs/mail/"
		cpCmd2 := exec.Command("cp", "-rf", srcFolder, destFolder)
		err = cpCmd2.Run()
		check(err)
		srcFolder = "defaults/mail/dovecot"
		cpCmd3 := exec.Command("cp", "-rf", srcFolder, destFolder)
		err = cpCmd3.Run()
		check(err)
		srcFolder = "defaults/mail/mailman"
		cpCmd4 := exec.Command("cp", "-rf", srcFolder, destFolder)
		err = cpCmd4.Run()
		check(err)
		srcFolder = "defaults/mail/roundcube"
		cpCmd5 := exec.Command("cp", "-rf", srcFolder, destFolder)
		err = cpCmd5.Run()
		check(err)

		os.MkdirAll("data/mail/spool", 0755)
		os.MkdirAll("data/mail/mailboxes", 0755)
		os.MkdirAll("data/mail/varmailman", 0755)
		os.MkdirAll("data/mail/libmailman", 0755)
		os.MkdirAll("data/mail/roundcube", 0755)

	}
}

func generarConfig() {
	vdomains := ""
	vmailboxs := ""

	// Remplazar main.cf cada vez
	srcFolder := "defaults/mail/postfix/conf/main.cf"
	destFolder := "configs/mail/postfix/conf/main.cf"
	cpCmd2 := exec.Command("cp", "-rf", srcFolder, destFolder)
	err := cpCmd2.Run()
	check(err)

	virtual := ""
	listdomains := "relay_domains = "
	transport := ""
	// Generar virtual_host y vmailbox
	for _, dominio := range model.Mgr.GetAllDominios() {
		vdomains = vdomains + dominio.Nombre + "\n"
		transport = transport + "listas." + dominio.Nombre + " mailman: \n"
		listdomains = listdomains + "listas." + dominio.Nombre + ","
		for _, cuenta := range dominio.Cuentas {
			vmailboxs = vmailboxs + cuenta.Nombre + "@" + dominio.Nombre + " " + dominio.Nombre + "/" + cuenta.Nombre + "/ \n"
			if len(cuenta.Renvio) > 0 {
				virtual = virtual + cuenta.Nombre + "@" + dominio.Nombre + " " + cuenta.Renvio + "\n"
			}
		}
		if len(dominio.CuentaDefecto) > 0 {
			virtual = virtual + "@" + dominio.Nombre + " " + dominio.CuentaDefecto + "@" + dominio.Nombre + " \n"
		}
	}

	maincf, err := os.OpenFile("configs/mail/postfix/conf/main.cf", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	maincf.WriteString(listdomains + "\n")
	maincf.Sync()
	maincf.Close()

	err = ioutil.WriteFile("configs/mail/postfix/conf/virtual_domains", []byte(vdomains), 0644)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("configs/mail/postfix/conf/virtual", []byte(virtual), 0644)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("configs/mail/postfix/conf/transport", []byte(transport), 0644)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("configs/mail/postfix/conf/vmailbox", []byte(vmailboxs), 0644)
	if err != nil {
		panic(err)
	}

	// Generar passwd y shadow
	for _, dominio := range model.Mgr.GetAllDominios() {
		os.MkdirAll("configs/mail/dovecot/conf/auth/"+dominio.Nombre, 0755)
		passwd := ""
		shadow := ""
		sieve := "require [\"vacation\"]; \n"
		for _, cuenta := range model.Mgr.GetCuentas(strconv.Itoa(int(dominio.ID))) {
			passwd = passwd + cuenta.Nombre + "@" + dominio.Nombre + "::101:104:\"" + cuenta.NombreReal + "\":/var/mail/vhosts/" + dominio.Nombre + "/" + cuenta.Nombre
			if cuenta.Cuota != 0 {
				passwd = passwd + "::userdb_quota_rule=*:bytes=" + strconv.Itoa(cuenta.Cuota) + "M \n"
			} else {
				passwd = passwd + "::\n"
			}
			shadow = shadow + cuenta.Nombre + "@" + dominio.Nombre + ":{DIGEST-MD5}" + cuenta.Password + "\n"
			if cuenta.Autoresponder.Activado == true {
				inicio, _ := time.Parse("2006-01-02", cuenta.Autoresponder.FechaIncio)
				fin, _ := time.Parse("2006-01-02", cuenta.Autoresponder.FechaFin)
				if time.Now().After(inicio) && time.Now().Before(fin) {
					// days define cada cuanto volver a repsonder a las mismas direcciones
					sieve = sieve + "vacation \n:days 1 \n:subject " + cuenta.Autoresponder.Asunto +
						"\n" + ":addresses [\"" + cuenta.Nombre + "@" + dominio.Nombre + "\"]\n" +
						"\"" + strings.Replace(cuenta.Autoresponder.Mensaje, ",", "\n", -1) + "\";\n"
				}
			}
			cuenta.Estado = 2
			model.Mgr.UpdateCuenta(&cuenta)
		}
		err := ioutil.WriteFile("configs/mail/dovecot/conf/auth/"+dominio.Nombre+"/passwd", []byte(passwd), 0644)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile("configs/mail/dovecot/conf/auth/"+dominio.Nombre+"/shadow", []byte(shadow), 0644)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile("configs/mail/dovecot/conf/auth/"+dominio.Nombre+"/autoreply.sieve", []byte(sieve), 0644)
		if err != nil {
			panic(err)
		}

	}

	srcFolder = "defaults/mail/mailman/conf/mm_cfg.py"
	destFolder = "configs/mail/mailman/conf/mm_cfg.py"
	cpCmd3 := exec.Command("cp", "-rf", srcFolder, destFolder)
	err = cpCmd3.Run()
	check(err)

	mmcfg := ""
	mmrun := ""
	for _, dominio := range model.Mgr.GetAllDominios() {
		mmcfg = mmcfg + "POSTFIX_STYLE_VIRTUAL_DOMAINS=['" + dominio.Nombre + "'] \n" +
			"add_virtualhost('" + dominio.Nombre + "', 'listas." + dominio.Nombre + " ')\n"

		for _, lista := range dominio.Listas {
			if lista.Estado == 1 {
				mmrun = mmrun + "newlist --emailhost=listas." + dominio.Nombre + " -a  " + lista.Nombre + " " + lista.EmailAdmin + " " + lista.Password + "\n"
				lista.Estado = 2
				model.Mgr.UpdateLista(&lista)
			} else if lista.Estado == 3 {
				mmrun = mmrun + "rmlist " + lista.Nombre
				model.Mgr.RemoveLista(strconv.Itoa(int(lista.ID)))
			}

		}
	}

	mmcfgf, err := os.OpenFile("configs/mail/mailman/conf/mm_cfg.py", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	mmcfgf.WriteString(mmcfg)
	mmcfgf.Sync()
	mmcfgf.Close()

	mmrunf, err := os.OpenFile("configs/mail/mailman/conf/add_lists.sh", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	mmrunf.WriteString(mmrun)
	mmrunf.Sync()
	mmrunf.Close()


	srcFolder = "defaults/mail/postfix/conf/amavisd.conf"
	destFolder = "configs/mail/postfix/conf/amavisd.conf"
	cpCmd4 := exec.Command("cp", "-rf", srcFolder, destFolder)
	err = cpCmd4.Run()
	check(err)

	amavisconf := "$sa_kill_level_deflt = { \n "

	for _, dominio := range model.Mgr.GetAllDominios() {
		if dominio.FiltroSpam == "alto" {
			amavisconf = amavisconf + "'***@" + dominio.Nombre + "' => '0',\n"
		} else if dominio.FiltroSpam == "medio" {
			amavisconf = amavisconf + "'***@" + dominio.Nombre + "' => '6',\n"
		} else if dominio.FiltroSpam == "bajo" {
			amavisconf = amavisconf + "'***@" + dominio.Nombre + "' => '12',\n"
		} else {
			amavisconf = amavisconf + "'***@" + dominio.Nombre + "' => '10000.0',\n"
		}
	}

	amavisconf = amavisconf + "'.' => '6.3, \n}; \n$sa_tag2_level_defl = $sa_kill_level_deflt; \n"

	amaviscfg, err := os.OpenFile("configs/mail/postfix/conf/amavisd.conf", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	amaviscfg.WriteString(amavisconf)
	amaviscfg.Sync()
	amaviscfg.Close()

}

func prepararImagenRoundcube() {
	cmdString := "docker images -q dp-img-mail-roundcube"
	imgCmd := exec.Command("/bin/sh", "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	imgCmd.Stderr = &stderr
	imgCmd.Stdout = &out
	err := imgCmd.Run()
	checkCmd(err, stderr)
	if len(out.String()) == 0 {
		log.Printf("Construyendo la imagen")
		cmdString := "cd configs/mail/roundcube/; docker build -t dp-img-mail-roundcube ."
		buildCmd := exec.Command("/bin/sh", "-c", cmdString)
		var out2 bytes.Buffer
		var stderr2 bytes.Buffer
		buildCmd.Stdout = &out2
		buildCmd.Stderr = &stderr2
		err := buildCmd.Run()
		log.Printf(out2.String())
		checkCmd(err, stderr2)
		log.Printf("Imagen creada")
	} else {
		log.Printf("La imagen ya existe, salteando paso.")
	}
}


func prepararImagenPostfix() {
	cmdString := "docker images -q dp-img-mail-postfix"
	imgCmd := exec.Command("/bin/sh", "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	imgCmd.Stderr = &stderr
	imgCmd.Stdout = &out
	err := imgCmd.Run()
	checkCmd(err, stderr)
	if len(out.String()) == 0 {
		log.Printf("Construyendo la imagen")
		cmdString := "cd configs/mail/postfix/; docker build -t dp-img-mail-postfix ."
		buildCmd := exec.Command("/bin/sh", "-c", cmdString)
		var out2 bytes.Buffer
		var stderr2 bytes.Buffer
		buildCmd.Stdout = &out2
		buildCmd.Stderr = &stderr2
		err := buildCmd.Run()
		log.Printf(out2.String())
		checkCmd(err, stderr2)
		log.Printf("Imagen creada")
	} else {
		log.Printf("La imagen ya existe, salteando paso.")
	}
}

func prepararImagenMailman() {
	cmdString := "docker images -q dp-img-mail-mailman"
	imgCmd := exec.Command("/bin/sh", "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	imgCmd.Stderr = &stderr
	imgCmd.Stdout = &out
	err := imgCmd.Run()
	checkCmd(err, stderr)
	if len(out.String()) == 0 {
		log.Printf("Construyendo la imagen")
		cmdString := "cd configs/mail/mailman/; docker build -t dp-img-mail-mailman ."
		buildCmd := exec.Command("/bin/sh", "-c", cmdString)
		var out2 bytes.Buffer
		var stderr2 bytes.Buffer
		buildCmd.Stdout = &out2
		buildCmd.Stderr = &stderr2
		err := buildCmd.Run()
		log.Printf(out2.String())
		checkCmd(err, stderr2)
		log.Printf("Imagen creada")
	} else {
		log.Printf("La imagen ya existe, salteando paso.")
	}
}

func prepararImagenDovecot() {
	cmdString := "docker images -q dp-img-mail-dovecot"
	imgCmd := exec.Command("/bin/sh", "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	imgCmd.Stderr = &stderr
	imgCmd.Stdout = &out
	err := imgCmd.Run()
	checkCmd(err, stderr)
	if len(out.String()) == 0 {
		log.Printf("Construyendo la imagen")
		cmdString := "cd configs/mail/dovecot/; docker build -t dp-img-mail-dovecot ."
		buildCmd := exec.Command("/bin/sh", "-c", cmdString)
		var out2 bytes.Buffer
		var stderr2 bytes.Buffer
		buildCmd.Stdout = &out2
		buildCmd.Stderr = &stderr2
		err := buildCmd.Run()
		log.Printf(out2.String())
		checkCmd(err, stderr2)
		log.Printf("Imagen creada")
	} else {
		log.Printf("La imagen ya existe, salteando paso.")
	}
}

func correrContenedorMailman() {
	wd, _ := os.Getwd()
	cmdString := "cd configs/mail/mailman; docker stop dp-mail-mailman; docker rm dp-mail-mailman; docker run -d -p 3880:80 -v " + wd + "/configs/mail/mailman/conf:/etc/mailman:ro -v " + wd + "/data/mail/varmailman:/var/lib/mailman -v" + wd + "/data/mail/libmailman:/usr/lib/mailman" + " --name dp-mail-mailman dp-img-mail-mailman"
	tarCmd := exec.Command("/bin/sh", "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	tarCmd.Stderr = &stderr
	tarCmd.Stdout = &out
	err := tarCmd.Run()
	checkCmd(err, stderr)
	log.Printf("Contenedor iniciado dp-mail-mailman " + out.String())
	time.Sleep(2 * time.Second)
}

func correrContenedorDovecot() {
	wd, _ := os.Getwd()
	cmdString := "cd configs/mail/dovecot; docker stop dp-mail-dovecot; docker rm dp-mail-dovecot; docker run -d -p 110:110 -p 143:143 -v " + wd + "/configs/mail/dovecot/conf:/etc/dovecot:ro -v " + wd + "/data/mail/mailboxes:/var/mail" + " --name dp-mail-dovecot dp-img-mail-dovecot"
	tarCmd := exec.Command("/bin/sh", "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	tarCmd.Stderr = &stderr
	tarCmd.Stdout = &out
	err := tarCmd.Run()
	checkCmd(err, stderr)
	log.Printf("Contenedor iniciado dp-mail-dovecot " + out.String())
	time.Sleep(2 * time.Second)
}

func correrContenedorPostfix() {
	wd, _ := os.Getwd()
	cmdString := "cd configs/mail/postfix; docker stop dp-mail-postfix; docker rm dp-mail-postfix; docker run -d -p 25:25 -v " + wd + "/configs/mail/postfix/conf:/etc/postfix -v " + wd + "/data/mail/mailboxes:/var/mail -v " + wd + "/data/mail/spool:/var/spool/postfix" + " --name dp-mail-postfix --volumes-from dp-mail-mailman --link dp-mail-dovecot:dp-mail-dovecot  dp-img-mail-postfix"
	tarCmd := exec.Command("/bin/sh", "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	tarCmd.Stderr = &stderr
	tarCmd.Stdout = &out
	err := tarCmd.Run()
	checkCmd(err, stderr)
	log.Printf("Contenedor iniciado dp-mail-postfix " + out.String())
}


func correrContenedorRoundcube() {
	wd, _ := os.Getwd()
	cmdString := "cd configs/mail/roundcube; docker stop dp-mail-roundcube; docker rm dp-mail-roudcube; docker run -d -p 9080:80 -v " + wd + "/data/mail/roundcube:/usr/share/webapps/roundcube/config " + " --name dp-mail-roundcube --link dp-mail-dovecot:dp-mail-dovecot --link dp-mail-postfix:dp-mail-postfix  dp-img-mail-roundcube"
	tarCmd := exec.Command("/bin/sh", "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	tarCmd.Stderr = &stderr
	tarCmd.Stdout = &out
	err := tarCmd.Run()
	checkCmd(err, stderr)
	log.Printf("Contenedor iniciado dp-mail-roundcube " + out.String())
}

func RunMailWorker() {
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
				prepararImagenMailman()
				prepararImagenRoundcube()

				dominio.Estado = 2
				model.Mgr.UpdateDominio(&dominio)
			} else if dominio.Estado == 2 || dominio.Estado == 4 {
				if ! isRunning("dp-mail-dovecot") {
					correrContenedorDovecot()
				} else {
					reiniciarContenedor("dp-mail-dovecot")
				}
				if ! isRunning("dp-mail-mailman") {
					correrContenedorMailman()
				} else {
					reiniciarContenedor("dp-mail-mailman")
				}
				if ! isRunning("dp-mail-postfix") {
					correrContenedorPostfix()
				} else {
					reiniciarContenedor("dp-mail-postfix")
				}
				if ! isRunning("dp-mail-roundcube") {
					correrContenedorRoundcube()
				} else {
					reiniciarContenedor("dp-mail-roundcube")
				}
				dominio.Estado = 3
				model.Mgr.UpdateDominio(&dominio)
			} else if dominio.Estado == 3 {
				if ! isRunning("dp-mail-mailman") || ! isRunning("dp-mail-roundcube")|| ! isRunning("dp-mail-dovecot") || ! isRunning("dp-mail-postfix") {
					dominio.Estado = 4
					model.Mgr.UpdateDominio(&dominio)
				}
			} else if dominio.Estado == 5 {

				model.Mgr.RemoveDominio(strconv.Itoa(int(dominio.ID)))
				generarConfig()
				if ! isRunning("dp-mail-dovecot") {
					correrContenedorDovecot()
				} else {
					reiniciarContenedor("dp-mail-dovecot")
				}
				if ! isRunning("dp-mail-mailman") {
					correrContenedorMailman()
				} else {
					reiniciarContenedor("dp-mail-mailman")
				}
				if ! isRunning("dp-mail-postfix") {
					correrContenedorPostfix()
				} else {
					reiniciarContenedor("dp-mail-postfix")
				}


			}

			if (time.Now().Hour() == 0) && (time.Now().Minute() == 0 ) {
				//Reconfiguro cuando sean las 00:00 para configurar los mensajes de vacaciones

				dominio.Estado = 1
				time.Sleep(60 * time.Second)
			}
		}
		time.Sleep(2 * time.Second)
	}

}
