package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/lifegoeson/blockchain-explorer/common"
	"time"
)

var db *sql.DB

func sqlOpen() {
	var err error
	db, err = sql.Open("postgres", "port=5432 user=postgres password=123456 dbname=fabricexplorer sslmode=disable")
	common.CheckErr(err)
}
func main() {

	sqlOpen()

	//getTxCountByBlockNum("",4)

	getBlockAndTxList("111",time.Date(2021,time.June,23,10,0,0,0,time.UTC),time.Date(2021,time.June,23,10,0,0,0,time.UTC), "")
}
