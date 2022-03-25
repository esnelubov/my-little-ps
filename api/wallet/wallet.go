package wallet

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"my-little-ps/cache_maps/currency"
	wc "my-little-ps/controllers/wallet"
	"my-little-ps/logger"
	"my-little-ps/models"
)

type API struct {
	logger           *logger.Log
	walletController *wc.Controller
	currencyCache    *currency.CacheMap
}

func New(logger *logger.Log, walletController *wc.Controller, currencyCache *currency.CacheMap) *API {
	return &API{
		logger:           logger,
		walletController: walletController,
		currencyCache:    currencyCache,
	}
}

type AddWalletRequest struct {
	Name     string
	Country  string
	City     string
	Currency string
}

func (a *API) ValidateAddWallet(req *AddWalletRequest) (err error) {
	if req.Name == "" || req.Country == "" || req.City == "" || req.Currency == "" {
		return errors.New("all fields should be filled")
	}

	if !a.currencyCache.HasCurrency(req.Currency) {
		return fmt.Errorf("currency is not allowed")
	}

	return nil
}

type AddWalletResponse struct {
	WalletID uint
}

func (a *API) AddWallet(req *AddWalletRequest) (resp *AddWalletResponse, err error) {
	a.logger.Debugf("Received the AddWallet request: %+v", req)

	if err = a.ValidateAddWallet(req); err != nil {
		return nil, err
	}

	newWallet := models.Wallet{
		Model:    gorm.Model{},
		Name:     req.Name,
		Country:  req.Country,
		City:     req.City,
		Currency: req.Currency,
		Balance:  0,
	}

	err = a.walletController.AddWallet(&newWallet)
	if err != nil {
		return nil, err
	}

	resp = &AddWalletResponse{
		WalletID: newWallet.ID,
	}

	a.logger.Debugf("Replying to the AddWallet request: %+v, with: %+v", req, resp)

	return
}
