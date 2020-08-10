package plugins

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	dbHistory "github.com/KuChainNetwork/kuchain/plugins/db_history"
	"github.com/KuChainNetwork/kuchain/plugins/test"
	"github.com/KuChainNetwork/kuchain/plugins/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// TODO: use a goroutine
var (
	plugins *Plugins
)

func InitPlugins(ctx Context, cfgs []BaseCfg) error {
	plugins = NewPlugins(ctx.Logger().With("module", "plugins"))
	for _, cfg := range cfgs {
		initPlugin(ctx, cfg, plugins)
	}

	plugins.Start()

	return nil
}

func StopPlugins(ctx Context) {
	if plugins != nil {
		plugins.Stop(ctx)
	}
}

func initPlugin(ctx Context, cfg BaseCfg, plugins *Plugins) {
	switch cfg.Name {
	case test.PluginName:
		plugins.RegPlugin(ctx, test.NewTestPlugin(ctx, cfg))
	case dbHistory.PluginName:
		plugins.RegPlugin(ctx, dbHistory.New(ctx, cfg))
	}
}

// HandleEvent plugins handler Events
func HandleEvent(ctx sdk.Context, evts types.ReqEvents) {

	if plugins == nil {
		return
	}

	for _, evt := range evts.Events {
		if ctx.BlockHeight() <= 0 {
			bz, _ := json.Marshal(evt)

			sh256 := sha256.New()
			sh256.Write(bz)
			hCode := hex.EncodeToString(sh256.Sum([]byte("")))

			plugins.EmitEvent(evt, evts.BlockHeight, hCode)
		}
	}
}
func HandleEventFromBlock(ctx sdk.Context, evts types.ReqEvents) {

	if plugins == nil {
		return
	}

	for _, evt := range evts.Events {
		bz, _ := json.Marshal(evt)

		sh256 := sha256.New()
		sh256.Write(bz)
		hCode := hex.EncodeToString(sh256.Sum([]byte("")))
		plugins.EmitEvent(evt, evts.BlockHeight, hCode)
	}
}

// HandleTx handler tx for each plugins
func HandleTx(ctxSdk sdk.Context, tx types.ReqTx) {
	if plugins == nil {
		return
	}

	plugins.EmitTx(tx)
}

// HandleBeginBlock emit begin block req to history plugin
func HandleBeginBlock(ctx sdk.Context, req types.ReqBlock) {

	if plugins == nil {
		return
	}

	plugins.EmitPluginMsg(types.NewMsgBeginBlock(req))
}

// HandleEndBlock emit end block req to history plugin
func HandleEndBlock(ctx sdk.Context, req abci.RequestEndBlock) {
	if plugins == nil {
		return
	}

	plugins.EmitPluginMsg(types.NewMsgEndBlock(req))
}
