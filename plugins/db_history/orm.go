package dbHistory

import (
	"github.com/KuChainNetwork/kuchain/plugins/db_history/chaindb"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/tendermint/tendermint/libs/log"
)

// createSchema creates database schema for User and Story models.
func createSchema(db *pg.DB, logger log.Logger) error {
	if err := chaindb.RegOrm(db, logger); err != nil {
		return err
	}

	models := []interface{}{
		(*SyncState)(nil),
	}

	for _, model := range models {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			logger.Debug("createSchema", "model", model)
			return err
		}
	}
	return nil
}
