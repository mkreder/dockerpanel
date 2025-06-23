package controller

import (
	"net/http"
	"strings"
	"github.com/mkreder/dockerpanel/login"
	"github.com/mkreder/dockerpanel/templates"
	"github.com/mkreder/dockerpanel/model"
	"github.com/mkreder/dockerpanel/tools"
)

func ZoneHandler(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		zones := model.Mgr.GetAllZones()
		templates.WriteZoneTemplate(w,zones,"")
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func AddZone(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		r.ParseForm()
		domain := strings.Join(r.Form["domain"],"")
		id := strings.Join(r.Form["id"],"")
		var err error

		if ( len(id) == 0 ) && ( model.Mgr.CheckIfZoneExists(domain) ){
			templates.WriteZoneTemplate(w, model.Mgr.GetAllZones(),"Zone " + domain + " already exists")
		} else {
			var zone model.Zone
			if len(id) != 0 {
				zone = model.Mgr.GetZone(id)
			}
			zone.Domain = domain
			zone.Email = strings.Join(r.Form["email"],"")

			if len(id) == 0{
				ip := tools.GetIPAddress()

				var recordA model.Record
				recordA.Type = "A"
				recordA.Name = domain + "."
				recordA.Value = ip

				var recordWww model.Record
				recordWww.Type = "A"
				recordWww.Name = "www"
				recordWww.Value = ip

				var recordMx model.Record
				recordMx.Type = "MX"
				recordMx.Name = domain + "."
				recordMx.Value = "mail." + domain + "."
				recordMx.Priority = "10"

				var recordNs model.Record
				recordNs.Type = "NS"
				recordNs.Name = domain + "."
				recordNs.Value = "ns1." + domain + "."

				var recordANs model.Record
				recordANs.Type = "A"
				recordANs.Name = "ns1"
				recordANs.Value = ip

				var recordAMx model.Record
				recordAMx.Type = "A"
				recordAMx.Name = "mail"
				recordAMx.Value = ip

				zone.Records = append(zone.Records, recordA)
				zone.Records = append(zone.Records, recordWww)
				zone.Records = append(zone.Records, recordMx)
				zone.Records = append(zone.Records, recordNs)
				zone.Records = append(zone.Records, recordAMx)
				zone.Records = append(zone.Records, recordANs)

			}


			zone.Status = 1

			if len(id) == 0 {
				err = model.Mgr.AddZone(&zone)
				if err != nil {
					templates.WriteZoneTemplate(w, model.Mgr.GetAllZones(),"Error adding zone")
				} else {
					http.Redirect(w,r,"/dns",http.StatusSeeOther)
				}
			} else {
				err = model.Mgr.UpdateZone(&zone)
				if err != nil {
					templates.WriteZoneTemplate(w, model.Mgr.GetAllZones(),"Error updating zone")
				} else {
					http.Redirect(w,r,"/dns",http.StatusSeeOther)
				}
			}
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}

func RemoveZone(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		id := r.URL.Query().Get("id")
		zone := model.Mgr.GetZone(id)
		zone.Status = 5
		err := model.Mgr.UpdateZone(&zone)
		if err != nil {
			templates.WriteZoneTemplate(w,model.Mgr.GetAllZones(),"Error deleting zone")
		} else {
			http.Redirect(w,r,"/dns",http.StatusSeeOther)
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}