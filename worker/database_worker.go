
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

func calculatePort(id uint) string{
	if ( id < 10 ){
		return "300" + strconv.Itoa(int(id))
	} else if (id < 100){
		return "30" + strconv.Itoa(int(id))
	} else if (id < 1000 ){
		return "3" + strconv.Itoa(int(id))
	} else {
		log.Panic("We don't support more than 1000 databases")
	}
	return ""
}

func createDatabaseConfigDirectory(name string){
	if _ , err := os.Stat("configs/db/Dockerfile"); os.IsNotExist(err){
		log.Printf("Creating configuration files")
		srcFolder := "defaults/db"
		os.MkdirAll("configs/db",0755)
		destFolder := "configs/"
		cpCmd := exec.Command("cp", "-rf", srcFolder, destFolder)
		err := cpCmd.Run()
		check(err)
		os.MkdirAll("data/database/" + name,0755)
	}
}

func createPMAConfig(){
	config := "<?php \n$cfg['blowfish_secret'] = 'abcdkfj2dsfsdfdfasjlkdjlkdsjfkdsjfadslkfajdkfjkldsjfdaflkdsafdx83kx';\n$i = 0;\n"
	for _, database := range model.Mgr.GetAllDatabases(){
		config = config + "$i++; \n" +
			"$cfg['Servers'][$i]['verbose'] = '"+ database.Name  +  "';\n" +
			"$cfg['Servers'][$i]['host'] = 'dp-database-" + database.Name + "';\n" +
		    "$cfg['Servers'][$i]['port'] = '3306';\n" +
		    "$cfg['Servers'][$i]['connect_type'] = 'tcp';\n" +
		    "$cfg['Servers'][$i]['extension'] = 'mysqli';\n" +
			"$cfg['Servers'][$i]['auth_type'] = 'cookie';\n" +
			"$cfg['Servers'][$i]['AllowNoPassword'] = false;\n"
	}

	err := ioutil.WriteFile("configs/db/phpmyadmin/config.inc.php", []byte(config), 0644)
	check(err)
}

func createSQL(db model.Database){
	sql := ""
	var user model.DatabaseUser
	dbas := model.Mgr.GetAllDatabaseAssociations()
	for _ , dba := range dbas {
		if dba.DatabaseID == db.ID && ( dba.Status == 1 || dba.Status == 3) {
			user = model.Mgr.GetDatabaseUser(strconv.Itoa(int(dba.DatabaseUserID)))
			for _, ip := range db.IPs {
				if user.Status == 1 && ip.Status == 1 && dba.Status == 1 {
					sql = sql + "GRANT ALL ON *.* TO '" + user.Name + "'@'" + ip.Value + "' IDENTIFIED BY '" + user.Password + "' ;\n"
				} else {
					sql = sql + "REVOKE ALL ON *.* FROM '" + user.Name + "'@'" + ip.Value + "' ;\n"
				}
			}
			ip := "172.17.0.0/255.255.255.0"
			if user.Status == 1 && dba.Status == 1 {
				sql = sql + "GRANT ALL ON *.* TO '" + user.Name + "'@'" + ip + "' IDENTIFIED BY '" + user.Password + "' ;\n"
			} else {
				sql = sql + "REVOKE ALL ON *.* FROM '" + user.Name + "'@'" + ip + "' ;\n"
			}

			if dba.Status == 1 {
				dba.Status = 2
				model.Mgr.UpdateDatabaseAssociation(&dba)
			}


			if dba.Status == 3 {
				model.Mgr.RemoveDatabaseAssociation(dba)
			}

		}

	}

	for _, ubd := range model.Mgr.GetUsersOfDatabase(strconv.Itoa(int(db.ID))){
		if ubd.Status == 2 {
			model.Mgr.RemoveDatabaseUser(strconv.Itoa(int(ubd.ID)))
		}
	}

	for _, ip := range db.IPs {
		if ip.Status == 2 {
			model.Mgr.RemoveAssociationIP(&db,&ip)
		}
	}

	sql = sql + "FLUSH PRIVILEGES; \n"
	os.MkdirAll("configs/db/" + db.Name + "/conf",0755)
	log.Printf(sql)
	err := ioutil.WriteFile("configs/db/" + db.Name + "/conf/userdata.sql", []byte(sql), 0644)
	check(err)
}


func buildPHPMyAdminImage(){
	createPMAConfig()
	log.Printf("Building the image")
	cmdString := "cd configs/db/phpmyadmin" + "; docker stop dp-db-phpmyadmin ; docker rm dp-db-phpmyadmin;  docker rmi dp-img-db-phpmyadmin; docker build -t dp-img-db-phpmyadmin" + " ."
	buildCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var out2 bytes.Buffer
	var stderr2 bytes.Buffer
	buildCmd.Stdout = &out2
	buildCmd.Stderr = &stderr2
	err := buildCmd.Run()
	log.Printf(out2.String())
	checkCmd(err,stderr2)
	log.Printf("Image created")
}

