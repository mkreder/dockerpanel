package model

import (
	"github.com/jinzhu/gorm"
)

//DBUser Username for DBs
type DBUser struct {
	gorm.Model
	Name           string
	Password       string
	DBAssociations []DBAssociation
	// 1: activado 2: desactivado
	Status int
}

func (mgr *manager) UpdateDBUser(dbu *DBUser) (err error) {
	return mgr.db.Save(&dbu).Error
}

func (mgr *manager) AddDBUser(dbu *DBUser) (err error) {
	return mgr.db.Create(dbu).Error
}

func (mgr *manager) CheckIfDBUserExists(name string) bool {
	var uftp DBUser
	exists := mgr.db.First(&uftp, "name = ?", name).RecordNotFound()
	return !exists
}

func (mgr *manager) GetAllDBUsers() []DBUser {
	var uftps []DBUser
	mgr.db.Preload("DBAssociations").Find(&uftps)
	return uftps
}

func (mgr *manager) GetDBUser(id string) DBUser {
	var uftp DBUser
	mgr.db.Preload("DBAssociations").First(&uftp, id)
	return uftp
}

func (mgr *manager) RemoveDBUser(id string) (err error) {
	return mgr.db.Delete(DBUser{}, "id == ?", id).Error
}

func (mgr *manager) GetDBUserbyDB(dbid string) []DBUser {
	var ubds []DBUser
	//TODO: Review this
	mgr.db.Joins("left join asociacion_bds on asociacion_bds.usuario_bd_id = usuario_bds.id").Where("asociacion_bds.bd_id = ?", dbid).Find(&ubds)
	return ubds
}
