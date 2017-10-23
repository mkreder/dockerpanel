package model

import (
	"github.com/jinzhu/gorm"
)

type Registro struct {
	gorm.Model
	Tipo string
	Nombre string
	Valor string
	Prioridad string
	ZonaID uint
}


func (mgr *manager) AddRegistros(registros []Registro) {
	var zona Zona
	for  _ , registro := range registros {
		mgr.db.Create(&registro)
	}
	mgr.db.First(&zona,registros[0].ZonaID)
	zona.Estado = 1
	mgr.db.Save(&zona)
}

func (mgr *manager) GetRegistros(zonaid string) []Registro{
	var registros []Registro
	mgr.db.Where("zona_id = ?", zonaid).Find(&registros)
	return registros
}

func (mgr *manager) CheckIfRegistroExists(nombre string,tipo string, valor string, prioridad string, zonaid string) bool{
	var registro Registro
	exists := mgr.db.First(&registro,"nombre = ? AND tipo = ? AND valor = ? AND prioridad = ? AND zona_id = ?",nombre,tipo,valor,prioridad,zonaid).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetRegistro(id string) Registro {
	var registro Registro
	mgr.db.First(&registro,id)
	return registro
}

func (mgr *manager) AddRegistro(registro *Registro) (err error) {
	var zona Zona
	mgr.db.First(&zona,registro.ZonaID)
	zona.Estado = 1
	mgr.db.Save(&zona)
	return mgr.db.Create(registro).Error
}

func (mgr *manager) UpdateRegistro(registro *Registro) (err error) {
	var zona Zona
	mgr.db.First(&zona,registro.ZonaID)
	zona.Estado = 1
	mgr.db.Save(&zona)
	return mgr.db.Save(&registro).Error
}

func (mgr *manager) RemoveRegistro(id string) (err error) {
	var zona Zona
	var registro Registro
	mgr.db.First(&registro,id)
	mgr.db.First(&zona,registro.ZonaID)
	zona.Estado = 1
	mgr.db.Save(&zona)
	return mgr.db.Delete(Registro{}, "id == ?", id).Error
}