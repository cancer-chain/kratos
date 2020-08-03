package chaindb

import (
	"fmt"
	"github.com/KuChainNetwork/kuchain/plugins/db_history/types"
	"github.com/KuChainNetwork/kuchain/utils/eventutil"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/tendermint/tendermint/libs/log"
)

type EventValidator struct {
	tableName struct{} `pg:"Validator,alias:Validator"` // default values are the same

	ID int // both "Id" and "ID" are detected as primary key

	Height            int64  `pg:"default:0",json:"height"`
	Address           string `pg:"unique:as",json:"Address"`
	Sender            string `pg:"unique:as",json:"sender"`
	ConsensusPubkey   string `json:"ConsensusPubkey"`
	Jailed            string `json:"Jailed"`
	Status            string `json:"Status"`
	Tokens            string `json:"Tokens"`
	DelegatorShares   string `json:"DelegatorShares"`
	Description       string `json:"Description"`
	UnbondingHeight   string `json:"UnbondingHeight"`
	UnbondingTime     string `json:"UnbondingTime"`
	Commission        string `json:"commission_rate"`
	MinSelfDelegation string `json:"min_self_delegation"`
	Time              string `json:"block_time"`
}

func vExec(db *pg.DB, model EventValidator, logger log.Logger) error {
	var m EventValidator
	err := orm.NewQuery(db, &m).Where(fmt.Sprintf("Address='%s' ", model.Address)).Select()
	if err != nil {
		logger.Debug("vExec1", "model", model)
		err = orm.Insert(db, &model)
	} else {
		logger.Debug("vExec2", "model", model)
		_, err = orm.NewQuery(db, &model).Where(fmt.Sprintf("Address='%s' ", model.Address)).Update()
	}
	return err
}

func vExecUpdate(db *pg.DB, addr string, logger log.Logger) error {
	var m EventValidator
	err := orm.NewQuery(db, &m).Where(fmt.Sprintf("Address='%s' ", addr)).Select()
	if err != nil {
		logger.Debug("vExecUpdate1", "model", m)
		err = orm.Insert(db, &m)
	} else {
		m.Jailed = "true"
		m.Status = "true"
		logger.Debug("vExecUpdate2", "model", m)
		_, err = orm.NewQuery(db, &m).Where(fmt.Sprintf("Address='%s' ", addr)).Update()
	}
	return err
}
func vExecUpdate2(db *pg.DB, sender string, logger log.Logger) error {
	var m EventValidator
	err := orm.NewQuery(db, &m).Where(fmt.Sprintf("Address='%s' ", sender)).Select()
	if err != nil {
		logger.Debug("vExecUpdate21", "model", m)
		err = orm.Insert(db, &m)
	} else {
		m.Jailed = "true"
		logger.Debug("vExecUpdate22", "model", m)
		_, err = orm.NewQuery(db, &m).Where(fmt.Sprintf("Address='%s' ", sender)).Update()
	}
	return err
}

func EventValidatorAdd(db *pg.DB, logger log.Logger, evt *types.Event) {
	var Validator EventValidator
	err := eventutil.UnmarshalKVMap(evt.Attributes, &Validator)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}
	Validator.Description = ToJson(Validator.Description)
	Validator.Commission = ToJson(Validator.Commission)

	err = vExec(db, Validator, logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}

}

func EventValidatorSlash(db *pg.DB, logger log.Logger, evt *types.Event) {
	var Validator EventValidator
	err := eventutil.UnmarshalKVMap(evt.Attributes, &Validator)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}
	Validator.Description = ToJson(Validator.Description)
	Validator.Commission = ToJson(Validator.Commission)

	err = vExecUpdate(db, Validator.Address, logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}
}

func EventValidatorUnjail(db *pg.DB, logger log.Logger, evt *types.Event) {
	var Validator EventValidator
	err := eventutil.UnmarshalKVMap(evt.Attributes, &Validator)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}
	Validator.Description = ToJson(Validator.Description)
	Validator.Commission = ToJson(Validator.Commission)

	err = vExecUpdate2(db, Validator.Sender, logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}
}
