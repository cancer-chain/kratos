package chaindb_test

import (
	"fmt"
	"github.com/KuChainNetwork/kuchain/plugins/db_history/config"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/stretchr/testify/require"

	"testing"
	"time"
)

func TestTx(t *testing.T) {
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

	type TestTab struct {
		tableName struct{} `pg:"TestTab,alias:TestTab"` // default values are the same
		ID        int64    // both "Id" and "ID" are detected as primary key

		Height int64 `pg:"default:0" json:"height"`
	}

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

	m := TestTab{
		Height: 1,
	}

	i := int64(0)
	tx, _ := db.Begin()

	for ; i < 1000000; i++ {
		m.Height = i
		err := orm.NewQuery(db, &m).Where(fmt.Sprintf("height='%d'", m.Height-1)).Select()
		if err != nil {

			tx2, _ := db.Begin()
			orm.Insert(db, &m)
			tx2.Commit()
		} else {
			require.NoError(t, err)

			tx2, _ := db.Begin()
			_, err = orm.NewQuery(db, &m).Where(fmt.Sprintf("height='%d'", m.Height-1)).Update()
			require.NoError(t, err)
			tx2.Commit()

			fmt.Sprintf("111222")
		}
	}
	tx.Commit()

	for ; i < 2000000; i++ {
		m.Height = i
		err := orm.NewQuery(db, &m).Where(fmt.Sprintf("height='%d'", m.Height-1)).Select()
		if err != nil {
			tx2, _ := db.Begin()
			orm.Insert(db, &m)
			tx2.Commit()
		} else {
			require.NoError(t, err)

			tx2, _ := db.Begin()
			_, err = orm.NewQuery(db, &m).Where(fmt.Sprintf("height='%d'", m.Height-1)).Update()

			tx2.Commit()
			require.NoError(t, err)

			if i == 2345 {
				panic(m)
			}
		}
	}
	tx.Commit()

}
