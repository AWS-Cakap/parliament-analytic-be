package models

import "gorm.io/gorm"

type Party struct {
	gorm.Model
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"type:varchar(100);unique;not null"`
	EstablishDate string `gorm:"type:date"`
	Logo          string `gorm:"type:varchar(255)"`
}
