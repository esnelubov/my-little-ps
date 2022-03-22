package models

import "gorm.io/gorm"

type OutOperation struct {
	gorm.Model

	OperationId    string
	TransactionId  string
	OriginWalletId uint
	TargetWalletId uint
	Amount         int64
	Currency       string
	Status         string
}
