package report

import (
	"errors"
	oc "my-little-ps/controllers/operation"
	wc "my-little-ps/controllers/wallet"
	"strconv"
	"time"
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

type GetOperationsRequest struct {
	WalletId string
	Offset   string
	Limit    string
	From     string
	To       string
}

type GetOperationsDecodedRequest struct {
	WalletId uint
	Offset   int64
	Limit    int64
	From     time.Time
	To       time.Time
}

type GetOperationsResponse struct {
	Operations []string
	Limit      int64
	Offset     int64
	From       time.Time
	To         time.Time
}

func (a *API) ParseGetOperations(req *GetOperationsRequest) (resp *GetOperationsDecodedRequest, err error) {
	var (
		walletId uint64
		offset   int64
		limit    int64
		from     time.Time
		to       time.Time
	)

	if req.WalletId == "" {
		return nil, errors.New("wallet id should be filled")
	}

	walletId, err = strconv.ParseUint(req.WalletId, 10, 64)
	if err != nil {
		return
	}

	resp = &GetOperationsDecodedRequest{
		WalletId: uint(walletId),
	}

	offset, err = strconv.ParseInt(req.Offset, 10, 64)
	if err == nil {
		resp.Offset = offset
	}

	limit, err = strconv.ParseInt(req.Limit, 10, 64)
	if err == nil {
		resp.Limit = limit
	}

	from, err = time.Parse(time.RFC3339, req.From)
	if err == nil {
		resp.From = from
	}

	to, err = time.Parse(time.RFC3339, req.To)
	if err == nil {
		resp.To = to
	}

	return resp, nil
}

func (a *API) GetOperations(req *GetOperationsRequest) (resp *GetOperationsResponse, err error) {
	var (
		reqDecoded *GetOperationsDecodedRequest
	)

	if reqDecoded, err = a.ParseGetOperations(req); err != nil {
		return nil, err
	}
}
