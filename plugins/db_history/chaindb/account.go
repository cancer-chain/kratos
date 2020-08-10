package chaindb

import (
	"fmt"
	"github.com/KuChainNetwork/kuchain/plugins/db_history/types"
	"github.com/KuChainNetwork/kuchain/utils/eventutil"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/tendermint/tendermint/libs/log"
)

const tableName = "account"

type Account struct {
	Height    int64  `json:"height"`
	AccountId string `json:"account"`
	Creator   string `json:"creator"`
	Auth      string `json:"auth"`
	Time      string `json:"block_time"`
}

type CreateAccountModel struct {
	tableName struct{} `pg:"account,alias:account"` // default values are the same
	ID        int64    // both "Id" and "ID" are detected as primary key

	Height    int64  `pg:",type:bigint,default:0" json:"height"`
	AccountId string `pg:"unique:as" json:"account_id"`
	Creator   string `json:"creator"`
	Auth      string `json:"auth"`
	Time      string `json:"time"`
}

func makeAccountAddSql(msg Account) CreateAccountModel {
	q := CreateAccountModel{
		Height:    msg.Height,
		AccountId: msg.AccountId,
		Creator:   msg.Creator,
		Auth:      msg.Auth,
		Time:      msg.Time,
	}

	return q
}

func accExec(db *pg.DB, model CreateAccountModel, logger log.Logger) error {
	var m CreateAccountModel
	err := orm.NewQuery(db, &m).Where(fmt.Sprintf(" account_id='%s'", model.AccountId)).Select()
	if err != nil {
		logger.Debug("accExec1", "model", model)
		err = orm.Insert(db, &model)
	} else {
		if len(m.Creator) > 0 && len(model.Creator) <= 0 {
			model.Creator = m.Creator
		}
		if len(m.Auth) > 0 && len(model.Auth) <= 0 {
			model.Auth = m.Auth
		}
		logger.Debug("accExec2", "model", model)
		_, err = orm.NewQuery(db, &model).Where(fmt.Sprintf(" account_id='%s'", model.AccountId)).Update()
	}

	return err
}

func EventAccountAdd(db *pg.DB, logger log.Logger, evt *types.Event) {
	var AccountMsg Account
	err := eventutil.UnmarshalKVMap(evt.Attributes, &AccountMsg)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}

	m := makeAccountAddSql(AccountMsg)
	err = accExec(db, m, logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}
}

func EventAccountUpdate(db *pg.DB, logger log.Logger, evt *types.Event) {
	var AccountMsg Account
	err := eventutil.UnmarshalKVMap(evt.Attributes, &AccountMsg)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}

	m := makeAccountAddSql(AccountMsg)
	err = accExec(db, m, logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}
}
