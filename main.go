package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log"
	"my-little-ps/app"
	"my-little-ps/config"
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
	a.Get("/hello", routes.Hello)
	//a.Get("/allbooks", routes.AllBooks)
	//a.Get("/book/:id", routes.GetBook)
	//a.Post("/book", routes.AddBook)
	//a.Put("/book", routes.Update)
	//a.Delete("/book", routes.Delete)
}

func main() {
	Conf = config.New("settings")
	Log = logger.New()
	DB = database.New(Conf)
	controllers.Setup(DB)
	_ = controllers.Operation.DB
	_ = controllers.Wallet.DB

	shutdownTimeoutSec := Conf.GetDurationSec("shutdownTimeoutSec")

	a := app.New(fiber.Config{
		Prefork:      Conf.GetBool("prefork"),
		IdleTimeout:  Conf.GetDurationSec("idleTimeoutSec"),
		ReadTimeout:  Conf.GetDurationSec("readTimeoutSec"),
		WriteTimeout: Conf.GetDurationSec("writeTimeoutSec"),
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
