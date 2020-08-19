package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	ReqBeginBlock
}

// NewMsgBeginBlock create begin block msg for plugin handler
func NewMsgBeginBlock(req ReqBeginBlock) *MsgBeginBlock {
	return &MsgBeginBlock{
		ReqBeginBlock: req,
	}
}

// MsgEndBlock end block msg for plugin handler
type MsgEndBlock struct {
	ReqEndBlock
}

// NewMsgEndBlock create end block msg for plugin
func NewMsgEndBlock(req ReqEndBlock) *MsgEndBlock {
	return &MsgEndBlock{
		ReqEndBlock: req,
	}
}
