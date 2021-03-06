package keeper

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
)

type (
	AssetTransfer = types.AssetTransfer
	Context       = types.Context
)

var (
	CoinAccountsFromDenom = types.CoinAccountsFromDenom
	CoinDenom             = types.CoinDenom
)

type (
	Coins    = types.Coins
	Coin     = types.Coin
	DecCoins = types.DecCoins
	DecCoin  = types.DecCoin
)

var (
	NewDec       = types.NewDec
	NewInt64Coin = types.NewInt64Coin
)
