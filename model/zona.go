package model

import (
	"github.com/jinzhu/gorm"
	"strconv"
)

type Zona struct {
	gorm.Model
	Dominio string
	Email string
	Estado int
	Registros []Registro
}


func (mgr *manager) AddZona(zona *Zona) (err error) {
	return mgr.db.Create(zona).Error
}

func (mgr *manager) UpdateZona(zona *Zona) (err error) {
	return mgr.db.Save(&zona).Error
}

func (mgr *manager) CheckIfZonaExists(dominio string) bool{
	var zona Zona
	exists := mgr.db.First(&zona,"dominio = ?",dominio).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetAllZonas() []Zona {
	var zonas []Zona
	mgr.db.Preload("Registros").Find(&zonas)
	return zonas
}

func (mgr *manager) GetZona(id string) Zona {
	var zona Zona
	mgr.db.Preload("Registros").First(&zona,id)
	return zona
}

func (mgr *manager) RemoveZona(id string) (err error) {
	for _ , registro := range mgr.GetZona(id).Registros {
		mgr.RemoveRegistro(strconv.Itoa(int(registro.ID)))
	}
	return mgr.db.Delete(Zona{}, "id == ?", id).Error
}