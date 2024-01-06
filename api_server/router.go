package api_server

import (
	"context"
	"embed"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"go-rustdesk-server/api_server/graph"
	"go-rustdesk-server/data_server"
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
func authErr() graphql.ResponseHandler {
	return func(ctx context.Context) *graphql.Response {
		return graphql.ErrorResponse(ctx, "auth error")
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL playground", "/query")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
func graphqlHandler() gin.HandlerFunc {
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	h.Use(handler.OperationFunc(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		operationContext := graphql.GetOperationContext(ctx)
		switch operationContext.OperationName {
		case "login":
			fallthrough
		case "IntrospectionQuery":
			return next(ctx)
		}
		token := operationContext.Headers.Get("Token")
		db, err := data_server.GetDataSever()
		if err != nil {
			return authErr()
		}
		user, err := db.CheckToken(token)
		if err != nil {
			return authErr()
		}
		ctx = context.WithValue(ctx, "user", user)
		return next(ctx)
	}))
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
