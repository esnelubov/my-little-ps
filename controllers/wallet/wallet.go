package wallet

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

func (c *Controller) AddWallet(wallet *models.Wallet) error {
	return c.DB.Create(wallet)
}
