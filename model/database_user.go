package model

import (
	"github.com/jinzhu/gorm"
)

type DatabaseUser struct {
	gorm.Model
	Name  string
	Password string
	DatabaseAssociations []DatabaseAssociation
	// 1: enabled 2: disabled
	Status int
}

func (mgr *manager) UpdateDatabaseUser(dbu *DatabaseUser) (err error) {
	return mgr.db.Save(&dbu).Error
}


func (mgr *manager) AddDatabaseUser(dbu *DatabaseUser) (err error) {
	return mgr.db.Create(dbu).Error
}


func (mgr *manager) CheckIfDatabaseUserExists(name string) bool{
	var dbu DatabaseUser
	exists := mgr.db.First(&dbu,"name = ?",name).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetAllDatabaseUsers() []DatabaseUser {
	var dbus []DatabaseUser
	mgr.db.Preload("DatabaseAssociations").Find(&dbus)
	return dbus
}

func (mgr *manager) GetDatabaseUser(id string) DatabaseUser {
	var dbu DatabaseUser
	mgr.db.Preload("DatabaseAssociations").First(&dbu,id)
	return dbu
}

func (mgr *manager) RemoveDatabaseUser(id string) (err error) {
	return mgr.db.Delete(DatabaseUser{}, "id == ?", id).Error
}


func (mgr *manager) GetUsersOfDatabase(dbid string) []DatabaseUser{
	var dbus []DatabaseUser
	mgr.db.Joins("left join database_associations on database_associations.database_user_id = database_users.id").Where("database_associations.database_id = ?", dbid).Find(&dbus)
	return dbus
}
