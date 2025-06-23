package model

import (
	"github.com/jinzhu/gorm"
	"log"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/mattn/go-sqlite3"
)

type Manager interface {
	AddWeb(web *Web) error
	CheckIfWebExists(domain string) bool
    GetAllWebs() []Web
	RemoveWeb(id string) (err error)
	UpdateWeb(web *Web) (err error)
	GetWeb(id string) Web

	GetUser(email string) User
	UpdatePassword(user string, hash string) (err error)

	AddZone(zone *Zone) error
	CheckIfZoneExists(domain string) bool
	GetAllZones() []Zone
	RemoveZone(id string) (err error)
	UpdateZone(zone *Zone) (err error)
	GetZone(id string) Zone

	AddRecords(records []Record)
	GetRecords(zoneid string) []Record
	CheckIfRecordExists(name string,recordtype string, value string, priority string, zoneid string ) bool
	GetRecord(id string) Record
	AddRecord(record *Record) (err error)
	UpdateRecord(record *Record) (err error)
	RemoveRecord(id string) (err error)

	AddFTPUser(ftpu *FTPUser)  (err error)
	UpdateFTPUser(ftpu *FTPUser) (err error)
	CheckIfFTPUserExists(name string, webid string) bool
	GetAllFTPUsers() []FTPUser
	GetFTPUser(id string) FTPUser
	RemoveFTPUser(id string) (err error)
	UpdateFtpConfig(anonWrite int, anonRead int, status int) (err error)
	GetFtpConfig() FtpConfig

    AddDatabase(db *Database) (err error)
	UpdateDatabase(db *Database) (err error)
	CheckIfDatabaseExists(name string) bool
	GetAllDatabases() []Database
	GetDatabase(id string) Database
	RemoveDatabase(id string) (err error)
	UpdateIP(ip IP,db Database)

	UpdateDatabaseUser(dbu *DatabaseUser) (err error)
	AddDatabaseUser(dbu *DatabaseUser) (err error)
	CheckIfDatabaseUserExists(name string) bool
	GetAllDatabaseUsers() []DatabaseUser
	GetDatabaseUser(id string) DatabaseUser
	RemoveDatabaseUser(id string) (err error)
	RemoveAssociationIP(db *Database, ip *IP) (err error)
	GetUsersOfDatabase(dbid string) []DatabaseUser

	AddDatabaseAssociation(dba *DatabaseAssociation) (err error)
	UpdateDatabaseAssociation(dba *DatabaseAssociation) (err error)
	CheckIfDatabaseAssociationExists(dbid string, dbuserid string) bool
	GetAllDatabaseAssociations() []DatabaseAssociation
	GetDatabaseAssociation(dbid string, dbuserid string) DatabaseAssociation
	RemoveDatabaseAssociation(dba DatabaseAssociation)

	AddDomain(domain *Domain) (err error)
	UpdateDomain(domain *Domain) (err error)
	CheckIfDomainExists(name string) bool
	GetAllDomains() []Domain
	GetDomain(id string) Domain
	RemoveDomain(id string) (err error)

	AddAccount(account *Account) (err error)
	UpdateAccount(account *Account) (err error)
	CheckIfAccountExists(name string, domainid string) bool
	GetAllAccounts() []Account
	GetAccount(id string) Account
	RemoveAccount(id string) (err error)
	GetAccounts(domainid string) []Account

	AddList(list *List) (err error)
	UpdateList(list *List) (err error)
	CheckIfListExists(name string, domainid string) bool
	GetAllLists() []List
	GetList(id string) List
	RemoveList(id string) (err error)
	GetLists(domainid string) []List

	AddAutoresponder(aresponder *Autoresponder) (err error)
	UpdateAutoresponder(aresponder *Autoresponder) (err error)
	CheckIfAutoresponderExists(accountid string) bool
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
	mgr.db.AutoMigrate(&User{})
	mgr.db.AutoMigrate(&Zone{})
	mgr.db.AutoMigrate(&Record{})
	mgr.db.AutoMigrate(&FTPUser{})
	mgr.db.AutoMigrate(&FtpConfig{})
	mgr.db.AutoMigrate(&DatabaseUser{})
	mgr.db.AutoMigrate(&DatabaseAssociation{})
	mgr.db.AutoMigrate(&Database{})
	mgr.db.AutoMigrate(&IP{})
	mgr.db.AutoMigrate(&Domain{})
	mgr.db.AutoMigrate(&Account{})
	mgr.db.AutoMigrate(&List{})
	mgr.db.AutoMigrate(&Autoresponder{})

	// Add default User
	var usr User

	if mgr.db.First(&usr, "email = ?", "admin@admin.com").RecordNotFound() {
		usr.Email = "admin@admin.com"
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		usr.Password = string(hash)
		mgr.db.Create(&usr)
	}

	var ftpConfig FtpConfig
	if mgr.db.First(&ftpConfig, "id = 1").RecordNotFound() {
		mgr.UpdateFtpConfig(0, 0,1)
	}
}

