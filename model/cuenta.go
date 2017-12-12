package model

import (
	"github.com/jinzhu/gorm"
	"strconv"
)

type Cuenta struct {
	gorm.Model
	Nombre  string
	NombreReal	string
	Password string
	Cuota int
	Estado int
	DominioID uint
	Autoresponder Autoresponder
	Renvio string
}

func (mgr *manager) AddCuenta(cuenta *Cuenta) (err error) {
	return mgr.db.Create(cuenta).Error
}

func (mgr *manager) UpdateCuenta(cuenta *Cuenta) (err error) {
	return mgr.db.Save(&cuenta).Error
}

func (mgr *manager) CheckIfCuentaExists(nombre string, dominioid string) bool{
	var cuenta Cuenta
	exists := mgr.db.First(&cuenta,"nombre = ? AND dominio_id = ?",nombre,dominioid).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetAllCuentas() []Cuenta {
	var cuentas []Cuenta
	mgr.db.Preload("Autoresponder").Find(&cuentas)
	return cuentas
}

func (mgr *manager) GetCuenta(id string) Cuenta {
	var cuenta Cuenta
	mgr.db.Preload("Autoresponder").First(&cuenta,id)
	return cuenta
}

func (mgr *manager) RemoveCuenta(id string) (err error) {
	cuenta := Mgr.GetCuenta(id)
	autoresponderid := strconv.Itoa(int(cuenta.Autoresponder.ID))
//	mgr.db.Model(&cuenta).Association("Autoresponder").Delete(&cuenta.Autoresponder)
	Mgr.RemoveAutoresponder(autoresponderid)
	return mgr.db.Delete(Cuenta{}, "id == ?", id).Error
}

func (mgr *manager) GetCuentas(dominioid string) []Cuenta {
	var cuentas []Cuenta
	mgr.db.Preload("Autoresponder").Where("dominio_id = ?", dominioid).Find(&cuentas)
	return cuentas
}

