package api

import (
	"my-little-ps/api/currency"
	"my-little-ps/api/report"
	"my-little-ps/api/transaction"
	"my-little-ps/api/wallet"
	ch "my-little-ps/cache_maps/currency"
	cc "my-little-ps/controllers/currency"
	oc "my-little-ps/controllers/operation"
	wc "my-little-ps/controllers/wallet"
	"my-little-ps/logger"
)

var (
	Wallet      *wallet.API
	Transaction *transaction.API
	Currency    *currency.API
	Report      *report.API
)

func Setup(logger *logger.Log, walletController *wc.Controller, operationController *oc.Controller, currencyController *cc.Controller, currencyCache *ch.CacheMap) {
	Currency = currency.New(logger, currencyController, currencyCache)
	Wallet = wallet.New(logger, walletController, currencyCache)
	Transaction = transaction.New(logger, walletController, operationController, currencyCache)
	Report = report.New(logger, walletController, operationController, currencyCache)
}
