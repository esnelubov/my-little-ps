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

func (c *Controller) UpdateFromFloat(rates map[string]float64) error {
	var (
		record  models.Currency
		records = make([]interface{}, len(rates))
	)

	for name, rate := range rates {
		record = models.Currency{
			Name:    name,
			USDRate: int64(rate * 1000000),
		}

		records = append(records, record)
	}

	return c.DB.SaveTx(records...)
}
