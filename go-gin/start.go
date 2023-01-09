package go_gin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Run() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/info", func(c *gin.Context) {
			c.String(http.StatusOK, "No info")
		})
	}
	r.Run()
}
