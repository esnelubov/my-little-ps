package models

import "gorm.io/gorm"

type Operation struct {
	gorm.Model

	UID      string
	WalletId uint
	Type     string
	Amount   int64
	Status   string
}
