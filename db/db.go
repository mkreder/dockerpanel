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

func (mgr *manager) AddWeb(web *model.Web) (err error) {
	mgr.db.Create(web)
	if errs := mgr.db.GetErrors(); len(errs) > 0 {
		err = errs[0]
	}
	return
}

func (mgr *manager) CheckIfWebExists(dominio string) bool{
	var web model.Web
	exists := mgr.db.First(&web,"dominio = ?",dominio).RecordNotFound()
	return ! exists
}

func (mgr *manager) migrate(){
	mgr.db.AutoMigrate(&model.Web{})
}

