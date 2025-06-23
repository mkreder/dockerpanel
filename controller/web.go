package controller

import "net/http"
import (
	"github.com/mkreder/dockerpanel/templates"
	"strings"
	"github.com/mkreder/dockerpanel/model"

	"github.com/mkreder/dockerpanel/login"
	"bytes"
	"io"
)

func WebHandler(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		webs := model.Mgr.GetAllWebs()
		templates.WriteWebTemplate(w,webs,"")
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func AddWeb(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		r.ParseMultipartForm(32  << 20)
		domain := strings.Join(r.Form["domain"],"")
		id := strings.Join(r.Form["id"],"")
		var err error

		if ( len(id) == 0 ) && ( model.Mgr.CheckIfWebExists(domain) ){
			templates.WriteWebTemplate(w, model.Mgr.GetAllWebs(),"Website " + domain + " already exists")
		} else {
			var web model.Web
			if (len(id) != 0) {
				web = model.Mgr.GetWeb(id)
			}
			web.Domain = domain

			cgi := strings.Join(r.Form["cgi"],"")
			if len(cgi) != 0 {
				web.CGI = true
			} else {
				web.CGI = false
			}

			ssl := strings.Join(r.Form["ssl"],"")
			if len(ssl) != 0 {
				web.SSL = true
				pem, _, _ := r.FormFile("pem")
				if pem != nil {
					buf := bytes.NewBuffer(nil)
					io.Copy(buf, pem)
					pem.Close()
					web.SSLCert = string(buf.Bytes())
				}
			} else {
				web.SSL = false
			}

			php := strings.Join(r.Form["php"],"")
			if len(php) != 0 {
				web.PHP = true
			} else {
				web.PHP = false
			}

			phpVersion := strings.Join(r.Form["phpVersion"],"")
			web.PHPVersion = phpVersion


			phpMethod := strings.Join(r.Form["phpMethod"],"")
			web.PHPMethod = phpMethod


			webserver := strings.Join(r.Form["webserver"],"")
			web.WebServer = webserver




			web.Status = 1
			if len(id) == 0 {
				err = model.Mgr.AddWeb(&web)
				if err != nil {
					templates.WriteWebTemplate(w, model.Mgr.GetAllWebs(),"Error adding website")
				} else {
					http.Redirect(w,r,"/web",http.StatusSeeOther)
				}
			} else {
				err = model.Mgr.UpdateWeb(&web)
				if err != nil {
					templates.WriteWebTemplate(w, model.Mgr.GetAllWebs(),"Error updating website")
				} else {
					http.Redirect(w,r,"/web",http.StatusSeeOther)
				}
			}
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}

func RemoveWeb(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		id := r.URL.Query().Get("id")
		if len(id) == 0 {
			http.Redirect(w,r,"/web",http.StatusSeeOther)
			return
		}
		err := model.Mgr.RemoveWeb(id)
		if err != nil {
			templates.WriteWebTemplate(w,model.Mgr.GetAllWebs(),"Error deleting website")
		} else {
			http.Redirect(w,r,"/web",http.StatusSeeOther)
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}