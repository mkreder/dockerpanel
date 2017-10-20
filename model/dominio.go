package model

import "github.com/jinzhu/gorm"

type Dominio struct {
	gorm.Model
	Nombre string
	Estado int
	FiltroSpam string
	CuentaDefecto string
	Cuentas []Cuenta
	Listas []Lista
}

func (mgr *manager) AddDominio(dominio *Dominio) (err error) {
	return mgr.db.Create(dominio).Error
}

func (mgr *manager) UpdateDominio(dominio *Dominio) (err error) {
	return mgr.db.Save(&dominio).Error
}

func (mgr *manager) CheckIfDominioExists(nombre string) bool{
	var dominio Dominio
	exists := mgr.db.First(&dominio,"nombre = ?",nombre).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetAllDominios() []Dominio {
	var dominios []Dominio
	mgr.db.Preload("Cuentas").Preload("Listas").Find(&dominios)
	return dominios
}

func (mgr *manager) GetDominio(id string) Dominio {
	var dominio Dominio
	mgr.db.Preload("Cuentas").Preload("Listas").First(&dominio,id)
	return dominio
}

func (mgr *manager) RemoveDominio(id string) (err error) {
	return mgr.db.Delete(Dominio{}, "id == ?", id).Error
}