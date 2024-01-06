package api_server

import (
	"fmt"
	logs "github.com/danbai225/go-logs"
	"github.com/gin-gonic/gin"
	"go-rustdesk-server/common"
	"go-rustdesk-server/web"
	"net/http"
)

func Start() {
	router := gin.Default()
	// 前端路由组
	frontendGroup := router.Group("/")
	{
		frontendGroup.GET("/", func(ctx *gin.Context) {
			file, _ := web.Dist.ReadFile("/dist/index.html")
			ctx.Data(200, "text/html; charset=utf-8", file)
		})
		_ = walkDir(web.Dist, "dist", frontendGroup, func(ctx *gin.Context) {
			ctx.FileFromFS("/dist"+ctx.Request.URL.Path, http.FS(web.Dist))
		})
	}
	//graphql
	router.GET("/graphql", playgroundHandler())
	router.POST("/query", graphqlHandler())
	// 启动Web服务器，并指定端口
	err := router.Run(fmt.Sprintf(":%d", common.Conf.WebPort))
	if err != nil {
		logs.Err(err)
	}
}
