package chaindb

import (
	"fmt"
	"github.com/KuChainNetwork/kuchain/plugins/db_history/types"
	"github.com/KuChainNetwork/kuchain/utils/eventutil"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/tendermint/tendermint/libs/log"
)

type Delegation struct {
	Height    string `json:"height"`
	Validator string `json:"validator"`
	Delegator string `json:"delegator"`
	Amount    string `json:"amount"`
	Time      string `json:"block_time"`
}

type CreateDelegationModel struct {
	tableName struct{} `pg:"Delegation,alias:Delegation"` // default values are the same

	ID          int    // both "Id" and "ID" are detected as primary key
	Height      string `pg:"default:0" json:"height"`
	Validator   string `pg:"unique:vd" json:"validator"`
	Delegator   string `pg:"unique:vd" json:"delegator"`
	Amount      int64  `pg:"default:0" json:"amount"`
	AmountFloat int64  `pg:"default:0" json:"amount_float"`
	Symbol      string `json:"symbol"`
	Time        string `json:"time"`
}

func makeDelegationSql(msg Delegation) CreateDelegationModel {
	coin, _ := NewCoin(msg.Amount)

	q := CreateDelegationModel{
		Height:      msg.Height,
		Validator:   msg.Validator,
		Delegator:   msg.Delegator,
		Amount:      coin.Amount,
		AmountFloat: coin.AmountFloat,
		Symbol:      coin.Symbol,
		Time:        msg.Time,
	}
	return q
}

func dExec(db *pg.DB, model CreateDelegationModel, logger log.Logger) error {
	var m CreateDelegationModel
	err := orm.NewQuery(db, &m).Where(fmt.Sprintf("validator='%s' and delegator='%s'", model.Validator, model.Delegator)).Select()
	if err != nil {
		logger.Debug("dExec1", "model", model)
		err = orm.Insert(db, &model)
	} else {
		model.Amount, model.AmountFloat = CoinAdd(model.Amount, model.AmountFloat, m.Amount, m.AmountFloat)
		logger.Debug("dExec2", "model", model)
		_, err = orm.NewQuery(db, &model).Where(fmt.Sprintf("validator='%s' and delegator='%s'", model.Validator, model.Delegator)).Update()
	}

	if err == nil {
		_, err = orm.NewQuery(db, &model).
			Where(fmt.Sprintf("validator='%s' and delegator='%s'", model.Validator, model.Delegator)).
			Set(fmt.Sprintf("amount=%d, amount_float=%d", model.Amount, model.AmountFloat)).Update()
	}
	return err
}

func EventDelegationAdd(db *pg.DB, logger log.Logger, evt *types.Event) {
	var msg Delegation
	err := eventutil.UnmarshalKVMap(evt.Attributes, &msg)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}

	q := makeDelegationSql(msg)
	err = dExec(db, q, logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}
}

func EventDelegationDel(db *pg.DB, logger log.Logger, evt *types.Event) {
	var msg Delegation
	err := eventutil.UnmarshalKVMap(evt.Attributes, &msg)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}

	q := makeDelegationSql(msg)
	err = dExec(db, q, logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}
}
