package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/lifegoeson/blockchain-explorer/common"
	"github.com/lifegoeson/blockchain-explorer/defaultclient"
)

var db *sql.DB

func sqlOpen() {
	var err error
	db, err = sql.Open("postgres", "port=5432 user=hppoc password=password dbname=fabricexplorer sslmode=disable")
	common.CheckErr(err)
}
func main() {
	sqlOpen()
	defaultclient.GetInstance()
	//sdk := dc.DefaultFabSdk
	//syncBlocks(sdk,"mychannel","e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	initChannels()
	//time.Sleep(100000000)
}

