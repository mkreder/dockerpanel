package controller

import "net/http"
import (
	"github.com/mkreder/dockerpanel/templates"
	"github.com/mkreder/dockerpanel/model"

	"github.com/mkreder/dockerpanel/login"
	"strings"
)

func MailHandler(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		dominios := model.Mgr.GetAllDominios()
		templates.WriteMailTemplate(w,dominios,"")
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func AddDominio(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		r.ParseForm()
		nombre := strings.Join(r.Form["nombre"],"")
		id := strings.Join(r.Form["id"],"")
		var err error

		if ( len(id) == 0 ) && ( model.Mgr.CheckIfDominioExists(nombre) ){
			templates.WriteMailTemplate(w, model.Mgr.GetAllDominios(),"El dominio de e-mail " + nombre + " ya existe")
		} else {
			var dominio model.Dominio
			if (len(id) != 0) {
				dominio = model.Mgr.GetDominio(id)
			}
			dominio.Nombre = nombre
			dominio.FiltroSpam = strings.Join(r.Form["filtro"],"")

			dominio.Estado = 1
			if len(id) == 0 {
				err = model.Mgr.AddDominio(&dominio)
				if err != nil {
					templates.WriteMailTemplate(w, model.Mgr.GetAllDominios(),"Error al cargar el dominio")
				} else {
					templates.WriteMailTemplate(w, model.Mgr.GetAllDominios(),"")
				}
			} else {
				err = model.Mgr.UpdateDominio(&dominio)
				if err != nil {
					templates.WriteMailTemplate(w, model.Mgr.GetAllDominios(),"Error al actualizar el dominio")
				} else {
					templates.WriteMailTemplate(w, model.Mgr.GetAllDominios(),"")
				}
			}
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}

func RemoveDominio(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		id := r.URL.Query().Get("id")
		err := model.Mgr.RemoveDominio(id)
		if err != nil {
			templates.WriteMailTemplate(w, model.Mgr.GetAllDominios(),"Error al borrar el dominio")
		} else {
			templates.WriteMailTemplate(w, model.Mgr.GetAllDominios(),"")
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}