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
		templates.WriteBDTemplate(w,bds,ubds,"")
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
			templates.WriteBDTemplate(w, model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),"La base de datos " + nombre + " ya existe")
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
					templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),"Error al crear base de datos")
				} else {
					templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),"")
				}
			} else {
				err = model.Mgr.UpdateBD(&bd)
				if err != nil {
					templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),"Error al actualizar base de datos")
				} else {
					templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),"")
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
		err := model.Mgr.RemoveBD(id)
		if err != nil {
			templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),"Error al borrar base de datos")
		} else {
			templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),"")
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func AddUsuarioBD(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		r.ParseForm()
		nombre := strings.Join(r.Form["nombre"],"")
		id := strings.Join(r.Form["id"],"")
		var err error

		if ( len(id) == 0 ) && ( model.Mgr.CheckIfUsuarioBDExists(nombre) ){
			templates.WriteBDTemplate(w, model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),"El usuario " + nombre + " ya existe")
		} else {
			var ubd model.UsuarioBD
			if (len(id) != 0) {
				ubd = model.Mgr.GetUsuarioBD(id)
			}

			ubd.Nombre = nombre
			password := strings.Join(r.Form["password"],"")
			if len(password) > 0 {
				ubd.Password = password
			}

			if len(id) == 0 {
				err = model.Mgr.AddUsuarioBD(&ubd)
				if err != nil {
					templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),"Error al crear usuario de base de datos")
				} else {
					templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),"")
				}
			} else {
				err = model.Mgr.UpdateUsuarioBD(&ubd)
				if err != nil {
					templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),"Error al actualizar el usuario de base de datos")
				} else {
					templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),"")
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
		err := model.Mgr.RemoveUsuarioBD(id)
		if err != nil {
			templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),"Error al borrar usuario de base de datos")
		} else {
			templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),"")
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func AssociateBD(w http.ResponseWriter, r *http.Request){
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		r.ParseForm()
		userid := strings.Join(r.Form["userid"],"")
		bdid := strings.Join(r.Form["bdid"],"")
		user := model.Mgr.GetUsuarioBD(userid)
		newbd := model.Mgr.GetBD(bdid)
		exists := 0
		for _ , bd := range user.BDs {
			if bd.ID == newbd.ID {
				templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),"La base seleccionada ya esta asociada al usuario")
				exists = 1
			}
		}
		if exists != 1 {
			user.BDs = append(user.BDs,newbd)
			newbd.Estado = 1
			model.Mgr.UpdateBD(&newbd)
			err := model.Mgr.UpdateUsuarioBD(&user)
			if err != nil {
				templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),"Error al actualizar el usuario de base de datos")
			} else {
				templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),"")
			}
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}

func DisassociateBD(w http.ResponseWriter, r *http.Request){
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		r.ParseForm()
		userid := strings.Join(r.Form["userid"],"")
		bdidstr := strings.Join(r.Form["bdid"],"")
		bdid, _ := strconv.Atoi(bdidstr)
		bd := model.Mgr.GetBD(bdidstr)
		bd.Estado = 1
		model.Mgr.UpdateBD(&bd)
		user := model.Mgr.GetUsuarioBD(userid)
		for _ , bd := range user.BDs {
			if int(bd.ID) == bdid {
				err := model.Mgr.RemoveAssociationUBD(&user,&bd)
				if err != nil {
					templates.WriteBDTemplate(w, model.Mgr.GetAllBDs(), model.Mgr.GetAllUsuarioBDs(), "Error al borrar la asociacion")
				} else {
					templates.WriteBDTemplate(w, model.Mgr.GetAllBDs(), model.Mgr.GetAllUsuarioBDs(), "")
				}
			}
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
		bd := model.Mgr.GetBD(bdid)
		newip := strings.Join(r.Form["ip"],"")
		exists := 0
		for _ , ip := range bd.IPs {
			if ip.Valor == newip {
				templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),"La IP ingresada ya esta asociada a esta base")
				exists = 1
			}
		}
		if exists != 1 {
			var ip model.IP
			ip.Valor = newip
			bd.Estado = 1
			bd.IPs = append(bd.IPs,ip)
			err := model.Mgr.UpdateBD(&bd)
			if err != nil {
				templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),"Error al agregar IP a la de base de datos")
			} else {
				templates.WriteBDTemplate(w,model.Mgr.GetAllBDs(),model.Mgr.GetAllUsuarioBDs(),"")
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
		rip := strings.Join(r.Form["ip"],"")
		bd := model.Mgr.GetBD(bdid)
		for _ , ip := range bd.IPs {
			if ip.Valor == rip {
				bd.Estado = 1
				model.Mgr.UpdateBD(&bd)
				err := model.Mgr.RemoveAssociationIP(&bd,&ip)
				if err != nil {
					templates.WriteBDTemplate(w, model.Mgr.GetAllBDs(), model.Mgr.GetAllUsuarioBDs(), "Error al borrar la asociación")
				} else {
					templates.WriteBDTemplate(w, model.Mgr.GetAllBDs(), model.Mgr.GetAllUsuarioBDs(), "")
				}
			}
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}