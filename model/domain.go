package model

import "github.com/jinzhu/gorm"

type Domain struct {
	gorm.Model
	Name string
	Status int
	SpamFilter string
	DefaultAccount string
	Accounts []Account
	Lists []List
}

func (mgr *manager) AddDomain(domain *Domain) (err error) {
	return mgr.db.Create(domain).Error
}

func (mgr *manager) UpdateDomain(domain *Domain) (err error) {
	return mgr.db.Save(&domain).Error
}

func (mgr *manager) CheckIfDomainExists(name string) bool{
	var domain Domain
	exists := mgr.db.First(&domain,"name = ?",name).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetAllDomains() []Domain {
	var domains []Domain
	mgr.db.Preload("Accounts").Preload("Lists").Find(&domains)
	return domains
}

func (mgr *manager) GetDomain(id string) Domain {
	var domain Domain
	mgr.db.Preload("Accounts").Preload("Lists").First(&domain,id)
	return domain
}

func (mgr *manager) RemoveDomain(id string) (err error) {
	return mgr.db.Delete(Domain{}, "id == ?", id).Error
}