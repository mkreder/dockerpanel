package model

import (
	"github.com/jinzhu/gorm"
)

type UsuarioBD struct {
	gorm.Model
	Nombre  string
	Password string
	Estado int
	BDs []BD `gorm:"many2many:usuariobd_bds;"`
}

func (mgr *manager) UpdateUsuarioBD(uftp *UsuarioBD) (err error) {
	return mgr.db.Save(&uftp).Error
}


func (mgr *manager) AddUsuarioBD(uftp *UsuarioBD) (err error) {
	return mgr.db.Create(uftp).Error
}


func (mgr *manager) CheckIfUsuarioBDExists(nombre string) bool{
	var uftp UsuarioBD
	exists := mgr.db.First(&uftp,"nombre = ?",nombre).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetAllUsuarioBDs() []UsuarioBD {
	var uftps []UsuarioBD
	mgr.db.Preload("BDs").Find(&uftps)
	return uftps
}

func (mgr *manager) GetUsuarioBD(id string) UsuarioBD {
	var uftp UsuarioBD
	mgr.db.Preload("BDs").First(&uftp,id)
	return uftp
}

func (mgr *manager) RemoveUsuarioBD(id string) (err error) {
	return mgr.db.Delete(UsuarioBD{}, "id == ?", id).Error
}

func (mgr *manager) RemoveAssociationUBD(usuario *UsuarioBD, bd *BD) (err error){
	return mgr.db.Model(&usuario).Association("BDs").Delete(&bd).Error
}
