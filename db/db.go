package db

import (
	"github.com/jinzhu/gorm"
	"log"
	"github.com/mkreder/dockerpanel/model"
	_ "github.com/mattn/go-sqlite3"
)

type Manager interface {
	AddWeb(web *model.Web) error
	CheckIfWebExists(dominio string) bool
    GetAllWebs() []model.Web
	RemoveWeb(id string) (err error)
	UpdateWeb(web *model.Web) (err error)
	GetWeb(id string) model.Web

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
