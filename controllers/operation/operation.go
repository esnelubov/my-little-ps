package operation

import (
	"github.com/gofiber/fiber/v2/utils"
	"gorm.io/gorm"
	"my-little-ps/constants"
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

// ExternalIn receive money from external sources
func (c *Controller) ExternalIn(transactionId string, targetWalletID uint, amount int64, currency string) error {
	newOp := &models.InOperation{
		Model:          gorm.Model{},
		OperationId:    utils.UUIDv4(),
		TransactionId:  transactionId,
		OriginWalletId: 0,
		TargetWalletId: targetWalletID,
		Amount:         amount,
		Currency:       currency,
		Status:         constants.OpStatusNew,
	}

	return c.DB.Create(newOp)
}

// InternalIn receive money to wallet
func (c *Controller) InternalIn(transactionId string, originWalletID uint, targetWalletID uint, amount int64, currency string) error {
	newOp := &models.InOperation{
		Model:          gorm.Model{},
		OperationId:    utils.UUIDv4(),
		TransactionId:  transactionId,
		OriginWalletId: originWalletID,
		TargetWalletId: targetWalletID,
		Amount:         amount,
		Currency:       currency,
		Status:         constants.OpStatusNew,
	}

	return c.DB.Create(newOp)
}

// InternalOut send money from wallet
func (c *Controller) InternalOut(transactionId string, originWalletID uint, targetWalletID uint, amount int64, currency string) error {
	newOp := &models.OutOperation{
		Model:          gorm.Model{},
		OperationId:    utils.UUIDv4(),
		TransactionId:  transactionId,
		OriginWalletId: originWalletID,
		TargetWalletId: targetWalletID,
		Amount:         amount,
		Currency:       currency,
		Status:         constants.OpStatusNew,
	}

	return c.DB.Create(newOp)
}
