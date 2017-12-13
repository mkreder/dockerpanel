package worker

import (
	"log"
	"bytes"
	"os"
	"github.com/mkreder/dockerpanel/model"
	"os/exec"
	"time"
)

func generarDockerFile(web model.Web){
	dockerfile := ""
	phppkg := ""
	phpdir := ""
	fpmver := ""
	dockerfile = dockerfile + "FROM ubuntu:trusty \nRUN apt-get update\nRUN mkdir /scripts && echo \"#!/bin/sh\" >> /scripts/start.sh && chmod +x /scripts/start.sh \n"
	if ( web.PHP == true ) && ( web.PHPversion != "5.5" ) {
		dockerfile = dockerfile + "RUN apt-get --assume-yes true install software-properties-common\n"+
			"RUN LC_ALL=C.UTF-8 add-apt-repository -y ppa:ondrej/php\n" +
		    "RUN apt-get update \n"
	}
	if  web.PHP == true {
		if web.PHPversion == "5.5" {
			phppkg = "php5"
			phpdir = "php5"
			fpmver = "php5"
		} else if web.PHPversion == "5.6" {
			phppkg = "php5.6"
			phpdir = "php/5.6"
			fpmver = "php\\/php5.6"
		} else if web.PHPversion == "7.0" {
			phppkg = "php7.0"
			phpdir = "php/7.0"
			fpmver = "php\\/php7.0"
		} else if web.PHPversion == "7.1" {
			phppkg = "php7.1"
			phpdir = "php/7.1"
			fpmver = "php\\/php7.1"
		}
		dockerfile = dockerfile + "RUN apt-get --assume-yes true install " + phppkg + " " + phppkg + "-mysql" + "\n"
	}
	if web.Webserver == "nginx"{
		if  web.PHP == true {
			dockerfile = dockerfile + "RUN service apache2 stop && update-rc.d -f apache2 remove && apt-get remove apache\n"
		}
		dockerfile = dockerfile + "RUN apt-get --asume-yes true install nginx \n"
	} else {
		dockerfile = dockerfile + "RUN apt-get --asume-yes true install apache2 \n"
	}
	if web.PHPmethod == "fpm" {
		if web.Webserver == "apache" {
			dockerfile = dockerfile + "RUN apt-get --asume-yes true install " + phppkg + "-fpm libapache2-mod-fastcgi \n"
			dockerfile = dockerfile + "RUN a2enmod actions fastcgi alias \n"
			dockerfile = dockerfile + "RUN a2enmod actions fastcgi alias \n"
			dockerfile = dockerfile + "COPY fpm-vhost.conf /etc/apache2/sites-available/web.conf \n"
			dockerfile = dockerfile + "RUN sed -i \"s/exampe.com/" + web.Dominio + "/g\"  /etc/apache2/sites-available/web.conf \n"
			dockerfile = dockerfile + "RUN sed -i \"s/PHPVER/" + phppkg + "/g\"  /etc/apache2/sites-available/web.conf \n"
		} else {
			dockerfile = dockerfile + "RUN apt-get --asume-yes true install " + phppkg +  "-fpm\n"
			dockerfile = dockerfile + "COPY fpm-nginx.conf /etc/nginx/sites-available/default \n"
			dockerfile = dockerfile + "RUN sed -i \"s/PHPVER/" + fpmver + "/g\"  /etc/nginx/sites-available/default \n"
			dockerfile = dockerfile + "RUN sed -i \"/pathinfo/s/;cgi/cgi/\" -i /etc/" + phpdir + "/fpm/php.ini \n"
		}
		dockerfile = dockerfile + "RUN echo \"/etc/init.d/" + phppkg + "-fpm start\" >> /scripts/start.sh \n"
	} else if web.PHPmethod == "fcgi" {
		dockerfile = dockerfile + "RUN apt-get --asume-yes true install apache2-suexec libapache2-mod-fcgid " + phppkg + "-cgi \n"
		dockerfile = dockerfile + "RUN a2dismod " + phppkg + "\n"
		dockerfile = dockerfile + "RUN a2enmod rewrite && a2enmod suexec && a2enmod include && a2enmod fcgid \n"
		dockerfile = dockerfile + "RUN sed -i \"/pathinfo/s/;cgi/cgi/\" -i /etc/" + phpdir + "/cgi/php.ini \n"
		dockerfile = dockerfile + "RUN echo \"AddHandler fcgid-script .fcgi .php \" > /etc/apache2/mods-enabled/fcgid.conf \n"
		dockerfile = dockerfile + "RUN sed '2 i PHP_Fix_Pathinfo_Enable 1' /etc/apache2/mods-available/fcgid.conf -i \n"
		dockerfile = dockerfile + "RUN groupadd web && useradd -s /bin/false -d /var/www/html -m -g web web \n"
		dockerfile = dockerfile + "RUN mkdir -p /var/www/php-fcgi-scripts/web \n"
		dockerfile = dockerfile + "RUN echo \"#!/bin/sh\" > /var/www/php-fcgi-scripts/web/php-fcgi-starter \n"
		dockerfile = dockerfile + "RUN echo \"PHPRC=/etc/" + phppkg +"/cgi/\" >> /var/www/php-fcgi-scripts/web/php-fcgi-starter\n"
		dockerfile = dockerfile + "RUN echo \"export PHPRC  PHP_FCGI_MAX_REQUESTS=5000 PHP_FCGI_CHILDREN=8 \">> /var/www/php-fcgi-scripts/web/php-fcgi-starter\n"
		dockerfile = dockerfile + "RUN echo \"exec /usr/lib/cgi-bin/php\" >> /var/www/php-fcgi-scripts/web/php-fcgi-starter\n"
		dockerfile = dockerfile + "RUN chmod 755 /var/www/php-fcgi-scripts/web/php-fcgi-starter \n"
		dockerfile = dockerfile + "RUN chown -R web:web /var/www/php-fcgi-scripts/web\n"
		dockerfile = dockerfile + "COPY fcgi-vhost.conf /etc/apache2/sites-available/web.conf \n"
		dockerfile = dockerfile + "RUN sed -i \"s/exampe.com/" + web.Dominio + "/g\"  /etc/apache2/sites-available/web.conf \n"
	} else {
		// no tengo PHP
		if web.Webserver == "apache" {
			dockerfile = dockerfile + "COPY nophp-vhost.conf /etc/apache2/sites-available/web.conf \n"
			dockerfile = dockerfile + "RUN sed -i \"s/exampe.com/" + web.Dominio + "/g\"  /etc/apache2/sites-available/web.conf \n"
		} else {
			dockerfile = dockerfile + "COPY nophp-nginx.conf  /etc/nginx/sites-available/\n"
		}
	}

	if (web.CGI == true ) && (web.Webserver == "apache") {
		dockerfile = dockerfile + "COPY cgi.conf /root/cgi.conf \n"
		dockerfile = dockerfile + "RUN cat /root/cgi.conf >>  /etc/apache2/sites-available/web.conf \n"
		dockerfile = dockerfile + "RUN rm -rf /root/cgi.conf \n"
	}


	if (web.CGI == true ) && (web.Webserver == "nginx") {
		dockerfile = dockerfile + "RUN apt-get --asume-yes true install fcgiwrap\n"
		dockerfile = dockerfile + "RUN echo \"/etc/init.d/fcgiwrap start\" >> /scripts/start.sh \n"
		dockerfile = dockerfile + "COPY cgi-nginx.conf /root/cgi-nginx.conf\n"
		dockerfile = dockerfile + "RUN cat /root/cgi-nginx.conf >> /etc/nginx/sites-available/default \n"
		dockerfile = dockerfile + "RUN rm -rf /root/cgi-nginx.conf \n"
	}

	if web.Webserver == "nginx"{
		dockerfile = dockerfile + "RUN echo \"}\" >> /etc/nginx/sites-available/default \n"
		dockerfile = dockerfile + "RUN sed '2 i daemon off;' -i /etc/nginx/nginx.conf \n"
		dockerfile = dockerfile + "RUN echo \"/usr/sbin/nginx\" >> /scripts/start.sh \n"
	} else {
		dockerfile = dockerfile + "RUN echo \"</VirtualHost>\" >> /etc/apache2/sites-available/web.conf \n"
		dockerfile = dockerfile + "RUN a2ensite web.conf \n"
		dockerfile = dockerfile + "RUN echo \"apachectl -e info -DFOREGROUND\" >> /scripts/start.sh \n"

	}
	dockerfile = dockerfile + "ENTRYPOINT [\"/scripts/start.sh\"]\n"
	dockerfile = dockerfile + "EXPOSE 80 \n"

	f, err := os.Create("configs/web/" + web.Dominio + "/Dockerfile")
	check(err)
	_, _ = f.WriteString (dockerfile)
	f.Sync()
	f.Close()
}


