package plugin

import (
	"encoding/json"
	"github.com/KuChainNetwork/kuchain/plugins"
	"github.com/KuChainNetwork/kuchain/plugins/types"
	"github.com/KuChainNetwork/kuchain/x/staking"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block

func PluginsGetValidatorByConsAddr(ctx sdk.Context, consAcc sdk.ConsAddress, k staking.Keeper) staking.ValidatorI {
	validator := k.ValidatorByConsAddr(ctx, consAcc)
	return validator
}

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k staking.Keeper, codec *codec.Codec) {
	proposerValidator := PluginsGetValidatorByConsAddr(ctx, ctx.BlockHeader().ProposerAddress, k)
	bz, _ := json.Marshal(proposerValidator)

	time := types.GetTxInfo(ctx, req.Header.Height, codec, plugins.HandleTx, plugins.HandleEventFromBlock)
	plugins.HandleBeginBlock(ctx,
		types.ReqBlock{
			RequestBeginBlock: req,
			ValidatorInfo:     string(bz),
			Time:              time,
		},
	)

	ctx.Logger().Debug("BeginBlocker",
		"proposerValidator:", proposerValidator, "proposer:", string(bz))
}

func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	plugins.HandleEndBlock(ctx, req)

	return []abci.ValidatorUpdate{}
}
