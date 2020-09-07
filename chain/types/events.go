package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
	"sync/atomic"
)

var sum uint64 = 0

func NewEvent(ctx sdk.Context, ty string, attrs ...sdk.Attribute) sdk.Event {
	e := sdk.Event{Type: ty}

	for _, attr := range attrs {
		e.Attributes = append(e.Attributes, attr.ToKVPair())
	}

	kv := sdk.Attribute{
		Key:   "height",
		Value: fmt.Sprintf("%d", ctx.BlockHeight()),
	}
	e.Attributes = append(e.Attributes, kv.ToKVPair())

	kv = sdk.Attribute{
		Key:   "event_id",
		Value: strconv.FormatUint(atomic.AddUint64(&sum, 1), 10),
	}
	e.Attributes = append(e.Attributes, kv.ToKVPair())

	kv = sdk.Attribute{
		Key:   "block_time",
		Value: ctx.BlockTime().Format("2006-01-02T15:04:05.999999999Z"),
	}

	e.Attributes = append(e.Attributes, kv.ToKVPair())
	return e
}

func EventToString(e sdk.Event) string {
	log := ""

	for _, n := range e.Attributes {
		log += string(n.Key) + ","
		log += string(n.Value) + "\n"
	}

	return log
}
func EventsToString(es sdk.Events) string {
	log := ""

	for _, e := range es {
		log += EventToString(e) + "\n\n"
	}

	return log
}
