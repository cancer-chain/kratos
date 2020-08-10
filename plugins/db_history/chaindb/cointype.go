package chaindb

import (
	"fmt"
	"github.com/KuChainNetwork/kuchain/plugins/db_history/types"
	"github.com/KuChainNetwork/kuchain/utils/eventutil"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/tendermint/tendermint/libs/log"
)

type EventCoinType struct {
	Max               string `json:"max"`
	Init              string `json:"init"`
	Amount            string `json:"amount"`
	Supply            string `json:"supply"`
	Module            string `json:"module"`
	CanLock           string `json:"canLock"`
	Creator           string `json:"creator"`
	Symbol            string `json:"symbol"`
	IssueCreateHeight string `json:"issueCreateHeight"`
	Height            int64  `json:"height"`
	CanIssue          string `json:"canIssue"`
	IssueToHeight     string `json:"issueToHeight"`
	Desc              string `json:"desc"`
	Time              string `json:"block_time"`
}

type CreateCoinTypeModel struct {
	tableName struct{} `pg:"Coins,alias:Coins"` // default values are the same

	ID int // both "Id" and "ID" are detected as primary key

	Max               string `json:"max"`
	Init              string `json:"init"`
	Amount            int64  `pg:"default:0" json:"amount"`
	AmountFloat       int64  `pg:"default:0" json:"amount_float"`
	Module            string `json:"module"`
	CanLock           string `json:"can_lock"`
	Creator           string `pg:"unique:cs" json:"creator"`
	Symbol            string `pg:"unique:cs" json:"symbol"`
	IssueCreateHeight string `json:"issue_create_height"`
	Height            int64  `pg:"default:0" json:"height"`
	CanIssue          string `json:"can_issue"`
	IssueToHeight     string `json:"issue_to_height"`
	Desc              string `json:"desc"`
	Time              string `json:"time"`
}

func makeCtpSql(model EventCoinType) CreateCoinTypeModel {
	coin, _ := NewCoin(model.Supply)
	if len(model.Supply) <= 0 {
		coin, _ = NewCoin(model.Amount)
	}

	q := CreateCoinTypeModel{
		Max:               model.Max,
		Init:              model.Init,
		Amount:            coin.Amount,
		AmountFloat:       coin.AmountFloat,
		Symbol:            model.Symbol,
		Module:            model.Module,
		CanLock:           model.CanLock,
		Creator:           model.Creator,
		IssueCreateHeight: model.IssueCreateHeight,
		Height:            model.Height,
		CanIssue:          model.CanIssue,
		IssueToHeight:     model.IssueToHeight,
		Desc:              model.Desc,
		Time:              model.Time,
	}

	if len(model.Symbol) <= 0 {
		q.Symbol = coin.Symbol
	}

	return q
}

func etExec(db *pg.DB, model CreateCoinTypeModel, logger log.Logger) error {
	var m CreateCoinTypeModel
	err := orm.NewQuery(db, &m).Where(fmt.Sprintf("Symbol='%s' and creator='%s'", model.Symbol, model.Creator)).Select()
	if err != nil {
		logger.Debug("etExec1", "model", model)
		err = orm.Insert(db, &model)
	} else {
		if len(m.Max) > 0 {
			model.Max = m.Max
		}
		if len(m.Init) > 0 {
			model.Init = m.Init
		}
		if len(m.CanLock) > 0 {
			model.CanLock = m.CanLock
		}
		if len(m.Creator) > 0 {
			model.Creator = m.Creator
		}
		if len(m.CanIssue) > 0 {
			model.CanIssue = m.CanIssue
		}
		if len(m.Module) > 0 {
			model.Module = m.Module
		}
		if len(m.IssueCreateHeight) > 0 {
			model.IssueCreateHeight = m.IssueCreateHeight
		}
		if len(m.IssueToHeight) > 0 {
			model.IssueToHeight = m.IssueToHeight
		}
		if len(m.Desc) > 0 {
			model.Desc = m.Desc
		}
		if len(m.Symbol) > 0 {
			model.Symbol = m.Symbol
		}

		model.Amount, model.AmountFloat = CoinAdd(model.Amount, model.AmountFloat, m.Amount, m.AmountFloat)

		logger.Debug("etExec2", "model", model)
		_, err = orm.NewQuery(db, &model).Where(fmt.Sprintf("Symbol='%s' and creator='%s'", model.Symbol, model.Creator)).Update()
	}
	if err == nil {
		_, err = orm.NewQuery(db, &model).
			Where(fmt.Sprintf("Symbol='%s' and creator='%s'", model.Symbol, model.Creator)).
			Set(fmt.Sprintf("amount=%d, amount_float=%d", model.Amount, model.AmountFloat)).Update()
	}
	return err
}

func EventCoinTypeAdd(db *pg.DB, logger log.Logger, evt *types.Event) {
	var CoinTypeMsg EventCoinType
	err := eventutil.UnmarshalKVMap(evt.Attributes, &CoinTypeMsg)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}

	m := makeCtpSql(CoinTypeMsg)
	logger.Debug("EventCoinTypeAdd", "\n", *evt, "\n", CoinTypeMsg, "\n", m)

	tx, _ := db.Begin()
	err = etExec(db, m, logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}
	tx.Commit()
}

func EventCoinTypeModifySupply(db *pg.DB, logger log.Logger, evt *types.Event) {
	var CoinTypeMsg EventCoinType
	err := eventutil.UnmarshalKVMap(evt.Attributes, &CoinTypeMsg)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}

	m := makeCtpSql(CoinTypeMsg)
	logger.Debug("EventCoinTypeModifySupply", "\n", *evt, "\n", CoinTypeMsg, "\n", m)

	tx, _ := db.Begin()
	err = etExec(db, m, logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}
	tx.Commit()
}
