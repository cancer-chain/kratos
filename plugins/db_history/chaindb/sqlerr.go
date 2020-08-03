package chaindb

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/tendermint/tendermint/libs/log"
	"time"
)

type ErrMsg struct {
	tableName struct{} `pg:"ErrMsg,alias:ErrMsg"` // default values are the same
	ID        int      // both "Id" and "ID" are detected as primary key

	Message string `json:"message"`
	Time    string `json:"time"`
}

func NewErrMsg(err error) ErrMsg {
	e := ErrMsg{
		Message: err.Error(),
		Time:    time.Now().String(),
	}
	return e
}

func EventErr(db *pg.DB, logger log.Logger, errIfo ErrMsg) {
	err := orm.Insert(db, &errIfo)
	if err != nil && logger != nil {
		logger.Error("ErrTableAdd add table error", "err", err.Error())
	}
}
