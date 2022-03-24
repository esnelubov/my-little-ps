package currency

import (
	"fmt"
	cc "my-little-ps/controllers/currency"
	"strconv"
)

type API struct {
	currencyController *cc.Controller
}

func New(currencyController *cc.Controller) *API {
	return &API{
		currencyController: currencyController,
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
		amount      int64
		hasCurrency bool
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

	hasCurrency, err = a.currencyController.HasCurrency(req.From)
	if err != nil {
		return
	}

	if !hasCurrency {
		return nil, fmt.Errorf("currency is not allowed: %s", req.From)
	}

	hasCurrency, err = a.currencyController.HasCurrency(req.To)
	if err != nil {
		return
	}

	if !hasCurrency {
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

	resp.Amount, err = a.currencyController.Convert(reqDec.From, reqDec.To, reqDec.Amount)
	if err != nil {
		return
	}

	return
}
