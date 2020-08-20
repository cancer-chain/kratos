package asset

import (
	"encoding/json"
	"github.com/KuChainNetwork/kuchain/chain/types"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/plugins"
	types2 "github.com/KuChainNetwork/kuchain/plugins/types"
	assettypes "github.com/KuChainNetwork/kuchain/x/asset/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

// InitGenesis account genesis init
func InitGenesis(ctx sdk.Context, ak Keeper, bz json.RawMessage) {
	logger := ak.Logger(ctx)

	var data GenesisState
	ModuleCdc.MustUnmarshalJSON(bz, &data)

	logger.Debug("init genesis", "module", ModuleName, "data", data)

	for _, a := range data.GenesisCoins {
		logger.Info("init genesis asset coin", "accountID", a.GetCreator(), "coins", a.GetSymbol(), "maxsupply:", a.GetMaxSupply())

		initSupply := types.NewCoin(a.GetMaxSupply().Denom, sdk.ZeroInt())

		err := ak.Create(ctx, a.GetCreator(), a.GetSymbol(), a.GetMaxSupply(), true, true, true, 0, initSupply, []byte{}) // TODO: genesis coins support opt
		if err != nil {
			panic(err)
		}

		ctx.EventManager().EmitEvent(
			chainTypes.NewEvent(ctx,
				assettypes.EventTypeCreate,
				sdk.NewAttribute(sdk.AttributeKeyModule, assettypes.AttributeValueCategory),
				sdk.NewAttribute(assettypes.AttributeKeyCreator, a.GetCreator().String()),
				sdk.NewAttribute(assettypes.AttributeKeySymbol, a.GetSymbol().String()),
				sdk.NewAttribute(assettypes.AttributeKeySupply, initSupply.String()),
				sdk.NewAttribute(assettypes.AttributeKeyMaxSupply, a.GetMaxSupply().String()),
				sdk.NewAttribute(assettypes.AttributeKeyCanIssue, strconv.FormatBool(true)),
				sdk.NewAttribute(assettypes.AttributeKeyCanLock, strconv.FormatBool(true)),
				sdk.NewAttribute(assettypes.AttributeKeyIssueCreateHeight, strconv.Itoa(1)),
				sdk.NewAttribute(assettypes.AttributeKeyIssueToHeight, strconv.Itoa(0)),
				sdk.NewAttribute(assettypes.AttributeKeyInit, initSupply.String()),
				sdk.NewAttribute(assettypes.AttributeKeyDescription, a.GetDescription()),
				sdk.NewAttribute(assettypes.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
			),
		)
	}

	for _, a := range data.GenesisAssets {
		logger.Info("init genesis account asset", "accountID", a.GetID(), "coins", a.GetCoins())
		err := ak.GenesisCoins(ctx, a.GetID(), a.GetCoins())
		if err != nil {
			panic(err)
		}

		for _, coin := range a.GetCoins() {
			creator, symbol, err := types.CoinAccountsFromDenom(coin.Denom)
			if err != nil {
				panic(err)
			}

			ctx.EventManager().EmitEvent(
				chainTypes.NewEvent(ctx,
					assettypes.EventTypeIssue,
					sdk.NewAttribute(sdk.AttributeKeyModule, assettypes.AttributeValueCategory),
					sdk.NewAttribute(assettypes.AttributeKeyCreator, creator.String()),
					sdk.NewAttribute(assettypes.AttributeKeySymbol, symbol.String()),
					sdk.NewAttribute(assettypes.AttributeKeyAmount, coin.String()),
				),
			)

			ctx.EventManager().EmitEvent(
				chainTypes.NewEvent(ctx,
					assettypes.EventTypeTransfer,
					sdk.NewAttribute(sdk.AttributeKeyModule, chainTypes.KuCodeSpace),
					sdk.NewAttribute(assettypes.AttributeKeyFrom, creator.String()),
					sdk.NewAttribute(assettypes.AttributeKeyTo, a.GetID().String()),
					sdk.NewAttribute(assettypes.AttributeKeyAmount, coin.String()),
				),
			)
		}
	}

	plugins.HandleEvent(ctx, types2.ReqEvents{BlockHeight: ctx.BlockHeight(), Events: ctx.EventManager().Events()})
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, ak Keeper) GenesisState {
	return GenesisState{}
}

// GenesisBalancesIterator implements genesis account iteration.
type GenesisBalancesIterator struct{}

// IterateGenesisBalances iterates over all the genesis accounts found in
// appGenesis and invokes a callback on each genesis account. If any call
// returns true, iteration stops.
func (GenesisBalancesIterator) IterateGenesisBalances(
	cdc *codec.Codec, appState types.AppGenesisState, cb func(GenesisAsset) (stop bool),
) {
	var gs GenesisState
	err := types.LoadGenesisStateFromBytes(cdc, appState, ModuleName, &gs)
	if err != nil {
		panic(err)
	}

	for _, a := range gs.GenesisAssets {
		if cb(a) {
			break
		}
	}
}
