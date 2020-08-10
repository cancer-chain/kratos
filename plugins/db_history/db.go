package dbHistory

import (
	"sync"

	"github.com/KuChainNetwork/kuchain/plugins/db_history/chaindb"
	"github.com/KuChainNetwork/kuchain/plugins/db_history/config"
	"github.com/go-pg/pg/v10"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

type dbWork struct {
	msg interface{}
}

func (w dbWork) IsStopped() bool {
	return w.msg == nil
}

func (w dbWork) IsEndBlock() (bool, int64) {
	if msg, ok := w.msg.(abci.RequestEndBlock); ok {
		return true, msg.Height
	}

	return false, 0
}

type dbService struct {
	logger      log.Logger
	database    *pg.DB
	errDatabase *pg.DB

	stat dbMsgs4Block
	sync *SyncState

	dbChan chan dbWork
	wg     sync.WaitGroup
}

func (dbs *dbService) GetBD() *pg.DB {
	return dbs.database
}

// NewDB create a connection commit event to db
func NewDB(cfg config.Cfg, logger log.Logger) *dbService {
	res := &dbService{
		database: pg.Connect(&pg.Options{
			Addr:     cfg.DB.Address,
			User:     cfg.DB.User,
			Password: cfg.DB.Password,
			Database: cfg.DB.Database,
		}),

		errDatabase: pg.Connect(&pg.Options{
			Addr:     cfg.DB.Address,
			User:     cfg.DB.User,
			Password: cfg.DB.Password,
			Database: cfg.DB.Database,
		}),

		logger: logger,
		dbChan: make(chan dbWork, 512),
	}

	if err := createSchema(res.database, res.logger); err != nil {
		panic(err)
	}

	res.sync = NewChainSyncStat(res.database, logger)
	res.stat = NewDBMsgs4Block(res.sync.BlockNum)

	return res
}

func (db *dbService) Start() error {
	db.logger.Info("Starting database service")

	db.wg.Add(1)
	go func() {
		defer db.wg.Done()

		for {
			work, ok := <-db.dbChan
			if !ok {
				db.logger.Info("msg channel closed")
				return
			}

			if work.IsStopped() {
				db.logger.Info("db service stopped")
				return
			}

			if ok, height := work.IsEndBlock(); ok {
				if err := UpdateChainSyncStat(db.database, db.logger, height, db.sync.ChainID); err != nil {
					db.logger.Error("UpdateChainSyncStat error", "err", err)
				}
				// no need process
				continue
			}

			if err := db.Process(&work); err != nil {
				db.logger.Error("db process error", "err", err)
			}
		}
	}()
	return nil
}

func (db *dbService) Process(work *dbWork) error {
	if err := chaindb.Process(db.database, db.logger, work.msg); err != nil {
		return err
	}

	return nil
}

func (db *dbService) Emit(work dbWork) {
	db.dbChan <- work
}

func (db *dbService) Stop() error {
	db.logger.Info("Stopping database service")

	db.dbChan <- dbWork{}
	db.wg.Wait()

	db.logger.Info("Database service stopped")

	db.database.Close()

	db.logger.Info("Database connection closed")
	return nil
}
