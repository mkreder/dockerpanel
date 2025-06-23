package controller

import (
	"net/http"
	"time"
)
import (
	"github.com/mkreder/dockerpanel/templates"
	"github.com/mkreder/dockerpanel/model"

	"github.com/mkreder/dockerpanel/login"
	"strconv"
	"strings"
	"github.com/mkreder/dockerpanel/tools"
)

func AccountHandler(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		domainid := r.URL.Query().Get("domainid")
		accounts := model.Mgr.GetAccounts(domainid)
		domain := model.Mgr.GetDomain(domainid)
		templates.WriteAccountTemplate(w,accounts,domain,"")
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func updateDomain(domainid string, accountname string, defaultaccount string){
	for {
		domain := model.Mgr.GetDomain(domainid)
		if ( domain.Status == 3 ){
			domain.Status = 1
			if len(defaultaccount) != 0 {
				domain.DefaultAccount = accountname
			}

			model.Mgr.UpdateDomain(&domain)
			break
		}
		time.Sleep(2 * time.Second)
	}
}

func AddAccount(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		r.ParseForm()
		domainid := strings.Join(r.Form["domainid"],"")
		name := strings.Join(r.Form["name"],"")
		id := strings.Join(r.Form["id"],"")

		if ( len(id) == 0 ) && ( model.Mgr.CheckIfAccountExists(name,domainid)){
			templates.WriteAccountTemplate(w, model.Mgr.GetAccounts(domainid),model.Mgr.GetDomain(domainid), "Email account " + name + " already exists")
		} else {
			var account model.Account
			if (len(id) != 0) {
				account = model.Mgr.GetAccount(id)
			}
			domain := model.Mgr.GetDomain(domainid)

			account.Name = name
			account.RealName = strings.Join(r.Form["realName"],"")
			quota, _ := strconv.Atoi(strings.Join(r.Form["quota"],""))
			account.Quota = quota

			password := strings.Join(r.Form["password"],"")
			if len(password) > 0 {
				account.Password = tools.GetMD5Hash(name + ":" + domain.Name + ":" + password)
			}

			arenabled := strings.Join(r.Form["arenabled"],"")
			if len(arenabled) != 0 {
				if len(id) == 0 {
					var autoresponder model.Autoresponder
					autoresponder.Enabled = true
					autoresponder.StartDate = strings.Join(r.Form["startDate"],"")
					autoresponder.EndDate = strings.Join(r.Form["endDate"],"")
					message := strings.Replace(strings.Replace(strings.Join(r.Form["message"],""),"\n",";",-1),"\r","",-1)
					autoresponder.Message = message
					autoresponder.Subject = strings.Join(r.Form["subject"],"")
					account.Autoresponder = autoresponder
				} else {
					account.Autoresponder.Enabled = true
					account.Autoresponder.StartDate = strings.Join(r.Form["startDate"],"")
					account.Autoresponder.EndDate = strings.Join(r.Form["endDate"],"")
					message := strings.Replace(strings.Replace(strings.Join(r.Form["message"],""),"\n",";",-1),"\r","",-1)
					account.Autoresponder.Message = message
					account.Autoresponder.Subject = strings.Join(r.Form["subject"],"")
				}

			} else {
				if len(id) == 0 {
					var autoresponder model.Autoresponder
					autoresponder.Enabled = false
					account.Autoresponder = autoresponder
				} else {
					account.Autoresponder.Enabled = false
					account.Autoresponder.EndDate = ""
					account.Autoresponder.StartDate = ""
					account.Autoresponder.Subject = "Out of office"
					account.Autoresponder.Message = ""
				}
			}

			forwardactive := strings.Join(r.Form["forwardactive"],"")
			if len(forwardactive) != 0 {
				account.Forward = strings.Join(r.Form["forwardAddress"],"")
			} else {
				account.Forward = ""
			}
			account.Status = 1
			defaultaccount := strings.Join(r.Form["defaultaccount"],"")

			go updateDomain(domainid,account.Name,defaultaccount)

			d64, err := strconv.ParseUint(domainid,10,32)
			account.DomainID = uint(d64)

			if len(id) == 0 {
				err = model.Mgr.AddAccount(&account)
				if err != nil {
					templates.WriteAccountTemplate(w, model.Mgr.GetAccounts(domainid),model.Mgr.GetDomain(domainid), "Error adding email account")
				} else {
					http.Redirect(w,r,"/editAccounts?domainid=" + domainid,http.StatusSeeOther)
				}
			} else {
				err = model.Mgr.UpdateAccount(&account)
				if err != nil {
					templates.WriteAccountTemplate(w, model.Mgr.GetAccounts(domainid),model.Mgr.GetDomain(domainid), "Error updating email account")
				} else {
					http.Redirect(w,r,"/editAccounts?domainid=" + domainid,http.StatusSeeOther)
				}
			}
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}

func RemoveAccount(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		id := r.URL.Query().Get("id")
		domainid := r.URL.Query().Get("domainid")
		go updateDomain(domainid,"","")
		err := model.Mgr.RemoveAccount(id)
		if err != nil {
			templates.WriteAccountTemplate(w, model.Mgr.GetAccounts(domainid),model.Mgr.GetDomain(domainid), "Error deleting email account")
		} else {
			http.Redirect(w,r,"/editAccounts?domainid=" + domainid,http.StatusSeeOther)
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}