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

