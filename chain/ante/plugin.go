package ante

import (
	"time"

	"github.com/KuChainNetwork/kuchain/plugins"
	"github.com/KuChainNetwork/kuchain/plugins/types"
	"github.com/KuChainNetwork/kuchain/x/plugin/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PluginHandlerDecorator struct {
}

func NewPluginHandlerDecorator() PluginHandlerDecorator {
	return PluginHandlerDecorator{}
}

func (isd PluginHandlerDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	ctx.Logger().Debug("plugin ante handler")

	// no need to increment sequence on CheckTx or RecheckTx
	if ctx.IsCheckTx() && !simulate {
		return next(ctx, tx, simulate)
	}

	if ctx.BlockHeight() == 0 {
		if std, ok := tx.(StdTx); ok {
			tx := types.RebuildTx(ctx, std, keeper.PCdec, ctx.BlockHeight(), time.Time{}, []byte(""), []byte(""))
			plugins.HandleTx(ctx, types.ReqTx{Txm: tx})
		}
	}

	return next(ctx, tx, simulate)
}
