package report

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/utils"
	"my-little-ps/cache_maps/currency"
	"my-little-ps/constants"
	oc "my-little-ps/controllers/operation"
	wc "my-little-ps/controllers/wallet"
	"my-little-ps/logger"
	"my-little-ps/models"
	"os"
	"sort"
	"strconv"
	"time"
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

type GetOperationsRequest struct {
	WalletId string
	Offset   string
	Limit    string
	From     string
	To       string
}

type GetOperationsDecodedRequest struct {
	WalletId uint
	Offset   int
	Limit    int
	From     time.Time
	To       time.Time
}

type GetOperationsResponse struct {
	Operations []*Operation
	Limit      int
	Offset     int
	From       string
	To         string
}

type Operation struct {
	Type      string
	Amount    int64
	Currency  string
	CreatedAt time.Time
}

func (a *API) ParseGetOperations(req *GetOperationsRequest) (resp *GetOperationsDecodedRequest, err error) {
	var (
		walletId uint64
		offset   int64
		limit    int64
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

	err = a.walletController.CheckWallets(resp.WalletId)
	if err != nil {
		return
	}

	offset, err = parseInt64(req.Offset, 0)
	if err != nil {
		return
	}
	resp.Offset = int(offset)

	limit, err = parseInt64(req.Limit, 1000)
	if err != nil {
		return
	}
	resp.Limit = int(limit)

	resp.From, err = parseRFC3339Time(req.From, time.Now().AddDate(0, -1, 0))
	if err != nil {
		return
	}

	resp.To, err = parseRFC3339Time(req.To, time.Now())
	if err != nil {
		return
	}

	return resp, nil
}

func parseInt64(val string, def int64) (int64, error) {
	if val == "" {
		return def, nil
	}

	return strconv.ParseInt(val, 10, 64)
}

func parseRFC3339Time(val string, def time.Time) (time.Time, error) {
	if val == "" {
		return def, nil
	}

	return time.Parse(time.RFC3339, val)
}

func (a *API) GetOperations(req *GetOperationsRequest) (resp *GetOperationsResponse, err error) {
	a.logger.Debugf("Received the GetOperations request: %+v", req)

	var (
		reqDec *GetOperationsDecodedRequest
		ops    []*Operation
	)

	if reqDec, err = a.ParseGetOperations(req); err != nil {
		return nil, err
	}

	ops, err = a.getSortedOperations(reqDec.WalletId, reqDec.From, reqDec.To, reqDec.Offset, reqDec.Limit)
	if err != nil {
		return
	}

	resp = &GetOperationsResponse{
		Operations: ops,
		Limit:      reqDec.Limit,
		Offset:     reqDec.Offset,
		From:       reqDec.From.Format(time.RFC3339),
		To:         reqDec.To.Format(time.RFC3339),
	}

	a.logger.Debugf("Replying to the GetOperations request: %+v, with %d operations", req, len(resp.Operations))

	return
}

func (a *API) GetOperationsCSV(req *GetOperationsRequest) (filename string, err error) {
	a.logger.Debugf("Received the GetOperationsCSV request: %+v", req)

	var (
		reqDec   *GetOperationsDecodedRequest
		opsChunk []*Operation
		f        *os.File
	)

	if reqDec, err = a.ParseGetOperations(req); err != nil {
		return
	}

	filename = utils.UUIDv4()

	f, err = os.CreateTemp("", filename)
	defer f.Close()

	_, err = f.WriteString("type,amount,currency,created_at\n")
	if err != nil {
		return
	}

	for {
		opsChunk, err = a.getSortedOperations(reqDec.WalletId, reqDec.From, reqDec.To, reqDec.Offset, reqDec.Limit)
		if err != nil {
			return
		}

		if len(opsChunk) == 0 {
			break
		}

		for _, op := range opsChunk {
			_, err = f.WriteString(fmt.Sprintf("%s,%d,%s,%s\n", op.Type, op.Amount, op.Currency, op.CreatedAt.Format(time.RFC3339Nano)))
			if err != nil {
				return
			}
		}

		reqDec.Offset += reqDec.Limit
	}

	a.logger.Debugf("Prepared CSV report %s for the GetOperationsCSV request: %+v", f.Name(), req)

	return f.Name(), nil
}

func (a *API) getSortedOperations(walletId uint, from time.Time, to time.Time, offset int, limit int) (ops []*Operation, err error) {
	var (
		inOps  []*models.InOperation
		outOps []*models.OutOperation
	)

	inOps, err = a.operationController.GetInOperations(walletId, from, to, offset, limit)
	if err != nil {
		return
	}

	outOps, err = a.operationController.GetOutOperations(walletId, from, to, offset, limit)
	if err != nil {
		return
	}

	for _, op := range inOps {
		ops = append(ops, &Operation{Type: constants.OpIn, Amount: op.Amount, Currency: op.Currency, CreatedAt: op.CreatedAt})
	}

	for _, op := range outOps {
		ops = append(ops, &Operation{Type: constants.OpOut, Amount: op.Amount, Currency: op.Currency, CreatedAt: op.CreatedAt})
	}

	sort.Slice(ops, func(i, j int) bool {
		return ops[i].CreatedAt.Before(ops[j].CreatedAt)
	})

	return
}

type GetOperationsTotalRequest struct {
	WalletId string
	Currency string
	From     string
	To       string
}

type GetOperationsTotalDecodedRequest struct {
	WalletId uint
	Currency string
	From     time.Time
	To       time.Time
}

type GetOperationsTotalResponse struct {
	Currency  string
	AmountIn  int64
	AmountOut int64
}

func (a *API) ParseGetOperationsTotal(req *GetOperationsTotalRequest) (resp *GetOperationsTotalDecodedRequest, err error) {
	var (
		walletId uint64
	)

	if req.WalletId == "" {
		return nil, errors.New("wallet id should be filled")
	}

	walletId, err = strconv.ParseUint(req.WalletId, 10, 64)
	if err != nil {
		return
	}

	resp = &GetOperationsTotalDecodedRequest{
		WalletId: uint(walletId),
	}

	err = a.walletController.CheckWallets(resp.WalletId)
	if err != nil {
		return
	}

	if req.Currency == "" {
		req.Currency = constants.USD
	}

	if !a.currencyCache.HasCurrency(req.Currency) {
		return nil, errors.New("currency is not allowed")
	}

	resp.Currency = req.Currency

	resp.From, err = parseRFC3339Time(req.From, time.Now().AddDate(0, -1, 0))
	if err != nil {
		return
	}

	resp.To, err = parseRFC3339Time(req.To, time.Now())
	if err != nil {
		return
	}

	return resp, nil
}

func (a *API) GetOperationsTotal(req *GetOperationsTotalRequest) (resp *GetOperationsTotalResponse, err error) {
	a.logger.Debugf("Received the GetOperationsTotal request: %+v", req)

	var (
		reqDec *GetOperationsTotalDecodedRequest
		inOps  []*models.InOperation
		outOps []*models.OutOperation
		amount int64
		offset = 0
		limit  = 1000
	)

	if reqDec, err = a.ParseGetOperationsTotal(req); err != nil {
		return
	}

	resp = &GetOperationsTotalResponse{
		Currency:  reqDec.Currency,
		AmountIn:  0,
		AmountOut: 0,
	}

	offset = 0

	for {
		inOps, err = a.operationController.GetInOperations(reqDec.WalletId, reqDec.From, reqDec.To, offset, limit)
		if err != nil {
			return
		}

		if len(inOps) == 0 {
			break
		}

		for _, op := range inOps {
			amount, err = a.currencyCache.Convert(op.Currency, resp.Currency, op.Amount)
			if err != nil {
				return
			}

			resp.AmountIn += amount
		}

		offset += limit
	}

	offset = 0

	for {
		outOps, err = a.operationController.GetOutOperations(reqDec.WalletId, reqDec.From, reqDec.To, offset, limit)
		if err != nil {
			return
		}

		if len(outOps) == 0 {
			break
		}

		for _, op := range outOps {
			amount, err = a.currencyCache.Convert(op.Currency, resp.Currency, op.Amount)
			if err != nil {
				return
			}

			resp.AmountOut += amount
		}

		offset += limit
	}

	a.logger.Debugf("Replying to the GetOperationsTotal request: %+v, with: %+v", req, resp)

	return
}
