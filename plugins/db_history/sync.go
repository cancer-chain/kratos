package dbHistory

import (
	"fmt"
	"github.com/KuChainNetwork/kuchain/plugins/db_history/types"
	types2 "github.com/KuChainNetwork/kuchain/plugins/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	btypes "github.com/tendermint/tendermint/types"
)

// dbMsgs4Block all msgs for a block, plugin will commit for all
type dbMsgs4Block struct {
	beginReq types2.ReqBeginBlock

	endReq types2.ReqEndBlock

	skip   bool
	events map[string]types.Event
	txs    []types2.ReqTx
	msgs   []sdk.Msg
}

func NewDBMsgs4Block(startHeight int64) dbMsgs4Block {
	return dbMsgs4Block{
		beginReq: types2.ReqBeginBlock{
			RequestBeginBlock: btypes.Block{
				Header: btypes.Header {
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
	return d.beginReq.RequestBeginBlock.Header.Height
}

func (d *dbMsgs4Block) Begin(ctx types.Context, req types2.ReqBeginBlock) {

	height := d.BlockHeight()
	reqHeight := req.RequestBeginBlock.Header.Height

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
		panic(fmt.Errorf("block height no match in begin %d %s", height, req.RequestBeginBlock.Header.LastBlockID.String()))
	}

	d.beginReq = req
}

func (d *dbMsgs4Block) End(ctx types.Context, req types2.ReqEndBlock) {
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
