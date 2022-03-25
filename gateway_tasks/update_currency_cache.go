package gateway_tasks

func (t *Tasks) UpdateCurrencyCache() {
	_ = t.CurrenciesCache.UpdateCurrencies()
}
