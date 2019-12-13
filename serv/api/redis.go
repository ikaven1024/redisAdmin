package api

import (
	"github.com/gin-gonic/gin"
	"github/ikaven/redisGoAdmin/redis"
)

func installRedisApi(router *gin.Engine) {
	a := &redisApi{
		cli: redis.NewRedisCli("192.168.11.128:6379"),
	}

	router.POST("/api/redis/keys/nodes", a.getKeyTreeNodes)
	router.GET("/api/redis/keys/summary", a.getKeySummary)
	router.GET("/api/redis/keys/value", a.getValueOfKey)
}

type redisApi struct {
	cli *redis.Client
}

// listKeySubMenus return sub nodes for key tree.
func (r redisApi) getKeyTreeNodes(c *gin.Context) {

	var body struct {
		Prefix string `json:"prefix"`
	}

	if err := c.ShouldBind(&body); err != nil {
		responseError(c, err)
	}

	children, err := r.cli.KeyMenus(body.Prefix)
	if err != nil {
		responseError(c, err)
		return
	}
	response(c, gin.H{
		"prefix":   body.Prefix,
		"children": children,
	})
}

// getKeySummary return ttl and type of key
func (r redisApi) getKeySummary(c *gin.Context) {
	var body struct {
		Key string `json:"key"`
	}
	if err := c.ShouldBind(&body); err != nil {
		responseError(c, err)
		return
	}

	ttl, err := r.cli.TTL(body.Key).Result()
	if err != nil {
		responseError(c, err)
		return
	}

	typ, err := r.cli.Type(body.Key).Result()
	if err != nil {
		responseError(c, err)
		return
	}

	response(c, gin.H{
		"key":  body.Key,
		"ttl":  ttl.String(),
		"type": typ,
	})
}

// getValueOfKey return value of key.
// support kinds of key, type is set in `Type`
func (r redisApi) getValueOfKey(c *gin.Context) {

	var body struct {
		Type     string `json:"type"`
		Key      string `json:"key"`
		Match    string `json:"match"`
		PageSize int64  `json:"pageSize"`
		PageNo   int64  `json:"pageNo"`
		Cursor   uint64 `json:"cursor"`
	}
	if err := c.ShouldBind(&body); err != nil {
		responseError(c, err)
		return
	}

	values, err := r.cli.GetValue(body.Key, redis.GetValueOpts{
		Type:     body.Type,
		Match:    body.Match,
		PageNo:   body.PageNo,
		PageSize: body.PageSize,
		Cursor:   body.Cursor,
	})
	if err != nil {
		responseError(c, err)
		return
	}
	response(c, gin.H{
		"key":    body.Key,
		"values": values,
	})
}
