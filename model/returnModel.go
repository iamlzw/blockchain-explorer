package model

import "time"

type GetTxCountByBlockNumResultModel struct {
	BlockNum int64
	TxCount int64
}

type GetTransactionByIDResultModel struct {
	Tx Transaction
	ChannelName string
}

type GetBlockActivityListResultModel struct {
	BlockNum int64
	TxCount int64
	DataHash string
	BlockHash string
	PreHash string
	CreateAt time.Time
	TxHash string
	ChannelName string
}

type GetTxListResultModel struct {
	CreatorMspId  string
	TxHash        string
	Type        string
	ChaincodeName string
	CreateAt      time.Time
	ChannelName   string
}

type GetBlockAndTxListResultModel struct {
	ChannelName string
	BlockNum int64
	TxCount      string
	DataHash        string
	BlockHash string
	PreHash      string
	CreateAt   time.Time
	TxHash string
}

type GetChainsInfoResultModel struct {
	Id int64
	ChannelName string
	Blocks int64
	ChannelGenesisHash string
	Trans int64
	CreateAt time.Time
	ChannelHash string
}

type GetChannelRefPeersResultModel struct{
	Id int64
	CreateAt time.Time
	PeerId string
	ChannelId string
	PeerType string
}

type GetPeerDataResultModel struct {
	ChannelName string
	Requests string
	ChannelGenesisHash string
	ServerHostName string
	MSPId string
	PeerType string
}

type GetTxPerChaincodeGenerateResultModel struct {
	ChaincodeName string
	ChannelName string
	Version string
	ChannelGenesisHash string
	Path string
	TxCount int
}

type GetOrgDataResultModel struct {

}

type GetTxCountByOrgResultModel struct {
	Name string
	Value int
}

type GetTxOrBlockByDateResultModel struct {
	Time string
	TxCount int
}
