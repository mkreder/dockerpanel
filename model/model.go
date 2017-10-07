package model

import "github.com/jinzhu/gorm"

type Web struct {
	gorm.Model
	Dominio string
	CGI bool
	SSL bool
	Python bool
	Ruby bool
	Perl bool
	PHP bool
	Status int
}

type User struct {
	gorm.Model
	Email string
	Password string
}