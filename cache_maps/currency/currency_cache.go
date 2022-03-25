package currency

import (
	"fmt"
	"my-little-ps/constants"
	cc "my-little-ps/controllers/currency"
	"my-little-ps/logger"
	"my-little-ps/models"
	"sync"
)

type CacheMap struct {
	logger             *logger.Log
	m                  sync.Map
	currencyController *cc.Controller
}

func New(logger *logger.Log, currencyController *cc.Controller) *CacheMap {
	return &CacheMap{
		logger:             logger,
		m:                  sync.Map{},
		currencyController: currencyController,
	}
}

func (c *CacheMap) UpdateCurrencies() (err error) {
	c.logger.Debug("Fetching currency rates to app cache")

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
	c.logger.Debugf("Checking if currency %s exists", name)

	_, ok := c.m.Load(name)

	return ok
}

func (c *CacheMap) getRate(currency string) (rate int64, err error) {
	c.logger.Debugf("Getting rate for %s", currency)

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

	c.logger.Debugf("Rate of %s is %d", currency, rate)

	return
}

func (c *CacheMap) Convert(from string, to string, amount int64) (result int64, err error) {
	c.logger.Debugf("Converting %d %s to %s", amount, from, to)

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

	c.logger.Debugf("%d %s is %d %s", amount, from, result, to)

	return result, nil
}
