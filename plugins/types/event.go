package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Event struct {
	BlockHeight int64
	Type        string
	Attributes  map[string]string

	HashCode string
}

func FromSdkEvent(evt sdk.Event, height int64, hashCode string) Event {
	res := Event{
		BlockHeight: height,
		Type:        evt.Type,
		Attributes:  make(map[string]string, len(evt.Attributes)),
		HashCode:    hashCode,
	}

	for _, attr := range evt.Attributes {
		res.Attributes[string(attr.Key)] = string(attr.Value)
	}

	return res
}
