package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

type PluginMsgHandler func(ctx Context, msg sdk.Msg)
type PluginTxHandler func(ctx Context, tx ReqTx)
type PluginEvtHandler func(ctx Context, evt Event)

type Plugin interface {
	Init(Context) error
	Start(Context) error
	Stop(Context) error

	OnBlockBegin(ctx Context, req ReqBeginBlock)
	OnBlockEnd(ctx Context, req ReqEndBlock)

	EvtHandler() PluginEvtHandler
	MsgHandler() PluginMsgHandler
	TxHandler() PluginTxHandler

	Logger() log.Logger
	Name() string
}
