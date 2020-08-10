package types

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"time"
)

type ReqBlock struct {
	abci.RequestBeginBlock
	ValidatorInfo string
	Time          time.Time
}
