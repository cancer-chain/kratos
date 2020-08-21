package types

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/KuChainNetwork/kuchain/chain/constants"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/kv"
	"time"

	chaintype "github.com/KuChainNetwork/kuchain/chain/types"
	ptypes "github.com/KuChainNetwork/kuchain/x/plugin/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	rpcclient "github.com/tendermint/tendermint/rpc/client/local"
	btypes "github.com/tendermint/tendermint/types"
)

type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type EventLog struct {
	Type       string      `json:"type"`
	Attributes []Attribute `json:"attributes"`
}

type eventsLog struct {
	MsgIndex int64      `json:"msg_index"`
	Log      string     `json:"log"`
	Events   []EventLog `json:"events"`
}

type rawLog struct {
	Code      uint32      `json:"code"`
	Data      string      `json:"data"`
	Log       []eventsLog `json:"log"`
	Info      string      `json:"info"`
	GasWanted int64       `json:"gas_wanted"`
	GasUsed   int64       `json:"gas_used"`
	CodeSpace string      `json:"code_space"`
}

type fee struct {
	Amount string `json:"amount" yaml:"amount"`
	Gas    uint64 `json:"gas" yaml:"gas"`
	Payer  string `json:"payer" yaml:"payer"`
}

func (f fee) ToString() string {
	bz, _ := json.Marshal(f)
	return string(bz)
}

type Txm struct {
	Height     int64                    `json:"height"`
	TxHash     []byte                   `json:"tx_hash"`
	Msgs       []json.RawMessage        `json:"msg"`
	Fee        fee                      `json:"fee"`
	Signatures []chaintype.StdSignature `json:"signatures"`
	Memo       string                   `json:"memo"`
	RawLog     rawLog                   `json:"raw_log"`
	Time       string                   `json:"time"`
	Senders    []string                 `json:"senders"`
}

type ReqTx struct {
	Txm
}

func PrettifyJSON(ctx sdk.Context, tx chaintype.StdTx, Cdc *codec.Codec) ([]json.RawMessage, error) {
	alias := struct {
		Msgs []json.RawMessage `json:"msgs"`
	}{
		Msgs: make([]json.RawMessage, 0, len(tx.Msgs)),
	}

	for _, msg := range tx.Msgs {
		if msg, ok := msg.(chaintype.Prettifier); ok {
			raw, err := msg.PrettifyJSON(Cdc)
			if err != nil {
				return nil, sdkerrors.Wrapf(err, "prettify json to msg")
			}
			alias.Msgs = append(alias.Msgs, raw)
		}
	}

	return alias.Msgs, nil
}

func RebuildTx(ctx sdk.Context, stdTx chaintype.StdTx, Cdc *codec.Codec,
	Height int64, Time time.Time, hash []byte, rawLog json.RawMessage) (btx Txm) {

	if Cdc == nil {
		return
	}

	json.Unmarshal(rawLog, &btx.RawLog)

	btx.Msgs, _ = PrettifyJSON(ctx, stdTx, Cdc)
	btx.Height = Height
	btx.Time = Time.Format("2006-01-02T15:04:05.999999999Z")
	btx.Memo = stdTx.Memo
	btx.TxHash = hash
	btx.Signatures = stdTx.Signatures
	btx.Senders = GetSenders(stdTx, Cdc)
	btx.Fee = fee{
		Amount: stdTx.Fee.Amount.String(),
		Gas:    stdTx.Fee.Gas,
		Payer:  stdTx.Fee.Payer.String(),
	}

	ctx.Logger().Debug("RebuildTx",
		"hash", base64.StdEncoding.EncodeToString(hash), "btx", btx)

	return
}

func makeEventForTxm(aEvent abci.Event) sdk.Event {
	return sdk.Event{
		Type:       aEvent.Type,
		Attributes: aEvent.Attributes,
	}
}
func makeFeeEvent(stdTx chaintype.StdTx, height int64, time2 time.Time) (Event sdk.Event) {

	Event.Type = "payfee"
	Event.Attributes = append(Event.Attributes, kv.Pair{
		Key:   []byte("amount"),
		Value: []byte(stdTx.Fee.Amount.String()),
	})
	Event.Attributes = append(Event.Attributes, kv.Pair{
		Key:   []byte("from"),
		Value: []byte(stdTx.Fee.Payer.String()),
	})
	Event.Attributes = append(Event.Attributes, kv.Pair{
		Key:   []byte("to"),
		Value: []byte(constants.GetFeeCollector().String()),
	})
	Event.Attributes = append(Event.Attributes, kv.Pair{
		Key:   []byte("height"),
		Value: []byte(fmt.Sprintf("%d", height)),
	})
	Event.Attributes = append(Event.Attributes, kv.Pair{
		Key:   []byte("block_time"),
		Value: []byte(time2.Format("2006-01-02T15:04:05.999999999Z")),
	})

	return
}

