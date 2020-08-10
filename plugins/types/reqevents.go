package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type ReqEvents struct {
	BlockHeight int64      `json:"block_height"`
	Events      sdk.Events `json:"events"`
}
