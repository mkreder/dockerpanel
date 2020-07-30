package model

import (
	"github.com/jinzhu/gorm"
)

//DB struck to represent a database (DB)
type DB struct {
	gorm.Model
	Name           string
	Status         int
	DBAssociations []DBAssociation
	IPs            []IP
}

func (mgr *manager) AddDB(db *DB) (err error) {
	return mgr.db.Create(db).Error
}

func (mgr *manager) UpdateDB(db *DB) (err error) {
	return mgr.db.Save(&db).Error
}

func (mgr *manager) CheckIfDBExists(name string) bool {
	var db DB
	notexists := mgr.db.First(&db, "name = ?", name).RecordNotFound()
	return !notexists
}

func (mgr *manager) GetAllDBs() []DB {
	var dbs []DB
	mgr.db.Preload("DBAssociations").Preload("IPs").Find(&dbs)
	return dbs
}

func (mgr *manager) GetDB(id string) DB {
	var db DB
	mgr.db.Preload("DBAssociations").Preload("IPs").First(&db, id)
	return db
}

func (mgr *manager) RemoveDB(id string) (err error) {
	return mgr.db.Delete(DB{}, "id == ?", id).Error
}

func (mgr *manager) RemoveAssociationIP(db *DB, ip *IP) (err error) {
	error := mgr.db.Model(&db).Association("IPs").Delete(&ip).Error
	mgr.db.Delete(IP{}, "value == ?", ip.Value)
	return error
}
