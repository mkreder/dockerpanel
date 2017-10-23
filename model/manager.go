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
	CheckIfRegistroExists(nombre string,tipo string, valor string, prioridad string, zonaid string ) bool
	GetRegistro(id string) Registro
	AddRegistro(registro *Registro) (err error)
	UpdateRegistro(registro *Registro) (err error)
	RemoveRegistro(id string) (err error)

	AddUsuarioFtp(uftp *UsuarioFTP)  (err error)
	UpdateUsuarioFtp(ftp *UsuarioFTP) (err error)
	CheckIfUsuarioFtpExists(nombre string, webid string) bool
	GetAllUsuarioFtps() []UsuarioFTP
	GetUsuarioFtp(id string) UsuarioFTP
	RemoveUsuarioFtp(id string) (err error)
	UpdateFtpConfig(anonWrite int, anonRead int) (err error)
	GetFtpConfig() FtpConfig

    AddBD(bd *BD) (err error)
	UpdateBD(bd *BD) (err error)
	CheckIfBDExists(nombre string) bool
	GetAllBDs() []BD
	GetBD(id string) BD
	RemoveBD(id string) (err error)

	UpdateUsuarioBD(uftp *UsuarioBD) (err error)
	AddUsuarioBD(uftp *UsuarioBD) (err error)
	CheckIfUsuarioBDExists(nombre string) bool
	GetAllUsuarioBDs() []UsuarioBD
	GetUsuarioBD(id string) UsuarioBD
	RemoveUsuarioBD(id string) (err error)
	RemoveAssociationUBD(usuario *UsuarioBD, bd *BD) (err error)
	RemoveAssociationIP(bd *BD, ip *IP) (err error)

	AddDominio(dominio *Dominio) (err error)
	UpdateDominio(dominio *Dominio) (err error)
	CheckIfDominioExists(nombre string) bool
	GetAllDominios() []Dominio
	GetDominio(id string) Dominio
	RemoveDominio(id string) (err error)

	AddCuenta(cuenta *Cuenta) (err error)
	UpdateCuenta(cuenta *Cuenta) (err error)
	CheckIfCuentaExists(nombre string, dominioid string) bool
	GetAllCuentas() []Cuenta
	GetCuenta(id string) Cuenta
	RemoveCuenta(id string) (err error)
	GetCuentas(dominioid string) []Cuenta

	AddLista(lista *Lista) (err error)
	UpdateLista(lista *Lista) (err error)
	CheckIfListaExists(nombre string, dominioid string) bool
	GetAllListas() []Lista
	GetLista(id string) Lista
	RemoveLista(id string) (err error)
	GetListas(dominioid string) []Lista

	AddAutoresponder(aresponder *Autoresponder) (err error)
	UpdateAutoresponder(aresponder *Autoresponder) (err error)
	CheckIfAutoresponderExists(cuentaid string) bool
	GetAllAutoresponders() []Autoresponder
	GetAutoresponder(id string) Autoresponder
	RemoveAutoresponder(id string) (err error)

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

func (mgr *manager) migrate() {
	mgr.db.AutoMigrate(&Web{})
	mgr.db.AutoMigrate(&Usuario{})
	mgr.db.AutoMigrate(&Zona{})
	mgr.db.AutoMigrate(&Registro{})
	mgr.db.AutoMigrate(&UsuarioFTP{})
	mgr.db.AutoMigrate(&FtpConfig{})
	mgr.db.AutoMigrate(&UsuarioBD{})
	mgr.db.AutoMigrate(&BD{})
	mgr.db.AutoMigrate(&IP{})
	mgr.db.AutoMigrate(&Dominio{})
	mgr.db.AutoMigrate(&Cuenta{})
	mgr.db.AutoMigrate(&Lista{})
	mgr.db.AutoMigrate(&Autoresponder{})

	// Add default Usuario
	var usr Usuario

	if mgr.db.First(&usr, "email = ?", "admin@admin.com").RecordNotFound() {
		usr.Email = "admin@admin.com"
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		usr.Password = string(hash)
		mgr.db.Create(&usr)
	}

	var ftpConfig FtpConfig
	if mgr.db.First(&ftpConfig, "id = 1").RecordNotFound() {
		mgr.UpdateFtpConfig(0, 0)
	}
}

