package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log"
	"my-little-ps/common/cache_maps/currency"
	"my-little-ps/common/config"
	"my-little-ps/common/controllers"
	"my-little-ps/common/database"
	"my-little-ps/common/logger"
	"my-little-ps/common/scheduler"
	"my-little-ps/common/tasks"
	"my-little-ps/gateway/api"
	"my-little-ps/gateway/app"
	"my-little-ps/gateway/routes"
	"os"
	"os/signal"
	"syscall"
)

var (
	Conf            config.IConfig
	Log             *logger.Log
	DB              *database.DB
	Scheduler       *scheduler.Scheduler
	CurrenciesCache *currency.CacheMap
	Tasks           *tasks.Tasks
)

var SetUpDependencies = func() {
	Conf = config.New("settings")
	Log = logger.New(Conf)
	DB = database.New(Conf)
	Scheduler = scheduler.New(Log, Conf)
	controllers.Setup(Log, DB)
	CurrenciesCache = currency.New(Log, controllers.Currency)
	Tasks = tasks.New(Log, CurrenciesCache)
	api.Setup(Log, controllers.Wallet, controllers.Operation, controllers.Currency, CurrenciesCache)
}

func setUpRoutes(a *app.App) {
	fmt.Println("Setting up routes")

	a.Post("/wallet", routes.AddWallet)
	a.Post("/receive_amount", routes.ReceiveAmount)
	a.Post("/transfer_amount", routes.TransferAmount)
	a.Post("/update_currencies", routes.UpdateCurrencies)
	a.Get("/operations/:walletId", routes.GetOperations)
	a.Get("/operations/file/:walletId", routes.GetOperationsCSV)
	a.Get("/operations/total/:walletId", routes.GetOperationsTotal)
	a.Get("/convert_amount/:amount/:from/:to", routes.ConvertAmount)
}

func setUpScheduler() {
	fmt.Println("Setting up the task scheduler")

	Scheduler.AddTask("@every 5m", Tasks.UpdateCurrencyCache)

	fmt.Println("Pre running some tasks")
	Tasks.UpdateCurrencyCache()
}

func main() {
	fmt.Println("Starting the gateway...")

	SetUpDependencies()

	setUpScheduler()
	Scheduler.Start()

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

	fmt.Printf("Shutting down in %v...\n", shutdownTimeoutSec)
	_ = a.GracefulShutdown(shutdownTimeoutSec)
	_ = Scheduler.GracefulShutdown(shutdownTimeoutSec)
	_ = Log.Sync()
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	return ctx.Status(200).JSON(&routes.ResponseData{Payload: &routes.PayloadError{Error: err.Error()}})
}
