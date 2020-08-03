package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// MsgEvent event msg for plugin handler
type MsgEvent struct {
	Evt Event
}

// NewMsgEvent new msg event
func NewMsgEvent(evt sdk.Event, height int64, hashCode string) *MsgEvent {
	return &MsgEvent{
		Evt: FromSdkEvent(evt, height, hashCode),
	}
}

// MsgStdTx stdTx msg for plugin handler
type MsgStdTx struct {
	Tx ReqTx
}

// NewMsgStdTx creates a new msg
func NewMsgStdTx(tx ReqTx) *MsgStdTx {
	return &MsgStdTx{
		Tx: tx, // no need deep copy as it will not be changed
	}
}

// MsgBeginBlock begin block msg for plugin handler
type MsgBeginBlock struct {
	ReqBlock
}

// NewMsgBeginBlock create begin block msg for plugin handler
func NewMsgBeginBlock(req ReqBlock) *MsgBeginBlock {
	return &MsgBeginBlock{
		ReqBlock: req,
	}
}

// MsgEndBlock end block msg for plugin handler
type MsgEndBlock struct {
	abci.RequestEndBlock
}

// NewMsgEndBlock create end block msg for plugin
func NewMsgEndBlock(req abci.RequestEndBlock) *MsgEndBlock {
	return &MsgEndBlock{
		RequestEndBlock: req,
	}
}
