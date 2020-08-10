package chaindb

import (
	"encoding/json"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/tendermint/tendermint/libs/log"
)

type Msg struct {
	Action string            `json:"action"`
	To     string            `json:"to"`
	Auth   []json.RawMessage `json:"auth"`
	Data   json.RawMessage   `json:"data"`
	From   string            `json:"from"`
	Amount []json.RawMessage `json:"amount"`
	Router string            `json:"router"`
}

type CreateTxMsgsModel struct {
	tableName struct{} `pg:"txmsgs,alias:txmsgs"` // default values are the same
	ID        int64    // both "Id" and "ID" are detected as primary key
	Height    int64    `pg:"default:0" json:"height"`
	TxId      int64    `pg:"btree:t" json:"tx_id"`
	Action    string   `json:"action"`
	To        string   `json:"to"`
	Auth      string   `json:"auth"`
	Data      string   `json:"data"`
	From      string   `json:"from"`
	Amount    string   `json:"amount"`
	Symbol    string   `json:"symbol"`
	Router    string   `json:"router"`
	Sender    string   `json:"sender"`
	Time      string   `json:"time"`
}

func buildTxMsg(logger log.Logger, m json.RawMessage, tx *txInDB, uid int64, sender string) (iMsg CreateTxMsgsModel) {
	var msg Msg
	json.Unmarshal(m, &msg)

	logger.Debug("InsertTxMsgs ", "msg", msg)

	iMsg.TxId = uid
	iMsg.Time = tx.Time
	iMsg.Height = tx.Height
	iMsg.Data = string(msg.Data)
	iMsg.To = msg.To
	iMsg.From = msg.From
	iMsg.Action = msg.Action
	iMsg.Router = msg.Router
	iMsg.Sender = sender

	bz, _ := json.Marshal(msg.Auth)
	iMsg.Auth = string(bz)

	for _, ad := range msg.Amount {
		var adn map[string]interface{}
		json.Unmarshal(ad, &adn)

		amount, ok := adn["amount"]
		if ok {
			iMsg.Amount += amount.(string) + " "
		}
		deNo, ok := adn["denom"]
		if ok {
			iMsg.Symbol += deNo.(string) + " "
		}
	}

	logger.Debug("buildTxMsg", "iMsg", iMsg)

	return
}

func InsertTxMsgs(db *pg.DB, logger log.Logger, tx *txInDB, tx_ *pg.Tx, uid int64) bool {

	for _, m := range tx.Msgs {
		iMsg := buildTxMsg(logger, m, tx, uid, "")
		err := orm.Insert(db, &iMsg)
		if err != nil {
			EventErr(db, logger, NewErrMsg(err))
		}
	}
	return true
}
