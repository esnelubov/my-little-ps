package wallet

import (
	"fmt"
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
	return c.DB.Create(&wallet)
}

func (c *Controller) HasWallet(id uint) (bool, error) {
	return c.DB.Has(&models.Wallet{}, map[string]interface{}{"id = ?": id})
}

func (c *Controller) CheckWallets(ids ...uint) (err error) {
	var hasWallet bool

	for _, id := range ids {
		hasWallet, err = c.HasWallet(id)

		if err != nil {
			return err
		}

		if !hasWallet {
			return fmt.Errorf("no wallet with id: %d", id)
		}
	}

	return nil
}
