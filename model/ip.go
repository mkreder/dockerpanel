package model

import "github.com/jinzhu/gorm"

type IP struct {
	gorm.Model
	Valor  string
	BdID uint
}

