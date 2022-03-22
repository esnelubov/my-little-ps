package currency

import (
	cc "my-little-ps/controllers/currency"
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
