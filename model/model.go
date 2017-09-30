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
	Status int
}
