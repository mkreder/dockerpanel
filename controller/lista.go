package controller

import "net/http"
import (
	"github.com/mkreder/dockerpanel/templates"
	"strings"
	"github.com/mkreder/dockerpanel/model"

	"github.com/mkreder/dockerpanel/login"
	"strconv"
)

func ListaHandler(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		dominioid := r.URL.Query().Get("dominioid")
		listas := model.Mgr.GetListas(dominioid)
		templates.WriteListaTemplate(w,listas,dominioid,"")
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func AddLista(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		r.ParseForm()
		dominioid := strings.Join(r.Form["dominioid"],"")
		nombre := strings.Join(r.Form["nombre"],"")
		id := strings.Join(r.Form["id"],"")

		if ( len(id) == 0 ) && ( model.Mgr.CheckIfListaExists(nombre,dominioid)){
			templates.WriteListaTemplate(w, model.Mgr.GetListas(dominioid),dominioid,"La lista de correo " + nombre + " ya existe")
		} else {
			var lista model.Lista
			if (len(id) != 0) {
				lista = model.Mgr.GetLista(id)
			}
			emailAdmin := strings.Join(r.Form["emailAdmin"],"")
			lista.Nombre = nombre
			lista.EmailAdmin = emailAdmin
			password := strings.Join(r.Form["password"],"")
			if len(password) > 0 {
				lista.Password = password
			}
			lista.Estado = 1

			d64, err := strconv.ParseUint(dominioid,10,32)

			lista.DominioID = uint(d64)

			dominio := model.Mgr.GetDominio(dominioid)
			dominio.Estado = 1
			model.Mgr.UpdateDominio(&dominio)

			if len(id) == 0 {
				err = model.Mgr.AddLista(&lista)
				if err != nil {
					templates.WriteListaTemplate(w, model.Mgr.GetListas(dominioid),dominioid,"Error al agregar la lista de correo")
				} else {
					templates.WriteListaTemplate(w, model.Mgr.GetListas(dominioid),dominioid,"")
				}
			} else {
				err = model.Mgr.UpdateLista(&lista)
				if err != nil {
					templates.WriteListaTemplate(w, model.Mgr.GetListas(dominioid),dominioid,"Error al modificar la lista de correo")
				} else {
					templates.WriteListaTemplate(w, model.Mgr.GetListas(dominioid),dominioid,"")
				}
			}
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}

func RemoveLista(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		id := r.URL.Query().Get("id")
		dominioid := r.URL.Query().Get("dominioid")
		dominio := model.Mgr.GetDominio(dominioid)
		dominio.Estado = 1
		model.Mgr.UpdateDominio(&dominio)
		lista := model.Mgr.GetLista(id)
		lista.Estado = 3
		err := model.Mgr.UpdateLista(&lista)
		if err != nil {
			templates.WriteListaTemplate(w, model.Mgr.GetListas(dominioid),dominioid,"Error al borrar la lista de correo")
		} else {
			templates.WriteListaTemplate(w, model.Mgr.GetListas(dominioid),dominioid,"")
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}