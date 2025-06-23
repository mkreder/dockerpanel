package model

import (
	"github.com/jinzhu/gorm"
)

type FTPUser struct {
	gorm.Model
	Name  string
	Password string
	Status int
	WebID  uint
}

func (mgr *manager) UpdateFTPUser(ftpu *FTPUser) (err error) {
	return mgr.db.Save(&ftpu).Error
}


func (mgr *manager) AddFTPUser(ftpu *FTPUser) (err error) {
	return mgr.db.Create(ftpu).Error
}


func (mgr *manager) CheckIfFTPUserExists(name string, webid string) bool{
	var ftpu FTPUser
	exists := mgr.db.First(&ftpu,"name = ? AND web_id = ?",name,webid).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetAllFTPUsers() []FTPUser {
	var ftpus []FTPUser
	mgr.db.Find(&ftpus)
	return ftpus
}

func (mgr *manager) GetFTPUser(id string) FTPUser {
	var ftpu FTPUser
	mgr.db.First(&ftpu,id)
	return ftpu
}

func (mgr *manager) RemoveFTPUser(id string) (err error) {
	return mgr.db.Delete(FTPUser{}, "id == ?", id).Error
}