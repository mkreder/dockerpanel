package controller

import "net/http"
import (
	"github.com/mkreder/dockerpanel/templates"
	"github.com/mkreder/dockerpanel/model"
	"strings"
	"github.com/mkreder/dockerpanel/login"
	"strconv"
)

func FtpHandler(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		webs := model.Mgr.GetAllWebs()
		ftpUsers := model.Mgr.GetAllFTPUsers()
		templates.WriteFtpTemplate(w,ftpUsers,webs,model.Mgr.GetFtpConfig(),"")
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func AddFTPUser(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		r.ParseForm()
		domain := strings.Join(r.Form["domain"],"")
		name := strings.Join(r.Form["name"],"")
		id := strings.Join(r.Form["id"],"")

		if ( len(id) == 0 ) && ( model.Mgr.CheckIfFTPUserExists(name,domain) ) {
			templates.WriteFtpTemplate(w, model.Mgr.GetAllFTPUsers(),model.Mgr.GetAllWebs(),model.Mgr.GetFtpConfig(),"FTP user " + name + "@" + domain +" already exists")
		} else {
			var ftpUser model.FTPUser
			if (len(id) != 0) {
				ftpUser = model.Mgr.GetFTPUser(id)
			}

			ftpUser.Name = name
			z64, err := strconv.ParseUint(domain,10,32)
			ftpUser.WebID = uint(z64)
			password := strings.Join(r.Form["password"],"")
			if len(password) > 0 {
				ftpUser.Password = password
			}

			ftpUser.Status = 1
			if len(id) == 0 {
				err = model.Mgr.AddFTPUser(&ftpUser)
				if err != nil {
					templates.WriteFtpTemplate (w, model.Mgr.GetAllFTPUsers(), model.Mgr.GetAllWebs(),model.Mgr.GetFtpConfig(),"Error adding FTP user")
				} else {
					http.Redirect(w,r,"/ftp",http.StatusSeeOther)
				}
			} else {
				err = model.Mgr.UpdateFTPUser(&ftpUser)
				if err != nil {
					templates.WriteFtpTemplate (w, model.Mgr.GetAllFTPUsers(), model.Mgr.GetAllWebs(),model.Mgr.GetFtpConfig(),"Error updating FTP user")
				} else {
					http.Redirect(w,r,"/ftp",http.StatusSeeOther)
				}
			}
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}

func RemoveFTPUser(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		id := r.URL.Query().Get("id")
		if len(id) == 0 {
			http.Redirect(w,r,"/ftp",http.StatusSeeOther)
			return
		}
		user := model.Mgr.GetFTPUser(id)
		user.Status = 3
		err := model.Mgr.UpdateFTPUser(&user)
		if err != nil {
			templates.WriteFtpTemplate(w,model.Mgr.GetAllFTPUsers(),model.Mgr.GetAllWebs(),model.Mgr.GetFtpConfig(),"Error deleting FTP user")
		} else {
			http.Redirect(w,r,"/ftp",http.StatusSeeOther)
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func FtpConfigHandler(w http.ResponseWriter, r *http.Request){
	userName := login.GetUserName(r)
	if userName != "" {
		r.ParseForm()
		write, _ := strconv.Atoi(strings.Join(r.Form["anonWrite"],""))
		read, _ := strconv.Atoi(strings.Join(r.Form["anonRead"],""))

		err := model.Mgr.UpdateFtpConfig(write,read,1)
		if err != nil {
			templates.WriteFtpTemplate(w,model.Mgr.GetAllFTPUsers(),model.Mgr.GetAllWebs(), model.Mgr.GetFtpConfig(), "Error saving FTP configuration")
		} else {
			http.Redirect(w,r,"/ftp",http.StatusSeeOther)
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}