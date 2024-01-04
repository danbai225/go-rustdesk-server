package api_server

import (
	"fmt"
	logs "github.com/danbai225/go-logs"
	"github.com/gin-gonic/gin"
	"go-rustdesk-server/common"
	"go-rustdesk-server/data_server"
)

func Start() {
	router := gin.Default()
	apiGroup := router.Group("/api/v1")
	apiGroup.POST("/login", login)
	// 启动Web服务器，并指定端口
	err := router.Run(fmt.Sprintf(":%d", common.Conf.WebPort))
	if err != nil {
		logs.Err(err)
	}
}
func login(ctx *gin.Context) {
	db, err := data_server.GetDataSever()
	if err != nil {
		logs.Err(err)
		ctx.JSON(500, gin.H{"msg": "server error"})
		return
	}
	name, err := db.GetUserByName("admin")
	if err != nil {
		logs.Err(err)
		ctx.JSON(500, gin.H{"msg": "server error"})
		return
	}
	if name == nil {
		ctx.JSON(200, gin.H{"msg": "not found"})
		return
	}
	if name.Password != "123" {
		ctx.JSON(200, gin.H{"msg": "password error"})
		return
	}
}
