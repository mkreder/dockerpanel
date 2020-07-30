package model

import (
	"github.com/jinzhu/gorm"
)

//DBAssociation Struct to represent associations between DBs and Users
type DBAssociation struct {
	gorm.Model
	DBID        uint
	UsuarioBDID uint
	// 1: add 2: ok 3: del
	Status int
}

func (mgr *manager) AddDBAssociation(dba *DBAssociation) (err error) {
	return mgr.db.Create(dba).Error
}

func (mgr *manager) UpdateDBAssociation(dba *DBAssociation) (err error) {
	return mgr.db.Save(&dba).Error
}

func (mgr *manager) CheckIfDbAssociationExists(dbid string, dbuserid string) bool {
	var dba DBAssociation
	// TODO: Rewiew this query
	exists := mgr.db.First(&dba, "db_id = ? AND user_db_id = ? ", dbid, dbuserid).RecordNotFound()
	return !exists
}

func (mgr *manager) GetDBAssociations() []DBAssociation {
	var dbas []DBAssociation
	mgr.db.Find(&dbas)
	return dbas
}

func (mgr *manager) GetDBAssociation(dbid string, dbuserid string) DBAssociation {
	var dba DBAssociation
	// TODO: Rewiew this query
	mgr.db.First(&dba, "db_id = ? AND db_user_id = ? ", dbid, dbuserid)
	return dba
}

func (mgr *manager) RemoveDBAssociation(dba DBAssociation) {
	mgr.db.Delete(&dba)
}
