package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ikaven1024/redisAdmin/redis_server"
	"github.com/ikaven1024/redisAdmin/repository"
)

type redisServerApi struct {
	manager *redis_server.Manager
}

func installRedisServerApi(router *gin.Engine, manager *redis_server.Manager) {
	a := &redisServerApi{
		manager: manager,
	}

	router.GET("/api/redisServers", a.list)
	router.POST("/api/redisServers", a.add)
	router.PUT("/api/redisServers/:id", a.update)
	router.DELETE("/api/redisServers/:id", a.delete)
}

func (a redisServerApi) list(c *gin.Context) {
	data, err := a.manager.List()
	if err != nil {
		_ = c.Error(NewServerError("查询Redis服务错误", err.Error()))
		return
	}

	type vo struct {
		ID        uint                   `json:"id"`
		Name      string                 `json:"name"`
		Mode      redis_server.RedisMode `json:"mode"`
		Addresses []string               `json:"addresses"`
		Password  string                 `json:"password"`
	}

	res := make([]vo, 0, len(data))
	for _, d := range data {
		res = append(res, vo{
			ID:        d.ID,
			Name:      d.Name,
			Mode:      d.Mode,
			Addresses: d.Addresses.Data,
			Password:  d.Password,
		})
	}

	c.JSON(200, result{
		Data: res,
	})
}

func (a redisServerApi) add(c *gin.Context) {
	data, ok := a.getRedisServer(c)
	if !ok {
		return
	}

	if err := a.manager.Create(&data); err != nil {
		_ = c.Error(NewServerError("添加Redis服务错误", err.Error()))
		return
	}

	c.JSON(http.StatusOK, result{
		Message: "创建成功",
		Data:    data,
	})
}

func (a redisServerApi) update(c *gin.Context) {
	id, ok := a.getID(c)
	if !ok {
		return
	}

	data, ok := a.getRedisServer(c)
	if !ok {
		return
	}
	data.ID = uint(id)

	if err := a.manager.Update(&data); err != nil {
		_ = c.Error(NewServerError("更新Redis服务错误", err.Error()))
		return
	}

	c.JSON(http.StatusOK, result{
		Message: "创建成功",
		Data:    data,
	})
}

func (a redisServerApi) delete(c *gin.Context) {
	id, ok := a.getID(c)
	if !ok {
		return
	}

	if err := a.manager.Delete(id); err != nil {
		_ = c.Error(NewBadRequestError("删除服务失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, result{
		Message: "删除成功",
	})
}

func (a redisServerApi) getRedisServer(c *gin.Context) (redis_server.RedisServer, bool) {
	var body struct {
		ID        uint                   `json:"id"`
		Name      string                 `json:"name"`
		Mode      redis_server.RedisMode `json:"mode"`
		Addresses []string               `json:"addresses"`
		Password  string                 `json:"password"`
	}
	err := c.ShouldBind(&body)
	if err != nil {
		_ = c.Error(NewBadRequestError("参数错误", err.Error()))
		return redis_server.RedisServer{}, false
	}
	s := redis_server.RedisServer{
		Name:      body.Name,
		Mode:      body.Mode,
		Addresses: repository.NewStringSlice(body.Addresses),
		Password:  body.Password,
	}
	s.ID = body.ID
	return s, true
}

func (a redisServerApi) getID(c *gin.Context) (uint, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(NewBadRequestError("参数id不合法", err.Error()))
		return 0, false
	}
	return uint(id), true
}
