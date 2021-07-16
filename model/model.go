package model

import "time"

type Block struct {
	Id           int64     `json:"id"`
	BlockNum     int64     `json:"blocknum"`
	DataHash     string    `json:"datahash"`
	PreHash      string    `json:"prehash"`
	TxCount      int64     `json:"txcount"`
	CreateAt     time.Time `json:"createat"`
	PrevBlockHash string    `json:"prev_blockhash"`
	BlockHash string `json:"blockhash"`
	ChannelGenesisHash string `json:"channel_genesis_hash"`
}

type Chaincode struct {
	Id int64 	`json:"id"`
	Name string 	`json:"name"`
	Version string `json:"version"`
	Path string 	`json:"path"`
	ChannelGenesisHash string `json:"channel_genesis_hash"`
	TxCount int64 `json:"txcount"`
	CreateAt     time.Time `json:"createat"`
}

type Channel struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	Blocks int64 `json:"blocks"`
	Trans int64 `json:"trans"`
	CreateAt     time.Time `json:"createat"`
	ChannelGenesisHash string `json:"channel_genesis_hash"`
	ChannelHash string `json:"channel_hash"`
	ChannelConfig []byte `json:"channel_config"`
	ChannelBlock []byte `json:"channel_block"`
	ChannelTx []byte `json:"channel_tx"`
	ChannelVersion interface{} `json:"channel_version"`
}

type Orderer struct {
	Id int64 `json:"id"`
	Requests string `json:"requests"`
	ServerHostName string `json:"server_hostname"`
	CreateAt     time.Time `json:"createat"`
}

type Peer struct {
	Id int64 `json:"id"`
	Org int64 `json:"org"`
	ChannelGenesisHash string `json:"channel_genesis_hash"`
	MspId string `json:"mspid"`
	Requests string `json:"requests"`
	Events string `json:"events"`
	ServerHostName string `json:"server_hostname"`
	CreateAt     time.Time `json:"createat"`
	PeerType string `json:"peer_type"`
}

type PeerRefChaincode struct {
	Id int64 `json:"id"`
	PeerId string `json:"peerid"`
	ChaincodeId string `json:"chaincodeid"`
	CCVersion string 	`json:"cc_version"`
	ChannelId string `json:"channelid"`
	CreateAt     time.Time `json:"createat"`
}

type PeerRefChannel struct {
	Id int64 `json:"id"`
	PeerId string `json:"peerid"`
	ChannelId string `json:"channelid"`
	PeerType string `json:"peer_type"`
	CreateAt     time.Time `json:"createat"`
}

type Transaction struct {
	Id int64 `json:"id"`
	BlockId int64 `json:"blockid"`
	TxHash string `json:"txhash"`
	CreateAt     time.Time `json:"createat"`
	ChaincodeName string `json:"chaincodename"`
	Status int32 `json:"status"`
	CreatorMspId string `json:"creator_msp_id"`
	EndorserMspId string `json:"endorser_msp_id"`
	ChaincodeId string `json:"chaincode_id"`
	Type string `json:"type"`
	ReadSet string `json:"read_set"`
	WriteSet string `json:"write_set"`
	ChannelGenesisHash string `json:"channel_genesis_hash"`
	ValidationCode string `json:"validation_code"`
	EnvelopeSignature string `json:"envelope_signature"`
	PayloadExtension string `json:"payload_extension"`
	CreatorIdBytes string `json:"creator_id_bytes"`
	CreatorNonce string `json:"creator_nonce"`
	ChaincodeProposalInput string `json:"chaincode_proposal_input"`
	TxResponse string `json:"tx_response"`
	PayloadProposalHash string `json:"payload_proposal_hash"`
	EndorserIdBytes string `json:"endorser_id_bytes"`
	EndorserSignature string `json:"endorser_signature"`
}

type WriteLock struct {
	WriteBlock int64 `json:"write_block"`
}
