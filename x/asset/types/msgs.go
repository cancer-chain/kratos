package types

import (
	"github.com/KuChainNetwork/kuchain/chain/msg"
	"github.com/KuChainNetwork/kuchain/chain/types"
)

// RouterKey is they name of the asset module
const RouterKey = ModuleName

var (
	RouterKeyName                 = types.MustName(RouterKey)
	_, _, _, _, _ types.KuMsgData = (*MsgCreateCoinData)(nil), (*MsgIssueCoinData)(nil), (*MsgBurnCoinData)(nil), (*MsgLockCoinData)(nil), (*MsgUnlockCoinData)(nil)
)

type (
	KuMsg = types.KuMsg
)

// NewMsgTransfer create msg transfer
func NewMsgTransfer(auth types.AccAddress, from types.AccountID, to types.AccountID, amount Coins) KuMsg {
	return *msg.MustNewKuMsg(
		RouterKeyName,
		msg.WithAuth(auth),
		msg.WithTransfer(from, to, amount),
	)
}

type MsgCreateCoin struct {
	types.KuMsg
}

type MsgCreateCoinData struct {
	Symbol        Name   `json:"symbol" yaml:"symbol"`                             // Symbol coin symbol name
	Creator       Name   `json:"creator" yaml:"creator"`                           // Creator coin creator account name
	MaxSupply     Coin   `json:"max_supply" yaml:"max_supply"`                     // MaxSupply coin max supply limit
	CanIssue      bool   `json:"can_issue,omitempty" yaml:"can_issue"`             // CanIssue if the coin can issue after create
	CanLock       bool   `json:"can_lock,omitempty" yaml:"can_lock"`               // CanLock if the coin can lock by user
	IssueToHeight int64  `json:"issue_to_height,omitempty" yaml:"issue_to_height"` // IssueToHeight if this is not zero, creator only can issue this
	InitSupply    Coin   `json:"init_supply" yaml:"init_supply"`                   // InitSupply coin init supply, if issue_to_height is not zero, this will be the start supply for issue
	Desc          []byte `json:"desc" yaml:"desc"`                                 // Description
}

func (MsgCreateCoinData) Type() types.Name { return types.MustName("create@asset") }

func (msg MsgCreateCoinData) Sender() AccountID {
	return NewAccountIDFromName(msg.Creator)
}

// NewMsgCreate new create coin msg
func NewMsgCreate(auth types.AccAddress, creator types.Name, symbol types.Name, maxSupply types.Coin, canIssue, canLock bool, issue2Height int64, initSupply types.Coin, desc []byte) MsgCreateCoin {
	return MsgCreateCoin{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgCreateCoinData{
				Creator:       creator,
				Symbol:        symbol,
				MaxSupply:     maxSupply,
				CanIssue:      canIssue,
				CanLock:       canLock,
				IssueToHeight: issue2Height,
				InitSupply:    initSupply,
				Desc:          desc,
			}),
		),
	}
}

type MsgIssueCoin struct {
	types.KuMsg
}

type MsgIssueCoinData struct {
	Symbol  Name `json:"symbol" yaml:"symbol"`   // Symbol coin symbol name
	Creator Name `json:"creator" yaml:"creator"` // Creator coin creator account name
	Amount  Coin `json:"amount" yaml:"amount"`   // MaxSupply coin max supply limit
}

// Type imp for data KuMsgData
func (MsgIssueCoinData) Type() types.Name { return types.MustName("issue") }

func (msg MsgIssueCoinData) Sender() AccountID {
	return NewAccountIDFromName(msg.Creator)
}

// NewMsgIssue new issue msg
func NewMsgIssue(auth types.AccAddress, creator, symbol types.Name, amount types.Coin) MsgIssueCoin {
	return MsgIssueCoin{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgIssueCoinData{
				Creator: creator,
				Symbol:  symbol,
				Amount:  amount,
			}),
		),
	}
}

type MsgBurnCoin struct {
	types.KuMsg
}

type MsgBurnCoinData struct {
	Id     AccountID `json:"id" yaml:"id"`         // Symbol coin symbol name
	Amount Coin      `json:"amount" yaml:"amount"` // MaxSupply coin max supply limit
}

// Type imp for data KuMsgData
func (MsgBurnCoinData) Type() types.Name { return types.MustName("burn") }

func (msg MsgBurnCoinData) Sender() AccountID {
	return msg.Id
}

// NewMsgBurn new issue msg
func NewMsgBurn(auth types.AccAddress, id types.AccountID, amount types.Coin) MsgIssueCoin {
	return MsgIssueCoin{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgBurnCoinData{
				Id:     id,
				Amount: amount,
			}),
		),
	}
}

// MsgLockCoin msg to lock coin
type MsgLockCoin struct {
	types.KuMsg
}

type MsgLockCoinData struct {
	Id                AccountID `json:"id" yaml:"id"`                                         // Id lock account
	Amount            Coins     `json:"amount" yaml:"amount"`                                 // Amount coins to lock
	UnlockBlockHeight int64     `json:"unlockBlockHeight,omitempty" yaml:"unlockBlockHeight"` // UnlockBlockHeight the block height the coins unlock
}

// Type imp for data KuMsgData
func (m *MsgLockCoinData) Type() types.Name { return types.MustName("lock@coin") }

func (m MsgLockCoinData) Sender() AccountID {
	return m.Id
}

// NewMsgLockCoin create new lock coin msg
func NewMsgLockCoin(auth types.AccAddress, id types.AccountID, amount types.Coins, unlockBlockHeight int64) MsgLockCoin {
	return MsgLockCoin{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgLockCoinData{
				Id:                id,
				Amount:            amount,
				UnlockBlockHeight: unlockBlockHeight,
			}),
		),
	}
}

// MsgUnlockCoin msg to unlock coin
type MsgUnlockCoin struct {
	types.KuMsg
}

type MsgUnlockCoinData struct {
	Id     AccountID `json:"id" yaml:"id"`         // Id lock account
	Amount Coins     `json:"amount" yaml:"amount"` // Amount coins to lock
}

// Type imp for data KuMsgData
func (m *MsgUnlockCoinData) Type() types.Name { return types.MustName("unlock@coin") }

func (m MsgUnlockCoinData) Sender() AccountID {
	return m.Id
}

// NewMsgUnlockCoin create new lock coin msg
func NewMsgUnlockCoin(auth types.AccAddress, id types.AccountID, amount types.Coins) MsgUnlockCoin {
	return MsgUnlockCoin{
		*msg.MustNewKuMsg(
			RouterKeyName,
			msg.WithAuth(auth),
			msg.WithData(Cdc(), &MsgUnlockCoinData{
				Id:     id,
				Amount: amount,
			}),
		),
	}
}
