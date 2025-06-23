package model

import (
	"github.com/jinzhu/gorm"
)

type Record struct {
	gorm.Model
	Type string
	Name string
	Value string
	Priority string
	ZoneID uint
}


func (mgr *manager) AddRecords(records []Record) {
	var zone Zone
	for  _ , record := range records {
		mgr.db.Create(&record)
	}
	mgr.db.First(&zone,records[0].ZoneID)
	zone.Status = 1
	mgr.db.Save(&zone)
}

func (mgr *manager) GetRecords(zoneid string) []Record{
	var records []Record
	mgr.db.Where("zone_id = ?", zoneid).Find(&records)
	return records
}

func (mgr *manager) CheckIfRecordExists(name string,recordtype string, value string, priority string, zoneid string) bool{
	var record Record
	exists := mgr.db.First(&record,"name = ? AND type = ? AND value = ? AND priority = ? AND zone_id = ?",name,recordtype,value,priority,zoneid).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetRecord(id string) Record {
	var record Record
	mgr.db.First(&record,id)
	return record
}

func (mgr *manager) AddRecord(record *Record) (err error) {
	var zone Zone
	mgr.db.First(&zone,record.ZoneID)
	zone.Status = 1
	mgr.db.Save(&zone)
	return mgr.db.Create(record).Error
}

func (mgr *manager) UpdateRecord(record *Record) (err error) {
	var zone Zone
	mgr.db.First(&zone,record.ZoneID)
	zone.Status = 1
	mgr.db.Save(&zone)
	return mgr.db.Save(&record).Error
}

func (mgr *manager) RemoveRecord(id string) (err error) {
	var zone Zone
	var record Record
	mgr.db.First(&record,id)
	mgr.db.First(&zone,record.ZoneID)
	zone.Status = 1
	mgr.db.Save(&zone)
	return mgr.db.Delete(Record{}, "id == ?", id).Error
}