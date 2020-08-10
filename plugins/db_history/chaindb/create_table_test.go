package chaindb_test

import (
	"fmt"
	"github.com/KuChainNetwork/kuchain/plugins/db_history/chaindb"
	"github.com/KuChainNetwork/kuchain/plugins/db_history/config"
	"github.com/KuChainNetwork/kuchain/utils/eventutil"
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

func TestCreateAccTable2(t *testing.T) {
	conf := config.Cfg{DB: config.DBCfg{
		Address:  "192.168.1.200:5432",
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

		Symbol:  "kcs",
		Account: "kuchain",
		Time:    "1234234234",
	}

	err := db.Insert(&m)
	require.NoError(t, err)

	m2 := chaindb.CreateAccCoinsModel{
		Height:      11100,
		Amount:      5000,
		AmountFloat: 4000,

		Symbol:  "kcs",
		Account: "kuchain",
		Time:    "xxxxxxxx",
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

func TestCreateAccTable(t *testing.T) {

	Attributes := make(map[string]string)
	Attributes["height"] = "100"
	Attributes["amount"] = "87917289371289"
	Attributes["amount_float"] = "9789"
	Attributes["symbol"] = "ktcs"
	Attributes["account"] = "dbxxx"
	Attributes["time"] = "2018-12-01"

	m2 := chaindb.CreateAccCoinsModel{}

	err := eventutil.UnmarshalKVMap(Attributes, &m2)
	require.NoError(t, err)

	conf := config.Cfg{DB: config.DBCfg{
		Address:  "192.168.1.200:5432",
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

	fmt.Println(m2)

	err = orm.NewQuery(db, &m2).Where(fmt.Sprintf("Symbol='%s' and account='%s'", m2.Symbol, m2.Account)).Select()
	if err != nil {
		orm.Insert(db, &m2)
	} else {
		require.NoError(t, err)

		m2.Amount += m2.Amount
		m2.AmountFloat += m2.AmountFloat

		_, err = orm.NewQuery(db, &m2).Where(fmt.Sprintf("Symbol='%s' and account='%s'", m2.Symbol, m2.Account)).Update()
		require.NoError(t, err)
	}

	time.Sleep(1 * time.Second)
}

func TestCreateAccTable3(t *testing.T) {

	Attributes := make(map[string]string)
	Attributes["max"] = "100"
	Attributes["init"] = "87917289371289"
	Attributes["amount"] = "232329789"
	Attributes["amount_float"] = "9789"
	Attributes["module"] = "module"
	Attributes["can_lock"] = "can_lock"
	Attributes["creator"] = "creator"
	Attributes["symbol"] = "ktcs"
	Attributes["issue_create_height"] = "issue_create_height"
	Attributes["height"] = "123123123"
	Attributes["can_issue"] = "can_issue"
	Attributes["issue_to_height"] = "issue_to_height"
	Attributes["desc"] = "desc"
	Attributes["time"] = "2018-12-01"

	m2 := chaindb.CreateCoinTypeModel{}

	err := eventutil.UnmarshalKVMap(Attributes, &m2)
	require.NoError(t, err)

	conf := config.Cfg{DB: config.DBCfg{
		Address:  "192.168.1.200:5432",
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
		(*chaindb.CreateCoinTypeModel)(nil),
	}

	for _, model := range models {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		require.NoError(t, err)
	}

	time.Sleep(1 * time.Second)

	fmt.Println(m2)
	orm.Insert(db, &m2)
}

func TestCreateAccTable4(t *testing.T) {

	type TestTab struct {
		tableName struct{} `pg:"TestTab,alias:TestTab"` // default values are the same

		ID int // both "Id" and "ID" are detected as primary key

		Num1 int64 `pg:"btree:num1"  json:"num_1" `
		Num2 int64 ` json:"num_2" pg:",type:bigint,default:0"`
		Num3 int64 ` json:"num_3" pg:",type:bigint,default:0"`
	}

	conf := config.Cfg{DB: config.DBCfg{
		Address:  "192.168.1.200:5432",
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
		(*TestTab)(nil),
	}

	for _, model := range models {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		require.NoError(t, err)
	}

	time.Sleep(1 * time.Second)

	m2 := TestTab{
		ID: 10000,
	}

	fmt.Println(m2)
	orm.Insert(db, &m2)

	m2.Num2 = 20000000
	m2.Num1 = 0
	m2.Num3 = 30000

	_, err := orm.NewQuery(db, &m2).Where(fmt.Sprintf("id='%d' ", m2.ID)).Update(&m2)
	//_, err = orm.NewQuery(db, &m2).Where(fmt.Sprintf("id='%d' ", m2.ID)).Set(fmt.Sprintf("num1=%d", 0)).Update()
	require.NoError(t, err)
}
