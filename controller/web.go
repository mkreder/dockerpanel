package controller

import "net/http"
import (
	"github.com/mkreder/dockerpanel/templates"
	"github.com/mkreder/dockerpanel/db"
	"strings"
	"log"
	"github.com/mkreder/dockerpanel/model"
	"strconv"
)

func WebHandler(w http.ResponseWriter, r *http.Request) {
	templates.WriteWebTemplate(w)
}

func AddWeb(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dominio := strings.Join(r.Form["dominio"],"")
	if ( db.Mgr.CheckIfWebExists(dominio) ){
		log.Printf("El dominio %s ya existe",dominio)
	} else {
		log.Printf("El dominio %s no existe ",dominio)
		var err error
		web := model.Web{}
		web.Dominio = dominio

		cgi := strings.Join(r.Form["cgi"],"")
		if len(cgi) != 0 {
			web.CGI, err = strconv.ParseBool(cgi)
			if err != nil {
				log.Fatal(err)
			}
		}

		ssl := strings.Join(r.Form["ssl"],"")
		if len(ssl) != 0 {
			web.SSL, err = strconv.ParseBool(ssl)
			if err != nil {
				log.Fatal(err)
			}
		}

		python := strings.Join(r.Form["python"],"")
		if len(python) != 0 {
			web.Python, err = strconv.ParseBool(python)
			if err != nil {
				log.Fatal(err)
			}
		}

		perl := strings.Join(r.Form["perl"],"")
		if len(perl) != 0 {
			web.Perl, err = strconv.ParseBool(perl)
			if err != nil {
				log.Fatal(err)
			}
		}

		ruby := strings.Join(r.Form["ruby"],"")
		if len(ruby) != 0 {
			web.Ruby, err = strconv.ParseBool(ruby)
			if err != nil {
				log.Fatal(err)
			}
		}

		db.Mgr.AddWeb(&web)

	}
	templates.WriteWebTemplate(w)

}