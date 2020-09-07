package chaindb

import (
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	ChainIdx = 1
)

var SyncBlockHeight int64

// SyncState sync state in pg database
type SyncState struct {
	tableName struct{} `pg:"sync_stat,alias:sync_stat"` // default values are the same

	ID       int // both "Id" and "ID" are detected as primary key
	BlockNum int64
	ChainID  string `pg:",unique"`
}

func SelectSyncStat(db *pg.DB, logger log.Logger) (*SyncState, error) {
	stat := &SyncState{
		ID: ChainIdx,
	}
	if err := db.Select(stat); err != nil {
		return stat, err
	}
	return stat, nil
}

func NewChainSyncStat(db *pg.DB, logger log.Logger) *SyncState {
	stat, err := SelectSyncStat(db, logger)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			// need init
			if err := db.Insert(stat); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	return stat
}

func UpdateChainSyncStat(db *pg.DB, logger log.Logger, num int64) (*SyncState, error) {
	stat, err := SelectSyncStat(db, logger)
	if err != nil {
		return nil, errors.Wrapf(err, "get sync stat err")
	}

	logger.Info("UpdateChainSyncStat get sync stat", "BlockNum", stat.BlockNum)

	stat.BlockNum = num

	err = db.Update(stat)
	return stat, err
}
