package main

import (
	"flag"
	"fmt"
	"log"
	"my-little-ps/common/config"
	"my-little-ps/common/controllers"
	"my-little-ps/common/database"
	"my-little-ps/common/logger"
	"my-little-ps/wlt_balancer/worker"
	"os"
	"os/signal"
	"syscall"
)

var (
	Conf     config.IConfig
	Log      *logger.Log
	DB       *database.DB
	Balancer *worker.Balancer
)

var SetUpDependencies = func() {
	Conf = config.New("settings")
	Log = logger.New(Conf)
	DB = database.New(Conf)
	controllers.Setup(Log, DB)
	Balancer = worker.New(Log, Conf, DB, controllers.Wallet, controllers.Operation)
}

func main() {
	workerNumber := flag.Int("number", 0, "balancer will partition wallets for the given number of workers")
	flag.Parse()

	if *workerNumber == 0 {
		log.Fatal("positive '--number' is required")
	}

	fmt.Printf("\"Starting to balance wallets for %d workers  ...\n", *workerNumber)

	SetUpDependencies()

	shutdownTimeoutSec := Conf.GetDurationSec("shutdownTimeoutSec")

	// Run in a different goroutine
	go func() {
		if err := Balancer.Run(*workerNumber); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	_ = <-c // This blocks the main thread until an interrupt is received

	fmt.Printf("Shutting down in %v...\n", shutdownTimeoutSec)
	_ = Log.Sync()
}
