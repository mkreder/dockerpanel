package model

import (
	"github.com/jinzhu/gorm"
)

type Usuario struct {
	gorm.Model
	Email string
	Password string
}

func (mgr *manager) GetUsuario(email string) Usuario {
	var usr Usuario
	mgr.db.First(&usr,"email = ?",email)
	return usr
}

func (mgr *manager) UpdatePassword(usuario string, hash string) (err error) {
	var usr Usuario
	mgr.db.First(&usr,"email = ?",usuario)
	usr.Password = hash
	return mgr.db.Save(&usr).Error
}
