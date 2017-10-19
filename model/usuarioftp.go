package model

import (
	"github.com/jinzhu/gorm"
)

type UsuarioFTP struct {
	gorm.Model
	Nombre  string
	Password string
	Estado int
	WebID  uint
}

func (mgr *manager) UpdateUsuarioFtp(uftp *UsuarioFTP) (err error) {
	return mgr.db.Save(&uftp).Error
}


func (mgr *manager) AddUsuarioFtp(uftp *UsuarioFTP) (err error) {
	return mgr.db.Create(uftp).Error
}


func (mgr *manager) CheckIfUsuarioFtpExists(nombre string, webid string) bool{
	var uftp UsuarioFTP
	exists := mgr.db.First(&uftp,"nombre = ? AND web_id = ?",nombre,webid).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetAllUsuarioFtps() []UsuarioFTP {
	var uftps []UsuarioFTP
	mgr.db.Find(&uftps)
	return uftps
}

func (mgr *manager) GetUsuarioFtp(id string) UsuarioFTP {
	var uftp UsuarioFTP
	mgr.db.First(&uftp,id)
	return uftp
}

func (mgr *manager) RemoveUsuarioFtp(id string) (err error) {
	return mgr.db.Delete(UsuarioFTP{}, "id == ?", id).Error
}