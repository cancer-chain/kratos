package assetSync

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/valyala/fasthttp"
)

// AssetSync
type AssetSync struct {
	host string
	port int
}

// Sync sync user coins
func (object *AssetSync) Sync(user, denom string,
	timeout time.Duration) (err error, coins chainTypes.Coins) {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()
	req.SetRequestURI(fmt.Sprintf("http://%s:%d/assets/coins/%s",
		object.host,
		object.port,
		user))
	if err = fasthttp.DoTimeout(req, res, timeout); nil != err {
		return
	}
	type Response struct {
		Height string `json:"height"`
		Result []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"result"`
	}
	var response Response
	if err = json.Unmarshal(res.Body(), &response); nil != err {
		return
	}
	for _, result := range response.Result {
		var amount int64
		amount, err = strconv.ParseInt(result.Amount, 10, 64)
		if nil != err {
			return
		}
		if denom != result.Denom {
			continue
		}
		coins = append(coins, chainTypes.NewInt64Coin(result.Denom, amount))
	}
	return
}

func New(host string, port int) *AssetSync {
	return &AssetSync{
		host: host,
		port: port,
	}
}
