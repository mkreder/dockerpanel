package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email    string
	Password string
}

func (mgr *manager) GetUser(email string) User {
	var usr User
	mgr.db.First(&usr, "email = ?", email)
	return usr
}

func (mgr *manager) UpdatePassword(usuario string, hash string) (err error) {
	var usr User
	mgr.db.First(&usr, "email = ?", usuario)
	usr.Password = hash
	return mgr.db.Save(&usr).Error
}
