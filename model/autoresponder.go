package model

import "github.com/jinzhu/gorm"

type Autoresponder struct {
	gorm.Model
	Activado bool
	Asunto string
	Mensaje string
	FechaIncio string
	FechaFin string
	CuentaID uint
}

func (mgr *manager) AddAutoresponder(aresponder *Autoresponder) (err error) {
	return mgr.db.Create(aresponder).Error
}

func (mgr *manager) UpdateAutoresponder(aresponder *Autoresponder) (err error) {
	return mgr.db.Save(&aresponder).Error
}

func (mgr *manager) CheckIfAutoresponderExists(cuentaid string) bool{
	var aresponder Autoresponder
	exists := mgr.db.First(&aresponder,"cuenta_id = ?",cuentaid).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetAllAutoresponders() []Autoresponder {
	var aresponders []Autoresponder
	mgr.db.Find(&aresponders)
	return aresponders
}

func (mgr *manager) GetAutoresponder(id string) Autoresponder {
	var aresponder Autoresponder
	mgr.db.First(&aresponder,id)
	return aresponder
}

func (mgr *manager) RemoveAutoresponder(id string) (err error) {
	return mgr.db.Delete(Autoresponder{}, "id == ?", id).Error
}