package api

import "github.com/gin-gonic/gin"

type userApi struct {
}

func installUserApi(router *gin.Engine) {
	a := &userApi{}

	router.POST("/api/user/login", a.login)
}

func (u userApi) login(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 20000,
		"data": map[string]interface{}{"token": "admin-token"},
	})
}
