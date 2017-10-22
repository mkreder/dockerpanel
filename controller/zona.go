package controller

import (
	"net/http"
	"strings"
	"github.com/mkreder/dockerpanel/login"
	"github.com/mkreder/dockerpanel/templates"
	"github.com/mkreder/dockerpanel/model"
	"github.com/mkreder/dockerpanel/tools"
)

func ZonaHandler(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		zonas := model.Mgr.GetAllZonas()
		templates.WriteZonaTemplate(w,zonas,"")
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func AddRegion(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		r.ParseForm()
		dominio := strings.Join(r.Form["dominio"],"")
		id := strings.Join(r.Form["id"],"")
		var err error

		if ( len(id) == 0 ) && ( model.Mgr.CheckIfZonaExists(dominio) ){
			templates.WriteZonaTemplate(w, model.Mgr.GetAllZonas(),"La zona " + dominio + " ya existe")
		} else {
			var zona model.Zona
			if len(id) != 0 {
				zona = model.Mgr.GetZona(id)
			}
			zona.Dominio = dominio
			zona.Email = strings.Join(r.Form["email"],"")

			if len(id) == 0{
				ip := tools.GetIPAddress()

				var registroA model.Registro
				registroA.Tipo = "A"
				registroA.Nombre = dominio + "."
				registroA.Valor = ip
				registroA.Prioridad = "0"

				var registroWww model.Registro
				registroWww.Tipo = "A"
				registroWww.Nombre = "www"
				registroWww.Valor = ip
				registroWww.Prioridad = "0"

				var registroMx model.Registro
				registroMx.Tipo = "MX"
				registroMx.Nombre = dominio + "."
				registroMx.Valor = "mail." + dominio + "."
				registroMx.Prioridad = "10"

				var registroNs model.Registro
				registroNs.Tipo = "NS"
				registroNs.Nombre = dominio + "."
				registroNs.Valor = "ns1." + dominio + "."
				registroNs.Prioridad = "0"

				var registroANs model.Registro
				registroANs.Tipo = "A"
				registroANs.Nombre = "ns1"
				registroANs.Valor = ip
				registroANs.Prioridad = "0"

				var registroAMx model.Registro
				registroAMx.Tipo = "A"
				registroAMx.Nombre = "mail"
				registroAMx.Valor = ip
				registroAMx.Prioridad = "0"

				zona.Registros = append(zona.Registros, registroA)
				zona.Registros = append(zona.Registros, registroWww)
				zona.Registros = append(zona.Registros, registroMx)
				zona.Registros = append(zona.Registros, registroNs)
				zona.Registros = append(zona.Registros, registroAMx)
				zona.Registros = append(zona.Registros, registroANs)

			}


			zona.Estado = 1

			if len(id) == 0 {
				err = model.Mgr.AddZona(&zona)
				if err != nil {
					templates.WriteZonaTemplate(w, model.Mgr.GetAllZonas(),"Error al agregar la zona")
				} else {
					templates.WriteZonaTemplate(w, model.Mgr.GetAllZonas(),"")
				}
			} else {
				err = model.Mgr.UpdateZona(&zona)
				if err != nil {
					templates.WriteZonaTemplate(w, model.Mgr.GetAllZonas(),"Error al actualizar la zona")
				} else {
					templates.WriteZonaTemplate(w, model.Mgr.GetAllZonas(),"")
				}
			}
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}

func RemoveZona(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		id := r.URL.Query().Get("id")
		err := model.Mgr.RemoveZona(id)
		if err != nil {
			templates.WriteZonaTemplate(w,model.Mgr.GetAllZonas(),"Error al borrar la zona")
		} else {
			templates.WriteZonaTemplate(w,model.Mgr.GetAllZonas(),"")
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}