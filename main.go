package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log"
	"my-little-ps/api"
	"my-little-ps/app"
	"my-little-ps/config"
	"my-little-ps/constants"
	"my-little-ps/controllers"
	"my-little-ps/database"
	"my-little-ps/logger"
	"my-little-ps/routes"
	"os"
	"os/signal"
	"syscall"
)

var (
	Conf config.IConfig
	Log  *logger.Log
	DB   *database.DB
)

func setUpRoutes(a *app.App) {
	a.Post("/wallet", routes.AddWallet)
	a.Post("/receive_amount", routes.ReceiveAmount)
	a.Post("/transfer_amount", routes.TransferAmount)
	a.Post("/update_currencies", routes.UpdateCurrencies)
	a.Get("/operations/:walletId", routes.GetOperations)
	a.Get("/operations_file/:walletId", routes.GetOperationsCSV)
	a.Get("/operations_total/:walletId", routes.GetOperationsTotal)
	a.Get("/convert_amount/:amount/:from/:to", routes.ConvertAmount)
}

func main() {
	Conf = config.New("settings")
	Log = logger.New()
	DB = database.New(Conf)
	constants.Setup()
	controllers.Setup(DB)
	api.Setup(controllers.Wallet, controllers.Operation, controllers.Currency)

	shutdownTimeoutSec := Conf.GetDurationSec("shutdownTimeoutSec")

	a := app.New(fiber.Config{
		Prefork:      Conf.GetBool("prefork"),
		IdleTimeout:  Conf.GetDurationSec("idleTimeoutSec"),
		ReadTimeout:  Conf.GetDurationSec("readTimeoutSec"),
		WriteTimeout: Conf.GetDurationSec("writeTimeoutSec"),
		ErrorHandler: ErrorHandler,
	})

	a.Use(recover.New())
	setUpRoutes(a)

	a.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	// Listen from a different goroutine
	go func() {
		if err := a.Listen(Conf.GetString("ip")); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	_ = <-c // This blocks the main thread until an interrupt is received

	fmt.Println("Shutting down...")
	_ = a.GracefulShutdown(shutdownTimeoutSec)
	_ = Log.Sync()
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	return ctx.Status(200).JSON(&routes.ResponseData{Payload: &routes.PayloadError{Error: err.Error()}})
}
