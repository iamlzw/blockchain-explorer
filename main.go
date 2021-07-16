package main

import (
	"bytes"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/lifegoeson/blockchain-explorer/common"
	"github.com/lifegoeson/blockchain-explorer/model"
	"github.com/spf13/viper"
	"io/ioutil"
	"time"
)

var db *sql.DB

func sqlOpen() {
	var err error
	db, err = sql.Open("postgres", "port=5432 user=hppoc password=password dbname=fabricexplorer sslmode=disable")
	common.CheckErr(err)
}
func main() {

	sqlOpen()

	//blknum,txcount := getTxCountByBlockNum("2dfaf3fa74316ef1b0b476d5535de673ab2516cab93664237bdf3e441558cf6d",972)
	//
	//fmt.Println(blknum,txcount)

	//getTransactionByIDResult := getTransactionByID("548b52e10ae039ad235fcdc900293b76fbeaa0ab62d0b5ed67f59fac60b8aae1")
	//
	//fmt.Println(getTransactionByIDResult)

	//startTime := time.Date(2021,time.July,9,16,44,10,0,time.UTC)
	//endTime := time.Date(2021,time.July,9,7,44,13,0,time.UTC)

	//getBlockActivityListResult := getBlockActivityList("2dfaf3fa74316ef1b0b476d5535de673ab2516cab93664237bdf3e441558cf6d")
	//blockAndtxs := getBlockAndTxList("2dfaf3fa74316ef1b0b476d5535de673ab2516cab93664237bdf3e441558cf6d",startTime,endTime,"")
	//chl := existChannel("mychannel1")
	//fmt.Println(chl)

	//block := model.Block{
	//	BlockNum: 100001,
	//	DataHash: "da6e26b6a5dc4b8511706e79ea3f3a0c3be1bf6dd9fec5d1cc840091988605",
	//	BlockHash: "da6e26b6a5dc4b8511706e79ea3f3a0c3be1bf6dd9fec5d1cc840091988605",
	//	PreHash: "da6e26b6a5dc4b8511706e79ea3f3a0c3be1bf6dd9fec5d1cc840091988605",
	//	TxCount: 100,
	//	CreateAt: time.Now(),
	//	PrevBlockHash: "da6e26b6a5dc4b8511706e79ea3f3a0c3be1bf6dd9fec5d1cc840091988605",
	//	ChannelGenesisHash: "2dfaf3fa74316ef1b0b476d5535de673ab2516cab93664237bdf3e441558cf6d",
	//}
	//fmt.Println(saveBlock(block))

	sdk := initSDK()
	//queryChaincodeInfo(sdk)
	//queryGenesisBlock(sdk)
	syncBlocks(sdk,"mychannel","e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	initChannels(sdk)
	//certByte2Pem()
	//getTxCountByBlockNum("",4)
	//discoveryTest("mychannel")
	//discoveryRaw()

	//tesetViper()
	//getBlockAndTxList("111",time.Date(2021,time.June,23,10,0,0,0,time.UTC),time.Date(2021,time.June,23,10,0,0,0,time.UTC), "")
}

func testSaveTransaction(){
	tx := model.Transaction{
		BlockId: 2587,
		TxHash: "c8d4b8f6c733b094af2290ff10173cd008b955f843069f9d1094a62e34873e23",
		CreateAt: time.Now(),
		ChaincodeName: "testchaincode",
		Status: 200,
		CreatorMspId: "Org1MSP",
	}
	saveTransaction(&tx)
}

func tesetViper(){
	blkFile, _ := ioutil.ReadFile("blockfiles/mychannel_0.json")
	v := viper.New()
	v.SetConfigType("json")
	v.MergeConfig(bytes.NewBuffer(blkFile))
	_ = v.MergeInConfig()
}
