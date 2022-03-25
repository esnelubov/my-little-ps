package currency

import (
	"fmt"
	"my-little-ps/cache_maps/currency"
	cc "my-little-ps/controllers/currency"
	"strconv"
)

type API struct {
	currencyController *cc.Controller
	currencyCache      *currency.CacheMap
}

func New(currencyController *cc.Controller, currencyCache *currency.CacheMap) *API {
	return &API{
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
	return nil
}

func (a *API) UpdateCurrencies(req *UpdateCurrenciesRequest) (resp *UpdateCurrenciesResponse, err error) {
	if err = a.ValidateUpdateCurrencies(req); err != nil {
		return nil, err
	}

	err = a.currencyController.UpdateFromFloat(req.Rates)
	if err != nil {
		return nil, err
	}

	return &UpdateCurrenciesResponse{}, nil
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
	Amount int64
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
	var (
		reqDec *ConvertAmountDecodedRequest
	)

	reqDec, err = a.ValidateConvertAmount(req)
	if err != nil {
		return
	}

	resp = &ConvertAmountResponse{}

	resp.Amount, err = a.currencyCache.Convert(reqDec.From, reqDec.To, reqDec.Amount)
	if err != nil {
		return
	}

	return
}
