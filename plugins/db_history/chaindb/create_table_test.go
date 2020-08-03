package chaindb_test

import (
	"fmt"
	"github.com/KuChainNetwork/kuchain/plugins/db_history/chaindb"
	"github.com/KuChainNetwork/kuchain/plugins/db_history/config"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
	"time"
)

type tLogger struct {
	log.Logger
}

func (t tLogger) Debug(msg string, keyvals ...interface{}) {
	fmt.Println(msg, keyvals)
}

func (t tLogger) Info(msg string, keyvals ...interface{}) {
	fmt.Println(msg, keyvals)
}
func (t tLogger) Error(msg string, keyvals ...interface{}) {
	fmt.Println(msg, keyvals)
}

func (t tLogger) With(keyvals ...interface{}) log.Logger {
	//fmt.Println(keyvals)
	return t
}

func TestCreateAccTable(t *testing.T) {
	conf := config.Cfg{DB: config.DBCfg{
		Address:  "192.168.1.78:5432",
		User:     "pguser",
		Password: "123456",
		Database: "kuchaindb",
	}}

	db := pg.Connect(&pg.Options{
		Addr:     conf.DB.Address,
		User:     conf.DB.User,
		Password: conf.DB.Password,
		Database: conf.DB.Database,
	})

	models := []interface{}{
		(*chaindb.CreateAccCoinsModel)(nil),
	}

	for _, model := range models {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		require.NoError(t, err)
	}

	time.Sleep(1 * time.Second)

	m := chaindb.CreateAccCoinsModel{
		Height:      0,
		Amount:      0,
		AmountFloat: 0,
		AmountStr:   "kuchan",
		Symbol:      "kcs",
		Account:     "kuchain",
		Time:        "1234234234",
	}

	err := db.Insert(&m)
	require.NoError(t, err)

	m2 := chaindb.CreateAccCoinsModel{
		Height:      11100,
		Amount:      5000,
		AmountFloat: 4000,
		AmountStr:   "kuchan",
		Symbol:      "kcs",
		Account:     "kuchain",
		Time:        "xxxxxxxx",
	}

	err = orm.NewQuery(db, &m).Where(fmt.Sprintf("Symbol='%s' and account='%s'", m2.Symbol, m2.Account)).Select()
	if err != nil {
		orm.Insert(db, &m2)
	} else {
		require.NoError(t, err)

		m2.Amount += m.Amount
		m2.AmountFloat += m.AmountFloat

		_, err = orm.NewQuery(db, &m2).Where(fmt.Sprintf("Symbol='%s' and account='%s'", m2.Symbol, m2.Account)).Update()
		require.NoError(t, err)
	}

	time.Sleep(1 * time.Second)
}
