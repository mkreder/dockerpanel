package model

import "github.com/jinzhu/gorm"

type List struct {
	gorm.Model
	Name  string
	AdminEmail string
	Password string
	Status int
	DomainID uint

}

func (mgr *manager) AddList(list *List) (err error) {
	return mgr.db.Create(list).Error
}

func (mgr *manager) UpdateList(list *List) (err error) {
	return mgr.db.Save(&list).Error
}

func (mgr *manager) CheckIfListExists(name string, domainid string) bool{
	var list List
	exists := mgr.db.First(&list,"name = ? AND domain_id = ?",name,domainid).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetAllLists() []List {
	var lists []List
	mgr.db.Find(&lists)
	return lists
}

func (mgr *manager) GetList(id string) List {
	var list List
	mgr.db.First(&list,id)
	return list
}

func (mgr *manager) RemoveList(id string) (err error) {
	return mgr.db.Delete(List{}, "id == ?", id).Error
}

func (mgr *manager) GetLists(domainid string) []List {
	var lists []List
	mgr.db.Where("domain_id = ?", domainid).Find(&lists)
	return lists
}