package model

import (
	"github.com/jinzhu/gorm"
)

type DatabaseAssociation struct {
	gorm.Model
	DatabaseID	uint
	DatabaseUserID uint
	// 1: add 2: ok 3: del
	Status int
}


func (mgr *manager) AddDatabaseAssociation(dba *DatabaseAssociation) (err error) {
	return mgr.db.Create(dba).Error
}

func (mgr *manager) UpdateDatabaseAssociation(dba *DatabaseAssociation) (err error) {
	return mgr.db.Save(&dba).Error
}

func (mgr *manager) CheckIfDatabaseAssociationExists(dbid string, dbuserid string) bool{
	var dba DatabaseAssociation
	exists := mgr.db.First(&dba,"database_id = ? AND database_user_id = ? ",dbid,dbuserid).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetAllDatabaseAssociations() []DatabaseAssociation {
	var dbas []DatabaseAssociation
	mgr.db.Find(&dbas)
	return dbas
}

func (mgr *manager) GetDatabaseAssociation(dbid string, dbuserid string) DatabaseAssociation {
	var dba DatabaseAssociation
	mgr.db.First(&dba,"database_id = ? AND database_user_id = ? ",dbid,dbuserid)
	return dba
}

func (mgr *manager) RemoveDatabaseAssociation(dba DatabaseAssociation) {
	mgr.db.Delete(&dba)
}


