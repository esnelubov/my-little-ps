package worker

import "my-little-ps/common/models"

type Chunk struct {
	InOperations  []*models.InOperation
	OutOperations []*models.OutOperation
	Wallet        *models.Wallet
}

type WalletToChunk map[uint]*Chunk

func (c WalletToChunk) AddInOperation(op *models.InOperation) {
	chunk, ok := c[op.TargetWalletId]
	if !ok {
		chunk = &Chunk{
			InOperations: []*models.InOperation{},
		}
	}
	chunk.InOperations = append(chunk.InOperations, op)
}

func (c WalletToChunk) AddOutOperation(op *models.OutOperation) {
	chunk, ok := c[op.OriginWalletId]
	if !ok {
		chunk = &Chunk{
			OutOperations: []*models.OutOperation{},
		}
	}
	chunk.OutOperations = append(chunk.OutOperations, op)
}

func (c WalletToChunk) AddWallet(w *models.Wallet) {
	chunk, ok := c[w.ID]
	if !ok {
		chunk = &Chunk{}
	}
	chunk.Wallet = w
}
