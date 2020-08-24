package chaindb

import (
	"fmt"
	"os"
	"sync"
	"time"

	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/asset/assetSync"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

func RegOrm(db *pg.DB, logger log.Logger, syncStatus bool) error {
	models := []interface{}{
		(*eventInDB)(nil),
		(*MessageInDB)(nil),
		(*KuTransferInDB)(nil),
		(*blockInDB)(nil),
		(*CreateAccCoinsModel)(nil),
		(*CreateAccountModel)(nil),
		(*BlockInfo)(nil),
		(*CreateCoinTypeModel)(nil),
		(*CreateDelegationModel)(nil),
		(*CreateDelegationChangeModel)(nil),
		(*CreateLockAccCoinsModel)(nil),
		(*EventValidator)(nil),
		(*CreateTxModel)(nil),
		(*CreateTxMsgsModel)(nil),
		(*ErrMsg)(nil),
		(*EventProposerRewardModel)(nil),
		(*EventUnBondModel)(nil),
	}

	for _, model := range models {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}

	startSync(db, syncStatus)

	return nil
}

func startSync(db *pg.DB, syncStatus bool) {
	if !syncStatus {
		return
	}

	startOnce.Do(func() {
		type Response struct {
			Account string           `json:"account"`
			Amount  chainTypes.Coins `json:"amount"`
			Error   error            `json:"error"`
		}
		m := &CreateAccCoinsModel{
			SyncState: 0,
		}

		if _, err := orm.NewQuery(db, m).
			Set("sync_state = ?", m.SyncState).
			Where("sync_state = ?", 1).
			Update(); nil != err {
			panic(err)
		}

		go func() {
			duration := 1 * time.Second
			tick := time.NewTimer(duration)
			assetSyncTool := assetSync.New("localhost", 80)

			for {
				select {
				case <-tick.C:
					var list []*CreateAccCoinsModel
					err := orm.NewQuery(db, &list).Where("sync_state=?", 0).Select()
					if nil != err {
						panic(err)
					}

					for 0 < len(list) {
						chunkSize := 128
						if chunkSize > len(list) {
							chunkSize = len(list)
						}

						var wg sync.WaitGroup
						wg.Add(chunkSize)

						for i := 0; i < chunkSize; i++ {
							m := list[i]
							m.SyncState = 1
							if _, err = orm.NewQuery(db, m).
								Set("sync_state = ?", m.SyncState).
								Where("account = ?", m.Account).
								Where("sync_state = ?", 0).
								Update(); nil != err {
								panic(err)
							}
							go func(m *CreateAccCoinsModel) {
								defer wg.Done()
								err, coins := assetSyncTool.Sync(m.Account, m.Symbol, 1*time.Second)
								if nil == err {
									if nil != coins && 0 < len(coins) {
										m.Amount = coins[0].Amount.QuoRaw(1000000000000000000).Int64()
										m.AmountFloat = coins[0].Amount.ModRaw(1000000000000000000).Int64()
									}
									m.SyncState = 2
									if _, err = orm.NewQuery(db, m).
										Set("sync_state = ?", m.SyncState).
										Set("amount = ?", m.Amount).
										Set("amount_float = ?", m.AmountFloat).
										Where("account = ?", m.Account).
										Where("sync_state = ?", 1).
										Update(); nil != err {
										panic(err)
									}
								} else {
									_, _ = fmt.Fprintln(os.Stderr, err)
									if _, err = orm.NewQuery(db, m).
										Set("sync_state = ?", 0).
										Where("account = ?", m.Account).
										Where("sync_state = ?", 1).
										Update(); nil != err {
										panic(err)
									}
								}
							}(m)
						}

						wg.Wait()
						list = list[chunkSize:]
					}

					tick.Reset(duration)
				}
			}
		}()
	})
}
