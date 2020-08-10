package chaindb

import (
	"github.com/KuChainNetwork/kuchain/plugins/db_history/types"
	"github.com/KuChainNetwork/kuchain/utils/eventutil"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/tendermint/tendermint/libs/log"
)

type DelegationChange struct {
	Height    string `json:"height"`
	Action    string `json:"action"`
	Hash      string `json:"hash"`
	Delegator string `json:"delegator"`
	Validator string `json:"validator"`
	Amount    string `json:"amount"`
	Fee       string `json:"fee"`
	Time      string `json:"time"`
}

type CreateDelegationChangeModel struct {
	tableName struct{} `pg:"delegation_change,alias:delegation_change"` // default values are the same
	ID        int64    // bot

	Height      string `json:"height"`
	Action      string `json:"action"`
	Hash        string `json:"hash"`
	Delegator   string `json:"delegator"`
	Validator   string `json:"validator"`
	Amount      int64  `pg:"default:0" json:"amount"`
	AmountFloat int64  `pg:"default:0" json:"amount_float"`
	Symbol      string `json:"symbol"`
	Fee         string `json:"fee"`
	Time        string `json:"time"`
}

func makeDelegationChangeSql(DeMsg DelegationChange) CreateDelegationChangeModel {

	coin, _ := NewCoin(DeMsg.Amount)

	q := CreateDelegationChangeModel{
		Height:      DeMsg.Height,
		Action:      DeMsg.Action,
		Hash:        DeMsg.Hash,
		Validator:   DeMsg.Validator,
		Delegator:   DeMsg.Delegator,
		Amount:      coin.Amount,
		AmountFloat: coin.AmountFloat,
		Symbol:      coin.Symbol,
		Fee:         DeMsg.Fee,
		Time:        DeMsg.Time,
	}
	return q
}

func EventDelegationChange(db *pg.DB, logger log.Logger, evt *types.Event) {
	var DMsg DelegationChange
	err := eventutil.UnmarshalKVMap(evt.Attributes, &DMsg)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}

	q := makeDelegationChangeSql(DMsg)
	logger.Debug("EventDelegationChange", "CreateDelegationChangeModel", q)
	err = orm.Insert(db, &q)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}
}
