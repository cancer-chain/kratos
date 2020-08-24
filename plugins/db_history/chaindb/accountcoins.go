package chaindb

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/KuChainNetwork/kuchain/plugins/db_history/types"
	"github.com/KuChainNetwork/kuchain/utils/eventutil"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/tendermint/tendermint/libs/log"
)

type EventAccCoins struct {
	Height               int64  `json:"height"`
	Amount               string `json:"amount"`
	Account              string `json:"creator"`
	From                 string `json:"from"`
	To                   string `json:"to"`
	DestinationValidator string `json:"destination_validator"`
	SourceValidator      string `json:"source_validator"`
	Creator              string `json:"creator"`
	Time                 string `json:"block_time"`
}

type CreateAccCoinsModel struct {
	tableName struct{} `pg:"AccCoins,alias:AccCoins"` // default values are the same
	ID        int64    // both "Id" and "ID" are detected as primary key

	Height      int64  `pg:"default:0" json:"height"`
	Amount      int64  `pg:"default:0" json:"amount"`
	AmountFloat int64  `pg:"default:0" json:"amount_float"`
	Symbol      string `pg:"unique:asy" json:"symbol"`
	Account     string `pg:"unique:asy" json:"account"`
	SyncState   int    `pg:"sync_state,notnull,default:0" json:"sync_state"`
	Time        string `json:"time"`
}

const minCoinSymbol = len("0s\t")

var (
	startOnce sync.Once
)

func MakeCoinSql(msg EventAccCoins, n ...int32) CreateAccCoinsModel {
	coin, _ := NewCoin(msg.Amount)

	m := CreateAccCoinsModel{
		Height:      msg.Height,
		Amount:      coin.Amount,
		AmountFloat: coin.AmountFloat,
		Symbol:      coin.Symbol,
		Account:     msg.Account,
		Time:        msg.Time,
	}

	if len(n) > 0 && n[0] < 0 {
		m.Amount = m.Amount * -1
		m.AmountFloat = m.AmountFloat * -1
	}

	return m
}

func acExec(db *pg.DB, model CreateAccCoinsModel, logger log.Logger) error {
	var m CreateAccCoinsModel
	err := orm.NewQuery(db, &m).Where(fmt.Sprintf("Symbol='%s' and account='%s'", model.Symbol, model.Account)).Select()
	if err != nil {
		logger.Debug("acExec1", "model", model)
		err = orm.Insert(db, &model)
	} else {
		logger.Debug("1 acExec2", "model", model, "m", m)
		model.Amount, model.AmountFloat = CoinAdd(model.Amount, model.AmountFloat, m.Amount, m.AmountFloat)
		logger.Debug("2 acExec2", "model", model)
		_, err = orm.NewQuery(db, &model).Where(fmt.Sprintf("Symbol='%s' and account='%s'", model.Symbol, model.Account)).Update()
	}
	if err == nil {
		_, err = orm.NewQuery(db, &model).
			Where(fmt.Sprintf("Symbol='%s' and account='%s'", model.Symbol, model.Account)).
			Set(fmt.Sprintf("amount=%d, amount_float=%d", model.Amount, model.AmountFloat)).Update()
	}
	return err
}

func EventAccCoinsAdd(db *pg.DB, logger log.Logger, evt *types.Event) {
	tx, _ := db.Begin()
	var AccMsg EventAccCoins
	err := eventutil.UnmarshalKVMap(evt.Attributes, &AccMsg)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}

	AccMsg.Account = AccMsg.Creator
	if len(AccMsg.Amount) <= minCoinSymbol {
		EventErr(db, logger, NewErrMsgString(fmt.Sprintf("height: %d,amount is null", AccMsg.Height)))
		return
	}

	err = acExec(db, MakeCoinSql(AccMsg), logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}
	tx.Commit()
}

func EventAccCoinsMintAdd(db *pg.DB, logger log.Logger, evt *types.Event) {
	tx, _ := db.Begin()
	var AccMsg EventAccCoins
	err := eventutil.UnmarshalKVMap(evt.Attributes, &AccMsg)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}
	AccMsg.Account = AccMsg.To

	if len(AccMsg.Amount) <= minCoinSymbol {
		_, file, line, _ := runtime.Caller(1)
		EventErr(db, logger, NewErrMsgString(fmt.Sprintf("height: %d,amount is null， %s, %d ", AccMsg.Height, file, line)))
		return
	}

	err = acExec(db, MakeCoinSql(AccMsg), logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}
	tx.Commit()
}

func EventAccCoinsReduce(db *pg.DB, logger log.Logger, evt *types.Event) {
	tx, _ := db.Begin()
	var AccMsg EventAccCoins
	err := eventutil.UnmarshalKVMap(evt.Attributes, &AccMsg)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}
	AccMsg.Account = AccMsg.From

	if len(AccMsg.Amount) <= minCoinSymbol {
		_, file, line, _ := runtime.Caller(1)
		EventErr(db, logger, NewErrMsgString(fmt.Sprintf("height: %d,amount is null， %s, %d ", AccMsg.Height, file, line)))
		return
	}

	err = acExec(db, MakeCoinSql(AccMsg, -1), logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}
	tx.Commit()
}

func EventAccCoinsMove(db *pg.DB, logger log.Logger, evt *types.Event) {
	var AccMsg1 EventAccCoins
	err := eventutil.UnmarshalKVMap(evt.Attributes, &AccMsg1)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}

	var AccMsg2 EventAccCoins
	err = eventutil.UnmarshalKVMap(evt.Attributes, &AccMsg2)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}

	if len(AccMsg1.Amount) <= minCoinSymbol || len(AccMsg2.Amount) <= minCoinSymbol {
		_, file, line, _ := runtime.Caller(1)
		EventErr(db, logger, NewErrMsgString(fmt.Sprintf("height: %d,amount is null， %s, %d ", AccMsg1.Height, file, line)))
		return
	}

	AccMsg1.Account = AccMsg1.From
	AccMsg2.Account = AccMsg2.To

	tx, _ := db.Begin()
	err = acExec(db, MakeCoinSql(AccMsg1, -1), logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}
	err = acExec(db, MakeCoinSql(AccMsg2), logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}
	tx.Commit()
}

func EventAccCompleteReDelegateCoinsMove(db *pg.DB, logger log.Logger, evt *types.Event) {
	var AccMsg1 EventAccCoins
	err := eventutil.UnmarshalKVMap(evt.Attributes, &AccMsg1)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}

	var AccMsg2 EventAccCoins
	err = eventutil.UnmarshalKVMap(evt.Attributes, &AccMsg2)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}
	AccMsg1.Account = AccMsg1.SourceValidator
	AccMsg2.Account = AccMsg2.DestinationValidator

	if len(AccMsg1.Amount) <= minCoinSymbol || len(AccMsg2.Amount) <= minCoinSymbol {
		_, file, line, _ := runtime.Caller(1)
		EventErr(db, logger, NewErrMsgString(fmt.Sprintf("height: %d,amount is null， %s, %d ", AccMsg1.Height, file, line)))
		return
	}

	tx, _ := db.Begin()
	err = acExec(db, MakeCoinSql(AccMsg1, -1), logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}

	err = acExec(db, MakeCoinSql(AccMsg2), logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}

	tx.Commit()
}
