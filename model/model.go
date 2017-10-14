package model

import "github.com/jinzhu/gorm"

type Web struct {
	gorm.Model
	Dominio string
	CGI bool
	SSL bool
	PHP bool
	PHPversion string
	PHPmethod string
	Webserver string
	Estado int
}

type Registro struct {
	gorm.Model
	Tipo string
	Nombre string
	Valor string
	Prioridad string
	ZonaID uint
}

type Zona struct {
	gorm.Model
	Dominio string
	Email string
	Estado int
	Registros []Registro
}


type User struct {
	gorm.Model
	Email string
	Password string
}

type Ftp struct {
	gorm.Model
	Username string
	Password string
	Web Web
}