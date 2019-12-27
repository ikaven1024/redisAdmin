package api

import (
	"net/http"

	"github.com/ikaven1024/redisAdmin/redis_server"

	"github.com/gin-gonic/gin"
)

type Api struct {
	engine *gin.Engine
}

func InstallApi(engine *gin.Engine, serverManager *redis_server.Manager) http.Handler {
	engine.Use(errorHandle)

	installUserApi(engine)
	installRedisApi(engine, serverManager)
	installRedisServerApi(engine, serverManager)

	return engine
}

func (a *Api) Run(addr string) error {
	return a.engine.Run(addr)
}

type result struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Detail  string      `json:"detail,omitempty"`
}

func response(c *gin.Context, data interface{}) {
	c.JSON(200, &result{
		//Success: true,
		Data: data,
	})
}
