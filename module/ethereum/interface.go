package ethereum

import (
	"github.com/denniswon/reddio/api/model"
	"github.com/ethereum/go-ethereum/common"
)

//Module interface
type Module interface {
	GetLatestBlock() *model.Block
	GetTxByHash(hash common.Hash) *model.Transaction
	GetAddressBalance(address string) (string, error)
	TransferEth(privKey string, to string, amount int64) (string, error)
}
