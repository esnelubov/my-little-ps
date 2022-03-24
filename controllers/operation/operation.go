package operation

import (
	"github.com/gofiber/fiber/v2/utils"
	"gorm.io/gorm"
	"my-little-ps/constants"
	"my-little-ps/database"
	"my-little-ps/models"
	"time"
)

type Controller struct {
	DB *database.DB
}

func NewController(db *database.DB) *Controller {
	return &Controller{
		DB: db,
	}
}

// NewExternalIn receive money from external sources
func (c *Controller) NewExternalIn(transactionId string, targetWalletID uint, amount int64, currency string) error {
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

// NewInternalIn receive money to wallet
func (c *Controller) NewInternalIn(transactionId string, originWalletID uint, targetWalletID uint, amount int64, currency string) error {
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

// NewInternalOut send money from wallet
func (c *Controller) NewInternalOut(transactionId string, originWalletID uint, targetWalletID uint, amount int64, currency string) error {
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

func (c *Controller) GetInOperations(walletId uint, from time.Time, to time.Time, offset int, limit int) (records []*models.InOperation, err error) {
	err = c.DB.Find(&records, map[string]interface{}{"target_wallet_id = ?": walletId, "created_at >= ?": from, "created_at < ?": to, "offset": offset, "limit": limit})
	return
}

func (c *Controller) GetOutOperations(walletId uint, from time.Time, to time.Time, offset int, limit int) (records []*models.OutOperation, err error) {
	err = c.DB.Find(&records, map[string]interface{}{"origin_wallet_id = ?": walletId, "created_at >= ?": from, "created_at < ?": to, "offset": offset, "limit": limit})
	return
}
