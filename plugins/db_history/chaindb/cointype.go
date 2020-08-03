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
	Init              string `json:"int"`
	Amount            int64  `pg:"default:0",json:"amount"`
	AmountFloat       int64  `pg:"default:0",json:"amount_float"`
	Module            string `json:"module"`
	CanLock           string `json:"can_lock"`
	Creator           string `pg:"unique:cs",json:"creator"`
	Symbol            string `pg:"unique:cs",json:"symbol"`
	IssueCreateHeight string `json:"issue_create_height"`
	Height            int64  `json:"height"`
	CanIssue          string `json:"can_issue"`
	IssueToHeight     string `json:"issue_to_height"`
	Desc              string `json:"_desc"`
	Time              string `json:"time"`
}

func makeCtpSql(model EventCoinType) CreateCoinTypeModel {
	coin, _ := NewCoin(model.Amount)
	if len(model.Supply) > 0 {
		coin, _ = NewCoin(model.Amount)
	}

	if len(coin.Symbol) <= 0 {
		var coinMax Coin
		if len(model.Max) > 0 {
			coinMax, _ = NewCoin(model.Max)
			coin.Symbol = coinMax.Symbol
		}
	}

	q := CreateCoinTypeModel{
		Max:               model.Max,
		Init:              model.Init,
		Amount:            coin.Amount,
		AmountFloat:       coin.AmountFloat,
		Symbol:            coin.Symbol,
		Module:            model.Module,
		CanLock:           model.CanLock,
		Creator:           model.Creator,
		IssueCreateHeight: model.IssueCreateHeight,
		Height:            model.Height,
		CanIssue:          model.CanLock,
		IssueToHeight:     model.IssueToHeight,
		Desc:              model.Desc,
		Time:              model.Time,
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
		model.Amount, model.AmountFloat = CoinAdd(model.Amount, model.AmountFloat, m.Amount, m.AmountFloat)
		model.Max = m.Max
		model.Desc = m.Desc
		model.IssueToHeight = m.IssueToHeight
		model.IssueCreateHeight = m.IssueCreateHeight
		model.CanIssue = m.CanIssue
		model.Creator = m.Creator

		logger.Debug("etExec2", "model", model)
		_, err = orm.NewQuery(db, &model).Where(fmt.Sprintf("Symbol='%s' and creator='%s'", model.Symbol, model.Creator)).Update()
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
	err = etExec(db, m, logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}
}

func EventCoinTypeModifySupply(db *pg.DB, logger log.Logger, evt *types.Event) {
	var CoinTypeMsg EventCoinType
	err := eventutil.UnmarshalKVMap(evt.Attributes, &CoinTypeMsg)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}

	m := makeCtpSql(CoinTypeMsg)
	err = etExec(db, m, logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}
}
