package login

import (
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"github.com/gorilla/securecookie"
	"github.com/mkreder/dockerpanel/templates"
	"github.com/mkreder/dockerpanel/db"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))


func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
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

func GetUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
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
		user := db.Mgr.GetUser(email)
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass)); err != nil {
			templates.WriteLoginTemplate(w,"Contraseña invalida")
		} else {
			setSession(email, w)
			templates.WriteHomeTemplate(w)
		}

	}
}


func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	templates.WriteLoginTemplate(w,"")
}