func PrintEventsLog(ctx sdk.Context, events sdk.Events, Height int64, desc string) {
	logEvents := ""
	for _, e := range events {
		logEvents += e.Type + ","
		for _, ar := range e.Attributes {
			logEvents += " " + string(ar.Key) + ":" + string(ar.Value) + " "
		}
		logEvents += ";" + "\n"
	}
	ctx.Logger().Debug("PrintEventsLog", "block_height", Height, "events", logEvents, "desc", desc)
}

func GetBlockTxInfo(ctx sdk.Context, Height int64, Cdc *codec.Codec) (err error, block btypes.Block,
	rTxs []ReqTx, rEvents ReqEvents, rTxEvents []ReqEvents, rFeeEvents ReqEvents) {

	if ptypes.PNode == nil {
		err = sdkerrors.New("PNode is nil", 2, fmt.Sprintf("block height %d", Height))
		return
	}

	blockPtr := ptypes.PNode.BlockStore().LoadBlock(Height)
	if blockPtr == nil {
		err = sdkerrors.New("block is nil", 1, fmt.Sprintf("block height %d", Height))
		return
	}

	block = *blockPtr
	var events sdk.Events
	var feeEvents sdk.Events

	getEvent := func() () {
		ResTx, err := rpcclient.New(ptypes.PNode).BlockResults(&Height)
		if err != nil {
			ctx.Logger().Error("getTx", "err", err)
			return
		}
		for i := 0; i < len(ResTx.BeginBlockEvents); i++ {
			events = append(events, makeEventForTxm(ResTx.BeginBlockEvents[i]))
		}
		return
	}

	getTx := func() (raws []json.RawMessage, codes []uint32) {
		ResTx, err := rpcclient.New(ptypes.PNode).BlockResults(&Height)
		if err != nil {
			ctx.Logger().Error("getTx", "err", err)
			return
		}

		for i := 0; i < len(ResTx.TxsResults); i++ {
			tr := rawLog{
				Code:      ResTx.TxsResults[i].Code,
				Data:      string(ResTx.TxsResults[i].Data),
				Info:      ResTx.TxsResults[i].Info,
				GasWanted: ResTx.TxsResults[i].GasWanted,
				GasUsed:   ResTx.TxsResults[i].GasUsed,
				CodeSpace: ResTx.TxsResults[i].Codespace,
			}

			json.Unmarshal([]byte(ResTx.TxsResults[i].Log), &tr.Log)

			bz, _ := json.Marshal(tr)
			raws = append(raws, bz)
			codes = append(codes, tr.Code)

			var txEvents sdk.Events
			for _, l := range tr.Log {
				for _, e := range l.Events {
					evt := sdk.Event{
						Type: e.Type,
					}
					for _, n := range e.Attributes {
						kv := kv.Pair{
							Key:   []byte(n.Key),
							Value: []byte(n.Value),
						}
						evt.Attributes = append(evt.Attributes, kv)
					}
					txEvents = append(txEvents, evt)
				}
			}

			if tr.Code == 0 { // tx is suc
				ev := ReqEvents{
					BlockHeight: block.Height,
					Events:      txEvents,
				}
				rTxEvents = append(rTxEvents, ev)
			}
		}

		ctx.Logger().Debug("getTx", "block_height", Height, "raws", raws)
		return
	}

	getTxInfo := func() {
		raws, _ := getTx()

		for i := 0; i < len(block.Data.Txs); i++ {
			var stdTx chaintype.StdTx
			err = Cdc.UnmarshalBinaryLengthPrefixed(block.Data.Txs[i], &stdTx)
			if err == nil {
				rtx := ReqTx{Txm: RebuildTx(ctx, stdTx, Cdc, block.Height, block.Time, block.Data.Txs[i].Hash(), raws[i])}
				rTxs = append(rTxs, rtx)
				feeEvents = append(feeEvents, makeFeeEvent(stdTx, block.Height, block.Time))
			} else {
				ctx.Logger().Error("getTxInfo", "err", err)
			}
		}

		getEvent()

		rEvents = ReqEvents{
			BlockHeight: block.Height,
			Events:      events,
		}

		rFeeEvents = ReqEvents{
			BlockHeight: block.Height,
			Events:      feeEvents,
		}

		PrintEventsLog(ctx, events, Height, "Events")
		for _, nEvent := range rTxEvents {
			PrintEventsLog(ctx, nEvent.Events, Height, "txEvents")
		}
		PrintEventsLog(ctx, feeEvents, Height, "feeEvents")
	}

	getTxInfo()

	return
}

func GetSenders(tx chaintype.StdTx, Cdc *codec.Codec) (senders []string) {
	for _, msg := range tx.Msgs {
		if msg, ok := msg.(chaintype.KuMsgDataHandler); ok {
			sender, err := msg.GetSender(Cdc)
			if err != nil {
				panic(fmt.Sprintf("get sender failed: %v", err))
			}

			senders = append(senders, sender.String())
		}
	}

	return senders
}
