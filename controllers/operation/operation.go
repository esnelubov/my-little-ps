package operation

import (
	"github.com/gofiber/fiber/v2/utils"
	"gorm.io/gorm"
	"my-little-ps/constants"
	"my-little-ps/database"
	"my-little-ps/logger"
	"my-little-ps/models"
	"time"
)

type Controller struct {
	logger *logger.Log
	DB     *database.DB
}

func NewController(logger *logger.Log, db *database.DB) *Controller {
	return &Controller{
		logger: logger,
		DB:     db,
	}
}

// NewExternalIn receive money from external sources
func (c *Controller) NewExternalIn(transactionId string, targetWalletID uint, amount int64, currency string) error {
	c.logger.Debugf("Adding new operation for payment from an outside source. "+
		"Transaction: %s, target wallet: %d, amount: %d, currency: %s",
		transactionId, targetWalletID, amount, currency)

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
	c.logger.Debugf("Adding new operation for payment from our wallet. "+
		"Transaction: %s, origin wallet: %d, target wallet: %d, amount: %d, currency: %s",
		transactionId, originWalletID, targetWalletID, amount, currency)

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
	c.logger.Debugf("Adding new operation for payout to our wallet. "+
		"Transaction: %s, origin wallet: %d, target wallet: %d, amount: %d, currency: %s",
		transactionId, originWalletID, targetWalletID, amount, currency)

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
	c.logger.Debugf("Getting IN operations for wallet %d, from %v to %v, with offset %d and limit %d", walletId, from, to, offset, limit)

	err = c.DB.Find(&records, map[string]interface{}{"target_wallet_id = ?": walletId, "created_at >= ?": from, "created_at < ?": to, "offset": offset, "limit": limit})
	return
}

func (c *Controller) GetOutOperations(walletId uint, from time.Time, to time.Time, offset int, limit int) (records []*models.OutOperation, err error) {
	c.logger.Debugf("Getting OUT operations for wallet %d, from %v to %v, with offset %d and limit %d", walletId, from, to, offset, limit)
	
	err = c.DB.Find(&records, map[string]interface{}{"origin_wallet_id = ?": walletId, "created_at >= ?": from, "created_at < ?": to, "offset": offset, "limit": limit})
	return
}
