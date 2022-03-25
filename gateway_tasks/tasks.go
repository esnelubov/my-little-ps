package gateway_tasks

import "my-little-ps/cache_maps/currency"

type Tasks struct {
	CurrenciesCache *currency.CacheMap
}

func New(currencyCache *currency.CacheMap) *Tasks {
	return &Tasks{
		CurrenciesCache: currencyCache,
	}
}
