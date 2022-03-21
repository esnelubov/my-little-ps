package routes

import (
	"github.com/gofiber/fiber/v2"
	"my-little-ps/api"
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
