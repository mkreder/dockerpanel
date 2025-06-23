package model

import "github.com/jinzhu/gorm"

type Autoresponder struct {
	gorm.Model
	Enabled bool
	Subject string
	Message string
	StartDate string
	EndDate string
	AccountID uint
}

func (mgr *manager) AddAutoresponder(aresponder *Autoresponder) (err error) {
	return mgr.db.Create(aresponder).Error
}

func (mgr *manager) UpdateAutoresponder(aresponder *Autoresponder) (err error) {
	return mgr.db.Save(&aresponder).Error
}

func (mgr *manager) CheckIfAutoresponderExists(accountid string) bool{
	var aresponder Autoresponder
	exists := mgr.db.First(&aresponder,"account_id = ?",accountid).RecordNotFound()
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