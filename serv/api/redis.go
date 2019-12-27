package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ikaven1024/redisAdmin/redis_client"
	"github.com/ikaven1024/redisAdmin/redis_server"
	"strconv"
)

const ContextKeyServer = "server"

func installRedisApi(router *gin.Engine, manager *redis_server.Manager) {
	a := &redisApi{
		serverManager: manager,
	}

	router.GET("/api/redis/dbCount", a.getDbCount)
	router.GET("/api/redis/keys/treeNodes", a.getKeyTreeNodes)
	router.GET("/api/redis/keys/summary", a.getKeySummary)
	router.GET("/api/redis/keys/value", a.getValueOfKey)
}

type redisApi struct {
	serverManager *redis_server.Manager
}

func (r redisApi) getDbCount(c *gin.Context) {
	cli, ok := r.getRedisClient(c)
	if !ok {
		return
	}

	data, err := cli.ConfigGet("databases").Result()
	if err != nil {
		_ = c.Error(NewServerError("查询DBSize失败", err.Error()))
		return
	}
	numStr := data[1].(string)
	num, _ := strconv.Atoi(numStr)
	response(c, gin.H{
		"number": num,
	})
}

// listKeySubMenus return sub nodes for key tree.
func (r redisApi) getKeyTreeNodes(c *gin.Context) {
	cli, ok := r.getRedisClient(c)
	if !ok {
		return
	}

	prefix := c.Query("prefix")

	children, err := cli.KeyMenus(prefix)
	if err != nil {
		_ = c.Error(NewServerError("加载key失败", err.Error()))
		return
	}
	response(c, gin.H{
		"prefix":   prefix,
		"children": children,
	})
}

// getKeySummary return ttl and type of key
func (r redisApi) getKeySummary(c *gin.Context) {
	cli, ok := r.getRedisClient(c)
	if !ok {
		return
	}

	key := c.Query("key")
	if len(key) == 0 {
		_ = c.Error(NewServerError("参数错误", "丢失查询参数"+key))
	}

	ttl, err := cli.TTL(key).Result()
	if err != nil {
		_ = c.Error(NewServerError("查询key ttl失败", err.Error()))
		return
	}

	typ, err := cli.Type(key).Result()
	if err != nil {
		_ = c.Error(NewServerError("查询key类型失败", err.Error()))
		return
	}

	response(c, gin.H{
		"key":  key,
		"ttl":  ttl.String(),
		"type": typ,
	})
}

// getValueOfKey return value of key.
// support kinds of key, type is set in `Type`
func (r redisApi) getValueOfKey(c *gin.Context) {
	cli, ok := r.getRedisClient(c)
	if !ok {
		return
	}

	var param struct {
		Type     string `form:"type"`
		Key      string `form:"key"`
		Match    string `form:"match"`
		PageSize int64  `form:"pageSize"`
		PageNo   int64  `form:"pageNo"`
		Cursor   uint64 `form:"cursor"`
	}
	if err := c.ShouldBindQuery(&param); err != nil {
		_ = c.Error(NewBadRequestError("参数错误", err.Error()))
		return
	}

	result, err := cli.GetValue(param.Key, redis_client.GetValueOpts{
		Type:     param.Type,
		Match:    param.Match,
		PageNo:   param.PageNo,
		PageSize: param.PageSize,
		Cursor:   param.Cursor,
	})
	if err != nil {
		_ = c.Error(NewServerError("查询value失败", err.Error()))
		return
	}
	response(c, result)
}

func (r redisApi) getRedisClient(c *gin.Context) (*redis_client.Client, bool) {
	var param struct {
		ServerID uint `form:"serverId"`
		DB       int  `form:"db"`
	}

	if err := c.ShouldBindQuery(&param); err != nil {
		_ = c.Error(NewBadRequestError("参数错误", err.Error()))
		return nil, false
	}

	serv, err := r.serverManager.Get(param.ServerID)
	if err != nil {
		_ = c.Error(NewServerError("查询服务信息失败", err.Error()))
		return nil, false
	}

	var cli *redis_client.Client
	if serv.IsCluster() {
		cli = redis_client.NewRedisClusterCli(serv.Addresses.Data, serv.Password)
	} else {
		cli = redis_client.NewRedisCli(serv.Addresses.Data[0], serv.Password, param.DB)
	}
	return cli, true
}
