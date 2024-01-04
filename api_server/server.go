package api_server

import (
	"embed"
	"fmt"
	logs "github.com/danbai225/go-logs"
	"github.com/gin-gonic/gin"
	"go-rustdesk-server/common"
	"go-rustdesk-server/data_server"
	"go-rustdesk-server/web"
	"net/http"
	"path/filepath"
	"strings"
)

func walkDir(fs embed.FS, dir string, group *gin.RouterGroup, handlerFunc gin.HandlerFunc) error {
	entries, err := fs.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			// 递归遍历子目录
			err = walkDir(fs, filepath.Join(dir, entry.Name()), group, handlerFunc)
			if err != nil {
				return err
			}
		} else {
			// 处理文件
			path := filepath.Join(strings.ReplaceAll(dir, "dist", ""), entry.Name())
			group.GET(path, handlerFunc)
		}
	}
	return nil
}
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
	// 后端路由组
	apiGroup := router.Group("/api/v1")
	{
		apiGroup.GET("/login", login)
	}
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
