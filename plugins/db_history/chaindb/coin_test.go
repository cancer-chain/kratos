package chaindb

import (
	"testing"
)

func TestCoin(t *testing.T) {
	coin1, err := NewCoin("222222000000000000000001kuchain/kts")
	if err != nil {
		t.Fatal(err)
	}

	coin2, err := NewCoin("111111000000000000000002kuchain/kts")
	if err != nil {
		t.Fatal(err)
	}

	coin3, err := NewCoin("0kuchain/kts")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("-------- %s -------\n", coin1.PrintTotalAmount())
	t.Logf("-------- %s -------\n", coin2.PrintTotalAmount())
	t.Logf("-------- %d %d %s -------\n", coin3.Amount, coin3.AmountFloat, coin3.Symbol)

	coin4, err := NewCoin("338897515734811486kratos/kts")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("-------- %d %d %s -------\n", coin4.Amount, coin4.AmountFloat, coin4.Symbol)
}

func TestShortSymbol(t *testing.T) {
	coin1, err := NewCoin("222222000000000000000001kuchain/kts")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("-------- %s -------\n", coin1.GetShortSymbol())

	coin2, err := NewCoin("222222000000000000000001kts")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("-------- %s -------\n", coin2.GetShortSymbol())
}

func TestNoSymbol(t *testing.T) {
	coin2, err := NewCoin("222222000000000000000001")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("-------- %d ------- %d ---- %s\n", coin2.Amount, coin2.AmountFloat, coin2.Symbol)
}
