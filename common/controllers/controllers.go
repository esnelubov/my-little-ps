package controllers

import (
	"my-little-ps/common/controllers/currency"
	"my-little-ps/common/controllers/operation"
	"my-little-ps/common/controllers/wallet"
	"my-little-ps/common/database"
	"my-little-ps/common/logger"
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
