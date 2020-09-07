package chaindb

import (
	"fmt"
	"github.com/KuChainNetwork/kuchain/plugins/db_history/types"
	"github.com/KuChainNetwork/kuchain/utils/eventutil"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/tendermint/tendermint/libs/log"
)

type UnBond struct {
	Height    int64  `json:"height"`
	Validator string `json:"validator"`
	Delegator string `json:"delegator"`
	Amount    string `json:"amount"`
	Time      string `json:"block_time"`
}

type EventUnBondModel struct {
	tableName struct{} `pg:"un_bond,alias:un_bond"` // default values are the same

	ID int // both "Id" and "ID" are detected as primary key

	Height      int64  `pg:"default:0" json:"height"`
	Validator   string `pg:"unique:apu" json:"validator"`
	Delegator   string `pg:"unique:apu" json:"delegator"`
	Amount      int64  `pg:"default:0" json:"amount"`
	AmountFloat int64  `pg:"default:0" json:"amount_float"`
	AmountStr   string `json:"amount_str"`
	Symbol      string `json:"symbol"`
	Time        string `json:"time"`
}

func makeUnBondSql(msg UnBond) (EventUnBondModel, error) {
	coin, err := NewCoin(msg.Amount)
	if err != nil {
		return EventUnBondModel{}, err
	}

	q := EventUnBondModel{
		Height:      msg.Height,
		Validator:   msg.Validator,
		Delegator:   msg.Delegator,
		Amount:      coin.Amount,
		AmountFloat: coin.AmountFloat,
		AmountStr:   msg.Amount,
		Symbol:      coin.Symbol,
		Time:        msg.Time,
	}
	return q, nil
}

func unBondAddExec(db *pg.DB, model EventUnBondModel, logger log.Logger) error {
	var m EventUnBondModel
	err := orm.NewQuery(db, &m).Where(fmt.Sprintf("validator='%s' and delegator='%s'", model.Validator, model.Delegator)).Select()
	if err != nil {
		logger.Debug("unBondAddExec", "model", model)
		err = orm.Insert(db, &model)
	} else {
		model.Amount, model.AmountFloat = CoinAdd(model.Amount, model.AmountFloat, m.Amount, m.AmountFloat)
		logger.Debug("unBondAddExec2", "model", model)
		_, err = orm.NewQuery(db, &model).Where(fmt.Sprintf("validator='%s' and delegator='%s'", model.Validator, model.Delegator)).Update()
	}

	if err == nil {
		_, err = orm.NewQuery(db, &model).
			Where(fmt.Sprintf("validator='%s' and delegator='%s'", model.Validator, model.Delegator)).
			Set(fmt.Sprintf("amount=%d, amount_float=%d", model.Amount, model.AmountFloat)).Update()
	}
	return err
}

func EventUnBondAdd(db *pg.DB, logger log.Logger, evt *types.Event) {
	var unBond UnBond
	err := eventutil.UnmarshalKVMap(evt.Attributes, &unBond)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}

	unBondModel, err := makeUnBondSql(unBond)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}

	err = unBondAddExec(db, unBondModel, logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}
}

func unBondCompleteExec(db *pg.DB, model EventUnBondModel, logger log.Logger) error {
	var m EventUnBondModel
	err := orm.NewQuery(db, &m).Where(fmt.Sprintf("validator='%s' and delegator='%s'", model.Validator, model.Delegator)).Select()
	if err != nil {
		logger.Debug("unBondCompleteExec", "model", model)
	} else {
		model.Amount, model.AmountFloat = CoinSub(model.Amount, model.AmountFloat, m.Amount, m.AmountFloat)
		logger.Debug("unBondCompleteExec2", "model", model)
		_, err = orm.NewQuery(db, &model).Where(fmt.Sprintf("validator='%s' and delegator='%s'", model.Validator, model.Delegator)).Update()
	}

	if err == nil {
		_, err = orm.NewQuery(db, &model).
			Where(fmt.Sprintf("validator='%s' and delegator='%s'", model.Validator, model.Delegator)).
			Set(fmt.Sprintf("amount=%d, amount_float=%d", model.Amount, model.AmountFloat)).Update()
	}
	return err
}

func EventUnBondComplete(db *pg.DB, logger log.Logger, evt *types.Event) {
	var unBond UnBond
	err := eventutil.UnmarshalKVMap(evt.Attributes, &unBond)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}

	unBondModel, err := makeUnBondSql(unBond)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}

	err = unBondCompleteExec(db, unBondModel, logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}
}
