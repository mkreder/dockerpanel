package model

import (
	"github.com/jinzhu/gorm"
)

type BD struct {
	gorm.Model
	Nombre  string
	Estado int
	AsociacionBDs []AsociacionBD
	IPs []IP
}


func (mgr *manager) AddBD(bd *BD) (err error) {
	return mgr.db.Create(bd).Error
}

func (mgr *manager) UpdateBD(bd *BD) (err error) {
	return mgr.db.Save(&bd).Error
}

func (mgr *manager) CheckIfBDExists(nombre string) bool{
	var bd BD
	exists := mgr.db.First(&bd,"nombre = ?",nombre).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetAllBDs() []BD {
	var bds []BD
	mgr.db.Preload("AsociacionBDs").Preload("IPs").Find(&bds)
	return bds
}

func (mgr *manager) GetBD(id string) BD {
	var bd BD
	mgr.db.Preload("AsociacionBDs").Preload("IPs").First(&bd,id)
	return bd
}

func (mgr *manager) RemoveBD(id string) (err error) {
	return mgr.db.Delete(BD{}, "id == ?", id).Error
}

func (mgr *manager) RemoveAssociationIP(bd *BD, ip *IP) (err error){
	error := mgr.db.Model(&bd).Association("IPs").Delete(&ip).Error
	mgr.db.Delete(IP{},"valor == ?", ip.Valor)
	return error
}
