package models

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model

	Name     string
	Country  string
	City     string
	Currency string
	Balance  int64
	Worker   int32
}
