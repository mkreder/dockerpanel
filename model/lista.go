package model

import "github.com/jinzhu/gorm"

type Lista struct {
	gorm.Model
	Nombre  string
	Password string
	Estado int
	DominioID uint

}

func (mgr *manager) AddLista(lista *Lista) (err error) {
	return mgr.db.Create(lista).Error
}

func (mgr *manager) UpdateLista(lista *Lista) (err error) {
	return mgr.db.Save(&lista).Error
}

func (mgr *manager) CheckIfListaExists(nombre string, dominioid string) bool{
	var lista Lista
	exists := mgr.db.First(&lista,"nombre = ? AND dominio_id = ?",nombre,dominioid).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetAllListas() []Lista {
	var listas []Lista
	mgr.db.Find(&listas)
	return listas
}

func (mgr *manager) GetLista(id string) Lista {
	var lista Lista
	mgr.db.First(&lista,id)
	return lista
}

func (mgr *manager) RemoveLista(id string) (err error) {
	return mgr.db.Delete(Lista{}, "id == ?", id).Error
}

func (mgr *manager) GetListas(dominioid string) []Lista {
	var listas []Lista
	mgr.db.Where("dominio_id = ?", dominioid).Find(&listas)
	return listas
}