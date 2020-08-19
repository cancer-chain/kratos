package chaindb

import (
	"encoding/json"
	"fmt"
	"github.com/KuChainNetwork/kuchain/plugins/types"
	"github.com/go-pg/pg/v10"
	"github.com/tendermint/tendermint/libs/log"
	"sync/atomic"
)

type BlockInfo struct {
	tableName struct{} `pg:"BlockInfo,alias:BlockInfo"` // default values are the same

	ID int64 // both "Id" and "ID" are detected as primary key

	BlockHash                    string `pg:"unique:as" json:"block_hash"`
	BlockProposalValidator       string `json:"proposal_validator"`
	BlockProposalTenderValidator string `json:"proposal_tender_validator"`

	BlockIdPartsHeaderTotal       string `json:"block_id_partsheader_total"`
	BlockIdPartsHeaderHash        string `json:"block_id_partsheader_hash"`
	BlockHeaderVersionBlock       string `json:"block_header_version_block"`
	BlockHeaderVersionApp         string `json:"block_header_version_app"`
	BlockHeaderChainId            string `json:"block_header_chainid"`
	BlockHeaderHeight             string `json:"block_header_height"`
	BlockHeaderTime               string `json:"block_header_time"`
	BlockHeaderLastBlockIdHash    string `json:"block_header_lastblockid_hash"`
	BlockHeaderLastCommitHash     string `json:"block_header_lastcommithash"`
	BlockHeaderDataHash           string `json:"block_header_datahash"`
	BlockHeaderNextValidatorsHash string `json:"block_header_nextvalidatorshash"`
	BlockHeaderConsensusHash      string `json:"block_header_consensushash"`
	BlockHeaderAppHash            string `json:"block_header_apphash"`
	BlockHeaderLastResultsHash    string `json:"block_header_lastresultshash"`
	BlockHeaderEvidenceHash       string `json:"block_header_evidencehash"`
	BlockHeaderProposerAddress    string `json:"block_header_proposeraddress"`
	BlockHeaderProposer           string `json:"block_header_proposer"`
	BlockDataHash                 string `json:"block_evidence_hash"`
	BlockHeaderValidators         string `json:"block_header_validators"`
	BlockLastCommitVotes          string `json:"block_lastcommit_votes"`
	BlockLastCommitRound          string `json:"block_lastcommit_round"`
	BlockLastCommitInfo           string `json:"block_lastcommit_info"`
	Time                          string `json:"time"`
}

type blockInDB struct {
	tableName struct{} `pg:"block,alias:block"` // default values are the same

	ID int64 // both "Id" and "ID" are detected as primary key

	types.ReqBeginBlock
}

func newBlockInDB(tb types.ReqBeginBlock) *blockInDB {
	return &blockInDB{
		ReqBeginBlock: tb,
	}
}

