package worker

import (
	ch "my-little-ps/common/cache_maps/currency"
	"my-little-ps/common/config"
	"my-little-ps/common/constants"
	oc "my-little-ps/common/controllers/operation"
	wc "my-little-ps/common/controllers/wallet"
	"my-little-ps/common/database"
	"my-little-ps/common/logger"
	"my-little-ps/common/models"
	"my-little-ps/common/pool"
	"time"
)

type Processor struct {
	logger              *logger.Log
	config              config.IConfig
	db                  *database.DB
	pool                *pool.TaskPool
	walletController    *wc.Controller
	operationController *oc.Controller
	currencyCache       *ch.CacheMap
	shutDown            chan struct{}
}

func New(logger *logger.Log, config config.IConfig, db *database.DB, pool *pool.TaskPool, walletController *wc.Controller, operationController *oc.Controller, currencyCache *ch.CacheMap) *Processor {
	return &Processor{
		logger:              logger,
		config:              config,
		db:                  db,
		pool:                pool,
		walletController:    walletController,
		operationController: operationController,
		currencyCache:       currencyCache,
		shutDown:            make(chan struct{}, 1),
	}
}

func (p *Processor) Run(workerNumber int) (err error) {
	p.logger.Debug("Starting processing cycle")

	var (
		delay time.Duration
	)

	delay, err = time.ParseDuration(p.config.GetString("operationProcessorDelay"))
	if err != nil {
		return
	}

	for {
		select {
		case <-p.shutDown:
			return
		default:
			err = p.ProcessOperations(workerNumber)
			if err != nil {
				return
			}
		}
		time.Sleep(delay)
	}
}

func (p *Processor) Shutdown() {
	p.logger.Debug("Stopping processing cycle")
	close(p.shutDown)
}

func (p *Processor) ProcessOperations(workerNumber int) (err error) {
	p.logger.Debug("Processing new operations")

	return p.db.Transaction(func(tx *database.DB) (err error) {
		var (
			walletController    = p.walletController.WithTransaction(tx)
			operationController = p.operationController.WithTransaction(tx)
			limit               = p.config.GetInt("operationProcessorLimit")
			wallets             []*models.Wallet
			inOperations        []*models.InOperation
			outOperations       []*models.OutOperation
		)

		if workerNumber == 0 {
			wallets, inOperations, outOperations, err = getAllOperations(walletController, operationController, limit)
		} else {
			wallets, inOperations, outOperations, err = getOperationsForWorker(workerNumber, walletController, operationController, limit)
		}
		if err != nil {
			return
		}

		p.logger.Debugf("Found %d IN and %d OUT operations for %d wallets", len(inOperations), len(outOperations), len(wallets))

		chunks := WalletToChunk{}

		for _, op := range inOperations {
			chunks.AddInOperation(op)
		}

		for _, op := range outOperations {
			chunks.AddOutOperation(op)
		}

		for _, w := range wallets {
			chunks.AddWallet(w)
		}

		for _, c := range chunks {
			p.pool.RunTask(func() error {
				return p.ProcessChunk(c)
			})
		}

		err = p.pool.WaitTasks().Error
		if err != nil {
			return err
		}

		err = saveRecords(tx, wallets, inOperations, outOperations)
		if err != nil {
			return err
		}

		return
	})
}

// ProcessChunk processes all operations for a wallet
func (p *Processor) ProcessChunk(chunk *Chunk) (err error) {
	p.logger.Debugf("Processing %d IN and %d OUT operations for wallet %d", len(chunk.InOperations), len(chunk.OutOperations), chunk.Wallet.ID)

	var (
		convertedAmount int64
	)

	for _, op := range chunk.InOperations {
		convertedAmount, err = p.currencyCache.Convert(op.Currency, chunk.Wallet.Currency, op.Amount)
		if err != nil {
			return
		}

		chunk.Wallet.Balance += convertedAmount
		op.Status = constants.OpStatusSuccess
	}

	for _, op := range chunk.OutOperations {
		convertedAmount, err = p.currencyCache.Convert(op.Currency, chunk.Wallet.Currency, op.Amount)
		if err != nil {
			return
		}

		if convertedAmount > chunk.Wallet.Balance {
			p.logger.Debugf("OUT Operation %s is declined (not enough money)", op.OperationId)
			op.Status = constants.OpStatusDecline
		} else {
			chunk.Wallet.Balance -= convertedAmount
			op.Status = constants.OpStatusSuccess

			err = p.operationController.NewInternalIn(op.TransactionId, op.OriginWalletId, op.TargetWalletId, convertedAmount, op.Currency)
			if err != nil {
				return
			}
		}
	}

	p.logger.Debugf("Wallet %d has new balance: %d", chunk.Wallet.ID, chunk.Wallet.Balance)

	return
}

func getWalletIds(inOperations []*models.InOperation, outOperations []*models.OutOperation) (ids []uint) {
	idSet := map[uint]struct{}{}
	for _, op := range inOperations {
		idSet[op.TargetWalletId] = struct{}{}
	}

	for _, op := range outOperations {
		idSet[op.OriginWalletId] = struct{}{}
	}

	ids = make([]uint, len(idSet))
	i := 0
	for id, _ := range idSet {
		ids[i] = id
		i++
	}

	return
}

func getAllOperations(wc *wc.Controller, oc *oc.Controller, limit int) (wallets []*models.Wallet, inOperations []*models.InOperation, outOperations []*models.OutOperation, err error) {
	inOperations, err = oc.GetNewInOperations(limit)
	if err != nil {
		return
	}

	outOperations, err = oc.GetNewOutOperations(limit)
	if err != nil {
		return
	}

	ids := getWalletIds(inOperations, outOperations)

	wallets, err = wc.GetWalletsWithIds(ids)

	return
}

func getOperationsForWorker(wn int, wc *wc.Controller, oc *oc.Controller, limit int) (wallets []*models.Wallet, inOperations []*models.InOperation, outOperations []*models.OutOperation, err error) {
	wallets, err = wc.GetWalletsForWorker(wn)
	if err != nil {
		return
	}

	ids := make([]uint, len(wallets))
	for i, w := range wallets {
		ids[i] = w.ID
	}

	inOperations, err = oc.GetNewInOperationsForWallets(ids, limit)
	if err != nil {
		return
	}

	outOperations, err = oc.GetNewOutOperationsForWallets(ids, limit)
	if err != nil {
		return
	}

	return
}

func saveRecords(tx *database.DB, wallets []*models.Wallet, inOperations []*models.InOperation, outOperations []*models.OutOperation) (err error) {
	for _, op := range inOperations {
		err = tx.Save(op)
		if err != nil {
			return
		}
	}

	for _, op := range outOperations {
		err = tx.Save(op)
		if err != nil {
			return
		}
	}

	for _, w := range wallets {
		err = tx.Save(w)
		if err != nil {
			return
		}
	}

	return
}
