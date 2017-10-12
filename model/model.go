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
	Status int
}


type Region struct {
	gorm.Model
	Dominio string
	Email string
	Registros []Registro
}

type Registro struct {
	gorm.Model
	Tipo string
	Nombre string
	Valor string
}

type User struct {
	gorm.Model
	Email string
	Password string
}