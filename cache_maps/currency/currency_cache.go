package currency

import (
	"fmt"
	"my-little-ps/constants"
	cc "my-little-ps/controllers/currency"
	"my-little-ps/models"
	"sync"
)

type CacheMap struct {
	m                  sync.Map
	currencyController *cc.Controller
}

func New(currencyController *cc.Controller) *CacheMap {
	return &CacheMap{
		m:                  sync.Map{},
		currencyController: currencyController,
	}
}

func (c *CacheMap) UpdateCurrencies() (err error) {
	var (
		records []*models.Currency
	)

	records, err = c.currencyController.GetAllRecords()
	if err != nil {
		return
	}

	for _, r := range records {
		c.m.Store(r.Name, r.USDRate)
	}

	return
}

func (c *CacheMap) HasCurrency(name string) bool {
	_, ok := c.m.Load(name)

	return ok
}

func (c *CacheMap) getRate(currency string) (rate int64, err error) {
	var (
		rateIface interface{}
		ok        bool
	)

	rateIface, ok = c.m.Load(currency)
	if !ok {
		return 0, fmt.Errorf("no rate for currency %s", currency)
	}

	rate, ok = rateIface.(int64)
	if !ok {
		return 0, fmt.Errorf("invalid rate for currency %s", currency)
	}

	return
}

func (c *CacheMap) Convert(from string, to string, amount int64) (result int64, err error) {
	var (
		fromUSDRate int64
		toUSDRate   int64
	)

	if from == to {
		return amount, nil
	}

	fromUSDRate, err = c.getRate(from)
	if err != nil {
		return
	}

	toUSDRate, err = c.getRate(to)
	if err != nil {
		return
	}

	if from == constants.USD {
		result = amount * constants.RateMultiplier / toUSDRate
	} else if to == constants.USD {
		result = amount * fromUSDRate / constants.RateMultiplier
	} else {
		result = amount * fromUSDRate / toUSDRate
	}

	return result, nil
}
