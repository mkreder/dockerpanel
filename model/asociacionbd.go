package model

import (
	"github.com/jinzhu/gorm"
)

type AsociacionBD struct {
	gorm.Model
	BDID	uint
	UsuarioBDID uint
	// 1: add 2: ok 3: del
	Estado int
}


func (mgr *manager) AddAsociacionBD(adb *AsociacionBD) (err error) {
	return mgr.db.Create(adb).Error
}

func (mgr *manager) UpdateAsociacionBD(adb *AsociacionBD) (err error) {
	return mgr.db.Save(&adb).Error
}

func (mgr *manager) CheckIfAsociacionBDExists(bdid string, usuariobdid string) bool{
	var adb AsociacionBD
	exists := mgr.db.First(&adb,"bd_id = ? AND usuario_bd_id = ? ",bdid,usuariobdid).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetAllAsociacionBDs() []AsociacionBD {
	var abds []AsociacionBD
	mgr.db.Find(&abds)
	return abds
}

func (mgr *manager) GetAsociacionBD(bdid string, usuariobdid string) AsociacionBD {
	var adb AsociacionBD
	mgr.db.First(&adb,"bd_id = ? AND usuario_bd_id = ? ",bdid,usuariobdid)
	return adb
}

func (mgr *manager) RemoveAsociacionBD(adb AsociacionBD) {
	mgr.db.Delete(&adb)
}


