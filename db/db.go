package db

import (
	"github.com/jinzhu/gorm"
	"log"
	"github.com/mkreder/dockerpanel/model"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/mattn/go-sqlite3"
)

type Manager interface {
	AddWeb(web *model.Web) error
	CheckIfWebExists(dominio string) bool
    GetAllWebs() []model.Web
	RemoveWeb(id string) (err error)
	UpdateWeb(web *model.Web) (err error)
	GetWeb(id string) model.Web

	GetUser(email string) model.User
	UpdatePassword(user string, hash string) (err error)

	AddZona(zona *model.Zona) error
	CheckIfZonaExists(dominio string) bool
	GetAllZonas() []model.Zona
	RemoveZona(id string) (err error)
	UpdateZona(zona *model.Zona) (err error)
	GetZona(id string) model.Zona

	AddRegistros(registros []model.Registro)
	GetRegistros(zonaid string) []model.Registro
	CheckIfRegistroExists(nombre string,tipo string, valor string, prioridad string) bool
	GetRegistro(id string) model.Registro
	AddRegistro(registro *model.Registro) (err error)
	UpdateRegistro(registro *model.Registro) (err error)
	RemoveRegistro(id string) (err error)


	migrate()
	// Add other methods
}

type manager struct {
	db *gorm.DB
}

var Mgr Manager

func init() {
	db, err :=gorm.Open("sqlite3", "dockerpanel.db")
	//migrate()
	if err != nil {
		log.Fatal("Failed to init db:", err)
	}
	Mgr = &manager{db: db}
	Mgr.migrate()
}
func (mgr *manager) migrate(){
	mgr.db.AutoMigrate(&model.Web{})
	mgr.db.AutoMigrate(&model.User{})
	mgr.db.AutoMigrate(&model.Zona{})
	mgr.db.AutoMigrate(&model.Registro{})

	// Add default user
	var usr model.User

	if mgr.db.First(&usr,"email = ?","admin@admin.com").RecordNotFound(){
		usr.Email = "admin@admin.com"
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		usr.Password = string(hash)
		mgr.db.Create(&usr)
	}

}


// Web


func (mgr *manager) AddWeb(web *model.Web) (err error) {
	return mgr.db.Create(web).Error
}

func (mgr *manager) UpdateWeb(web *model.Web) (err error) {
	return mgr.db.Save(&web).Error
}

func (mgr *manager) CheckIfWebExists(dominio string) bool{
	var web model.Web
	exists := mgr.db.First(&web,"dominio = ?",dominio).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetAllWebs() []model.Web {
	var webs []model.Web
	mgr.db.Find(&webs)
	return webs
}

func (mgr *manager) GetWeb(id string) model.Web {
	var web model.Web
	mgr.db.First(&web,id)
	return web
}

func (mgr *manager) RemoveWeb(id string) (err error) {
	return mgr.db.Delete(model.Web{}, "id == ?", id).Error
}

// User

func (mgr *manager) GetUser(email string) model.User {
	var usr model.User
	mgr.db.First(&usr,"email = ?",email)
	return usr
}

func (mgr *manager) UpdatePassword(user string, hash string) (err error) {
	var usr model.User
	mgr.db.First(&usr,"email = ?",user)
	usr.Password = hash
	return mgr.db.Save(&usr).Error
}


// Zona

func (mgr *manager) AddZona(zona *model.Zona) (err error) {
	return mgr.db.Create(zona).Error
}

func (mgr *manager) UpdateZona(zona *model.Zona) (err error) {
	return mgr.db.Save(&zona).Error
}

func (mgr *manager) CheckIfZonaExists(dominio string) bool{
	var zona model.Zona
	exists := mgr.db.First(&zona,"dominio = ?",dominio).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetAllZonas() []model.Zona {
	var zonas []model.Zona
	mgr.db.Find(&zonas)
	return zonas
}

func (mgr *manager) GetZona(id string) model.Zona {
	var zona model.Zona
	mgr.db.First(&zona,id)
	return zona
}

func (mgr *manager) RemoveZona(id string) (err error) {
	return mgr.db.Delete(model.Zona{}, "id == ?", id).Error
}

func (mgr *manager) AddRegistros(registros []model.Registro) {
	for  _ , registro := range registros {
		mgr.db.Create(&registro)

	}
}

func (mgr *manager) GetRegistros(zonaid string) []model.Registro{
	var registros []model.Registro
	mgr.db.Where("zona_id = ?", zonaid).Find(&registros)
	return registros
}

func (mgr *manager) CheckIfRegistroExists(nombre string,tipo string, valor string, prioridad string) bool{
	var registro model.Registro
	exists := mgr.db.First(&registro,"nombre = ? AND tipo = ? AND valor = ? AND prioridad = ? ",nombre,tipo,valor,prioridad).RecordNotFound()
	return ! exists
}

func (mgr *manager) GetRegistro(id string) model.Registro {
	var registro model.Registro
	mgr.db.First(&registro,id)
	return registro
}

func (mgr *manager) AddRegistro(registro *model.Registro) (err error) {
	return mgr.db.Create(registro).Error
}

func (mgr *manager) UpdateRegistro(registro *model.Registro) (err error) {
	return mgr.db.Save(&registro).Error
}

func (mgr *manager) RemoveRegistro(id string) (err error) {
	return mgr.db.Delete(model.Registro{}, "id == ?", id).Error
}