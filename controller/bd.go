package controller

import "net/http"
import (
	"github.com/mkreder/dockerpanel/templates"
	"github.com/mkreder/dockerpanel/model"

	"github.com/mkreder/dockerpanel/login"
	"strings"
	"strconv"
)

func BDHandler(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		bds := model.Mgr.GetAllBDs()
		ubds := model.Mgr.GetAllUsuarioBDs()
		abds := model.Mgr.GetAllAsociacionBDs()
		templates.WriteBDTemplate(w,bds,ubds,abds,"")
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func AddBD(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		r.ParseForm()
		nombre := strings.Join(r.Form["nombreBD"],"")
		id := strings.Join(r.Form["idBD"],"")
		var err error

		if ( len(id) == 0 ) && ( model.Mgr.CheckIfBDExists(nombre) ){
			templates.WriteBDTemplate(w, model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),model.Mgr.GetAllAsociacionBDs(),"La base de datos " + nombre + " ya existe")
		} else {
			var bd model.BD
			if (len(id) != 0) {
				bd = model.Mgr.GetBD(id)
			}

			bd.Nombre = nombre
			bd.Estado = 1
			if len(id) == 0 {
				err = model.Mgr.AddBD(&bd)
				if err != nil {
					templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),model.Mgr.GetAllAsociacionBDs(),"Error al crear base de datos")
				} else {
					templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),model.Mgr.GetAllAsociacionBDs(),"")
				}
			} else {
				err = model.Mgr.UpdateBD(&bd)
				if err != nil {
					templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),model.Mgr.GetAllAsociacionBDs(),"Error al actualizar base de datos")
				} else {
					templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),model.Mgr.GetAllAsociacionBDs(),"")
				}
			}
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}

func RemoveBD(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		id := r.URL.Query().Get("id")
		bd := model.Mgr.GetBD(id)
		bd.Estado = 5
		err := model.Mgr.UpdateBD(&bd)
		if err != nil {
			templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),model.Mgr.GetAllAsociacionBDs(),"Error al borrar base de datos")
		} else {
			http.Redirect(w,r,"/bd",http.StatusSeeOther)
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func AddUsuarioBD(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		r.ParseForm()
		nombre := strings.Join(r.Form["nombre"], "")
		id := strings.Join(r.Form["id"], "")
		var err error

		if len(nombre) == 0 {
			templates.WriteBDTemplate(w, model.Mgr.GetAllBDs(), model.Mgr.GetAllUsuarioBDs(), model.Mgr.GetAllAsociacionBDs(), "")
			return
		}

		if ( len(id) == 0 ) && ( model.Mgr.CheckIfUsuarioBDExists(nombre) ){
			templates.WriteBDTemplate(w, model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),model.Mgr.GetAllAsociacionBDs(),"El usuario " + nombre + " ya existe")
		} else {
			var ubd model.UsuarioBD
			if (len(id) != 0) {
				ubd = model.Mgr.GetUsuarioBD(id)
			}

			ubd.Nombre = nombre
			ubd.Estado = 1
			password := strings.Join(r.Form["password"],"")
			if len(password) > 0 {
				ubd.Password = password
			}

			if len(id) == 0 {
				err = model.Mgr.AddUsuarioBD(&ubd)
				if err != nil {
					templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),model.Mgr.GetAllAsociacionBDs(),"Error al crear usuario de base de datos")
				} else {
					http.Redirect(w,r,"/bd",http.StatusSeeOther)
				}
			} else {
				err = model.Mgr.UpdateUsuarioBD(&ubd)
				if err != nil {
					templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),model.Mgr.GetAllAsociacionBDs(),"Error al actualizar el usuario de base de datos")
				} else {
					http.Redirect(w,r,"/bd",http.StatusSeeOther)
				}
			}
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}

