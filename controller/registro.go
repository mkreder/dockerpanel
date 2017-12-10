package controller

import (
	"net/http"
	"strings"
	"github.com/mkreder/dockerpanel/login"
	"github.com/mkreder/dockerpanel/templates"
	"github.com/mkreder/dockerpanel/model"
	"strconv"
)

func RegistroHandler(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		zonaid := r.URL.Query().Get("id")
		registros := model.Mgr.GetRegistros(zonaid)
		templates.WriteRegistroTemplate(w,registros,zonaid,"")
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func AddRegistro(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		r.ParseForm()
		tipo := strings.Join(r.Form["tipo"],"")
		nombre := strings.Join(r.Form["nombre"],"")
		valor := strings.Join(r.Form["valor"],"")
		prioridad := strings.Join(r.Form["prioridad"],"")
		zonaid := strings.Join(r.Form["zonaid"],"")
		id := strings.Join(r.Form["id"],"")

		if ( len(id) == 0 ) && ( model.Mgr.CheckIfRegistroExists(nombre,tipo,valor,prioridad,zonaid) ){
			templates.WriteRegistroTemplate(w, model.Mgr.GetRegistros(zonaid),zonaid,"El registro ya existe")
		} else {
			var registro model.Registro
			if len(id) != 0 {
				registro = model.Mgr.GetRegistro(id)
			}
			registro.Tipo = tipo
			registro.Nombre = nombre
			registro.Valor = valor
			registro.Prioridad = prioridad

			zid, _ := strconv.Atoi(zonaid)
			registro.ZonaID = uint(zid)

			if len(id) == 0 {
				err := model.Mgr.AddRegistro(&registro)
				if err != nil {
					templates.WriteRegistroTemplate(w, model.Mgr.GetRegistros(zonaid),zonaid,"Error al agregar el registro")
				} else {
					http.Redirect(w,r,"/editRegistros?id=" + zonaid,http.StatusSeeOther)
				}
			} else {
				err := model.Mgr.UpdateRegistro(&registro)
				if err != nil {
					templates.WriteRegistroTemplate(w, model.Mgr.GetRegistros(zonaid),zonaid,"Error al actualizar el registro")
				} else {
					http.Redirect(w,r,"/editRegistros?id=" + zonaid,http.StatusSeeOther)
				}
			}
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}

func RemoveRegistro(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		id := r.URL.Query().Get("id")
		zonaid := r.URL.Query().Get("zonaid")
		err := model.Mgr.RemoveRegistro(id)
		if err != nil {
			templates.WriteRegistroTemplate(w,model.Mgr.GetRegistros(zonaid),zonaid,"Error al borrar el registro")
		} else {
			http.Redirect(w,r,"/editRegistros?id=" + zonaid,http.StatusSeeOther)
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}