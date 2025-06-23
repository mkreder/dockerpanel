package controller

import "net/http"
import (
	"github.com/mkreder/dockerpanel/templates"
	"strings"
	"github.com/mkreder/dockerpanel/model"

	"github.com/mkreder/dockerpanel/login"
	"strconv"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		domainid := r.URL.Query().Get("domainid")
		lists := model.Mgr.GetLists(domainid)
		templates.WriteListTemplate(w,lists,domainid,"")
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func AddList(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		r.ParseForm()
		domainid := strings.Join(r.Form["domainid"],"")
		name := strings.Join(r.Form["name"],"")
		id := strings.Join(r.Form["id"],"")

		if ( len(id) == 0 ) && ( model.Mgr.CheckIfListExists(name,domainid)){
			templates.WriteListTemplate(w, model.Mgr.GetLists(domainid),domainid,"Mailing list " + name + " already exists")
		} else {
			var list model.List
			if (len(id) != 0) {
				list = model.Mgr.GetList(id)
			}
			adminEmail := strings.Join(r.Form["adminEmail"],"")
			list.Name = name
			list.AdminEmail = adminEmail
			password := strings.Join(r.Form["password"],"")
			if len(password) > 0 {
				list.Password = password
			}
			list.Status = 1

			d64, err := strconv.ParseUint(domainid,10,32)

			list.DomainID = uint(d64)

			domain := model.Mgr.GetDomain(domainid)
			domain.Status = 1
			model.Mgr.UpdateDomain(&domain)

			if len(id) == 0 {
				err = model.Mgr.AddList(&list)
				if err != nil {
					templates.WriteListTemplate(w, model.Mgr.GetLists(domainid),domainid,"Error adding mailing list")
				} else {
					http.Redirect(w,r,"/editLists?domainid=" + domainid,http.StatusSeeOther)
				}
			} else {
				err = model.Mgr.UpdateList(&list)
				if err != nil {
					templates.WriteListTemplate(w, model.Mgr.GetLists(domainid),domainid,"Error updating mailing list")
				} else {
					http.Redirect(w,r,"/editLists?domainid=" + domainid,http.StatusSeeOther)
				}
			}
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}

func RemoveList(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		id := r.URL.Query().Get("id")
		domainid := r.URL.Query().Get("domainid")
		domain := model.Mgr.GetDomain(domainid)
		domain.Status = 1
		model.Mgr.UpdateDomain(&domain)
		list := model.Mgr.GetList(id)
		list.Status = 3
		err := model.Mgr.UpdateList(&list)
		if err != nil {
			templates.WriteListTemplate(w, model.Mgr.GetLists(domainid),domainid,"Error deleting mailing list")
		} else {
			http.Redirect(w,r,"/editLists?domainid=" + domainid,http.StatusSeeOther)
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}