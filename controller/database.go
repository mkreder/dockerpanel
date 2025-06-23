package controller

import "net/http"
import (
	"github.com/mkreder/dockerpanel/templates"
	"github.com/mkreder/dockerpanel/model"

	"github.com/mkreder/dockerpanel/login"
	"strings"
	"strconv"
)

func DatabaseHandler(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		databases := model.Mgr.GetAllDatabases()
		dbUsers := model.Mgr.GetAllDatabaseUsers()
		dbAssociations := model.Mgr.GetAllDatabaseAssociations()
		templates.WriteDatabaseTemplate(w,databases,dbUsers,dbAssociations,"")
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func AddDatabase(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		r.ParseForm()
		name := strings.Join(r.Form["databaseName"],"")
		id := strings.Join(r.Form["databaseID"],"")
		var err error

		if ( len(id) == 0 ) && ( model.Mgr.CheckIfDatabaseExists(name) ){
			templates.WriteDatabaseTemplate(w, model.Mgr.GetAllDatabases(),model.Mgr.GetAllDatabaseUsers(),model.Mgr.GetAllDatabaseAssociations(),"Database " + name + " already exists")
		} else {
			var database model.Database
			if (len(id) != 0) {
				database = model.Mgr.GetDatabase(id)
			}

			database.Name = name
			database.Status = 1
			if len(id) == 0 {
				err = model.Mgr.AddDatabase(&database)
				if err != nil {
					templates.WriteDatabaseTemplate(w,model.Mgr.GetAllDatabases(),model.Mgr.GetAllDatabaseUsers(),model.Mgr.GetAllDatabaseAssociations(),"Error creating database")
				} else {
					templates.WriteDatabaseTemplate(w,model.Mgr.GetAllDatabases(),model.Mgr.GetAllDatabaseUsers(),model.Mgr.GetAllDatabaseAssociations(),"")
				}
			} else {
				err = model.Mgr.UpdateDatabase(&database)
				if err != nil {
					templates.WriteDatabaseTemplate(w,model.Mgr.GetAllDatabases(),model.Mgr.GetAllDatabaseUsers(),model.Mgr.GetAllDatabaseAssociations(),"Error updating database")
				} else {
					templates.WriteDatabaseTemplate(w,model.Mgr.GetAllDatabases(),model.Mgr.GetAllDatabaseUsers(),model.Mgr.GetAllDatabaseAssociations(),"")
				}
			}
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}

func RemoveDatabase(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		id := r.URL.Query().Get("id")
		database := model.Mgr.GetDatabase(id)
		database.Status = 5
		err := model.Mgr.UpdateDatabase(&database)
		if err != nil {
			templates.WriteDatabaseTemplate(w,model.Mgr.GetAllDatabases(),model.Mgr.GetAllDatabaseUsers(),model.Mgr.GetAllDatabaseAssociations(),"Error deleting database")
		} else {
			http.Redirect(w,r,"/db",http.StatusSeeOther)
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func AddDatabaseUser(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		r.ParseForm()
		name := strings.Join(r.Form["name"], "")
		id := strings.Join(r.Form["id"], "")
		var err error

		if len(name) == 0 {
			templates.WriteDatabaseTemplate(w, model.Mgr.GetAllDatabases(), model.Mgr.GetAllDatabaseUsers(), model.Mgr.GetAllDatabaseAssociations(), "")
			return
		}

		if ( len(id) == 0 ) && ( model.Mgr.CheckIfDatabaseUserExists(name) ){
			templates.WriteDatabaseTemplate(w, model.Mgr.GetAllDatabases(),model.Mgr.GetAllDatabaseUsers(),model.Mgr.GetAllDatabaseAssociations(),"User " + name + " already exists")
		} else {
			var dbUser model.DatabaseUser
			if (len(id) != 0) {
				dbUser = model.Mgr.GetDatabaseUser(id)
			}

			dbUser.Name = name
			dbUser.Status = 1
			password := strings.Join(r.Form["password"],"")
			if len(password) > 0 {
				dbUser.Password = password
			}

			if len(id) == 0 {
				err = model.Mgr.AddDatabaseUser(&dbUser)
				if err != nil {
					templates.WriteDatabaseTemplate(w,model.Mgr.GetAllDatabases(),model.Mgr.GetAllDatabaseUsers(),model.Mgr.GetAllDatabaseAssociations(),"Error creating database user")
				} else {
					http.Redirect(w,r,"/db",http.StatusSeeOther)
				}
			} else {
				err = model.Mgr.UpdateDatabaseUser(&dbUser)
				if err != nil {
					templates.WriteDatabaseTemplate(w,model.Mgr.GetAllDatabases(),model.Mgr.GetAllDatabaseUsers(),model.Mgr.GetAllDatabaseAssociations(),"Error updating database user")
				} else {
					http.Redirect(w,r,"/db",http.StatusSeeOther)
				}
			}
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}

func RemoveDatabaseUser(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		id := r.URL.Query().Get("id")
		if len(id) == 0 {
			templates.WriteDatabaseTemplate(w, model.Mgr.GetAllDatabases(), model.Mgr.GetAllDatabaseUsers(), model.Mgr.GetAllDatabaseAssociations(), "")
			return
		}
		user := model.Mgr.GetDatabaseUser(id)
		user.Status = 2
		err := model.Mgr.UpdateDatabaseUser(&user)
		if err != nil {
			templates.WriteDatabaseTemplate(w,model.Mgr.GetAllDatabases(),model.Mgr.GetAllDatabaseUsers(),model.Mgr.GetAllDatabaseAssociations(),"Error deleting database user")
		} else {
			http.Redirect(w,r,"/db",http.StatusSeeOther)
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func AssociateDatabase(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		r.ParseForm()
		userid := strings.Join(r.Form["userid"], "")
		dbid := strings.Join(r.Form["dbid"], "")
		if len(userid) == 0 {
			http.Redirect(w,r,"/db",http.StatusSeeOther)
			return
		}
		newDatabase := model.Mgr.GetDatabase(dbid)
		if model.Mgr.CheckIfDatabaseAssociationExists(dbid, userid) {
			templates.WriteDatabaseTemplate(w, model.Mgr.GetAllDatabases(), model.Mgr.GetAllDatabaseUsers(), model.Mgr.GetAllDatabaseAssociations(), "Selected database is already associated with the user")
		} else {
			var dba model.DatabaseAssociation
			dbidint, _ := strconv.Atoi(dbid)
			dba.DatabaseID = uint(dbidint)
			useridint, _ := strconv.Atoi(userid)
			dba.DatabaseUserID = uint(useridint)
			dba.Status = 1
			newDatabase.Status = 1
			err := model.Mgr.AddDatabaseAssociation(&dba)
			if err != nil {
				templates.WriteDatabaseTemplate(w, model.Mgr.GetAllDatabases(), model.Mgr.GetAllDatabaseUsers(),model.Mgr.GetAllDatabaseAssociations(), "Error updating database user")
			} else {
				model.Mgr.UpdateDatabase(&newDatabase)
				http.Redirect(w,r,"/db",http.StatusSeeOther)
			}
		}
	} else {
		templates.WriteLoginTemplate(w, "")
	}
}


func DisassociateDatabase(w http.ResponseWriter, r *http.Request){
	userName := login.GetUserName(r)
	if userName != "" {
		r.ParseForm()
		userid := strings.Join(r.Form["userid"],"")
		dbid := strings.Join(r.Form["dbid"],"")
		if len(userid) == 0 {
			http.Redirect(w,r,"/db",http.StatusSeeOther)
			return
		}
		database := model.Mgr.GetDatabase(dbid)
		dba := model.Mgr.GetDatabaseAssociation(dbid,userid)
		dba.Status = 3
		database.Status = 1
		model.Mgr.UpdateDatabase(&database)
		err := model.Mgr.UpdateDatabaseAssociation(&dba)
		if err != nil {
			templates.WriteDatabaseTemplate(w, model.Mgr.GetAllDatabases(), model.Mgr.GetAllDatabaseUsers(), model.Mgr.GetAllDatabaseAssociations(),"Error deleting association")
		} else {
			http.Redirect(w,r,"/db",http.StatusSeeOther)
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}

func AddDatabaseIP(w http.ResponseWriter, r *http.Request){
	userName := login.GetUserName(r)
	if userName != "" {
		r.ParseForm()
		dbid := strings.Join(r.Form["dbid"],"")
		if len(dbid) == 0 {
			http.Redirect(w,r,"/db",http.StatusSeeOther)
			return
		}
		database := model.Mgr.GetDatabase(dbid)
		newip := strings.Join(r.Form["ip"],"")
		exists := 0
		for _ , ip := range database.IPs {
			if ip.Value == newip {
				templates.WriteDatabaseTemplate(w,model.Mgr.GetAllDatabases(),model.Mgr.GetAllDatabaseUsers(),model.Mgr.GetAllDatabaseAssociations(),"The entered IP is already associated with this database")
				exists = 1
			}
		}
		if exists != 1 {
			var ip model.IP
			ip.Value = newip
			ip.Status = 1
			database.Status = 1
			database.IPs = append(database.IPs,ip)
			err := model.Mgr.UpdateDatabase(&database)
			if err != nil {
				templates.WriteDatabaseTemplate(w,model.Mgr.GetAllDatabases(),model.Mgr.GetAllDatabaseUsers(),model.Mgr.GetAllDatabaseAssociations(),"Error adding IP to database")
			} else {
				http.Redirect(w,r,"/db",http.StatusSeeOther)
			}
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}

func RemoveDatabaseIP(w http.ResponseWriter, r *http.Request){
	userName := login.GetUserName(r)
	if userName != "" {
		r.ParseForm()
		dbid := strings.Join(r.Form["dbid"],"")
		if len(dbid) == 0 {
			http.Redirect(w,r,"/db",http.StatusSeeOther)
			return
		}
		rip := strings.Join(r.Form["ip"],"")
		database := model.Mgr.GetDatabase(dbid)
		for _ , ip := range database.IPs {
			if ip.Value == rip {
				database.Status = 1
				ip.Status = 2
				model.Mgr.UpdateIP(ip,database)
			}
		}
		http.Redirect(w,r,"/db",http.StatusSeeOther)
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}