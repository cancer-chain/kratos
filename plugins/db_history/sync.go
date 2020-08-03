package dbHistory

import (
	"fmt"
	types2 "github.com/KuChainNetwork/kuchain/plugins/types"

	"github.com/KuChainNetwork/kuchain/plugins/db_history/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	ChainIdx = 1
)

// SyncState sync state in pg database
type SyncState struct {
	tableName struct{} `pg:"sync_stat,alias:sync_stat"` // default values are the same

	ID       int // both "Id" and "ID" are detected as primary key
	BlockNum int64
	ChainID  string `pg:",unique"`
}

func NewChainSyncStat(db *pg.DB, logger log.Logger) *SyncState {
	stat := &SyncState{
		ID: ChainIdx,
	}
	if err := db.Select(stat); err != nil {
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

func UpdateChainSyncStat(db *pg.DB, logger log.Logger, num int64, chainID string) error {
	stat := &SyncState{
		ID: ChainIdx,
	}
	err := db.Select(stat)
	if err != nil {
		return errors.Wrapf(err, "get sync stat err")
	}

	logger.Info("get sync stat", "num", stat.BlockNum)

	return db.Update(&SyncState{
		ID:       ChainIdx,
		BlockNum: num,
		ChainID:  chainID,
	})
}

// dbMsgs4Block all msgs for a block, plugin will commit for all
type dbMsgs4Block struct {
	beginReq types2.ReqBlock

	endReq abci.RequestEndBlock

	skip   bool
	events map[string]types.Event
	txs    []types2.ReqTx
	msgs   []sdk.Msg
}

func NewDBMsgs4Block(startHeight int64) dbMsgs4Block {
	return dbMsgs4Block{
		beginReq: types2.ReqBlock{
			RequestBeginBlock: abci.RequestBeginBlock{
				Header: abci.Header{
					Height: startHeight,
				},
			},
		},
		skip: false,

		events: make(map[string]types.Event),
		txs:    make([]types2.ReqTx, 0, 256),
		msgs:   make([]sdk.Msg, 0, 1024),
	}
}

func (d *dbMsgs4Block) BlockHeight() int64 {
	return d.beginReq.Header.Height
}

func (d *dbMsgs4Block) Begin(ctx types.Context, req types2.ReqBlock) {

	height := d.BlockHeight()
	reqHeight := req.Header.GetHeight()

	ctx.Logger().Debug("msgs begin block", "req", reqHeight, "curr", height)

	if reqHeight <= height && height > 0 {
		d.skip = true

		ctx.Logger().Debug("skip by heght")

		d.events = make(map[string]types.Event)
		d.txs = d.txs[0:0]
		d.msgs = d.msgs[0:0]

		return
	} else {
		d.skip = false
	}

	if (height + 1) != reqHeight {
		panic(fmt.Errorf("block height no match in begin %d %s", height, req.Header.LastBlockId.String()))
	}

	d.beginReq = req
}

func (d *dbMsgs4Block) End(ctx types.Context, req abci.RequestEndBlock) {
	height := d.BlockHeight()

	ctx.Logger().Debug("end for block", "height", height, "req", req.Height)

	d.events = make(map[string]types.Event)
	d.txs = d.txs[0:0]
	d.msgs = d.msgs[0:0]

	if req.Height < height {
		return
	}

	if height != req.Height {
		panic(fmt.Errorf("block height no match in end %d %d", height, req.Height))
	}

	d.endReq = req
}

func (d *dbMsgs4Block) AppendEvent(evt types.Event) {
	_, ok := d.events[evt.HashCode]
	if ok {
		return
	}
	d.events[evt.HashCode] = evt
}

func (d *dbMsgs4Block) AppendTx(tx types2.ReqTx) {

	d.txs = append(d.txs, tx)
}

func (d *dbMsgs4Block) AppendMsg(msg sdk.Msg) {
	d.msgs = append(d.msgs, msg)
}
