package model

import "github.com/jinzhu/gorm"

type IP struct {
	gorm.Model
	Value  string
	DatabaseID uint
	// 1: enabled 2: disabled
	Status int
}


func (mgr *manager) UpdateIP(ip IP,db Database)  {
	mgr.db.Exec("update ips set status = 2 where id = ?",ip.ID)
	mgr.db.Exec("update databases set status = 1 where id = ?",db.ID)
	mgr.db.Save(&ip)
}
