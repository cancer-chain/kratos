package plugin

import (
	"encoding/json"
	"github.com/KuChainNetwork/kuchain/plugins"
	"github.com/KuChainNetwork/kuchain/plugins/db_history/chaindb"
	"github.com/KuChainNetwork/kuchain/plugins/types"
	"github.com/KuChainNetwork/kuchain/x/staking"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block

func getValidatorByConsAddr(ctx sdk.Context, consAcc sdk.ConsAddress, k staking.Keeper) staking.ValidatorI {
	validator := k.ValidatorByConsAddr(ctx, consAcc)
	return validator
}

var storageBlockHeight int64

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k staking.Keeper, codec *codec.Codec) {
	ctx.Logger().Debug("EndBlocker", "SyncBlockHeight:", chaindb.SyncBlockHeight)

	if req.Header.Height < 2 {
		return
	}

	storageBlockHeight = chaindb.SyncBlockHeight
	err, block, rTxs, rEvents, rTxEvents, rFeeEvents := types.GetBlockTxInfo(ctx, storageBlockHeight+1, codec)
	if err == nil {
		proposerValidator := getValidatorByConsAddr(ctx, ctx.BlockHeader().ProposerAddress, k)
		bz, _ := json.Marshal(proposerValidator)

		plugins.HandleBeginBlock(
			ctx,
			types.ReqBeginBlock{
				RequestBeginBlock: block,
				Tx:                rTxs,
				Events:            rEvents,
				TxEvents:          rTxEvents,
				FeeEvents:         rFeeEvents,
				ValidatorInfo:     string(bz),
				Time:              block.Time,
			},
		)
		ctx.Logger().Debug("BeginBlocker",
			"proposerValidator:", proposerValidator, "proposer:", string(bz))
	} else {
		ctx.Logger().Error("BeginBlocker", "GetBlockTxInfo err:", err)
	}
}

func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	if req.Height < 2 {
		return []abci.ValidatorUpdate{}
	}

	plugins.HandleEndBlock(ctx, types.ReqEndBlock{Height: storageBlockHeight + 1})
	return []abci.ValidatorUpdate{}
}
