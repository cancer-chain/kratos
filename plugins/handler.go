package plugins

import (
	"errors"
	"sync"

	"github.com/KuChainNetwork/kuchain/plugins/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

type pluginMsg interface{}

// Plugins a handler for all plugins to reg
type Plugins struct {
	plugins     []Plugin
	txHandlers  []types.PluginTxHandler
	msgHandlers []types.PluginMsgHandler
	evtHandlers []types.PluginEvtHandler

	msgChan chan pluginMsg
	closed  bool
	logger  log.Logger
	wg      sync.WaitGroup
}

func NewPlugins(logger log.Logger) *Plugins {
	return &Plugins{
		msgChan: make(chan pluginMsg, 512),
		closed:  false,
		logger:  logger,
	}
}

func (p *Plugins) RegPlugin(ctx Context, plugin Plugin) {
	plugin.Logger().Info("init plugin", "name", plugin.Name())

	for _, p := range p.plugins {
		if p.Name() == plugin.Name() {
			panic(errors.New("plugin reg two times"))
		}
	}

	if err := plugin.Init(ctx); err != nil {
		panic(err)
	}

	p.plugins = append(p.plugins, plugin)

	if tx := plugin.TxHandler(); tx != nil {
		p.txHandlers = append(p.txHandlers, tx)
	}

	if msg := plugin.MsgHandler(); msg != nil {
		p.msgHandlers = append(p.msgHandlers, msg)
	}

	if evt := plugin.EvtHandler(); evt != nil {
		p.evtHandlers = append(p.evtHandlers, evt)
	}
}

func (p *Plugins) onTx(ctx types.Context, tx types.ReqTx) {
	for _, h := range p.txHandlers {
		h(ctx, tx)
	}
}

func (p *Plugins) onEvent(ctx types.Context, evt types.Event) {
	for _, h := range p.evtHandlers {
		h(ctx, evt)
	}
}

func (p *Plugins) onBeginBlock(ctx types.Context, msg *types.MsgBeginBlock) {
	ctx.Logger().Info("on begin block", "header", msg.Header)
	for _, plugin := range p.plugins {
		plugin.OnBlockBegin(ctx, msg.ReqBlock)
	}
}

func (p *Plugins) onEndBlock(ctx types.Context, msg *types.MsgEndBlock) {
	ctx.Logger().Info("on end block", "height", msg.GetHeight(), "pnum", len(p.plugins))
	for _, plugin := range p.plugins {
		plugin.OnBlockEnd(ctx, msg.RequestEndBlock)
	}
}

func (p *Plugins) Start() {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()

		ctx := types.NewContext(p.logger)
		for _, p := range p.plugins {
			if err := p.Start(ctx); err != nil {
				panic(err)
			}
		}

		for {
			msg, ok := <-p.msgChan
			if !ok {
				p.logger.Info("msg channel closed")
				return
			}

			if msg == nil {
				p.logger.Info("stop channel")
				return
			}

			ctx := NewContext(p.logger)

			switch msg := msg.(type) {
			case *types.MsgEvent:
				p.onEvent(ctx, msg.Evt)
			case *types.MsgStdTx:
				p.onTx(ctx, msg.Tx)
			case *types.MsgBeginBlock:
				p.onBeginBlock(ctx, msg)
			case *types.MsgEndBlock:
				p.onEndBlock(ctx, msg)
			}
		}
	}()
}

func (p *Plugins) EmitEvent(evt sdk.Event, height int64, hashCode string) {
	p.msgChan <- types.NewMsgEvent(evt, height, hashCode)
}

func (p *Plugins) EmitTx(tx types.ReqTx) {
	p.msgChan <- types.NewMsgStdTx(tx)
}

func (p *Plugins) EmitPluginMsg(msg interface{}) {
	p.msgChan <- msg
}

func (p *Plugins) Stop(ctx types.Context) {
	p.msgChan <- nil
	p.wg.Wait()

	for _, plg := range p.plugins {
		plg.Stop(ctx)
	}
}
