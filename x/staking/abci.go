package staking

import (
	"github.com/KuChainNetwork/kuchain/plugins"
	ptypes "github.com/KuChainNetwork/kuchain/plugins/types"
	"github.com/KuChainNetwork/kuchain/x/staking/keeper"
	"github.com/KuChainNetwork/kuchain/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker will persist the current header and validator set as a historical entry
// and prune the oldest entry based on the HistoricalEntries parameter
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	k.TrackHistoricalInfo(ctx)
}

// Called every block, update validator set
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	r := k.BlockValidatorUpdates(ctx)

	es := sdk.Events{}
	validators := k.GetAllValidators(ctx)
	for _, validator := range validators {
		es = append(es, keeper.MakeValidatorEvent(ctx, types.EventTypeEndValidator, validator))
	}

	plugins.HandleEventFromBlock(ctx, ptypes.ReqEvents{BlockHeight: ctx.BlockHeight(), Events: es})
	return r
}
