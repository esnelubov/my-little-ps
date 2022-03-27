package worker

import (
	"errors"
	"my-little-ps/common/config"
	oc "my-little-ps/common/controllers/operation"
	wc "my-little-ps/common/controllers/wallet"
	"my-little-ps/common/database"
	"my-little-ps/common/logger"
	"my-little-ps/common/models"
	"my-little-ps/wlt_balancer/partitioner"
	"time"
)

type Balancer struct {
	logger              *logger.Log
	config              config.IConfig
	db                  *database.DB
	walletController    *wc.Controller
	operationController *oc.Controller
}

func New(logger *logger.Log, config config.IConfig, db *database.DB, walletController *wc.Controller, operationController *oc.Controller) *Balancer {
	return &Balancer{
		logger:              logger,
		config:              config,
		db:                  db,
		walletController:    walletController,
		operationController: operationController,
	}
}

func (b *Balancer) Run(workerNumber int) (err error) {
	b.logger.Debug("Starting rebalancing cycle")

	var (
		delay time.Duration
	)

	delay, err = time.ParseDuration(b.config.GetString("walletBalancerDelay"))
	if err != nil {
		return
	}

	for {
		b.logger.Debug("Performing new rebalancing iteration")
		err = b.RebalanceWallets(workerNumber, delay)
		if err != nil {
			return
		}
		b.logger.Debug("Wallets successfully rebalanced")
		time.Sleep(delay)
	}
}

func (b *Balancer) RebalanceWallets(workerNumber int, duration time.Duration) (err error) {
	var (
		walletToOpCount  map[int64]int64
		opCountToWallets map[int64][]int64
		walletCount      int64
	)

	walletToOpCount, err = b.operationController.CountOperationsPerWallet(time.Now().Add(-duration))
	if err != nil {
		return err
	}

	walletCount, err = b.walletController.CountWallets()
	if err != nil {
		return err
	}

	opCountToWallets = make(map[int64][]int64)
	for walletId, opCount := range walletToOpCount {
		wallets, ok := opCountToWallets[opCount]
		if !ok {
			wallets = []int64{}
		}
		wallets = append(wallets, walletId)
		opCountToWallets[opCount] = wallets
	}

	opCounts := []int64{}
	for opCount, _ := range opCountToWallets {
		opCounts = append(opCounts, opCount)
	}

	inactiveWalletsCount := walletCount - int64(len(opCountToWallets))

	for inactiveWalletsCount > 0 {
		opCounts = append(opCounts, 0)
		inactiveWalletsCount -= 1
	}

	partitions := partitioner.KarmarkarKarp(opCounts, workerNumber)

	opCountToWorkers := make(map[int64][]int, len(opCounts))
	for i, partition := range partitions {
		wn := i + 1

		for _, opCount := range partition {
			workers, ok := opCountToWorkers[opCount]
			if !ok {
				workers = []int{}
			}
			workers = append(workers, wn)
			opCountToWorkers[opCount] = workers
		}
	}

	return b.UpdateWalletsWithWorkers(walletToOpCount, opCountToWorkers)
}

var NoMoreWallets = errors.New("No more wallets to update")

func (b *Balancer) UpdateWalletsWithWorkers(walletToOpCount map[int64]int64, opCountToWorkers map[int64][]int) (err error) {
	var (
		offset = 0
		limit  = b.config.GetInt("walletBalancerLimit")
	)

	for err == nil || err != NoMoreWallets {
		err = b.db.Transaction(func(tx *database.DB) (err error) {
			var (
				walletController = b.walletController.WithTransaction(tx)
				wallets          []*models.Wallet
				workerNumber     int
			)

			wallets, err = walletController.GetAllWallets(offset, limit)
			if err != nil {
				return
			}

			if len(wallets) == 0 {
				return NoMoreWallets
			}

			b.logger.Debugf("Updating workers for %d wallets", len(wallets))

			for _, wallet := range wallets {
				opCount, ok := walletToOpCount[int64(wallet.ID)]
				if !ok {
					opCount = 0
				}

				workers := opCountToWorkers[opCount]
				workerNumber, opCountToWorkers[opCount] = pop(workers)

				wallet.Worker = int32(workerNumber)
				err = tx.Save(wallet)
				if err != nil {
					return
				}
			}

			offset += limit

			return
		})
	}

	if err == NoMoreWallets {
		err = nil
	}

	return
}

func pop(stack []int) (element int, modified_stack []int) {
	element, modified_stack = stack[len(stack)-1], stack[:len(stack)-1]
	return
}
