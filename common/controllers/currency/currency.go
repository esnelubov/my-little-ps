package currency

import (
	"my-little-ps/common/constants"
	"my-little-ps/common/database"
	"my-little-ps/common/logger"
	"my-little-ps/common/models"
)

type Controller struct {
	logger          *logger.Log
	DB              *database.DB
	withTransaction bool
}

func NewController(logger *logger.Log, db *database.DB) *Controller {
	return &Controller{
		logger: logger,
		DB:     db,
	}
}

func (c *Controller) WithTransaction(db *database.DB) *Controller {
	return &Controller{
		logger:          c.logger,
		DB:              db,
		withTransaction: true,
	}
}

func (c *Controller) UpdateFromFloat(rates map[string]float64) (err error) {
	c.logger.Debugf("Updating the currencies table with values: %+v", rates)

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

	err = c.DB.Last(record, map[string]interface{}{"name = ?": name}, c.withTransaction)

	if err == c.DB.ErrRecordNotFound {
		record = &models.Currency{
			Name: name,
		}

		err = nil
	}

	return
}

func (c *Controller) GetAllRecords() (records []*models.Currency, err error) {
	c.logger.Debug("Getting all rates from the currencies table")

	err = c.DB.Find(&records, map[string]interface{}{}, c.withTransaction)
	return
}
