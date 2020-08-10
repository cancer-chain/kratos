package keeper

import (
	"github.com/KuChainNetwork/kuchain/chain/constants"
	chaintype "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (a AssetKeeper) PayFee(ctx sdk.Context, payer types.AccountID, fee types.Coins) error {
	ctx.Logger().Debug("pay fee", "payer", payer, "fee", fee)
	if err := a.CoinsToPower(ctx, payer, constants.GetFeeCollector(), fee); err != nil {
		return sdkerrors.Wrap(err, "pay fee")
	}

	ctx.EventManager().EmitEvent(
		chaintype.NewEvent(ctx,
			types.EventTypePayFee,
			sdk.NewAttribute(sdk.AttributeKeyAmount, fee.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, payer.String()),
			sdk.NewAttribute(types.AttributeKeyTo, constants.GetFeeCollector().String()),
		),
	)

	ctx.Logger().Debug("PayFee",
		"amount", fee.String(), "from", payer.String(), "to", constants.GetFeeCollector().String())

	return nil
}
