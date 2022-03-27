package operation

import (
	"github.com/gofiber/fiber/v2/utils"
	"gorm.io/gorm"
	"my-little-ps/common/constants"
	"my-little-ps/common/database"
	"my-little-ps/common/logger"
	"my-little-ps/common/models"
	"strings"
	"time"
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

	err = c.DB.Find(&records, map[string]interface{}{"target_wallet_id = ?": walletId, "created_at >= ?": from, "created_at < ?": to, "offset": offset, "limit": limit}, c.withTransaction)
	return
}

func (c *Controller) GetOutOperations(walletId uint, from time.Time, to time.Time, offset int, limit int) (records []*models.OutOperation, err error) {
	c.logger.Debugf("Getting OUT operations for wallet %d, from %v to %v, with offset %d and limit %d", walletId, from, to, offset, limit)

	err = c.DB.Find(&records, map[string]interface{}{"origin_wallet_id = ?": walletId, "created_at >= ?": from, "created_at < ?": to, "offset": offset, "limit": limit}, c.withTransaction)
	return
}

func (c *Controller) GetNewInOperations(limit int) (records []*models.InOperation, err error) {
	c.logger.Debugf("Getting all new IN operations (limit %d)", limit)
	err = c.DB.Find(&records, map[string]interface{}{"status = ?": constants.OpStatusNew, "limit": limit}, c.withTransaction)
	return
}

func (c *Controller) GetNewOutOperations(limit int) (records []*models.OutOperation, err error) {
	c.logger.Debugf("Getting all new OUT operations (limit %d)", limit)
	err = c.DB.Find(&records, map[string]interface{}{"status = ?": constants.OpStatusNew, "limit": limit}, c.withTransaction)
	return
}

func (c *Controller) GetNewInOperationsForWallets(wallets []uint, limit int) (records []*models.InOperation, err error) {
	c.logger.Debugf("Getting new IN operations (limit %d) for given wallets", limit)
	if len(wallets) == 0 {
		return []*models.InOperation{}, nil
	}

	err = c.DB.Find(&records, map[string]interface{}{"target_wallet_id IN ?": wallets, "status = ?": constants.OpStatusNew, "limit": limit}, c.withTransaction)
	return
}

func (c *Controller) GetNewOutOperationsForWallets(wallets []uint, limit int) (records []*models.OutOperation, err error) {
	c.logger.Debugf("Getting new OUT operations (limit %d) for given wallets", limit)
	if len(wallets) == 0 {
		return []*models.OutOperation{}, nil
	}

	err = c.DB.Find(&records, map[string]interface{}{"origin_wallet_id IN ?": wallets, "status = ?": constants.OpStatusNew, "limit": limit}, c.withTransaction)
	return
}

func (c *Controller) CountOperationsPerWallet(from time.Time) (counts map[int64]int64, err error) {
	c.logger.Debugf("Counting ALL operations from %v per wallet", from)

	var (
		sqlIn   = "SELECT target_wallet_id AS wallet_id, COUNT(*) AS op_count FROM {Table} WHERE created_at > ? GROUP BY target_wallet_id"
		sqlOut  = "SELECT origin_wallet_id AS wallet_id, COUNT(*) AS op_count FROM {Table} WHERE created_at > ? GROUP BY origin_wallet_id"
		results []map[string]interface{}
	)

	counts = make(map[int64]int64)

	countResults := func(results []map[string]interface{}) {
		for _, r := range results {
			count, ok := counts[r["wallet_id"].(int64)]
			if !ok {
				count = 0
			}

			counts[r["wallet_id"].(int64)] = r["op_count"].(int64) + count
		}
	}

	results = []map[string]interface{}{}
	sqlIn = strings.Replace(sqlIn, "{Table}", c.DB.TableName("in_operations"), -1)
	err = c.DB.Raw(&results, sqlIn, from)
	if err != nil {
		return
	}
	countResults(results)

	results = []map[string]interface{}{}
	sqlOut = strings.Replace(sqlOut, "{Table}", c.DB.TableName("out_operations"), -1)
	err = c.DB.Raw(&results, sqlOut, from)
	if err != nil {
		return
	}
	countResults(results)

	return
}
