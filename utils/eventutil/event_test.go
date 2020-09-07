package eventutil_test

import (
	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	wallet = simapp.NewWallet()
)

type testEvent struct {
	Name   types.Name       `json:"name"`
	Names  []types.Name     `json:"names"`
	Id     types.AccountID  `json:"id"`
	IdAddr types.AccountID  `json:"idaddr"`
	Auth   types.AccAddress `json:"auth"`
	Coin   types.Coin       `json:"coin"`
	Coins  types.Coins      `json:"coins"`
	Str    string           `json:"str"`
	Int    int64            `json:"int"`
	Ints   []int64          `json:"ints"`
	Strs   []string         `json:"strs"`
	UInt   uint64           `json:"uint"`
	Bool1  bool             `json:"b1"`
	Bool2  bool             `json:"b2"`
}

var (
	testCoin  = types.NewInt64Coin(constants.DefaultBondDenom, 100001)
	testCoins = types.NewCoins(
		types.NewInt64Coin(constants.DefaultBondDenom, 100),
		types.NewInt64Coin("foo/test1", 200),
		types.NewInt64Coin("foo/test2", 300),
		types.NewInt64Coin("foo/test3", 400),
	)

	testName   = types.MustName("adss@sdssd")
	testNames  = []types.Name{types.MustName("abcde"), types.MustName("aa@aaa"), types.MustName("kuc.ha@in")}
	testIDName = types.MustAccountID("adss@sd2s")
	testAuth   = wallet.NewAccAddress()
	testIDAddr = types.NewAccountIDFromAccAdd(wallet.NewAccAddress())

	testEvt sdk.Event = sdk.NewEvent("test",
		sdk.NewAttribute("name", testName.String()),
		sdk.NewAttribute("names", "abcde, aa@aaa,kuc.ha@in"),
		sdk.NewAttribute("id", testIDName.String()),
		sdk.NewAttribute("idaddr", testIDAddr.String()),
		sdk.NewAttribute("auth", testAuth.String()),
		sdk.NewAttribute("coin", testCoin.String()),
		sdk.NewAttribute("coins", testCoins.String()),
		sdk.NewAttribute("str", "test str for event"),
		sdk.NewAttribute("int", "1234567"),
		sdk.NewAttribute("uint", "7654321"),
		sdk.NewAttribute("strs", "test,str, for, event"),
		sdk.NewAttribute("ints", "-12,34,56,7"),
		sdk.NewAttribute("b1", "true"),
		sdk.NewAttribute("b2", "false"),
	)
)

type EventValidator1 struct {
	tableName struct{} `pg:"Validator,alias:Validator"` // default values are the same

	ID int // both "Id" and "ID" are detected as primary key

	Height  int64  `pg:"default:0" json:"height"`
	Address string `pg:"unique:as" json:"address"`
	Sender  string `pg:"unique:as" json:"sender"`
}

var (
	H    int64 = 11239
	Addr       = "abcde"
	Send       = "cdasd"

	VEvent = sdk.NewEvent("test1",
		sdk.NewAttribute("height", "11239"),
		sdk.NewAttribute("address", "abcde"),
		sdk.NewAttribute("sender", "cdasd"),
	)
)