func crearDirectorioConfigWeb(web model.Web) {
	if _, err := os.Stat("configs/web/" + web.Dominio); os.IsNotExist(err) {
		srcFolder := "defaults/web/default"
		os.MkdirAll("configs/web", 0755)
		destFolder := "configs/web/" + web.Dominio
		cpCmd := exec.Command("cp", "-rf", srcFolder, destFolder)
		err := cpCmd.Run()
		check(err)
		os.MkdirAll("data/web/"+web.Dominio, 0755)
		srcFolder = "configs/web/" + web.Dominio + "/info.php"
		destFolder = "data/web/" + web.Dominio + "/"
		cpCmd2 := exec.Command("cp", "-rf", srcFolder, destFolder)
		err = cpCmd2.Run()
		check(err)
	}
	if _, err := os.Stat("configs/web/loadbalancer/conf/" + web.Dominio + ".conf"); os.IsNotExist(err) {
		os.MkdirAll("configs/web/loadbalancer/conf/ssl", 0755)
		srcFolder := "defaults/web/loadbalancer/default.conf"
		destFolder := "configs/web/loadbalancer/conf/"
		cpCmd2 := exec.Command("cp", "-rf", srcFolder, destFolder)
		err = cpCmd2.Run()
		srcFolder = "defaults/web/loadbalancer/Dockerfile"
		destFolder = "configs/web/loadbalancer/"
		cpCmd2 = exec.Command("cp", "-rf", srcFolder, destFolder)
		err = cpCmd2.Run()
	}

}

