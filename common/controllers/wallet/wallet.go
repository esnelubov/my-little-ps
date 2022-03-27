package wallet

import (
	"fmt"
	"my-little-ps/common/database"
	"my-little-ps/common/logger"
	"my-little-ps/common/models"
	"strings"
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

func (c *Controller) AddWallet(wallet *models.Wallet) error {
	c.logger.Debugf("Adding wallet %+v", wallet)

	return c.DB.Create(wallet)
}

func (c *Controller) HasWallet(id uint) (bool, error) {
	c.logger.Debugf("Checking if wallet %d exists", id)

	return c.DB.Has(&models.Wallet{}, map[string]interface{}{"id = ?": id})
}

func (c *Controller) WalletsMustExist(ids ...uint) (err error) {
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

func (c *Controller) GetWalletsWithIds(ids []uint) (records []*models.Wallet, err error) {
	c.logger.Debug("Getting wallets with given ids")

	if len(ids) == 0 {
		return []*models.Wallet{}, nil
	}

	err = c.DB.Find(&records, map[string]interface{}{"id IN ?": ids}, c.withTransaction)
	return
}

func (c *Controller) GetWalletsForWorker(number int) (records []*models.Wallet, err error) {
	c.logger.Debugf("Getting wallets for worker %d", number)

	err = c.DB.Find(&records, map[string]interface{}{"worker = ?": number}, c.withTransaction)
	return
}

func (c *Controller) GetAllWallets(offset int, limit int) (records []*models.Wallet, err error) {
	c.logger.Debugf("Getting all wallets (offset %d, limit %d)", offset, limit)

	err = c.DB.Find(&records, map[string]interface{}{"offset": offset, "limit": limit}, c.withTransaction)
	return
}

func (c *Controller) CountWallets() (count int64, err error) {
	c.logger.Debug("Counting all wallets")
	var (
		sql = "SELECT COUNT(*) AS wallet_count FROM {Table}"
	)

	sql = strings.Replace(sql, "{Table}", c.DB.TableName("wallets"), -1)
	err = c.DB.Raw(&count, sql)
	return
}
