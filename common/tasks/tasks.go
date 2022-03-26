package tasks

import (
	"my-little-ps/common/cache_maps/currency"
	"my-little-ps/common/logger"
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
