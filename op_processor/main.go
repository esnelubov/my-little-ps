package main

import (
	"flag"
	"fmt"
	"log"
	"my-little-ps/common/cache_maps/currency"
	"my-little-ps/common/config"
	"my-little-ps/common/controllers"
	"my-little-ps/common/database"
	"my-little-ps/common/logger"
	"my-little-ps/common/pool"
	"my-little-ps/common/scheduler"
	"my-little-ps/common/tasks"
	"my-little-ps/op_processor/worker"
	"os"
	"os/signal"
	"syscall"
)

var (
	Conf            config.IConfig
	Log             *logger.Log
	DB              *database.DB
	Scheduler       *scheduler.Scheduler
	Pool            *pool.TaskPool
	CurrenciesCache *currency.CacheMap
	Tasks           *tasks.Tasks
	Processor       *worker.Processor
)

var SetUpDependencies = func() {
	Conf = config.New("settings")
	Log = logger.New(Conf)
	DB = database.New(Conf)
	Scheduler = scheduler.New(Log, Conf)
	Pool = pool.New(Log, Conf)
	controllers.Setup(Log, DB)
	CurrenciesCache = currency.New(Log, controllers.Currency)
	Tasks = tasks.New(Log, CurrenciesCache)
	Processor = worker.New(Log, Conf, DB, Pool, controllers.Wallet, controllers.Operation, CurrenciesCache)
}

func setUpScheduler() {
	fmt.Println("Setting up the task scheduler")

	Scheduler.AddTask("@every 5m", Tasks.UpdateCurrencyCache)

	fmt.Println("Pre running some tasks")
	Tasks.UpdateCurrencyCache()
}

func main() {
	workerNumber := flag.Int("number", 0, "worker will process only operations of wallets marked with the given number")
	flag.Parse()

	if *workerNumber == 0 {
		fmt.Println("Starting to process operations for ALL wallets")
	} else {
		fmt.Printf("\"Starting to process operations for wallets group %d ...\n", *workerNumber)
	}

	SetUpDependencies()

	setUpScheduler()
	Scheduler.Start()

	shutdownTimeoutSec := Conf.GetDurationSec("shutdownTimeoutSec")

	// Run in a different goroutine
	go func() {
		if err := Processor.Run(*workerNumber); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	_ = <-c // This blocks the main thread until an interrupt is received

	fmt.Printf("Shutting down in %v...\n", shutdownTimeoutSec)
	Processor.Shutdown()
	_ = Scheduler.GracefulShutdown(shutdownTimeoutSec)
	_ = Pool.GracefulShutdown(shutdownTimeoutSec)
	_ = Log.Sync()

}
