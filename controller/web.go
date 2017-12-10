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
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		webs := model.Mgr.GetAllWebs()
		templates.WriteWebTemplate(w,webs,"")
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func AddWeb(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		r.ParseMultipartForm(32  << 20)
		dominio := strings.Join(r.Form["dominio"],"")
		id := strings.Join(r.Form["id"],"")
		var err error

		if ( len(id) == 0 ) && ( model.Mgr.CheckIfWebExists(dominio) ){
			templates.WriteWebTemplate(w, model.Mgr.GetAllWebs(),"El sitio web " + dominio + " ya existe")
		} else {
			var web model.Web
			if (len(id) != 0) {
				web = model.Mgr.GetWeb(id)
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
				pem, _, _ := r.FormFile("pem")
				if pem != nil {
					buf := bytes.NewBuffer(nil)
					io.Copy(buf, pem)
					pem.Close()
					web.CertSSL = string(buf.Bytes())
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
			web.PHPversion = phpVersion


			phpMethod := strings.Join(r.Form["phpMethod"],"")
			web.PHPmethod = phpMethod


			webserver := strings.Join(r.Form["webserver"],"")
			web.Webserver = webserver




			web.Estado = 1
			if len(id) == 0 {
				err = model.Mgr.AddWeb(&web)
				if err != nil {
					templates.WriteWebTemplate(w, model.Mgr.GetAllWebs(),"Error al agregar el sitio web")
				} else {
					http.Redirect(w,r,"/web",http.StatusSeeOther)
				}
			} else {
				err = model.Mgr.UpdateWeb(&web)
				if err != nil {
					templates.WriteWebTemplate(w, model.Mgr.GetAllWebs(),"Error al actualizar el sitio web")
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
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		id := r.URL.Query().Get("id")
		if len(id) == 0 {
			http.Redirect(w,r,"/web",http.StatusSeeOther)
			return
		}
		err := model.Mgr.RemoveWeb(id)
		if err != nil {
			templates.WriteWebTemplate(w,model.Mgr.GetAllWebs(),"Error al borrar el sitio web")
		} else {
			http.Redirect(w,r,"/web",http.StatusSeeOther)
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}