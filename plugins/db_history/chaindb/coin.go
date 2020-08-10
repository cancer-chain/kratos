package chaindb

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

const (
	length18 = 18
	count18  = 1000000000000000000
)

type Coin struct {
	Amount      int64  `pg:"default:0" json:"amount"`
	AmountFloat int64  `pg:"default:0" json:"amount_float"`
	AmountStr   string `json:"amount_str"`
	Symbol      string `pg:"unique:as" json:"symbol"`
}

func NewCoin(coinStr string) (Coin, error) {
	amountStr, symbol, err := splitSymbol(coinStr)
	if err != nil {
		return Coin{}, errors.Wrap(err, "split symbol failed")
	}

	amount, amountFloat, err := parseAmountStr(amountStr)
	if err != nil {
		return Coin{}, errors.Wrap(err, "parse amount str failed")
	}

	return Coin{
		Amount:      amount,
		AmountFloat: amountFloat,
		AmountStr:   amountStr,
		Symbol:      symbol,
	}, nil
}

func NewNegativeCoin(coin Coin) Coin {
	return Coin{
		Amount:      -1 * coin.Amount,
		AmountFloat: -1 * coin.AmountFloat,
		AmountStr:   coin.AmountStr,
		Symbol:      coin.Symbol,
	}
}

func (c Coin) PrintTotalAmount() string {
	return c.AmountStr
}

func (c Coin) GetShortSymbol() string {
	s := strings.Split(c.Symbol, "/")
	if len(s) == 2 {
		return s[1]
	}
	return c.Symbol
}

func splitSymbol(coinStr string) (amountStr, symbol string, err error) {
	hasSymbol := false

	for i, v := range coinStr {
		if v < '0' || v > '9' {
			amountStr = coinStr[:i]
			symbol = coinStr[i:]
			hasSymbol = true
			break
		}
	}

	if len(amountStr) != 0 {
		return amountStr, symbol, nil
	}

	if hasSymbol {
		return "", "", errors.New("amount in coinStr is null")
	}

	return coinStr, "", nil
}

func parseAmountStr(amountStr string) (int64, int64, error) {
	if len(amountStr) <= length18 {
		amountFloat, err := strconv.ParseInt(amountStr, 10, 64)
		if err != nil {
			return 0, 0, err
		}

		return 0, amountFloat, nil
	}

	amount, err := strconv.ParseInt(amountStr[:len(amountStr)-length18], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	amountFloat, err := strconv.ParseInt(amountStr[len(amountStr)-length18:], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return amount, amountFloat, nil
}

func parseAmount(amount, amountFloat int64) (res string) {
	if amountFloat < count18 {
		res = getSA(amount) + getSF(amount, amountFloat)
	} else {
		res = getSA(amount+(amountFloat-count18)/count18) + getSF(amount, amountFloat-count18)
	}

	if len(strings.Trim(res, "0")) == 0 {
		return "0"
	}

	return res
}

func getSA(amount int64) string {
	if amount == 0 {
		return "0"
	}
	return strconv.FormatInt(amount, 10)
}

func getSF(amount, amountFloat int64) string {
	if amount == 0 {
		return strconv.FormatInt(amountFloat, 10)
	}

	return fmt.Sprintf("%0*d", 18, amountFloat)
}

func (c Coin) GetAmountSql(tabName string, amountName string, amountFloatName string) string {
	return fmt.Sprintf("%s=%s.%s-1+%d+(%s.%s+%d+%d)/%d", amountName, tabName, amountName, c.Amount, tabName, amountFloatName, count18, c.AmountFloat, count18)
}

func (c Coin) GetAmountFloatSql(talName string, amountFloatName string) string {
	return fmt.Sprintf("%s=(%s.%s+ %d +%d) %% %d", amountFloatName, talName, amountFloatName, c.AmountFloat, count18, count18)
}

func CoinAdd(a int64, af int64, b int64, bf int64) (int64, int64) {
	f := (count18 + af + bf) / count18
	af = (count18 + af + bf) % count18
	a = a - 1 + b + f

	return a, af
}

func CoinSub(a int64, af int64, b int64, bf int64) (int64, int64) {
	f := (count18 + af - bf) / count18
	af = (count18 + af - bf) % count18
	a = a - 1 - b + f

	return a, af
}
