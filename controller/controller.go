package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lifegoeson/blockchain-explorer/defaultclient"
	"github.com/lifegoeson/blockchain-explorer/service"
	"strconv"
	"strings"
	"time"
)

//func GetCurBlockNum(c *gin.Context) {
//	curBlockNum := service.GetCurBlockNum(c.Query("channelGenesisHash"))
//	fmt.Println(curBlockNum)
//	chlByte,err := json.Marshal(curBlockNum)
//	if err != nil {
//		c.Data(500,"json",[]byte("获取当前区块号失败"))
//	}
//	c.Data(200,"json",chlByte)
//}

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

//func GetTxList(c *gin.Context){
//	channelGenesisHash := c.GetString("channelGenesisHash")
//	blockNum := c.GetInt64("blocknum")
//	txId := c.GetString("txId")
//	from := c.GetTime("from")
//	to := c.GetTime("to")
//	organizations := c.GetString("organizations")
//	txInfos,err := service.GetTxList(channelGenesisHash,blockNum,txId,from,to,organizations)
//	//txInfosByte,err := json.Marshal(txInfos)
//	if err != nil {
//		c.Data(500,"json",[]byte("获取交易信息失败"))
//	}
//	c.JSON(200,txInfos)
//}

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
	channelGenesisHash := c.Query("channelGenesisHash")
	peerInfos:= service.GetPeerData(channelGenesisHash)
	ledgerHeight := service.GetCurBlockNum(channelGenesisHash)
	//p := make(map[string]interface{})
	var ps []map[string]interface{}
	for _,peer := range peerInfos {
		m := make(map[string]interface{})
		m["ChannelName"] = peer.ChannelName
		m["Requests"] = peer.Requests
		m["ChannelGenesisHash"] = peer.ChannelGenesisHash
		m["ServerHostName"] = peer.ServerHostName
		m["MSPId"] = peer.MSPId
		m["PeerType"] = peer.PeerType
		m["ledgerHeight"] = ledgerHeight
		m["low"] = 0
		m["Unsigned"] = true
		ps = append(ps,m)
	}
	c.JSON(200,ps)
}

func GetCurBlockNum(c *gin.Context){
	channelGenesisHash := c.Query("channelGenesisHash")
	ledgerHeight := service.GetCurBlockNum(channelGenesisHash)
	c.JSON(200,ledgerHeight)
}

func GetBlockAndTxList(c *gin.Context){
	channelGenesisHash := c.PostForm("channelGenesisHash")
	from := c.PostForm("from")
	from_int,_ := strconv.ParseInt(from, 10, 64)
	//fmt.Println(from)
	startTime :=time.Unix(from_int/1000,0)
	fmt.Println(startTime)   //打印结果：2017-04-11 13:30:39
	to := c.PostForm("to")
	to_int,_ := strconv.ParseInt(to, 10, 64)
	fmt.Println(to)
	endTime := time.Unix(to_int/1000,0)
	fmt.Println(endTime)
	orgs := c.PostForm("orgs")
	orgsarray := strings.Split(orgs,",")
	os := ""
	for i, org := range orgsarray{
		if i == len(orgsarray) - 1 {
			os += "'" + org + "'"
		}else {
			os += "'" + org + "',"
		}

	}
	currentPage := c.PostForm("current")
	pageSize := c.PostForm("pageSize")
	currentPage2Int,_ := strconv.Atoi(currentPage)
	pageSize2Int,_ := strconv.Atoi(pageSize)
	offset := (currentPage2Int - 1)  * pageSize2Int
	blockAndTxList := service.GetBlockAndTxListByPage(channelGenesisHash,startTime,endTime,orgs,pageSize2Int,offset)
	c.JSON(200,blockAndTxList)

}

func GetTxListByPage(c *gin.Context){
	channelGenesisHash := c.PostForm("channelGenesisHash")
	blockNum := c.PostForm("blocknum")
	blockNum_int,_ := strconv.Atoi(blockNum)
	from := c.PostForm("from")
	from_int,_ := strconv.ParseInt(from, 10, 64)
	//fmt.Println(from)
	startTime :=time.Unix(from_int/1000,0)
	fmt.Println(startTime)   //打印结果：2017-04-11 13:30:39
	to := c.PostForm("to")
	to_int,_ := strconv.ParseInt(to, 10, 64)
	fmt.Println(to)
	endTime := time.Unix(to_int/1000,0)
	fmt.Println(endTime)
	orgs := c.PostForm("orgs")
	orgsarray := strings.Split(orgs,",")
	os := ""
	for i, org := range orgsarray{
		if i == len(orgsarray) - 1 {
			os += "'" + org + "'"
		}else {
			os += "'" + org + "',"
		}
	}
	currentPage := c.PostForm("current")
	pageSize := c.PostForm("pageSize")
	currentPage2Int,_ := strconv.Atoi(currentPage)
	pageSize2Int,_ := strconv.Atoi(pageSize)
	offset := (currentPage2Int - 1)  * pageSize2Int
	blockAndTxList, _ := service.GetTxListByPage(channelGenesisHash, int64(blockNum_int),"",startTime,endTime,orgs,pageSize2Int,offset)
	c.JSON(200,blockAndTxList)
}



