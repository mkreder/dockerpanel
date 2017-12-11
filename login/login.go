package login

import (
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"github.com/gorilla/securecookie"
	"github.com/mkreder/dockerpanel/templates"
	"github.com/mkreder/dockerpanel/model"

)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))


func setSession(UsuarioName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": UsuarioName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func GetUNombreUsuario(request *http.Request) (UsuarioName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			UsuarioName = cookieValue["name"]
		}
	}
	return UsuarioName
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}



func LoginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	pass := r.FormValue("password")
	if len(email) == 0 {
		templates.WriteLoginTemplate(w, "")
	} else {
		Usuario := model.Mgr.GetUsuario(email)
		if err := bcrypt.CompareHashAndPassword([]byte(Usuario.Password), []byte(pass)); err != nil {
			templates.WriteLoginTemplate(w,"Contraseña invalida")
		} else {
			setSession(email, w)
			http.Redirect(w,r,"/",http.StatusSeeOther)
		}

	}
}


func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	templates.WriteLoginTemplate(w,"")
}

