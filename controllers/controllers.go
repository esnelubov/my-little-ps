package controllers

import (
	"my-little-ps/controllers/operation"
	"my-little-ps/controllers/wallet"
	"my-little-ps/database"
)

var (
	Operation *operation.Controller
	Wallet    *wallet.Controller
)

func Setup(DB *database.DB) {
	Operation = operation.NewController(DB)
	Wallet = wallet.NewController(DB)
}