func generarConfLLB(web model.Web){
	conf := "server { \n listen 80;\n"
	if web.SSL == true {
		conf = conf + " listen 443 ssl;\n"
		f, err := os.Create("configs/web/loadbalancer/conf/ssl/" + web.Dominio + ".pem")
		check(err)
		_, _ = f.WriteString (web.CertSSL)
		f.Sync()
		f.Close()
		conf = conf + " ssl_certificate /etc/nginx/conf.d/ssl/" + web.Dominio + ".pem;\n"
		conf = conf + " ssl_certificate_key /etc/nginx/conf.d/ssl/" + web.Dominio + ".pem;\n"
	}
	conf = conf + " server_name " + web.Dominio + " www." + web.Dominio + ";\n"
	conf = conf + " location / { \n  proxy_pass http://dp-web-" + web.Dominio +";\n }\n}\n"
	f, err := os.Create("configs/web/loadbalancer/conf/" + web.Dominio + ".conf")
	check(err)
	_, _ = f.WriteString (conf)
	f.Sync()
	f.Close()
}

func buildearContenedorWeb(web model.Web){
	cmdString := "docker images -q dp-img-web-" + web.Dominio
	imgCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	imgCmd.Stderr = &stderr
	imgCmd.Stdout = &out
	err := imgCmd.Run()
	checkCmd(err,stderr)
	if len(out.String()) == 0 {
		log.Printf("Construyendo la imagen")
		cmdString := "cd configs/web/" + web.Dominio + "; docker build -t dp-img-web-" + web.Dominio + " ."
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

func buildearContenedorLB(){
	cmdString := "docker images -q dp-img-web-loadbalancer"
	imgCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	imgCmd.Stderr = &stderr
	imgCmd.Stdout = &out
	err := imgCmd.Run()
	checkCmd(err,stderr)
	if len(out.String()) == 0 {
		log.Printf("Construyendo la imagen")
		cmdString := "cd configs/web/loadbalancer; docker build -t dp-img-web-loadbalancer ."
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

func correrContenedorWeb(web model.Web){
	wd, _ := os.Getwd()
	cmdString := "docker stop dp-web-" + web.Dominio + "; docker rm dp-web-" + web.Dominio + "; docker run -d  -v " + wd + "/data/web/" + web.Dominio + ":/var/www/html" + " --name dp-web-" + web.Dominio  + " dp-img-web-" + web.Dominio
	tarCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	tarCmd.Stderr = &stderr
	tarCmd.Stdout = &out
	err := tarCmd.Run()
	checkCmd(err,stderr)
	log.Printf("Contenedor iniciado dp-web" + web.Dominio + " " + out.String() )
}

func correrContenedorLB(){
	wd, _ := os.Getwd()
	link := ""
	for _ , web := range model.Mgr.GetAllWebs(){
		if web.Estado == 3 {
			link = link + " --link dp-web-" + web.Dominio + ":dp-web-" + web.Dominio
		}
	}
	cmdString := "docker stop dp-web-loadbalancer; docker rm dp-web-loadbalancer; docker run -d  -p 80:80 -p 443:443  -v" + wd + "/configs/web/loadbalancer/conf:/etc/nginx/conf.d:ro" + link + " --name dp-web-loadbalancer dp-img-web-loadbalancer"
	tarCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	tarCmd.Stderr = &stderr
	tarCmd.Stdout = &out
	err := tarCmd.Run()
	checkCmd(err,stderr)
	log.Printf("Contenedor (re)iniciado dp-web-loadbalancer " + out.String() )
}

func removerWeb(web model.Web){
	cmdString := "docker stop dp-web-" + web.Dominio +  "; docker rm dp-web-" + web.Dominio +  " ; docker rmi dp-img-web- " + web.Dominio
	rmCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	rmCmd.Stderr = &stderr
	rmCmd.Stdout = &out
	err := rmCmd.Run()
	checkCmd(err,stderr)
	os.RemoveAll("configs/web/loadbalancer/conf/" + web.Dominio + ".conf")
	os.RemoveAll("configs/web/loadbalancer/conf/ssl/" + web.Dominio + ".pem")
}

func RunWebWorker (){
	log.Printf("Iniciando Web worker")
	// Loop para siempre
	for {
		for _ , web := range model.Mgr.GetAllWebs() {
			if web.Estado == 1 {
				log.Printf("Trabajando en la web " + web.Dominio )

				crearDirectorioConfigWeb(web)
				generarConfLLB(web)
				generarDockerFile(web)
				buildearContenedorWeb(web)
				buildearContenedorLB()

				web.Estado = 2
				model.Mgr.UpdateWeb(&web)
			} else if web.Estado == 2 || web.Estado == 4 {
				if ! isRunning("dp-web-" + web.Dominio){
					correrContenedorWeb(web)
				} else {
					reiniciarContenedor("dp-web-" + web.Dominio)
				}
				web.Estado = 3
				model.Mgr.UpdateWeb(&web)
				generarConfLLB(web)
				correrContenedorLB()
			} else if web.Estado == 3 {
				if ! isRunning("dp-web-" + web.Dominio){
					web.Estado = 4
					model.Mgr.UpdateWeb(&web)
				}
			} else if web.Estado == 5 {
				removerWeb(web)
				generarConfLLB(web)
				correrContenedorLB()
			}
			if ! isRunning("dp-web-loadbalancer"){
				correrContenedorLB()
			}

		}
		time.Sleep(2 * time.Second)
	}
}