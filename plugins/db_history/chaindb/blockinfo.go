package chaindb

import (
	"encoding/json"
	"fmt"
	"github.com/KuChainNetwork/kuchain/plugins/types"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/tendermint/tendermint/libs/log"
)

type BlockInfo struct {
	tableName struct{} `pg:"BlockInfo,alias:BlockInfo"` // default values are the same

	ID int64 // both "Id" and "ID" are detected as primary key

	BlockHash                    string `pg:"unique:as",json:"block_hash"`
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
	BlockByzantineValidators      string `json:"block_byzantinevalidators"`
	BlockLastCommitVotes          string `json:"block_lastcommit_votes"`
	BlockLastCommitRound          string `json:"block_lastcommit_round"`
	BlockLastCommitInfo           string `json:"block_lastcommit_info"`
	Time                          string `json:"time"`
}

type blockInDB struct {
	tableName struct{} `pg:"block,alias:block"` // default values are the same

	ID int64 // both "Id" and "ID" are detected as primary key

	types.ReqBlock
}

func newBlockInDB(tb types.ReqBlock) *blockInDB {
	return &blockInDB{
		ReqBlock: tb,
	}
}

func InsertBlockInfo(db *pg.DB, logger log.Logger, bk *blockInDB) error {
	err := db.Insert(bk)
	if err != nil {
		panic(err)
	}

	msg := BlockInfo{
		BlockHash:                     Hash2Hex(bk.Hash),
		BlockIdPartsHeaderTotal:       "",
		BlockIdPartsHeaderHash:        "",
		BlockHeaderVersionBlock:       bk.Header.Version.String(),
		BlockHeaderVersionApp:         fmt.Sprintf("%d", bk.Header.Version.App),
		BlockHeaderChainId:            bk.Header.ChainID,
		BlockHeaderHeight:             fmt.Sprintf("%d", bk.Header.Height),
		BlockHeaderTime:               TimeFormat(bk.Header.Time),
		BlockHeaderLastBlockIdHash:    Hash2Hex(bk.Header.LastBlockId.Hash),
		BlockHeaderLastCommitHash:     Hash2Hex(bk.Header.LastCommitHash),
		BlockHeaderDataHash:           Hash2Hex(bk.Header.DataHash),
		BlockHeaderNextValidatorsHash: Hash2Hex(bk.Header.NextValidatorsHash),
		BlockHeaderConsensusHash:      Hash2Hex(bk.Header.ConsensusHash),
		BlockHeaderAppHash:            Hash2Hex(bk.Header.AppHash),
		BlockHeaderLastResultsHash:    Hash2Hex(bk.Header.LastResultsHash),
		BlockHeaderEvidenceHash:       Hash2Hex(bk.Header.EvidenceHash),
		BlockHeaderProposerAddress:    Hash2Hex(bk.Header.ProposerAddress),
		BlockHeaderProposer:           Hash2Hex(bk.Header.ProposerAddress),
		BlockDataHash:                 "",
		BlockLastCommitRound:          fmt.Sprintf("%d", bk.LastCommitInfo.Round),
		Time:                          TimeFormat(bk.Time),
	}

	bz, _ := json.Marshal(bk.ByzantineValidators)
	msg.BlockByzantineValidators = string(bz)

	bz, _ = json.Marshal(bk.LastCommitInfo.Votes)
	msg.BlockLastCommitVotes = string(bz)

	bz, _ = json.Marshal(bk.LastCommitInfo)
	msg.BlockLastCommitInfo = string(bz)

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

	err = orm.Insert(db, &msg)
	if err != nil {
		EventErr(db, logger, NewErrMsg(err))
	}
	return err
}
