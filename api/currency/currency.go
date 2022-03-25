package currency

import (
	"errors"
	"fmt"
	"my-little-ps/cache_maps/currency"
	cc "my-little-ps/controllers/currency"
	"my-little-ps/logger"
	"strconv"
)

type API struct {
	logger             *logger.Log
	currencyController *cc.Controller
	currencyCache      *currency.CacheMap
}

func New(logger *logger.Log, currencyController *cc.Controller, currencyCache *currency.CacheMap) *API {
	return &API{
		logger:             logger,
		currencyController: currencyController,
		currencyCache:      currencyCache,
	}
}

type UpdateCurrenciesRequest struct {
	Rates map[string]float64
}

type UpdateCurrenciesResponse struct {
}

func (a *API) ValidateUpdateCurrencies(req *UpdateCurrenciesRequest) error {
	if len(req.Rates) == 0 {
		return errors.New("rates field must be set")
	}

	return nil
}

func (a *API) UpdateCurrencies(req *UpdateCurrenciesRequest) (resp *UpdateCurrenciesResponse, err error) {
	a.logger.Debugf("Received the UpdateCurrencies request: %+v", req)

	if err = a.ValidateUpdateCurrencies(req); err != nil {
		return nil, err
	}

	err = a.currencyController.UpdateFromFloat(req.Rates)
	if err != nil {
		return nil, err
	}

	resp = &UpdateCurrenciesResponse{}

	a.logger.Debugf("Replying to the UpdateCurrencies request: %+v, with: %+v", req, resp)

	return
}

type ConvertAmountRequest struct {
	Amount string
	From   string
	To     string
}

type ConvertAmountDecodedRequest struct {
	Amount int64
	From   string
	To     string
}

type ConvertAmountResponse struct {
	Amount   int64
	Currency string
}

func (a *API) ValidateConvertAmount(req *ConvertAmountRequest) (resp *ConvertAmountDecodedRequest, err error) {
	var (
		amount int64
	)

	resp = &ConvertAmountDecodedRequest{
		From: req.From,
		To:   req.To,
	}

	amount, err = strconv.ParseInt(req.Amount, 10, 64)
	if err != nil {
		return
	}

	resp.Amount = amount

	if !a.currencyCache.HasCurrency(req.From) {
		return nil, fmt.Errorf("currency is not allowed: %s", req.From)
	}

	if !a.currencyCache.HasCurrency(req.To) {
		return nil, fmt.Errorf("currency is not allowed: %s", req.To)
	}

	return
}

func (a *API) ConvertAmount(req *ConvertAmountRequest) (resp *ConvertAmountResponse, err error) {
	a.logger.Debugf("Received the ConvertAmount request: %+v", req)

	var (
		reqDec *ConvertAmountDecodedRequest
	)

	reqDec, err = a.ValidateConvertAmount(req)
	if err != nil {
		return
	}

	resp = &ConvertAmountResponse{
		Currency: reqDec.To,
	}

	resp.Amount, err = a.currencyCache.Convert(reqDec.From, reqDec.To, reqDec.Amount)
	if err != nil {
		return
	}

	a.logger.Debugf("Replying to the ConvertAmount request: %+v, with: %+v", req, resp)

	return
}
