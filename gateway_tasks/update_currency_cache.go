package gateway_tasks

func (t *Tasks) UpdateCurrencyCache() {
	err := t.CurrenciesCache.UpdateCurrencies()

	if err != nil {
		t.logger.Errorf("Can't update currencies from db: %s", err.Error())
	}
}
