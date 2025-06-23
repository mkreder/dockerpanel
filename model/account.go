package model

import (
	"github.com/jinzhu/gorm"
	"strconv"
)

type Account struct {
	gorm.Model
	Name  string
	RealName	string
	Password string
	Quota int
	Status int
	DomainID uint
	Autoresponder Autoresponder
	Forward string
}

func (mgr *manager) AddAccount(account *Account) (err error) {
	return mgr.db.Create(account).Error
}

func (mgr *manager) UpdateAccount(account *Account) (err error) {
	return mgr.db.Save(&account).Error
}

func (mgr *manager) CheckIfAccountExists(name string, domainid string) bool{
	var account Account
	exists := mgr.db.First(&account,"name = ? AND domain_id = ?",name,domainid).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetAllAccounts() []Account {
	var accounts []Account
	mgr.db.Preload("Autoresponder").Find(&accounts)
	return accounts
}

func (mgr *manager) GetAccount(id string) Account {
	var account Account
	mgr.db.Preload("Autoresponder").First(&account,id)
	return account
}

func (mgr *manager) RemoveAccount(id string) (err error) {
	account := Mgr.GetAccount(id)
	autoresponderid := strconv.Itoa(int(account.Autoresponder.ID))
//	mgr.db.Model(&account).Association("Autoresponder").Delete(&account.Autoresponder)
	Mgr.RemoveAutoresponder(autoresponderid)
	return mgr.db.Delete(Account{}, "id == ?", id).Error
}

func (mgr *manager) GetAccounts(domainid string) []Account {
	var accounts []Account
	mgr.db.Preload("Autoresponder").Where("domain_id = ?", domainid).Find(&accounts)
	return accounts
}

