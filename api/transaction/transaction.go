package transaction

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/utils"
	"my-little-ps/constants"
	oc "my-little-ps/controllers/operation"
	wc "my-little-ps/controllers/wallet"
)

type API struct {
	walletController    *wc.Controller
	operationController *oc.Controller
}

func New(walletController *wc.Controller, operationController *oc.Controller) *API {
	return &API{
		walletController:    walletController,
		operationController: operationController,
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

func (a *API) ValidateReceiveAmount(req *ReceiveAmountRequest) error {
	if req.Amount <= 0 || req.Currency == "" {
		return errors.New("all fields should be filled")
	}

	err := a.walletController.CheckWallets(req.WalletId)
	if err != nil {
		return err
	}

	if !constants.AllowedCurrencies.Has(req.Currency) {
		return fmt.Errorf("allowed currencies: %v", constants.AllowedCurrencies)
	}

	return nil
}

func (a *API) ReceiveAmount(req *ReceiveAmountRequest) (resp *ReceiveAmountResponse, err error) {
	if err = a.ValidateReceiveAmount(req); err != nil {
		return nil, err
	}

	transactionId := utils.UUIDv4()

	err = a.operationController.ExternalIn(transactionId, req.WalletId, req.Amount, req.Currency)
	if err != nil {
		return nil, err
	}

	return &ReceiveAmountResponse{
		TransactionId: transactionId,
	}, nil
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

	if !constants.AllowedCurrencies.Has(req.Currency) {
		return fmt.Errorf("allowed currencies: %v", constants.AllowedCurrencies)
	}

	return nil
}

func (a *API) TransferAmount(req *TransferAmountRequest) (resp *TransferAmountResponse, err error) {
	if err = a.ValidateTransferAmount(req); err != nil {
		return nil, err
	}

	transactionId := utils.UUIDv4()

	err = a.operationController.InternalOut(transactionId, req.OriginWalletId, req.TargetWalletId, req.Amount, req.Currency)
	if err != nil {
		return nil, err
	}

	return &TransferAmountResponse{
		TransactionId: transactionId,
	}, nil
}
