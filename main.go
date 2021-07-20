package main

import (
	_ "database/sql"
	_ "github.com/lib/pq"
	_ "github.com/lifegoeson/blockchain-explorer/common"
	"github.com/lifegoeson/blockchain-explorer/defaultclient"
	"github.com/lifegoeson/blockchain-explorer/server"
	"github.com/lifegoeson/blockchain-explorer/service"
	"log"
)

func main() {
	service.SqlOpen()
	defaultclient.GetInstance()
	initChannels()
	r := server.InitRouter()
	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}

