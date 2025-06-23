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

func createDirectoryConfigMail() {
	if _, err := os.Stat("configs/mail/postfix/Dockerfile"); os.IsNotExist(err) {
		log.Printf("Creating configuration files")
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
	for _, domain := range model.Mgr.GetAllDomains() {
		vdomains = vdomains + domain.Name + "\n"
		transport = transport + "listas." + domain.Name + " mailman: \n"
		listdomains = listdomains + "listas." + domain.Name + ","
		for _, account := range domain.Accounts {
			vmailboxs = vmailboxs + account.Name + "@" + domain.Name + " " + domain.Name + "/" + account.Name + "/ \n"
			if len(account.Forward) > 0 {
				virtual = virtual + account.Name + "@" + domain.Name + " " + account.Forward + "\n"
			}
		}
		if len(domain.DefaultAccount) > 0 {
			virtual = virtual + "@" + domain.Name + " " + domain.DefaultAccount + "@" + domain.Name + " \n"
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
	for _, domain := range model.Mgr.GetAllDomains() {
		os.MkdirAll("configs/mail/dovecot/conf/auth/"+domain.Name, 0755)
		passwd := ""
		shadow := ""
		sieve := "require [\"vacation\"]; \n"
		for _, account := range model.Mgr.GetAccounts(strconv.Itoa(int(domain.ID))) {
			passwd = passwd + account.Name + "@" + domain.Name + "::101:104:\"" + account.RealName + "\":/var/mail/vhosts/" + domain.Name + "/" + account.Name
			if account.Quota != 0 {
				passwd = passwd + "::userdb_quota_rule=*:bytes=" + strconv.Itoa(account.Quota) + "M \n"
			} else {
				passwd = passwd + "::\n"
			}
			shadow = shadow + account.Name + "@" + domain.Name + ":{DIGEST-MD5}" + account.Password + "\n"
			if account.Autoresponder.Enabled == true {
				inicio, _ := time.Parse("2006-01-02", account.Autoresponder.StartDate)
				fin, _ := time.Parse("2006-01-02", account.Autoresponder.EndDate)
				if time.Now().After(inicio) && time.Now().Before(fin) {
					// days define cada cuanto volver a repsonder a las mismas direcciones
					sieve = sieve + "vacation \n  :days 1 \n  :subject \"" + account.Autoresponder.Subject +
						"\"\n" + "  :addresses [\"" + account.Name + "@" + domain.Name + "\"]\n" +
						"\"" + strings.Replace(account.Autoresponder.Message, ",", "\n", -1) + "\";\n"
				}
				os.MkdirAll("configs/mail/dovecot/conf/auth/" +domain.Name + "/" + account.Name,0755)
				err = ioutil.WriteFile("configs/mail/dovecot/conf/auth/"+domain.Name + "/" + account.Name + "/autoreply.sieve", []byte(sieve), 0644)
				if err != nil {
					panic(err)
				}
			}
	//		account.Status = 2
	//		model.Mgr.UpdateCuenta(&account)
		}
		err := ioutil.WriteFile("configs/mail/dovecot/conf/auth/"+domain.Name+"/passwd", []byte(passwd), 0644)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile("configs/mail/dovecot/conf/auth/"+domain.Name+"/shadow", []byte(shadow), 0644)
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
	for _, domain := range model.Mgr.GetAllDomains() {
		mmcfg = mmcfg + "POSTFIX_STYLE_VIRTUAL_DOMAINS=['" + domain.Name + "'] \n" +
			"add_virtualhost('" + domain.Name + "', 'listas." + domain.Name + " ')\n"

		for _, lista := range domain.Lists {
			if lista.Status == 1 {
				mmrun = mmrun + "newlist --emailhost=listas." + domain.Name + " -a  " + lista.Name + " " + lista.AdminEmail + " " + lista.Password + "\n"
				lista.Status = 2
				model.Mgr.UpdateList(&lista)
			} else if lista.Status == 3 {
				mmrun = mmrun + "rmlist " + lista.Name
				model.Mgr.RemoveList(strconv.Itoa(int(lista.ID)))
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

	for _, domain := range model.Mgr.GetAllDomains() {
		if domain.SpamFilter == "high" {
			amavisconf = amavisconf + "'***@" + domain.Name + "' => '0',\n"
		} else if domain.SpamFilter == "medium" {
			amavisconf = amavisconf + "'***@" + domain.Name + "' => '6',\n"
		} else if domain.SpamFilter == "low" {
			amavisconf = amavisconf + "'***@" + domain.Name + "' => '12',\n"
		} else {
			amavisconf = amavisconf + "'***@" + domain.Name + "' => '10000.0',\n"
		}
	}

	amavisconf = amavisconf + "'.' => '6.3,' \n}; \n \n"

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
		log.Printf("Building the image")
		cmdString := "cd configs/mail/roundcube/; docker build -t dp-img-mail-roundcube ."
		buildCmd := exec.Command("/bin/sh", "-c", cmdString)
		var out2 bytes.Buffer
		var stderr2 bytes.Buffer
		buildCmd.Stdout = &out2
		buildCmd.Stderr = &stderr2
		err := buildCmd.Run()
		log.Printf(out2.String())
		checkCmd(err, stderr2)
		log.Printf("Image created")
	} else {
		log.Printf("Image already exists, skipping step.")
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
		log.Printf("Building the image")
		cmdString := "cd configs/mail/postfix/; docker build -t dp-img-mail-postfix ."
		buildCmd := exec.Command("/bin/sh", "-c", cmdString)
		var out2 bytes.Buffer
		var stderr2 bytes.Buffer
		buildCmd.Stdout = &out2
		buildCmd.Stderr = &stderr2
		err := buildCmd.Run()
		log.Printf(out2.String())
		checkCmd(err, stderr2)
		log.Printf("Image created")
	} else {
		log.Printf("Image already exists, skipping step.")
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
		log.Printf("Building the image")
		cmdString := "cd configs/mail/mailman/; docker build -t dp-img-mail-mailman ."
		buildCmd := exec.Command("/bin/sh", "-c", cmdString)
		var out2 bytes.Buffer
		var stderr2 bytes.Buffer
		buildCmd.Stdout = &out2
		buildCmd.Stderr = &stderr2
		err := buildCmd.Run()
		log.Printf(out2.String())
		checkCmd(err, stderr2)
		log.Printf("Image created")
	} else {
		log.Printf("Image already exists, skipping step.")
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
		log.Printf("Building the image")
		cmdString := "cd configs/mail/dovecot/; docker build -t dp-img-mail-dovecot ."
		buildCmd := exec.Command("/bin/sh", "-c", cmdString)
		var out2 bytes.Buffer
		var stderr2 bytes.Buffer
		buildCmd.Stdout = &out2
		buildCmd.Stderr = &stderr2
		err := buildCmd.Run()
		log.Printf(out2.String())
		checkCmd(err, stderr2)
		log.Printf("Image created")
	} else {
		log.Printf("Image already exists, skipping step.")
	}
}

func runMailmanContainer() {
	wd, _ := os.Getwd()
	cmdString := "cd configs/mail/mailman; docker stop dp-mail-mailman; docker rm dp-mail-mailman; docker run -d -p 3880:80 -v " + wd + "/configs/mail/mailman/conf:/etc/mailman -v " + wd + "/data/mail/varmailman:/var/lib/mailman -v" + wd + "/data/mail/libmailman:/usr/lib/mailman" + " --name dp-mail-mailman dp-img-mail-mailman"
	tarCmd := exec.Command("/bin/sh", "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	tarCmd.Stderr = &stderr
	tarCmd.Stdout = &out
	err := tarCmd.Run()
	checkCmd(err, stderr)
	log.Printf("Container started dp-mail-mailman " + out.String())
	time.Sleep(2 * time.Second)
}

func runDovecotContainer() {
	wd, _ := os.Getwd()
	cmdString := "cd configs/mail/dovecot; docker stop dp-mail-dovecot; docker rm dp-mail-dovecot; docker run -d -p 110:110 -p 143:143 -v " + wd + "/configs/mail/dovecot/conf:/etc/dovecot:ro -v " + wd + "/data/mail/mailboxes:/var/mail" + " --name dp-mail-dovecot dp-img-mail-dovecot"
	tarCmd := exec.Command("/bin/sh", "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	tarCmd.Stderr = &stderr
	tarCmd.Stdout = &out
	err := tarCmd.Run()
	checkCmd(err, stderr)
	log.Printf("Container started dp-mail-dovecot " + out.String())
	time.Sleep(2 * time.Second)
}

func runPostfixContainer() {
	wd, _ := os.Getwd()
	cmdString := "cd configs/mail/postfix; docker stop dp-mail-postfix; docker rm dp-mail-postfix; docker run -d -p 25:25 -v " + wd + "/configs/mail/postfix/conf:/etc/postfix -v " + wd + "/data/mail/mailboxes:/var/mail -v " + wd + "/data/mail/spool:/var/spool/postfix" + " --name dp-mail-postfix --volumes-from dp-mail-mailman --link dp-mail-dovecot:dp-mail-dovecot  dp-img-mail-postfix"
	tarCmd := exec.Command("/bin/sh", "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	tarCmd.Stderr = &stderr
	tarCmd.Stdout = &out
	err := tarCmd.Run()
	checkCmd(err, stderr)
	log.Printf("Container started dp-mail-postfix " + out.String())
}


func runRoundcubeContainer() {
	wd, _ := os.Getwd()
	cmdString := "cd configs/mail/roundcube; docker stop dp-mail-roundcube; docker rm dp-mail-roudcube; docker run -d -p 9080:80 -v " + wd + "/data/mail/roundcube:/usr/share/webapps/roundcube/config " + " --name dp-mail-roundcube --link dp-mail-dovecot:dp-mail-dovecot --link dp-mail-postfix:dp-mail-postfix  dp-img-mail-roundcube"
	tarCmd := exec.Command("/bin/sh", "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	tarCmd.Stderr = &stderr
	tarCmd.Stdout = &out
	err := tarCmd.Run()
	checkCmd(err, stderr)
	log.Printf("Container started dp-mail-roundcube " + out.String())
}

func RunMailWorker() {
	log.Printf("Starting Mail worker")
	// Loop forever
	for {
		for _, domain := range model.Mgr.GetAllDomains() {
			if domain.Status == 1 {
				log.Printf("Working on domain " + domain.Name)

				createDirectoryConfigMail()
				generarConfig()
				prepararImagenDovecot()
				prepararImagenPostfix()
				prepararImagenMailman()
				prepararImagenRoundcube()

				domain.Status = 2
				model.Mgr.UpdateDomain(&domain)
			} else if domain.Status == 2 || domain.Status == 4 {
				if ! isRunning("dp-mail-dovecot") {
					runDovecotContainer()
				} else {
					restartContainer("dp-mail-dovecot")
				}
				if ! isRunning("dp-mail-mailman") {
					runMailmanContainer()
				} else {
					restartContainer("dp-mail-mailman")
				}
				if ! isRunning("dp-mail-postfix") {
					runPostfixContainer()
				} else {
					restartContainer("dp-mail-postfix")
				}
				if ! isRunning("dp-mail-roundcube") {
					runRoundcubeContainer()
				} else {
					restartContainer("dp-mail-roundcube")
				}
				domain.Status = 3
				model.Mgr.UpdateDomain(&domain)
			} else if domain.Status == 3 {
				if ! isRunning("dp-mail-mailman") || ! isRunning("dp-mail-roundcube")|| ! isRunning("dp-mail-dovecot") || ! isRunning("dp-mail-postfix") {
					domain.Status = 4
					model.Mgr.UpdateDomain(&domain)
				}
			} else if domain.Status == 5 {

				model.Mgr.RemoveDomain(strconv.Itoa(int(domain.ID)))
				generarConfig()
				if ! isRunning("dp-mail-dovecot") {
					runDovecotContainer()
				} else {
					restartContainer("dp-mail-dovecot")
				}
				if ! isRunning("dp-mail-mailman") {
					runMailmanContainer()
				} else {
					restartContainer("dp-mail-mailman")
				}
				if ! isRunning("dp-mail-postfix") {
					runPostfixContainer()
				} else {
					restartContainer("dp-mail-postfix")
				}


			}

			if (time.Now().Hour() == 0) && (time.Now().Minute() == 0 ) {
				//Reconfiguro cuando sean las 00:00 para configure los mensajes de vacaciones

				domain.Status = 1
				time.Sleep(60 * time.Second)
			}
		}
		time.Sleep(2 * time.Second)
	}

}