func RemoveUsuarioBD(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		id := r.URL.Query().Get("id")
		if len(id) == 0 {
			templates.WriteBDTemplate(w, model.Mgr.GetAllBDs(), model.Mgr.GetAllUsuarioBDs(), model.Mgr.GetAllAsociacionBDs(), "")
			return
		}
		usuario := model.Mgr.GetUsuarioBD(id)
		usuario.Estado = 2
		err := model.Mgr.UpdateUsuarioBD(&usuario)
		if err != nil {
			templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),model.Mgr.GetAllAsociacionBDs(),"Error al borrar usuario de base de datos")
		} else {
			http.Redirect(w,r,"/bd",http.StatusSeeOther)
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func AssociateBD(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		r.ParseForm()
		userid := strings.Join(r.Form["userid"], "")
		bdid := strings.Join(r.Form["bdid"], "")
		if len(userid) == 0 {
			http.Redirect(w,r,"/bd",http.StatusSeeOther)
			return
		}
		newbd := model.Mgr.GetBD(bdid)
		if model.Mgr.CheckIfAsociacionBDExists(bdid, userid) {
			templates.WriteBDTemplate(w, model.Mgr.GetAllBDs(), model.Mgr.GetAllUsuarioBDs(), model.Mgr.GetAllAsociacionBDs(), "La base seleccionada ya esta asociada al usuario")
		} else {
			var adb model.AsociacionBD
			bdidint, _ := strconv.Atoi(bdid)
			adb.BDID = uint(bdidint)
			useridint, _ := strconv.Atoi(userid)
			adb.UsuarioBDID = uint(useridint)
			adb.Estado = 1
			newbd.Estado = 1
			err := model.Mgr.AddAsociacionBD(&adb)
			if err != nil {
				templates.WriteBDTemplate(w, model.Mgr.GetAllBDs(), model.Mgr.GetAllUsuarioBDs(),model.Mgr.GetAllAsociacionBDs(), "Error al actualizar el usuario de base de datos")
			} else {
				model.Mgr.UpdateBD(&newbd)
				http.Redirect(w,r,"/bd",http.StatusSeeOther)
			}
		}
	} else {
		templates.WriteLoginTemplate(w, "")
	}
}


func DisassociateBD(w http.ResponseWriter, r *http.Request){
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		r.ParseForm()
		userid := strings.Join(r.Form["userid"],"")
		bdid := strings.Join(r.Form["bdid"],"")
		if len(userid) == 0 {
			http.Redirect(w,r,"/bd",http.StatusSeeOther)
			return
		}
		bd := model.Mgr.GetBD(bdid)
		abd := model.Mgr.GetAsociacionBD(bdid,userid)
		abd.Estado = 3
		bd.Estado = 1
		model.Mgr.UpdateBD(&bd)
		err := model.Mgr.UpdateAsociacionBD(&abd)
		if err != nil {
			templates.WriteBDTemplate(w, model.Mgr.GetAllBDs(), model.Mgr.GetAllUsuarioBDs(), model.Mgr.GetAllAsociacionBDs(),"Error al borrar la asociacion")
		} else {
			http.Redirect(w,r,"/bd",http.StatusSeeOther)
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}

func AddIP(w http.ResponseWriter, r *http.Request){
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		r.ParseForm()
		bdid := strings.Join(r.Form["bdid"],"")
		if len(bdid) == 0 {
			http.Redirect(w,r,"/bd",http.StatusSeeOther)
			return
		}
		bd := model.Mgr.GetBD(bdid)
		newip := strings.Join(r.Form["ip"],"")
		exists := 0
		for _ , ip := range bd.IPs {
			if ip.Valor == newip {
				templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),model.Mgr.GetAllAsociacionBDs(),"La IP ingresada ya esta asociada a esta base")
				exists = 1
			}
		}
		if exists != 1 {
			var ip model.IP
			ip.Valor = newip
			ip.Estado = 1
			bd.Estado = 1
			bd.IPs = append(bd.IPs,ip)
			err := model.Mgr.UpdateBD(&bd)
			if err != nil {
				templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),model.Mgr.GetAllAsociacionBDs(),"Error al agregar IP a la de base de datos")
			} else {
				http.Redirect(w,r,"/bd",http.StatusSeeOther)
			}
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}

func RemoveIP(w http.ResponseWriter, r *http.Request){
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		r.ParseForm()
		bdid := strings.Join(r.Form["bdid"],"")
		if len(bdid) == 0 {
			http.Redirect(w,r,"/bd",http.StatusSeeOther)
			return
		}
		rip := strings.Join(r.Form["ip"],"")
		bd := model.Mgr.GetBD(bdid)
		for _ , ip := range bd.IPs {
			if ip.Valor == rip {
				bd.Estado = 1
				ip.Estado = 2
				model.Mgr.UpdateIP(ip,bd)
			}
		}
		http.Redirect(w,r,"/bd",http.StatusSeeOther)
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}