func InsertBlockInfo(db *pg.DB, logger log.Logger, bk *blockInDB) error {

	err := db.Insert(bk)
	if err != nil {
		panic(err)
	}

	msg := BlockInfo{
		BlockHash:               Hash2Hex(bk.ReqBeginBlock.RequestBeginBlock.Hash()),
		BlockIdPartsHeaderTotal: "",
		BlockIdPartsHeaderHash:  Hash2Hex(bk.ReqBeginBlock.RequestBeginBlock.Header.Hash()),

		BlockHeaderVersionApp:         fmt.Sprintf("%d", bk.ReqBeginBlock.RequestBeginBlock.Header.Version.App),
		BlockHeaderVersionBlock:       fmt.Sprintf("%d", bk.ReqBeginBlock.RequestBeginBlock.Header.Version.Block),
		BlockHeaderChainId:            bk.ReqBeginBlock.RequestBeginBlock.Header.ChainID,
		BlockHeaderHeight:             fmt.Sprintf("%d", bk.ReqBeginBlock.RequestBeginBlock.Header.Height),
		BlockHeaderTime:               TimeFormat(bk.ReqBeginBlock.RequestBeginBlock.Header.Time),
		BlockHeaderLastBlockIdHash:    Hash2Hex(bk.ReqBeginBlock.RequestBeginBlock.Header.LastBlockID.Hash),
		BlockHeaderLastCommitHash:     Hash2Hex(bk.ReqBeginBlock.RequestBeginBlock.Header.LastCommitHash),
		BlockHeaderDataHash:           Hash2Hex(bk.ReqBeginBlock.RequestBeginBlock.Header.DataHash),
		BlockHeaderNextValidatorsHash: Hash2Hex(bk.ReqBeginBlock.RequestBeginBlock.Header.NextValidatorsHash),
		BlockHeaderConsensusHash:      Hash2Hex(bk.ReqBeginBlock.RequestBeginBlock.Header.ConsensusHash),
		BlockHeaderAppHash:            Hash2Hex(bk.ReqBeginBlock.RequestBeginBlock.Header.AppHash),
		BlockHeaderLastResultsHash:    Hash2Hex(bk.ReqBeginBlock.RequestBeginBlock.Header.LastResultsHash),
		BlockHeaderEvidenceHash:       Hash2Hex(bk.ReqBeginBlock.RequestBeginBlock.Header.EvidenceHash),
		BlockHeaderProposerAddress:    Hash2Hex(bk.ReqBeginBlock.RequestBeginBlock.Header.ProposerAddress),
		BlockHeaderProposer:           Hash2Hex(bk.ReqBeginBlock.RequestBeginBlock.Header.ProposerAddress),
		BlockDataHash:                 Hash2Hex(bk.ReqBeginBlock.RequestBeginBlock.Header.DataHash),
		BlockLastCommitRound:          fmt.Sprintf("%d", bk.ReqBeginBlock.RequestBeginBlock.LastCommit.Round),
		BlockHeaderValidators:         Hash2Hex(bk.ReqBeginBlock.RequestBeginBlock.Header.ValidatorsHash),
		Time:                          TimeFormat(bk.Time),
	}

	n := len(bk.RequestBeginBlock.LastCommit.Signatures)
	for i := 0; i < n; i++ {
		vote := bk.RequestBeginBlock.LastCommit.GetVote(i)
		bz, _ := json.Marshal(vote)
		msg.BlockLastCommitVotes = msg.BlockLastCommitVotes + string(bz) + ","
	}

	bz, _ := json.Marshal(bk.RequestBeginBlock.LastCommit)
	msg.BlockLastCommitInfo = string(bz)

	bz, _ = json.Marshal(bk.ReqBeginBlock.RequestBeginBlock.Header.Version)

	//get proposal
	var tMap map[string]interface{}
	json.Unmarshal([]byte(bk.ValidatorInfo), &tMap)

	operatorAccount, ok := tMap["operator_account"]
	if ok {
		msg.BlockProposalValidator = operatorAccount.(string)
	}

	consensusPubkey, ok := tMap["consensus_pubkey"]
	if ok {
		msg.BlockProposalTenderValidator = consensusPubkey.(string)
	}

	tx, _ := db.Begin()

	{ //blockinfo
		err = db.Insert(&msg)
		if err != nil {
			EventErr(db, logger, NewErrMsg(err))
		}
		logger.Debug("InsertBlockInfo", "blockinfo", msg)
	}
	{ //tx msg
		for _, txInBk := range bk.Tx {
			err, TxUid := InsertTxm(db, logger, txInDB{ReqTx: txInBk})
			if err != nil {
				EventErr(db, logger, NewErrMsg(err))
			}

			for _, m := range txInBk.Msgs {
				iMsg := buildTxMsg(logger, m, txInDB{ReqTx: txInBk}, TxUid, "")
				err := db.Insert(&iMsg)
				if err != nil {
					EventErr(db, logger, NewErrMsg(err))
				}
			}
		}
	}
	{ //events
		Events := makeEvent(bk.Events, logger)
		for _, evt := range Events {
			err = InsertEvent(db, logger, &evt)
			if err != nil {
				EventErr(db, logger, NewErrMsg(err))
			}
		}

		for _, txInBk := range bk.Tx {
			if txInBk.RawLog.Code == 0 {
				for i := 0; i < len(bk.TxEvents); i++ {
					txEvents := makeEvent(bk.TxEvents[i], logger)
					for _, evt := range txEvents {
						err = InsertEvent(db, logger, &evt)
						if err != nil {
							EventErr(db, logger, NewErrMsg(err))
						}
					}
				}
			}
		}

		feeEvents := makeEvent(bk.FeeEvents, logger)
		for _, evt := range feeEvents {
			err = InsertEvent(db, logger, &evt)
			if err != nil {
				EventErr(db, logger, NewErrMsg(err))
			}
		}
	}
	{ //SyncStat
		var stat *SyncState
		if stat, err = UpdateChainSyncStat(db, logger, bk.ReqBeginBlock.RequestBeginBlock.Header.Height); err != nil {
			logger.Error("UpdateChainSyncStat error", "err", err)
		}
		atomic.StoreInt64(&SyncBlockHeight, stat.BlockNum)
	}

	tx.Commit()

	return err
}
