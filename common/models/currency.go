package models

import "gorm.io/gorm"

type Currency struct {
	gorm.Model

	Name    string
	USDRate int64 //x * USDRate / 1000000
}
