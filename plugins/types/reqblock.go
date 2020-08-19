package types

import (
	//abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/types"
	"time"
)

type ReqBeginBlock struct {
	RequestBeginBlock types.Block `json:"request_begin_block"`
	Tx                []ReqTx     `json:"tx"`
	Events            ReqEvents   `json:"events"`
	TxEvents          []ReqEvents `json:"tx_events"`
	FeeEvents         ReqEvents   `json:"fee_events"`
	ValidatorInfo     string      `json:"validator_info"`
	Time              time.Time   `json:"time"`
}

type ReqEndBlock struct {
	Height int64
}