func buildDatabaseImage(){
	cmdString := "docker images -q dp-img-db"
	imgCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	imgCmd.Stderr = &stderr
	imgCmd.Stdout = &out
	err := imgCmd.Run()
	checkCmd(err,stderr)
	if len(out.String()) == 0 {
		log.Printf("Building the image")
		cmdString := "cd configs/db/" + "; docker build -t dp-img-db" + " ."
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

func runPMAContainer(){
	link := ""
	for _ , db := range model.Mgr.GetAllDatabases(){
		if db.Status == 3 {
			link = link + " --link dp-db-" + db.Name + ":dp-db-" + db.Name
		}
	}
	cmdString := "cd configs/db/phpmyadmin; docker stop dp-db-phpmyadmin; docker rm dp-db-phpmyadmin; docker run -d -p 58080:8080 " + link + " --name dp-db-phpmyadmin dp-img-db-phpmyadmin"
	runCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	runCmd.Stdout = &out
	runCmd.Stderr = &stderr
	err := runCmd.Run()
	checkCmd(err,stderr)
	log.Printf("Container started dp-db-phpmyadmin "  + out.String() )
}

func runDatabaseContainer(db model.Database){
	wd, _ := os.Getwd()
	cmdString := "cd configs/db/" + db.Name + "; docker stop dp-db-" + db.Name + "; docker rm dp-db-"+ db.Name +"; docker run -d -p " + calculatePort(db.ID) + ":3306"   +   " -v " + wd + "/configs/db/" + db.Name + "/conf:/conf:ro"  + " -v " + wd + "/data/db/" + db.Name + ":/var/lib/mysql" + " --name dp-db-" + db.Name + " dp-img-db"
	runCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	runCmd.Stdout = &out
	runCmd.Stderr = &stderr
	err := runCmd.Run()
	checkCmd(err,stderr)
	log.Printf("Container started dp-db-" + db.Name + " " + out.String() )
}

func executeSQL(db model.Database){
	cmdString := "docker exec dp-db-" + db.Name + " /scripts/loadconf.sh"
	execCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	execCmd.Stdout = &out
	execCmd.Stderr = &stderr
	err := execCmd.Run()
	checkCmd(err,stderr)
	log.Printf("SQL ejecutado dp-db-" + db.Name + " " + out.String() )
}

func deleteDB(db model.Database){
	cmdString := "docker stop dp-db-" + db.Name + " ; docker rm dp-db-" + db.Name
	execCmd := exec.Command("/bin/sh" , "-c", cmdString)
	var stderr bytes.Buffer
	var out bytes.Buffer
	execCmd.Stdout = &out
	execCmd.Stderr = &stderr
	err := execCmd.Run()
	checkCmd(err,stderr)
	model.Mgr.RemoveDatabase(strconv.Itoa(int(db.ID)))
	log.Printf("DB dp-db-" + db.Name + " borrada ")
}

func RunDBWorker() {
	log.Printf("Starting Database worker")
	// Loop forever
	for {
		for _ , db := range model.Mgr.GetAllDatabases() {
			if db.Status == 1 {
				log.Printf("Working on DB " + db.Name )

				if ! isRunning("dp-db-" + db.Name){
					createDatabaseConfigDirectory(db.Name)
				}
				createSQL(db)
				if ! isRunning("dp-db-" + db.Name) {
					buildDatabaseImage()
					buildPHPMyAdminImage()
				} else {
					executeSQL(db)
				}
				db.Status = 2
				model.Mgr.UpdateDatabase(&db)

			} else if db.Status == 2 || db.Status == 4 {
				if ! isRunning("dp-db-" + db.Name){
					runDatabaseContainer(db)
				}
				db.Status = 3
				model.Mgr.UpdateDatabase(&db)
				if ! isRunning("dp-db-phpmyadmin"){
					runPMAContainer()
				}
				time.Sleep(10*time.Second)
			} else if db.Status == 3 {
				if ! isRunning("dp-db-" + db.Name){
					db.Status = 4
					model.Mgr.UpdateDatabase(&db)
				}
			} else if db.Status == 5 {
				deleteDB(db)
			}
		}
		for _ , user := range model.Mgr.GetAllDatabaseUsers() {
			if user.Status == 2 {
				count := 0
				for _, dba := range model.Mgr.GetAllDatabaseAssociations(){
					if ( dba.DatabaseUserID == user.ID ){
						count = count + 1;
					}
				}
				if count == 0 {
					model.Mgr.RemoveDatabaseUser(strconv.Itoa(int(user.ID)))
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
}