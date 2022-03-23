package currency

import (
	"my-little-ps/database"
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

		record.USDRate = int64(rate * 1000000)

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
