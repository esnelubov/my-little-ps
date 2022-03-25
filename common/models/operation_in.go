package models

import "gorm.io/gorm"

type InOperation struct {
	gorm.Model

	OperationId    string
	TransactionId  string
	OriginWalletId uint
	TargetWalletId uint
	Amount         int64
	Currency       string
	Status         string
}
