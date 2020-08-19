package types

import (
	//abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/types"
	"time"
)

type ReqBeginBlock struct {
	RequestBeginBlock types.Block
	Tx                ReqTx
	Events            ReqEvents
	TxEvents          ReqEvents
	FeeEvents         ReqEvents

	ValidatorInfo string
	Time          time.Time
}

type ReqEndBlock struct {
	Height int64
}
