package transaction

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/utils"
	"my-little-ps/cache_maps/currency"
	oc "my-little-ps/controllers/operation"
	wc "my-little-ps/controllers/wallet"
	"my-little-ps/logger"
)

type API struct {
	logger              *logger.Log
	walletController    *wc.Controller
	operationController *oc.Controller
	currencyCache       *currency.CacheMap
}

func New(logger *logger.Log, walletController *wc.Controller, operationController *oc.Controller, currencyCache *currency.CacheMap) *API {
	return &API{
		logger:              logger,
		walletController:    walletController,
		operationController: operationController,
		currencyCache:       currencyCache,
	}
}

type ReceiveAmountRequest struct {
	WalletId uint
	Amount   int64
	Currency string
}

type ReceiveAmountResponse struct {
	TransactionId string
}

func (a *API) ValidateReceiveAmount(req *ReceiveAmountRequest) (err error) {
	if req.Amount <= 0 || req.Currency == "" {
		return errors.New("all fields should be filled")
	}

	err = a.walletController.CheckWallets(req.WalletId)
	if err != nil {
		return err
	}

	if !a.currencyCache.HasCurrency(req.Currency) {
		return fmt.Errorf("currency is not allowed")
	}

	return nil
}

func (a *API) ReceiveAmount(req *ReceiveAmountRequest) (resp *ReceiveAmountResponse, err error) {
	a.logger.Debugf("Received the ReceiveAmount request: %+v", req)

	if err = a.ValidateReceiveAmount(req); err != nil {
		return nil, err
	}

	transactionId := utils.UUIDv4()

	err = a.operationController.NewExternalIn(transactionId, req.WalletId, req.Amount, req.Currency)
	if err != nil {
		return nil, err
	}

	resp = &ReceiveAmountResponse{
		TransactionId: transactionId,
	}

	a.logger.Debugf("Replying to the ReceiveAmount request: %+v, with: %+v", req, resp)

	return
}

type TransferAmountRequest struct {
	OriginWalletId uint
	TargetWalletId uint
	Amount         int64
	Currency       string
}

type TransferAmountResponse struct {
	TransactionId string
}

func (a *API) ValidateTransferAmount(req *TransferAmountRequest) error {
	if req.Amount <= 0 || req.Currency == "" {
		return errors.New("all fields should be filled")
	}

	err := a.walletController.CheckWallets(req.OriginWalletId, req.TargetWalletId)
	if err != nil {
		return err
	}

	if !a.currencyCache.HasCurrency(req.Currency) {
		return fmt.Errorf("currency is not allowed")
	}

	return nil
}

func (a *API) TransferAmount(req *TransferAmountRequest) (resp *TransferAmountResponse, err error) {
	a.logger.Debugf("Received the TransferAmount request: %+v", req)

	if err = a.ValidateTransferAmount(req); err != nil {
		return nil, err
	}

	transactionId := utils.UUIDv4()

	err = a.operationController.NewInternalOut(transactionId, req.OriginWalletId, req.TargetWalletId, req.Amount, req.Currency)
	if err != nil {
		return nil, err
	}

	resp = &TransferAmountResponse{
		TransactionId: transactionId,
	}

	a.logger.Debugf("Replying to the TransferAmount request: %+v, with: %+v", req, resp)

	return
}
