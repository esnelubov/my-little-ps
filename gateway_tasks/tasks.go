package gateway_tasks

import (
	"my-little-ps/cache_maps/currency"
	"my-little-ps/logger"
)

type Tasks struct {
	logger          *logger.Log
	CurrenciesCache *currency.CacheMap
}

func New(logger *logger.Log, currencyCache *currency.CacheMap) *Tasks {
	return &Tasks{
		logger:          logger,
		CurrenciesCache: currencyCache,
	}
}
