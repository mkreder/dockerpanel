package model

import (
	"github.com/jinzhu/gorm"
)

type Database struct {
	gorm.Model
	Name  string
	Status int
	DatabaseAssociations []DatabaseAssociation
	IPs []IP
}


func (mgr *manager) AddDatabase(db *Database) (err error) {
	return mgr.db.Create(db).Error
}

func (mgr *manager) UpdateDatabase(db *Database) (err error) {
	return mgr.db.Save(&db).Error
}

func (mgr *manager) CheckIfDatabaseExists(name string) bool{
	var db Database
	exists := mgr.db.First(&db,"name = ?",name).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetAllDatabases() []Database {
	var dbs []Database
	mgr.db.Preload("DatabaseAssociations").Preload("IPs").Find(&dbs)
	return dbs
}

func (mgr *manager) GetDatabase(id string) Database {
	var db Database
	mgr.db.Preload("DatabaseAssociations").Preload("IPs").First(&db,id)
	return db
}

func (mgr *manager) RemoveDatabase(id string) (err error) {
	return mgr.db.Delete(Database{}, "id == ?", id).Error
}

func (mgr *manager) RemoveAssociationIP(db *Database, ip *IP) (err error){
	error := mgr.db.Model(&db).Association("IPs").Delete(&ip).Error
	mgr.db.Delete(IP{},"value == ?", ip.Value)
	return error
}
