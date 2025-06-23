package model

import (
	"github.com/jinzhu/gorm"
	"strconv"
)

type Zone struct {
	gorm.Model
	Domain string
	Email string
	Status int
	Records []Record
}


func (mgr *manager) AddZone(zone *Zone) (err error) {
	return mgr.db.Create(zone).Error
}

func (mgr *manager) UpdateZone(zone *Zone) (err error) {
	return mgr.db.Save(&zone).Error
}

func (mgr *manager) CheckIfZoneExists(domain string) bool{
	var zone Zone
	exists := mgr.db.First(&zone,"domain = ?",domain).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetAllZones() []Zone {
	var zones []Zone
	mgr.db.Preload("Records").Find(&zones)
	return zones
}

func (mgr *manager) GetZone(id string) Zone {
	var zone Zone
	mgr.db.Preload("Records").First(&zone,id)
	return zone
}

func (mgr *manager) RemoveZone(id string) (err error) {
	for _ , record := range mgr.GetZone(id).Records {
		mgr.RemoveRecord(strconv.Itoa(int(record.ID)))
	}
	return mgr.db.Delete(Zone{}, "id == ?", id).Error
}