package controller

import "net/http"
import (
	"github.com/mkreder/dockerpanel/templates"
	"github.com/mkreder/dockerpanel/db"
	"strings"
	"log"
	"github.com/mkreder/dockerpanel/model"

)

func WebHandler(w http.ResponseWriter, r *http.Request) {
	webs := db.Mgr.GetAllWebs()
	templates.WriteWebTemplate(w,webs)
}

func AddWeb(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dominio := strings.Join(r.Form["dominio"],"")
	id := strings.Join(r.Form["id"],"")

	if ( len(id) == 0 ) && ( db.Mgr.CheckIfWebExists(dominio) ){
		log.Printf("El dominio %s ya existe",dominio)
	} else {
		var web model.Web
		if (len(id) != 0) {
			web = db.Mgr.GetWeb(id)
		}
		web.Dominio = dominio

		cgi := strings.Join(r.Form["cgi"],"")
		if len(cgi) != 0 {
			web.CGI = true
		} else {
			web.CGI = false
		}

		ssl := strings.Join(r.Form["ssl"],"")
		if len(ssl) != 0 {
			web.SSL = true
		} else {
			web.SSL = false
		}


		php := strings.Join(r.Form["php"],"")
		if len(php) != 0 {
			web.PHP = true
		} else {
			web.PHP = false
		}

		python := strings.Join(r.Form["python"],"")
		if len(python) != 0 {
			web.Python = true
		} else {
			web.Python = false
		}

		perl := strings.Join(r.Form["perl"],"")
		if len(perl) != 0 {
			web.Perl = true
		} else {
			web.Perl = false
		}

		ruby := strings.Join(r.Form["ruby"],"")
		if len(ruby) != 0 {
			web.Ruby = true
		} else {
			web.Ruby = false
		}

		web.Status = 1
		if len(id) == 0 {
			db.Mgr.AddWeb(&web)
		} else {
			db.Mgr.UpdateWeb(&web)
		}


	}
	templates.WriteWebTemplate(w,db.Mgr.GetAllWebs())


}

func RemoveWeb(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	db.Mgr.RemoveWeb(id)
	templates.WriteWebTemplate(w,db.Mgr.GetAllWebs())
}