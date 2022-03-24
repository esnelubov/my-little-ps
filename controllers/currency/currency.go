package currency

import (
	"my-little-ps/constants"
	"my-little-ps/database"
	"my-little-ps/helpers/datatypes/set_of_string"
	"my-little-ps/models"
)

type Controller struct {
	DB *database.DB
}

func NewController(db *database.DB) *Controller {
	return &Controller{
		DB: db,
	}
}

func (c *Controller) UpdateFromFloat(rates map[string]float64) (err error) {
	var (
		record *models.Currency
	)

	for name, rate := range rates {
		record, err = c.GetOrInitRateRecord(name)
		if err != nil {
			return
		}

		record.USDRate = int64(rate * constants.RateMultiplier)

		if err = c.DB.Save(record); err != nil {
			return
		}
	}

	return nil
}

func (c *Controller) GetOrInitRateRecord(name string) (record *models.Currency, err error) {
	record = &models.Currency{}

	err = c.DB.Last(record, map[string]interface{}{"name = ?": name})

	if err == c.DB.ErrRecordNotFound {
		record = &models.Currency{
			Name: name,
		}

		err = nil
	}

	return
}

func (c *Controller) GetAllowedCurrencies() (set set_of_string.Type, err error) {
	currencies := []string{}

	err = c.DB.Raw(&currencies, "SELECT name FROM currencies")
	if err != nil {
		return
	}

	set = set_of_string.New(currencies...)

	return
}

func (c *Controller) HasCurrency(name string) (result bool, err error) {
	records := []*models.Currency{}

	err = c.DB.Find(&records, map[string]interface{}{"name = ?": name})
	if err != nil {
		return
	}

	return len(records) != 0, nil
}

func (c *Controller) Convert(from string, to string, amount int64) (result int64, err error) {
	var (
		records     = []*models.Currency{}
		fromUSDRate int64
		toUSDRate   int64
	)

	if from == to {
		return amount, nil
	}

	err = c.DB.Find(&records, map[string]interface{}{"name IN ?": []string{from, to}})
	if err != nil {
		return
	}

	for _, r := range records {
		if r.Name == from {
			fromUSDRate = r.USDRate
		}

		if r.Name == to {
			toUSDRate = r.USDRate
		}
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
