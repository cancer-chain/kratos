package dbHistory

import (
	"github.com/KuChainNetwork/kuchain/plugins/db_history/types"
	types2 "github.com/KuChainNetwork/kuchain/plugins/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (t *plugin) OnEvent(ctx types.Context, evt types.Event) {

	t.db.stat.AppendEvent(evt)
	t.logger.Info("on event", "type", evt.Type, "events num", len(t.db.stat.events))
}

func (t *plugin) OnTx(ctx types.Context, tx types2.ReqTx) {
	t.logger.Info("on tx", "tx", tx)
	t.db.stat.AppendTx(tx)
}

func (t *plugin) OnMsg(ctx types.Context, msg sdk.Msg) {
	t.logger.Info("on msg", "msg", msg)
	t.db.stat.AppendMsg(msg)
}

func (t *plugin) OnBlockBegin(ctx types.Context, req types2.ReqBeginBlock) {
	t.db.stat.Begin(ctx, req)
}

func (t *plugin) OnBlockEnd(ctx types.Context, req types2.ReqEndBlock) {
	logger := ctx.Logger()

	if t.db.stat.skip {
		return
	}

	logger.Info("on block end", "height", req.Height, "events num:", len(t.db.stat.events))

	t.db.Emit(dbWork{
		msg: t.db.stat.beginReq,
	})

	for _, evt := range t.db.stat.events {
		t.db.Emit(dbWork{
			msg: evt,
		})
		logger.Debug("OnBlockEnd", "evt", evt)
	}

	for _, tx := range t.db.stat.txs {
		t.db.Emit(dbWork{
			msg: tx,
		})
	}

	for _, msg := range t.db.stat.msgs {
		t.db.Emit(dbWork{
			msg: msg,
		})
	}

	t.db.Emit(dbWork{
		msg: req,
	})

	t.db.stat.End(ctx, req)
}
