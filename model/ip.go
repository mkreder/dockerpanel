package model

import "github.com/jinzhu/gorm"

type IP struct {
	gorm.Model
	Valor  string
	BdID uint
	// 1: activado 2: desactivado
	Estado int
}


func (mgr *manager) UpdateIP(ip IP,bd BD)  {
	mgr.db.Exec("update ips set Estado = 2 where id = ?",ip.ID)
	mgr.db.Exec("update bds set Estado = 1 where id = ?",bd.ID)
	mgr.db.Save(&ip)
}
