package api

import (
	"my-little-ps/api/wallet"
	wc "my-little-ps/controllers/wallet"
)

var (
	Wallet *wallet.API
)

func Setup(walletController *wc.Controller) {
	Wallet = wallet.New(walletController)
}
