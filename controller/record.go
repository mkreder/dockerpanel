package controller

import (
	"net/http"
	"strings"
	"github.com/mkreder/dockerpanel/login"
	"github.com/mkreder/dockerpanel/templates"
	"github.com/mkreder/dockerpanel/model"
	"strconv"
)

func RecordHandler(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		zoneid := r.URL.Query().Get("id")
		records := model.Mgr.GetRecords(zoneid)
		templates.WriteRecordTemplate(w,records,zoneid,"")
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}

func AddRecord(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		r.ParseForm()
		recordtype := strings.Join(r.Form["type"],"")
		name := strings.Join(r.Form["name"],"")
		value := strings.Join(r.Form["value"],"")
		priority := strings.Join(r.Form["priority"],"")
		zoneid := strings.Join(r.Form["zoneid"],"")
		id := strings.Join(r.Form["id"],"")

		if ( len(id) == 0 ) && ( model.Mgr.CheckIfRecordExists(name,recordtype,value,priority,zoneid) ){
			templates.WriteRecordTemplate(w, model.Mgr.GetRecords(zoneid),zoneid,"Record already exists")
		} else {
			var record model.Record
			if len(id) != 0 {
				record = model.Mgr.GetRecord(id)
			}
			record.Type = recordtype
			record.Name = name
			record.Value = value
			record.Priority = priority

			zid, _ := strconv.Atoi(zoneid)
			record.ZoneID = uint(zid)

			if len(id) == 0 {
				err := model.Mgr.AddRecord(&record)
				if err != nil {
					templates.WriteRecordTemplate(w, model.Mgr.GetRecords(zoneid),zoneid,"Error adding record")
				} else {
					http.Redirect(w,r,"/editRecords?id=" + zoneid,http.StatusSeeOther)
				}
			} else {
				err := model.Mgr.UpdateRecord(&record)
				if err != nil {
					templates.WriteRecordTemplate(w, model.Mgr.GetRecords(zoneid),zoneid,"Error updating record")
				} else {
					http.Redirect(w,r,"/editRecords?id=" + zoneid,http.StatusSeeOther)
				}
			}
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}

func RemoveRecord(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		id := r.URL.Query().Get("id")
		zoneid := r.URL.Query().Get("zoneid")
		err := model.Mgr.RemoveRecord(id)
		if err != nil {
			templates.WriteRecordTemplate(w,model.Mgr.GetRecords(zoneid),zoneid,"Error deleting record")
		} else {
			http.Redirect(w,r,"/editRecords?id=" + zoneid,http.StatusSeeOther)
		}
	} else {
		templates.WriteLoginTemplate(w,"")
	}

}