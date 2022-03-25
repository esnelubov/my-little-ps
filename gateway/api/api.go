package api

import (
	ch "my-little-ps/common/cache_maps/currency"
	cc "my-little-ps/common/controllers/currency"
	oc "my-little-ps/common/controllers/operation"
	wc "my-little-ps/common/controllers/wallet"
	"my-little-ps/common/logger"
	"my-little-ps/gateway/api/currency"
	"my-little-ps/gateway/api/report"
	"my-little-ps/gateway/api/transaction"
	"my-little-ps/gateway/api/wallet"
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
