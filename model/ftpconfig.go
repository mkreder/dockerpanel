package model

import (
	"github.com/jinzhu/gorm"
)

type FtpConfig struct {
	gorm.Model
	AnonWrite int
	AnonRead int
	Estado int
}

func (mgr *manager) UpdateFtpConfig(anonWrite int, anonRead int, estado int) (err error) {
	var ftpConfig FtpConfig
	if  mgr.db.First(&ftpConfig,"id = ?",1).RecordNotFound()  {
		ftpConfig.AnonRead = anonRead
		ftpConfig.AnonWrite = anonWrite
		ftpConfig.Estado = estado
		return mgr.db.Create(&ftpConfig).Error
	} else {
		ftpConfig.AnonRead = anonRead
		ftpConfig.AnonWrite = anonWrite
		ftpConfig.Estado = estado
		return mgr.db.Save(&ftpConfig).Error
	}
}


func (mgr *manager) GetFtpConfig() FtpConfig {
	var ftpConfig FtpConfig
	mgr.db.First(&ftpConfig,1)
	return ftpConfig
}