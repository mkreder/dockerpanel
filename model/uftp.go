package model

import (
	"github.com/jinzhu/gorm"
)

type UsuarioFtp struct {
	gorm.Model
	Nombre  string
	Password string
	Estado int
	WebID  uint
}

func (mgr *manager) UpdateUsuarioFtp(uftp *UsuarioFtp) (err error) {
	return mgr.db.Save(&uftp).Error
}


func (mgr *manager) AddUsuarioFtp(uftp *UsuarioFtp) (err error) {
	return mgr.db.Create(uftp).Error
}


func (mgr *manager) CheckIfUsuarioFtpExists(nombre string, webid string) bool{
	var uftp UsuarioFtp
	exists := mgr.db.First(&uftp,"nombre = ? AND web_id = ?",nombre,webid).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetAllUsuarioFtps() []UsuarioFtp {
	var uftps []UsuarioFtp
	mgr.db.Find(&uftps)
	return uftps
}

func (mgr *manager) GetUsuarioFtp(id string) UsuarioFtp {
	var uftp UsuarioFtp
	mgr.db.First(&uftp,id)
	return uftp
}

func (mgr *manager) RemoveUsuarioFtp(id string) (err error) {
	return mgr.db.Delete(UsuarioFtp{}, "id == ?", id).Error
}