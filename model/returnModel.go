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
