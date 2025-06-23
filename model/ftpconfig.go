package model

import (
	"github.com/jinzhu/gorm"
)

type FtpConfig struct {
	gorm.Model
	AnonWrite int
	AnonRead int
	Status int
}

func (mgr *manager) UpdateFtpConfig(anonWrite int, anonRead int, status int) (err error) {
	var ftpConfig FtpConfig
	if  mgr.db.First(&ftpConfig,"id = ?",1).RecordNotFound()  {
		ftpConfig.AnonRead = anonRead
		ftpConfig.AnonWrite = anonWrite
		ftpConfig.Status = status
		return mgr.db.Create(&ftpConfig).Error
	} else {
		ftpConfig.AnonRead = anonRead
		ftpConfig.AnonWrite = anonWrite
		ftpConfig.Status = status
		return mgr.db.Save(&ftpConfig).Error
	}
}


func (mgr *manager) GetFtpConfig() FtpConfig {
	var ftpConfig FtpConfig
	mgr.db.First(&ftpConfig,1)
	return ftpConfig
}