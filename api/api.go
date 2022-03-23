package api

import (
	"my-little-ps/api/currency"
	"my-little-ps/api/report"
	"my-little-ps/api/transaction"
	"my-little-ps/api/wallet"
	cc "my-little-ps/controllers/currency"
	oc "my-little-ps/controllers/operation"
	wc "my-little-ps/controllers/wallet"
)

var (
	Wallet      *wallet.API
	Transaction *transaction.API
	Currency    *currency.API
	Report      *report.API
)

func Setup(walletController *wc.Controller, operationController *oc.Controller, currencyController *cc.Controller) {
	Wallet = wallet.New(walletController)
	Transaction = transaction.New(walletController, operationController)
	Currency = currency.New(currencyController)
	Report = report.New(walletController, operationController)
}
