package chaindb

import (
	"fmt"
	"github.com/KuChainNetwork/kuchain/plugins/db_history/types"
	"github.com/KuChainNetwork/kuchain/utils/eventutil"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/tendermint/tendermint/libs/log"
	"strconv"
	"strings"
)

type ProposerReward struct {
	Height    int64  `json:"height"`
	Validator string `json:"validator"`
	Amount    string `json:"amount"`
	Time      string `json:"block_time"`
}

type EventProposerRewardModel struct {
	tableName struct{} `pg:"proposer_reward,alias:proposer_reward"` // default values are the same

	ID int // both "Id" and "ID" are detected as primary key

	Height      int64  `pg:"default:0" json:"height"`
	Validator   string `pg:"unique:apu" json:"validator"`
	Amount      int64  `pg:"default:0" json:"amount"`
	AmountFloat int64  `pg:"default:0" json:"amount_float"`
	AmountStr   string `json:"amount_str"`
	Symbol      string `json:"symbol"`
	Time        string `json:"time"`
}

func proposerRewardExec(db *pg.DB, model EventProposerRewardModel, logger log.Logger) error {
	var m EventProposerRewardModel

	err := orm.NewQuery(db, &m).Where(fmt.Sprintf("validator='%s'", model.Validator)).Select()
	if err != nil {
		logger.Debug("proposerRewardExec", "model", model)
		err = orm.Insert(db, &model)
	} else {
		//model.Amount, model.AmountFloat = CoinAdd(model.Amount, model.AmountFloat, m.Amount, m.AmountFloat)
		logger.Debug("proposerRewardExec2", "model", model)
		_, err = orm.NewQuery(db, &model).Where(fmt.Sprintf("validator='%s'", model.Validator)).Update()
	}

	if err == nil {
		_, err = orm.NewQuery(db, &model).
			Where(fmt.Sprintf("validator='%s'", model.Validator)).
			Set(fmt.Sprintf("amount=%d, amount_float=%d", model.Amount, model.AmountFloat)).Update()
	}
	return err
}

func EventProposerRewardAdd(db *pg.DB, logger log.Logger, evt *types.Event) {
	type attributes struct {
		Height    string `json:"height"`
		Validator string `json:"validator"`
		Amount    string `json:"amount"`
		Time      string `json:"block_time"`
	}
	var attr attributes

	err := eventutil.UnmarshalKVMap(evt.Attributes, &attr)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}

	var reward EventProposerRewardModel
	reward.Height, err = strconv.ParseInt(attr.Height, 10, 64)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}

	reward.AmountStr = attr.Amount
	reward.Validator = attr.Validator
	reward.Time = attr.Time

	amountStrList := strings.Split(attr.Amount, ".")
	amountStr := amountStrList[0]
	if len(amountStrList) > 1 {
		amountStr = amountStrList[len(amountStrList)-1]
	}

	_, symbol, err := splitSymbol(amountStr)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}

	reward.Symbol = symbol
	reward.Amount, reward.AmountFloat, err = parseAmountStr(amountStrList[0])
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
		return
	}

	err = proposerRewardExec(db, reward, logger)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}
}
