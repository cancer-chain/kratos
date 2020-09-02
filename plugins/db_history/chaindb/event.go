package chaindb

import (
	"github.com/KuChainNetwork/kuchain/plugins/db_history/types"
	"github.com/go-pg/pg/v10"
	"github.com/tendermint/tendermint/libs/log"
)

type eventInDB struct {
	tableName struct{} `pg:"events,alias:events"` // default values are the same

	ID          int64  // both "Id" and "ID" are detected as primary key
	Type        string `json:"type"`
	BlockHeight int64  `pg:"default:0" json:"BlockHeight"`
	Attributes  map[string]string
}

func InsertEvent(db *pg.DB, logger log.Logger, evt *types.Event) error {

	logger.Info("InsertEvent", "insert Type:", evt.Type, "height:", evt.BlockHeight, "event", evt)

	if evt.Type == "create" { //asset
		EventCoinTypeAdd(db, logger, evt)
	} else if evt.Type == "issue" {
		EventCoinTypeModifySupply(db, logger, evt)
		EventAccCoinsAdd(db, logger, evt)
	} else if evt.Type == "burn" {
		EventCoinTypeModifySupply(db, logger, evt)
		EventAccCoinsReduce(db, logger, evt)
	} else if evt.Type == "transfer" {
		EventAccCoinsMove(db, logger, evt)
	} else if evt.Type == "delegate" { //staking
		EventDelegationAdd(db, logger, evt)
		EventDelegationChange(db, logger, evt)
	} else if evt.Type == "edit_validator" || evt.Type == "create_validator" || evt.Type == "end_validator" {
		EventValidatorAdd(db, logger, evt)
	} else if evt.Type == "unbond" {
		EventDelegationDel(db, logger, evt)
		EventDelegationChange(db, logger, evt)
		EventUnBondAdd(db, logger, evt)
	} else if evt.Type == "complete_redelegation" {
		EventAccCompleteReDelegateCoinsMove(db, logger, evt)
		EventDelegationChange(db, logger, evt)
	} else if evt.Type == "SendCoinsFromModuleToAccount" {
		EventAccCoinsMove(db, logger, evt)
	} else if evt.Type == "SendCoinsFromModuleToModule" {
		EventAccCoinsMove(db, logger, evt)
	} else if evt.Type == "DelegateCoinsFromAccountToModule" {
		EventAccCoinsMove(db, logger, evt)
	} else if evt.Type == "UndelegateCoinsFromModuleToAccount" {
		EventAccCoinsMove(db, logger, evt)
	} else if evt.Type == "ModuleMintCoins" {
		EventAccCoinsMintAdd(db, logger, evt)
	} else if evt.Type == "account.create" || evt.Type == "initModuleAccount" {
		EventAccountAdd(db, logger, evt)
	} else if evt.Type == "account.authupdate" {
		EventAccountUpdate(db, logger, evt)
	} else if evt.Type == "lock" {
		EventLockAccCoinsAdd(db, logger, evt)
	} else if evt.Type == "unlock" {
		EventUnLockAccCoinsAdd(db, logger, evt)
	} else if evt.Type == "payfee" {
		EventAccCoinsMove(db, logger, evt)
	} else if evt.Type == "slash" {
		EventValidatorSlash(db, logger, evt)
	} else if evt.Type == "unjail" {
		EventValidatorUnjail(db, logger, evt)
	} else if evt.Type == "proposer_reward" {
		EventProposerRewardAdd(db, logger, evt)
	} else if evt.Type == "complete_unbonding" {
		EventUnBondComplete(db, logger, evt)
	} else if evt.Type == "redelegate" {
		EventReDelegate(db, logger, evt)
	}

	return db.Insert(&eventInDB{
		Type:        evt.Type,
		BlockHeight: evt.BlockHeight,
		Attributes:  evt.Attributes,
	})
}
