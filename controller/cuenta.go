package controller

import "net/http"
import (
	"github.com/mkreder/dockerpanel/templates"
	"github.com/mkreder/dockerpanel/model"

	"github.com/mkreder/dockerpanel/login"
	"strconv"
	"strings"
	"github.com/mkreder/dockerpanel/tools"
)

func CuentaHandler(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		dominioid := r.URL.Query().Get("dominioid")
		cuentas := model.Mgr.GetCuentas(dominioid)
		dominio := model.Mgr.GetDominio(dominioid)
		templates.WriteCuentaTemplate(w,cuentas,dominio,"")
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func AddCuenta(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		r.ParseForm()
		dominioid := strings.Join(r.Form["dominioid"],"")
		nombre := strings.Join(r.Form["nombre"],"")
		id := strings.Join(r.Form["id"],"")

		if ( len(id) == 0 ) && ( model.Mgr.CheckIfCuentaExists(nombre,dominioid)){
			templates.WriteCuentaTemplate(w, model.Mgr.GetCuentas(dominioid),model.Mgr.GetDominio(dominioid), "La cuenta de correo " + nombre + " ya existe")
		} else {
			var cuenta model.Cuenta
			if (len(id) != 0) {
				cuenta = model.Mgr.GetCuenta(id)
			}
			dominio := model.Mgr.GetDominio(dominioid)

			cuenta.Nombre = nombre
			cuenta.NombreReal = strings.Join(r.Form["nombreReal"],"")
			cuota, _ := strconv.Atoi(strings.Join(r.Form["cuota"],""))
			cuenta.Cuota = cuota

			password := strings.Join(r.Form["password"],"")
			if len(password) > 0 {
				cuenta.Password = tools.GetMD5Hash(nombre + ":" + dominio.Nombre + ":" + password)
			}



			aractivado := strings.Join(r.Form["aractivado"],"")
			if len(aractivado) != 0 {
				if len(id) == 0 {
					var autoresponder model.Autoresponder
					autoresponder.Activado = true
					autoresponder.FechaIncio = strings.Join(r.Form["fechaInicio"],"")
					autoresponder.FechaFin = strings.Join(r.Form["fechaFin"],"")
					mensaje := strings.Replace(strings.Replace(strings.Join(r.Form["mensaje"],""),"\n",";",-1),"\r","",-1)
					autoresponder.Mensaje = mensaje
					autoresponder.Asunto = strings.Join(r.Form["asunto"],"")
					cuenta.Autoresponder = autoresponder
				} else {
					cuenta.Autoresponder.Activado = true
					cuenta.Autoresponder.FechaIncio = strings.Join(r.Form["fechaInicio"],"")
					cuenta.Autoresponder.FechaFin = strings.Join(r.Form["fechaFin"],"")
					mensaje := strings.Replace(strings.Replace(strings.Join(r.Form["mensaje"],""),"\n",";",-1),"\r","",-1)
					cuenta.Autoresponder.Mensaje = mensaje
					cuenta.Autoresponder.Asunto = strings.Join(r.Form["asunto"],"")
				}

			} else {
				if len(id) == 0 {
					var autoresponder model.Autoresponder
					autoresponder.Activado = false
					cuenta.Autoresponder = autoresponder
				} else {
					cuenta.Autoresponder.Activado = false
					cuenta.Autoresponder.FechaFin = ""
					cuenta.Autoresponder.FechaIncio = ""
					cuenta.Autoresponder.Asunto = "Fuera de la oficina"
					cuenta.Autoresponder.Mensaje = ""
				}
			}

			renvioactivo := strings.Join(r.Form["renvioactivo"],"")
			if len(renvioactivo) != 0 {
				cuenta.Renvio = strings.Join(r.Form["direccionRenvio"],"")
			} else {
				cuenta.Renvio = ""
			}
			cuenta.Estado = 1
			dominio.Estado = 1

			cuentadefecto := strings.Join(r.Form["cuentadefecto"],"")
			if len(cuentadefecto) != 0 {
				dominio.CuentaDefecto = cuenta.Nombre
			}

			model.Mgr.UpdateDominio(&dominio)

			d64, err := strconv.ParseUint(dominioid,10,32)
			cuenta.DominioID = uint(d64)

			if len(id) == 0 {
				err = model.Mgr.AddCuenta(&cuenta)
				if err != nil {
					templates.WriteCuentaTemplate(w, model.Mgr.GetCuentas(dominioid),model.Mgr.GetDominio(dominioid), "Error al agregar cuenta de correo")
				} else {
					http.Redirect(w,r,"/editCuentas?dominioid=" + dominioid,http.StatusSeeOther)
				}
			} else {
				err = model.Mgr.UpdateCuenta(&cuenta)
				if err != nil {
					templates.WriteCuentaTemplate(w, model.Mgr.GetCuentas(dominioid),model.Mgr.GetDominio(dominioid), "Error al modificar cuenta de correo")
				} else {
					http.Redirect(w,r,"/editCuentas?dominioid=" + dominioid,http.StatusSeeOther)
				}
			}
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}

func RemoveCuenta(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		id := r.URL.Query().Get("id")
		dominioid := r.URL.Query().Get("dominioid")
		dominio := model.Mgr.GetDominio(dominioid)
		dominio.Estado = 1
		model.Mgr.UpdateDominio(&dominio)
		err := model.Mgr.RemoveCuenta(id)
		if err != nil {
			templates.WriteCuentaTemplate(w, model.Mgr.GetCuentas(dominioid),model.Mgr.GetDominio(dominioid), "Error al borrar cuenta de correo")
		} else {
			http.Redirect(w,r,"/editCuentas?dominioid=" + dominioid,http.StatusSeeOther)
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}