package model

import (
	"github.com/jinzhu/gorm"
)

type Ftp struct {
	gorm.Model
	Username string
	Password string
	//TODO: review Web association
	Web Web
}

func (mgr *manager) UpdateFtp(ftp *Ftp) (err error) {
	return mgr.db.Save(&ftp).Error
}