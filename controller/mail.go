package controller

import "net/http"
import (
	"github.com/mkreder/dockerpanel/templates"
	"github.com/mkreder/dockerpanel/model"

	"github.com/mkreder/dockerpanel/login"
	"strings"
)

func MailHandler(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		domains := model.Mgr.GetAllDomains()
		templates.WriteMailTemplate(w,domains,"")
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func AddDomain(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		r.ParseForm()
		name := strings.Join(r.Form["name"],"")
		id := strings.Join(r.Form["id"],"")
		var err error

		if ( len(id) == 0 ) && ( model.Mgr.CheckIfDomainExists(name) ){
			templates.WriteMailTemplate(w, model.Mgr.GetAllDomains(),"Email domain " + name + " already exists")
		} else {
			var domain model.Domain
			if (len(id) != 0) {
				domain = model.Mgr.GetDomain(id)
			}
			domain.Name = name
			domain.SpamFilter = strings.Join(r.Form["filter"],"")

			domain.Status = 1
			if len(id) == 0 {
				err = model.Mgr.AddDomain(&domain)
				if err != nil {
					templates.WriteMailTemplate(w, model.Mgr.GetAllDomains(),"Error loading domain")
				} else {
					http.Redirect(w,r,"/mail",http.StatusSeeOther)
				}
			} else {
				err = model.Mgr.UpdateDomain(&domain)
				if err != nil {
					templates.WriteMailTemplate(w, model.Mgr.GetAllDomains(),"Error updating domain")
				} else {
					templates.WriteMailTemplate(w, model.Mgr.GetAllDomains(),"")
				}
			}
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}

func RemoveDomain(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		id := r.URL.Query().Get("id")
		if len(id) == 0 {
			http.Redirect(w,r,"/mail",http.StatusSeeOther)
			return
		}
		err := model.Mgr.RemoveDomain(id)
		if err != nil {
			templates.WriteMailTemplate(w, model.Mgr.GetAllDomains(),"Error deleting domain")
		} else {
			http.Redirect(w,r,"/mail",http.StatusSeeOther)
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}