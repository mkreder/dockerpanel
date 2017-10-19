package controller

import "net/http"
import (
	"github.com/mkreder/dockerpanel/templates"
	"github.com/mkreder/dockerpanel/model"
	"strings"
	"github.com/mkreder/dockerpanel/login"
	"strconv"
)

func FtpHandler(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		webs := model.Mgr.GetAllWebs()
		uftps := model.Mgr.GetAllUsuarioFtps()
		templates.WriteFtpTemplate(w,uftps,webs,model.Mgr.GetFtpConfig(),"")
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func AddUsuarioFtp(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		r.ParseForm()
		dominio := strings.Join(r.Form["dominio"],"")
		nombre := strings.Join(r.Form["nombre"],"")
		id := strings.Join(r.Form["id"],"")

		if ( len(id) == 0 ) && ( model.Mgr.CheckIfUsuarioFtpExists(nombre,dominio) ) {
			templates.WriteFtpTemplate(w, model.Mgr.GetAllUsuarioFtps(),model.Mgr.GetAllWebs(),model.Mgr.GetFtpConfig(),"El usuario de FTP " + nombre + "@" + dominio +" ya existe")
		} else {
			var uftp model.UsuarioFTP
			if (len(id) != 0) {
				uftp = model.Mgr.GetUsuarioFtp(id)
			}

			uftp.Nombre = nombre
			z64, err := strconv.ParseUint(dominio,10,32)
			uftp.WebID = uint(z64)
			password := strings.Join(r.Form["password"],"")
			if len(password) > 0 {
				uftp.Password = password
			}

			uftp.Estado = 1
			if len(id) == 0 {
				err = model.Mgr.AddUsuarioFtp(&uftp)
				if err != nil {
					templates.WriteFtpTemplate (w, model.Mgr.GetAllUsuarioFtps(), model.Mgr.GetAllWebs(),model.Mgr.GetFtpConfig(),"Error al agregar el usuario FTP")
				} else {
					templates.WriteFtpTemplate (w, model.Mgr.GetAllUsuarioFtps(), model.Mgr.GetAllWebs(),model.Mgr.GetFtpConfig(),"")
				}
			} else {
				err = model.Mgr.UpdateUsuarioFtp(&uftp)
				if err != nil {
					templates.WriteFtpTemplate (w, model.Mgr.GetAllUsuarioFtps(), model.Mgr.GetAllWebs(),model.Mgr.GetFtpConfig(),"Error al actualizar el usuario FTP")
				} else {
					templates.WriteFtpTemplate (w, model.Mgr.GetAllUsuarioFtps(), model.Mgr.GetAllWebs(),model.Mgr.GetFtpConfig(),"")
				}
			}
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}

func RemoveUsuarioFtp(w http.ResponseWriter, r *http.Request) {
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		id := r.URL.Query().Get("id")
		err := model.Mgr.RemoveUsuarioFtp(id)
		if err != nil {
			templates.WriteFtpTemplate(w,model.Mgr.GetAllUsuarioFtps(),model.Mgr.GetAllWebs(),model.Mgr.GetFtpConfig(),"Error al borrar el usuario FTP")
		} else {
			templates.WriteFtpTemplate(w,model.Mgr.GetAllUsuarioFtps(),model.Mgr.GetAllWebs(),model.Mgr.GetFtpConfig(),"")
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func FtpConfigHandler(w http.ResponseWriter, r *http.Request){
	UsuarioName := login.GetUNombreUsuario(r)
	if UsuarioName != "" {
		r.ParseForm()
		write, _ := strconv.Atoi(strings.Join(r.Form["anonWrite"],""))
		read, _ := strconv.Atoi(strings.Join(r.Form["anonRead"],""))

		err := model.Mgr.UpdateFtpConfig(write,read)
		if err != nil {
			templates.WriteFtpTemplate(w,model.Mgr.GetAllUsuarioFtps(),model.Mgr.GetAllWebs(), model.Mgr.GetFtpConfig(), "Error al guardar configuración FTP")
		} else {
			templates.WriteFtpTemplate(w,model.Mgr.GetAllUsuarioFtps(),model.Mgr.GetAllWebs(), model.Mgr.GetFtpConfig(),"")
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}