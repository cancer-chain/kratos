package chaindb

import (
	"fmt"
	"github.com/KuChainNetwork/kuchain/plugins/db_history/types"
	"github.com/KuChainNetwork/kuchain/utils/eventutil"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/tendermint/tendermint/libs/log"
)

type EventLockAccCoins struct {
	Amount  string `json:"amount"`
	Height  int64  `json:"height"`
	Account string `json:"account"`
	From    string `json:"from"`
	Module  string `json:"module"`
	Time    string `json:"block_time"`
}

type CreateLockAccCoinsModel struct {
	tableName struct{} `pg:"LockAccCoins,alias:LockAccCoins"` // default values are the same

	ID int // both "Id" and "ID" are detected as primary key

	Amount      int64  `pg:"default:0" json:"amount"`
	AmountFloat int64  `pg:"default:0" json:"amount_float"`
	AmountStr   string `json:"amount_str"`
	Symbol      string `pg:"unique:as" json:"symbol"`
	Height      int64  `json:"height"`
	Account     string `pg:"unique:as" json:"account"`
	Time        string `json:"time"`
}

func MakeLockCoinSql(msg EventLockAccCoins, n ...int32) CreateLockAccCoinsModel {
	coin, _ := NewCoin(msg.Amount)

	m := CreateLockAccCoinsModel{
		Amount:      coin.Amount,
		AmountFloat: coin.AmountFloat,
		Symbol:      coin.Symbol,
		Height:      msg.Height,
		Account:     msg.Account,
		Time:        msg.Time,
	}

	if len(msg.From) > 0 {
		m.Account = msg.From
	}

	if len(n) > 0 && n[0] < 0 {
		m.Amount = m.Amount * -1
		m.AmountFloat = m.AmountFloat * -1
	}

	return m
}

func lExec(db *pg.DB, model CreateLockAccCoinsModel, logger log.Logger) error {
	var m CreateLockAccCoinsModel
	err := orm.NewQuery(db, &m).Where(fmt.Sprintf("Symbol='%s' and account='%s'", model.Symbol, model.Account)).Select()
	if err != nil {
		logger.Debug("lExec1", "model", model)
		err = orm.Insert(db, &model)
	} else {
		model.Amount, model.AmountFloat = CoinAdd(model.Amount, model.AmountFloat, m.Amount, m.AmountFloat)
		logger.Debug("lExec2", "model", model)
		_, err = orm.NewQuery(db, &model).Where(fmt.Sprintf("Symbol='%s' and account='%s'", model.Symbol, model.Account)).Update()
	}

	if err == nil {
		_, err = orm.NewQuery(db, &model).
			Where(fmt.Sprintf("Symbol='%s' and account='%s'", model.Symbol, model.Account)).
			Set(fmt.Sprintf("amount=%d, amount_float=%d", model.Amount, model.AmountFloat)).Update()
	}
	return err
}

func EventLockAccCoinsAdd(db *pg.DB, logger log.Logger, evt *types.Event) {
	var lockMsg EventLockAccCoins
	err := eventutil.UnmarshalKVMap(evt.Attributes, &lockMsg)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}

	m := MakeLockCoinSql(lockMsg)

	tx, _ := db.Begin()
	err = lExec(db, m, logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}
	tx.Commit()
}

func EventUnLockAccCoinsAdd(db *pg.DB, logger log.Logger, evt *types.Event) {
	var lockMsg EventLockAccCoins
	err := eventutil.UnmarshalKVMap(evt.Attributes, &lockMsg)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}

	m := MakeLockCoinSql(lockMsg, -1)

	tx, _ := db.Begin()
	err = lExec(db, m, logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}
	tx.Commit()
}
