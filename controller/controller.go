package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lifegoeson/blockchain-explorer/defaultclient"
	"github.com/lifegoeson/blockchain-explorer/service"
	"strconv"
)

func GetCurBlockNum(c *gin.Context) {
	curBlockNum := service.GetCurBlockNum(c.Query("channelGenesisHash"))
	fmt.Println(curBlockNum)
	chlByte,err := json.Marshal(curBlockNum)
	if err != nil {
		c.Data(500,"json",[]byte("获取当前区块号失败"))
	}
	c.Data(200,"json",chlByte)
}

func GetTxCountByBlockNum(c *gin.Context){
	blocknum,err := strconv.Atoi(c.Query("blocknum"))
	channelGenesisBlock := c.Query("channelGenesisHash")
	txCountInfo := service.GetTxCountByBlockNum(channelGenesisBlock, int64(blocknum))
	//fmt.Println(curBlockNum)
	txCountInfoByte,err := json.Marshal(txCountInfo)
	if err != nil {
		c.Data(500,"json",[]byte("获取区块交易数量失败"))
	}
	c.Data(200,"json",txCountInfoByte)
}

func GetTransactionByID(c *gin.Context){
	txhash := c.Query("txhash")
	txInfo := service.GetTransactionByID(txhash)
	txInfoByte,err := json.Marshal(txInfo)
	if err != nil {
		c.Data(500,"json",[]byte("获取交易信息失败"))
	}
	c.Data(200,"json",txInfoByte)
}

func GetBlockActivityList(c *gin.Context){
	channelGenesisBlock := c.GetString("channelGenesisHash")
	blks := service.GetBlockActivityList(channelGenesisBlock)
	blksBytes,err := json.Marshal(blks)
	if err != nil {
		c.Data(500,"json",[]byte("获取最近区块信息失败"))
	}
	c.Data(200,"json",blksBytes)
}

func GetTxList(c *gin.Context){
	channelGenesisHash := c.GetString("channelGenesisHash")
	blockNum := c.GetInt64("blocknum")
	txId := c.GetString("txId")
	from := c.GetTime("from")
	to := c.GetTime("to")
	organizations := c.GetString("organizations")
	txInfos,err := service.GetTxList(channelGenesisHash,blockNum,txId,from,to,organizations)
	//txInfosByte,err := json.Marshal(txInfos)
	if err != nil {
		c.Data(500,"json",[]byte("获取交易信息失败"))
	}
	c.JSON(200,txInfos)
}

func GetBaseInfos(c * gin.Context){
	m := make(map[string]interface{})
	//获取通道信息
	chls := service.GetChannelsInfo(defaultclient.GetInstance().DefaultServerName)
	peerInfos:= service.GetPeerData(defaultclient.GetInstance().DefaultChannelGenHash)
	ccs := service.GetChaincodeCount(defaultclient.GetInstance().DefaultChannelGenHash)
	txCount := service.GetTxCount(defaultclient.GetInstance().DefaultChannelGenHash)
	blkActivity := service.GetBlockActivityList(defaultclient.GetInstance().DefaultChannelGenHash)
	blkCount := service.GetBlockCount(defaultclient.GetInstance().DefaultChannelGenHash)

	for _, chl := range chls{
		if chl.ChannelGenesisHash == defaultclient.GetInstance().DefaultChannelGenHash {
			m["defaultchannel"] = chl
		}
	}
	m["chls"] = chls
	m["peers"] = peerInfos
	m["ccs"] = ccs
	m["txCount"] = txCount
	m["blkActivity"] = blkActivity
	m["blkCount"] = blkCount
	c.JSON(200,m)
}

func GetPeerInfos(c *gin.Context){
	//channelGenesisHash := c.GetString("channelGenesisHash")
	peerInfos:= service.GetPeerData(defaultclient.GetInstance().DefaultChannelGenHash)
	c.JSON(200,peerInfos)
}



