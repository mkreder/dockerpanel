package controller

import (
	"net/http"
	"strings"
	"github.com/mkreder/dockerpanel/login"
	"github.com/mkreder/dockerpanel/db"
	"github.com/mkreder/dockerpanel/templates"
	"golang.org/x/crypto/bcrypt"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		r.ParseForm()
		pass1 := strings.Join(r.Form["pass1"],"")
		if len(pass1) > 0 {
			hash, _ := bcrypt.GenerateFromPassword([]byte(pass1), bcrypt.DefaultCost)
			err := db.Mgr.UpdatePassword(login.GetUserName(r),string(hash))
			if err != nil {
				templates.WriteProfileTemplate(w,"Error al cambiar contraseña")
			} else {
				templates.WriteProfileTemplate(w,"Contraseña guardada")
			}
		} else {
			templates.WriteProfileTemplate(w,"")
		}

	} else {
		templates.WriteLoginTemplate(w,"")
	}

}