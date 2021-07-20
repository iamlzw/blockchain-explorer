package server

import (
	"github.com/gin-gonic/gin"
	"github.com/lifegoeson/blockchain-explorer/controller"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/block/count", controller.GetCurBlockNum)
	r.GET("/tx/count",controller.GetTxCountByBlockNum)
	r.GET("/tx/info",controller.GetTransactionByID)
	r.GET("/block/activity",controller.GetBlockActivityList)
	r.GET("/base/infos",controller.GetBaseInfos)
	return r
}