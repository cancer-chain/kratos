package account

import (
	"encoding/json"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/plugins"
	types2 "github.com/KuChainNetwork/kuchain/plugins/types"
	"github.com/KuChainNetwork/kuchain/x/account/exported"
	"github.com/KuChainNetwork/kuchain/x/account/types"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis account genesis init
func InitGenesis(ctx sdk.Context, ak Keeper, data json.RawMessage) {
	logger := ak.Logger(ctx)

	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)

	for _, a := range genesisState.Accounts {
		logger.Info("init genesis account", "name", a.GetName(), "auth", a.GetAuth())
		ak.SetAccount(ctx, ak.NewAccount(ctx, a))

		// ensure auth init
		ak.EnsureAuthInited(ctx, a.GetAuth())

		ctx.EventManager().EmitEvents(sdk.Events{
			chainTypes.NewEvent(ctx,
				types.EventTypeCreateAccount,
				sdk.NewAttribute(types.AttributeKeyCreator, types.AttributeValueCreator),
				sdk.NewAttribute(types.AttributeKeyAccount, a.GetID().String()),
				sdk.NewAttribute(types.AttributeKeyAuth, a.GetAuth().String()),
				sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
			),
		})

		if _, ok := a.GetID().ToName(); ok {
			ak.AddAccountByAuth(ctx, a.GetAuth(), a.GetName().String())
		}
	}

	plugins.HandleEvent(ctx, types2.ReqEvents{BlockHeight: ctx.BlockHeight(), Events: ctx.EventManager().Events()})
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, ak Keeper) GenesisState {
	var genAccounts exported.GenesisAccounts
	ak.IterateAccounts(ctx, func(account exported.Account) bool {
		genAccounts = append(genAccounts, account.(exported.GenesisAccount))
		return false
	})

	return GenesisState{
		Accounts: genAccounts,
	}
}
