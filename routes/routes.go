package routes

import (
	"github.com/gofiber/fiber/v2"
	"my-little-ps/api"
	"my-little-ps/api/currency"
	"my-little-ps/api/transaction"
	"my-little-ps/api/wallet"
)

type ResponseData struct {
	Payload interface{} `json:"payload"`
}

type PayloadError struct {
	Error string `json:"error"`
}

func AddWallet(c *fiber.Ctx) (err error) {
	var (
		req  *wallet.AddWalletRequest
		resp *wallet.AddWalletResponse
	)

	if err = c.BodyParser(&req); err != nil {
		return err
	}

	resp, err = api.Wallet.AddWallet(req)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(&ResponseData{Payload: resp})
}

func ReceiveAmount(c *fiber.Ctx) (err error) {
	var (
		req  *transaction.ReceiveAmountRequest
		resp *transaction.ReceiveAmountResponse
	)

	if err = c.BodyParser(&req); err != nil {
		return err
	}

	resp, err = api.Transaction.ReceiveAmount(req)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(&ResponseData{Payload: resp})
}

func TransferAmount(c *fiber.Ctx) (err error) {
	var (
		req  *transaction.TransferAmountRequest
		resp *transaction.TransferAmountResponse
	)

	if err = c.BodyParser(&req); err != nil {
		return err
	}

	resp, err = api.Transaction.TransferAmount(req)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(&ResponseData{Payload: resp})
}

func UpdateCurrencies(c *fiber.Ctx) (err error) {
	var (
		req  *currency.UpdateCurrenciesRequest
		resp *currency.UpdateCurrenciesResponse
	)

	if err = c.BodyParser(&req); err != nil {
		return err
	}

	resp, err = api.Currency.UpdateCurrencies(req)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(&ResponseData{Payload: resp})
}
