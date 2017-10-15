package model

import (
	"github.com/jinzhu/gorm"
	"log"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/mattn/go-sqlite3"
)

type Manager interface {
	AddWeb(web *Web) error
	CheckIfWebExists(dominio string) bool
    GetAllWebs() []Web
	RemoveWeb(id string) (err error)
	UpdateWeb(web *Web) (err error)
	GetWeb(id string) Web

	GetUsuario(email string) Usuario
	UpdatePassword(Usuario string, hash string) (err error)

	AddZona(zona *Zona) error
	CheckIfZonaExists(dominio string) bool
	GetAllZonas() []Zona
	RemoveZona(id string) (err error)
	UpdateZona(zona *Zona) (err error)
	GetZona(id string) Zona

	AddRegistros(registros []Registro)
	GetRegistros(zonaid string) []Registro
	CheckIfRegistroExists(nombre string,tipo string, valor string, prioridad string) bool
	GetRegistro(id string) Registro
	AddRegistro(registro *Registro) (err error)
	UpdateRegistro(registro *Registro) (err error)
	RemoveRegistro(id string) (err error)


	UpdateFtp(ftp *Ftp) (err error)

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
	mgr.db.AutoMigrate(&Web{})
	mgr.db.AutoMigrate(&Usuario{})
	mgr.db.AutoMigrate(&Zona{})
	mgr.db.AutoMigrate(&Registro{})

	// Add default Usuario
	var usr Usuario

	if mgr.db.First(&usr,"email = ?","admin@admin.com").RecordNotFound(){
		usr.Email = "admin@admin.com"
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		usr.Password = string(hash)
		mgr.db.Create(&usr)
	}

}
