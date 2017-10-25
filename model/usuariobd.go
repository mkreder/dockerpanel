package model

import (
	"github.com/jinzhu/gorm"
)

type UsuarioBD struct {
	gorm.Model
	Nombre  string
	Password string
	AsociacionBDs []AsociacionBD
	// 1: activado 2: desactivado
	Estado int
}

func (mgr *manager) UpdateUsuarioBD(ubd *UsuarioBD) (err error) {
	return mgr.db.Save(&ubd).Error
}


func (mgr *manager) AddUsuarioBD(ubd *UsuarioBD) (err error) {
	return mgr.db.Create(ubd).Error
}


func (mgr *manager) CheckIfUsuarioBDExists(nombre string) bool{
	var uftp UsuarioBD
	exists := mgr.db.First(&uftp,"nombre = ?",nombre).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetAllUsuarioBDs() []UsuarioBD {
	var uftps []UsuarioBD
	mgr.db.Preload("AsociacionBDs").Find(&uftps)
	return uftps
}

func (mgr *manager) GetUsuarioBD(id string) UsuarioBD {
	var uftp UsuarioBD
	mgr.db.Preload("AsociacionBDs").First(&uftp,id)
	return uftp
}

func (mgr *manager) RemoveUsuarioBD(id string) (err error) {
	return mgr.db.Delete(UsuarioBD{}, "id == ?", id).Error
}


func (mgr *manager) GetUsuariosDeBD(bdid string) []UsuarioBD{
	var ubds []UsuarioBD
	mgr.db.Joins("left join asociacion_bds on asociacion_bds.usuario_bd_id = usuario_bds.id").Where("asociacion_bds.bd_id = ?", bdid).Find(&ubds)
	return ubds
}
