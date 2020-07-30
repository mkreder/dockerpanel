package model

import "github.com/jinzhu/gorm"

//IP struct to represent an IP whitelisted to access a DB
type IP struct {
	gorm.Model
	Value string
	DbID  uint
	// 1: activado 2: desactivado
	Status int
}

func (mgr *manager) UpdateIP(ip IP, db DB) {
	mgr.db.Exec("update ips set Status = 2 where id = ?", ip.ID)
	mgr.db.Exec("update bds set Status = 1 where id = ?", db.ID)
	mgr.db.Save(&ip)
}
