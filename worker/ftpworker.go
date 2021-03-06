package worker

import (
	"github.com/mkreder/dockerpanel/model"
	"log"
	"os"
	"os/exec"
	"time"
	"bytes"
	"strconv"
)

func calcularUid(id uint) string{
	idstr := strconv.Itoa(int(id))
	if (id < 10){
		return "100" + idstr
	} else if (id < 100 ){
		return "10" + idstr
	} else if (id < 1000 ){
		return "1" + idstr
	} else {
		return "5000"
	}
}

func crearDirectorioConfigFTP(){
	if _ , err := os.Stat("configs/ftp/Dockerfile"); os.IsNotExist(err){
		os.MkdirAll("configs/ftp",0755)
	}
}

func generarConfigFTP(){
	srcFolder := "defaults/ftp/Dockerfile"
	destFolder := "configs/ftp/Dockerfile"
	cpCmd := exec.Command("cp", "-rf", srcFolder, destFolder)
	err := cpCmd.Run()
	check(err)
	dockerfile, err := os.OpenFile("configs/ftp/Dockerfile", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	conf := ""
	ftpConfig := model.Mgr.GetFtpConfig()
	if ftpConfig.AnonRead == 0 && ftpConfig.AnonWrite == 0 {
		conf = "RUN sed -i \"s/anonymous_enable=YES/anonymous_enable=NO/g\" /etc/vsftpd.conf\n"
	} else {
		if ftpConfig.AnonRead == 0 {
			// Permito escribir pero no leer
			conf = "RUN chmod 2733 /var/lib/ftp\n"
		}
		if ftpConfig.AnonWrite == 1 {
			conf = conf + "RUN sed -i \"s/#anon_upload_enable=YES/anon_upload_enable=YES/g\" /etc/vsftpd.conf\n"
			conf = conf + "RUN sed -i \"s/#anon_mkdir_write_enable=YES/#anon_mkdir_write_enable=YES/g\" /etc/vsftpd.conf\n"
		}
	}
	model.Mgr.UpdateFtpConfig(ftpConfig.AnonWrite,ftpConfig.AnonRead,2)

	for _, user := range model.Mgr.GetAllUsuarioFtps(){
		web := model.Mgr.GetWeb(strconv.Itoa(int(user.WebID)))
		conf = conf + "RUN mkdir -p /data/" + web.Dominio + "\n"
		conf = conf + "RUN useradd -g 1001 -u " + calcularUid(user.ID) +  " -d /data/" + web.Dominio +" " + user.Nombre + "\n"
		conf = conf + "RUN echo " + user.Nombre + ":" + user.Password + " | chpasswd \n"

		wd, _ := os.Getwd()
		destFolder := wd + "/data/web/" + web.Dominio
		uid, err := strconv.Atoi(calcularUid(user.ID))
		os.Chown(destFolder,uid,1001)
		check(err)
	}

	dockerfile.WriteString(conf)
	dockerfile.Sync()
	dockerfile.Close()
}

func buildearContenedorFTP(){
	log.Printf("Construyendo la imagen")
	cmdString := "cd configs/ftp/; docker build -t dp-img-ftp ."
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

func correrContenedorFTP(){
	wd, _ := os.Getwd()
	volume := ""
	for _ , web := range model.Mgr.GetAllWebs(){
		if web.Estado == 3 {
			volume = volume + " -v " + wd + "/data/web/" + web.Dominio + ":/data/" + web.Dominio + " "
		}
	}
	cmdString := "docker stop dp-ftp; docker rm dp-ftp; docker run -d -p 20:20 -p 21:21 -p 21100-21110:21100-21110" + volume + " --name dp-ftp dp-img-ftp"
	tarCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	tarCmd.Stderr = &stderr
	tarCmd.Stdout = &out
	err := tarCmd.Run()
	checkCmd(err,stderr)
	log.Printf("Contenedor iniciado dp-ftp " + out.String() )
}



func removerUsuarioFTP(uftp model.UsuarioFTP){
	model.Mgr.RemoveUsuarioFtp(strconv.Itoa(int(uftp.ID)))
	generarConfigFTP()
	buildearContenedorFTP()
	correrContenedorFTP()
}

func RunFTPWorker (){
	log.Printf("Iniciando FTP worker")
	// Loop para siempre
	for {
		//if model.Mgr.GetFtpConfig().Estado == 1 {
		if model.Mgr.GetFtpConfig().AnonRead == 1 || model.Mgr.GetFtpConfig().AnonWrite == 1 || ( len(model.Mgr.GetAllUsuarioFtps()) > 0 ){
			if model.Mgr.GetFtpConfig().Estado == 1{
				crearDirectorioConfigFTP()
				generarConfigFTP()
				buildearContenedorFTP()
				correrContenedorFTP()
			}
		}
		for _ , uftp := range model.Mgr.GetAllUsuarioFtps() {
			if uftp.Estado == 1 {
				log.Printf("Trabajando en el Usuario " + uftp.Nombre )

				crearDirectorioConfigFTP()
				generarConfigFTP()
				buildearContenedorFTP()
				correrContenedorFTP()

				uftp.Estado = 2
				model.Mgr.UpdateUsuarioFtp(&uftp)
			} else if uftp.Estado == 3 {
				removerUsuarioFTP(uftp)
			}
		}
		if model.Mgr.GetFtpConfig().AnonRead == 1 || model.Mgr.GetFtpConfig().AnonWrite == 1 || ( len(model.Mgr.GetAllUsuarioFtps()) > 0 ) {
			if ! isRunning("dp-ftp") {
				correrContenedorFTP()
			}
		}
		time.Sleep(2 * time.Second)
	}
}