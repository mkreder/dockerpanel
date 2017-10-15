package model

import (
	"github.com/jinzhu/gorm"
)

type Web struct {
	gorm.Model
	Dominio string
	CGI bool
	SSL bool
	PHP bool
	PHPversion string
	PHPmethod string
	Webserver string
	Estado int
}



func (mgr *manager) AddWeb(web *Web) (err error) {
	return mgr.db.Create(web).Error
}

func (mgr *manager) UpdateWeb(web *Web) (err error) {
	return mgr.db.Save(&web).Error
}

func (mgr *manager) CheckIfWebExists(dominio string) bool{
	var web Web
	exists := mgr.db.First(&web,"dominio = ?",dominio).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetAllWebs() []Web {
	var webs []Web
	mgr.db.Find(&webs)
	return webs
}

func (mgr *manager) GetWeb(id string) Web {
	var web Web
	mgr.db.First(&web,id)
	return web
}

func (mgr *manager) RemoveWeb(id string) (err error) {
	return mgr.db.Delete(Web{}, "id == ?", id).Error
}