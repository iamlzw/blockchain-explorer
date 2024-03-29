package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lifegoeson/blockchain-explorer/defaultclient"
	"github.com/lifegoeson/blockchain-explorer/model"
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

func GetBaseChannelInfo(c *gin.Context){
	m := make(map[string]interface{})
	chls := service.GetChannelsInfo(defaultclient.GetInstance().DefaultServerName)
	for _, chl := range chls{
		if chl.ChannelGenesisHash == defaultclient.GetInstance().DefaultChannelGenHash {
			m["defaultchannel"] = chl
		}
	}
	m["chls"] = chls
	c.JSON(200,m)
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
	channelGenesisHash := c.Query("channelGenesisHash")
	m := make(map[string]interface{})
	peerInfos:= service.GetPeerData(channelGenesisHash)
	ccs := service.GetChaincodeCount(channelGenesisHash)
	txCount := service.GetTxCount(channelGenesisHash)
	blkActivity := service.GetBlockActivityList(channelGenesisHash)
	blkCount := service.GetBlockCount(channelGenesisHash)
	txgroup := service.GetTxByOrg(channelGenesisHash)
	m["peers"] = peerInfos
	m["ccs"] = ccs
	m["txCount"] = txCount
	m["blkActivity"] = blkActivity
	m["blkCount"] = blkCount
	jsonByte,_ := json.Marshal(txgroup)
	m["txgroup"] = jsonByte
	c.JSON(200,m)
}

func GetTxCountGroup(c * gin.Context){
	channelGenesisHash := c.Query("channelGenesisHash")
	txgroup := service.GetTxByOrg(channelGenesisHash)
	c.JSON(200,txgroup)
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

func GetChannelInfo(c *gin.Context){
	channelGenesisHash := c.Query("channelGenesisHash")
	chl := service.GetChannelConfig(channelGenesisHash)
	c.JSON(200,chl)
}

func GetChaincodes(c *gin.Context){
	channelGenesisHash := c.Query("channelGenesisHash")
	chl := service.GetChaincodes(channelGenesisHash)
	c.JSON(200,chl)
}

func GetTxCountByMonth(c *gin.Context){
	channelGenesisHash := c.Query("channelGenesisHash")
	month := c.Query("channelGenesisHash")
	month2Int,err := strconv.Atoi(month)
	if err != nil  {
		fmt.Println(err)
	}
	txData := service.GetTxByMonth(channelGenesisHash,month2Int)
	c.JSON(200,txData)
}

func GetTxOrBlockCountByTime(c *gin.Context){
	channelGenesisHash := c.Query("channelGenesisHash")
	queryType := c.Query("queryType")
	t := ""
	timeType := c.Query("timeType")
	fmt.Println(queryType,channelGenesisHash,timeType)
	var data []model.GetTxOrBlockByDateResultModel
	if queryType == "block" && timeType == "hour"{
		data = service.GetBlockByHour(channelGenesisHash,t)
	} else if queryType == "block" && timeType == "min" {
		data = service.GetBlockByMin(channelGenesisHash,t)
	} else if queryType == "tx" && timeType == "hour" {
		data = service.GetTxByHour(channelGenesisHash,t)
	} else {
		data = service.GetTxByMin(channelGenesisHash,t)
	}
	c.JSON(200,data)
}



