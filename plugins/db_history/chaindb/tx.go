package chaindb

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	ptypes "github.com/KuChainNetwork/kuchain/plugins/types"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/tendermint/tendermint/libs/log"
	"strings"
)

type txInDB struct {
	ptypes.ReqTx
}
type CreateTxModel struct {
	tableName struct{} `pg:"tx,alias:tx"` // default values are the same

	TxUid      int64  `json:"tx_uid"`
	Height     int64  `json:"height"`
	TxHash     string `json:"tx_hash"`
	Msgs       string `json:"msg"`
	Fee        string `json:"fee"`
	Signatures string `json:"signatures"`
	Memo       string `json:"memo"`
	RawLog     string `json:"raw_log"`
	Senders    string `json:"senders"`
	Time       string `json:"time"`
}

func newTxInDB(tx ptypes.ReqTx) *txInDB {
	return &txInDB{
		ReqTx: tx,
	}
}

type Signature struct {
	PubKey    string `json:"pub_key"`
	Signature string `json:"signature"`
}

func makeTxmSql(tm ptypes.ReqTx) CreateTxModel {

	bz, _ := json.Marshal(tm.Msgs)
	Msg := string(bz)

	if len(Msg) <= 0 {
		Msg = "{}"
	}

	Hash := strings.ToUpper(hex.EncodeToString(tm.TxHash))
	Fee := tm.Fee.ToString()
	if len(Fee) <= 0 {
		Fee = "{}"
	}

	snowNode, _ := NewSnowNode(0)
	Uid := snowNode.Generate().Int64()

	type signature struct {
		PubKey    string
		Signature []byte
	}
	var tmpSignatures []signature
	for _, p := range tm.Signatures {
		tmpSignatures = append(
			tmpSignatures,
			signature{
				PubKey:    base64.StdEncoding.EncodeToString(p.PubKey.Bytes()),
				Signature: p.Signature,
			},
		)
	}

	bz, _ = json.Marshal(tmpSignatures)
	Sins := string(bz)

	if len(Sins) <= 0 {
		Sins = "{}"
	}

	bz, _ = json.Marshal(tm.RawLog)
	rawLog := string(bz)
	if len(rawLog) <= 0 {
		rawLog = "{}"
	}

	bz, _ = json.Marshal(tm.Senders)
	Sender := string(bz)
	if len(Sender) <= 0 {
		Sender = "{}"
	}

	q := CreateTxModel{
		TxUid:      Uid,
		Height:     tm.Height,
		TxHash:     Hash,
		Msgs:       Msg,
		Fee:        Fee,
		Signatures: Sins,
		Memo:       tm.Memo,
		RawLog:     rawLog,
		Senders:    Sender,
		Time:       tm.Time,
	}

	return q
}

func InsertTxm(db *pg.DB, logger log.Logger, tx txInDB) (error, int64) {

	q := makeTxmSql(tx.ReqTx)
	err := orm.Insert(db, &q)

	logger.Debug("InsertTxm", "txm", q)
	return err, q.TxUid
}

func getEvent0(tm ptypes.ReqTx, logger log.Logger) (Events []ptypes.Event) {
	for _, l := range tm.RawLog.Log {
		for _, e := range l.Events {
			evt := ptypes.Event{
				BlockHeight: tm.Height,
				HashCode:    strings.ToUpper(hex.EncodeToString(tm.TxHash)),
				Type:        e.Type,
			}

			evt.Attributes = make(map[string]string)
			for _, kv := range e.Attributes {
				evt.Attributes[kv.Key] = kv.Value
			}
			Events = append(Events, evt)
		}
	}

	logger.Debug("makeEvent",
		"evts", Events, "raw", tm.RawLog, "RawCode", tm.RawLog.Code, "height", tm.Height)
	return
}

func InsertTxm0(db *pg.DB, logger log.Logger, tx *txInDB) error {
	Events := getEvent0(tx.ReqTx, logger)

	tx_, _ := db.Begin()
	q := makeTxmSql(tx.ReqTx)
	err := orm.Insert(db, &q)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}

	for _, m := range tx.Msgs {
		iMsg := buildTxMsg(logger, m, *tx, 0, "")
		err := orm.Insert(db, &iMsg)
		if err != nil {
			EventErr(db, logger, NewErrMsg(err))
		}
	}

	if tx.RawLog.Code == 0 { //fee
		for _, evt := range Events {
			err = InsertEvent(db, logger, &evt)
			if err != nil {
				EventErr(db, logger, NewErrMsg(err))
			}
		}
	}

	tx_.Commit()

	return nil
}
