package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
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
	//chlInfos := service.GetCurBlockNum("2dfaf3fa74316ef1b0b476d5535de673ab2516cab93664237bdf3e441558cf6d")
	//chls := service.GetChannelsInfo("peer0.org1.example.ocm")
	//peerInfos, _ := service.GetChannelRefPeers(chls[0].ChannelName)
}


