package wallet

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	cc "my-little-ps/controllers/currency"
	wc "my-little-ps/controllers/wallet"
	"my-little-ps/models"
)

type API struct {
	walletController   *wc.Controller
	currencyController *cc.Controller
}

func New(walletController *wc.Controller, currencyController *cc.Controller) *API {
	return &API{
		walletController:   walletController,
		currencyController: currencyController,
	}
}

type AddWalletRequest struct {
	Name     string
	Country  string
	City     string
	Currency string
}

func (a *API) ValidateAddWallet(req *AddWalletRequest) (err error) {
	var (
		hasCurrency bool
	)

	if req.Name == "" || req.Country == "" || req.City == "" || req.Currency == "" {
		return errors.New("all fields should be filled")
	}

	hasCurrency, err = a.currencyController.HasCurrency(req.Currency)
	if err != nil {
		return err
	}

	if !hasCurrency {
		return fmt.Errorf("currency is not allowed")
	}

	return nil
}

type AddWalletResponse struct {
	WalletID uint
}

func (a *API) AddWallet(req *AddWalletRequest) (resp *AddWalletResponse, err error) {
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

	return &AddWalletResponse{
		WalletID: newWallet.ID,
	}, nil
}
