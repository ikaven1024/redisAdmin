package api

import (
	"github.com/gin-gonic/gin"
	"log"
)

func RunApi() {
	gin.ForceConsoleColor()

	router := gin.Default()
	installRedisApi(router)

	log.Println(router.Run(":80"))
}

type apiResult struct {
	Data    interface{} `json:"data,emit"`
	Success bool        `json:"success"`
	Message string      `json:"message,emit"`
}

func response(c *gin.Context, data interface{}) {
	c.JSON(200, &apiResult{
		Success: true,
		Data:    data,
	})
}

func responseError(c *gin.Context, err error) {
	_ = c.Error(err)
	c.JSON(200, &apiResult{
		Success: false,
		Message: err.Error(),
	})
}
