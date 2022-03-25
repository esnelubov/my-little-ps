package controllers

import (
	"my-little-ps/controllers/currency"
	"my-little-ps/controllers/operation"
	"my-little-ps/controllers/wallet"
	"my-little-ps/database"
	"my-little-ps/logger"
)

var (
	Operation *operation.Controller
	Wallet    *wallet.Controller
	Currency  *currency.Controller
)

func Setup(logger *logger.Log, db *database.DB) {
	Operation = operation.NewController(logger, db)
	Wallet = wallet.NewController(logger, db)
	Currency = currency.NewController(logger, db)
}
