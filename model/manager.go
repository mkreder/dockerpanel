package model

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Manager interface {
	AddWeb(web *Web) error
	CheckIfWebExists(dominio string) bool
	GetAllWebs() []Web
	RemoveWeb(id string) (err error)
	UpdateWeb(web *Web) (err error)
	GetWeb(id string) Web

	GetUser(email string) User
	UpdatePassword(Usuario string, hash string) (err error)

	AddZona(zona *Zona) error
	CheckIfZonaExists(dominio string) bool
	GetAllZonas() []Zona
	RemoveZona(id string) (err error)
	UpdateZona(zona *Zona) (err error)
	GetZona(id string) Zona

	AddRegistros(registros []Registro)
	GetRegistros(zonaid string) []Registro
	CheckIfRegistroExists(nombre string, tipo string, valor string, prioridad string, zonaid string) bool
	GetRegistro(id string) Registro
	AddRegistro(registro *Registro) (err error)
	UpdateRegistro(registro *Registro) (err error)
	RemoveRegistro(id string) (err error)

	AddUsuarioFtp(uftp *UsuarioFTP) (err error)
	UpdateUsuarioFtp(ftp *UsuarioFTP) (err error)
	CheckIfUsuarioFtpExists(nombre string, webid string) bool
	GetAllUsuarioFtps() []UsuarioFTP
	GetUsuarioFtp(id string) UsuarioFTP
	RemoveUsuarioFtp(id string) (err error)
	UpdateFtpConfig(anonWrite int, anonRead int, estado int) (err error)
	GetFtpConfig() FtpConfig

	AddDB(db *DB) (err error)
	UpdateDB(db *DB) (err error)
	CheckIfDBExists(name string) bool
	GetAllDBs() []DB
	GetDB(id string) DB
	RemoveDB(id string) (err error)
	RemoveAssociationIP(db *DB, ip *IP) (err error)

	UpdateIP(ip IP, bd BD)

	UpdateDBUser(dbu *DBUser) (err error)
	AddDBUser(dbu *DBUser) (err error)
	CheckIfDBUserExists(name string) bool
	GetAllDBUsers() []DBUser
	GetDBUser(id string) DBUser
	RemoveDBUser(id string) (err error)
	GetDBUserbyDB(dbid string) []DBUser

	AddDBAssociation(dba *DBAssociation) (err error)
	UpdateDBAssociation(dba *DBAssociation) (err error)
	CheckIfDBAssociationExists(dbid string, dbuserid string) bool
	GetDBAssociations() []DBAssociation
	GetDBAssociation(dbid string, dbuserid string) DBAssociation
	RemoveAsociacionBD(adb AsociacionBD)

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
	db, err := gorm.Open("sqlite3", "dockerpanel.db")
	//migrate()
	if err != nil {
		log.Fatal("Failed to init db:", err)
	}
	Mgr = &manager{db: db}
	Mgr.migrate()
}

func (mgr *manager) migrate() {
	mgr.db.AutoMigrate(&Web{})
	mgr.db.AutoMigrate(&User{})
	mgr.db.AutoMigrate(&Zona{})
	mgr.db.AutoMigrate(&Registro{})
	mgr.db.AutoMigrate(&UsuarioFTP{})
	mgr.db.AutoMigrate(&FtpConfig{})
	mgr.db.AutoMigrate(&DBUser{})
	mgr.db.AutoMigrate(&AsociacionBD{})
	mgr.db.AutoMigrate(&BD{})
	mgr.db.AutoMigrate(&IP{})
	mgr.db.AutoMigrate(&Dominio{})
	mgr.db.AutoMigrate(&Cuenta{})
	mgr.db.AutoMigrate(&Lista{})
	mgr.db.AutoMigrate(&Autoresponder{})

	// Add default Usuario
	var usr User

	if mgr.db.First(&usr, "email = ?", "admin@admin.com").RecordNotFound() {
		usr.Email = "admin@admin.com"
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		usr.Password = string(hash)
		mgr.db.Create(&usr)
	}

	var ftpConfig FtpConfig
	if mgr.db.First(&ftpConfig, "id = 1").RecordNotFound() {
		mgr.UpdateFtpConfig(0, 0, 1)
	}
}
