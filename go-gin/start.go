package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	v1 := r.Group("/v1/user")
	{
		v1.GET("/info", func(c *gin.Context) {
			c.String(http.StatusOK, "No info")
		})
		v1.POST("/create", func(c *gin.Context) {
			err := c.Request.ParseForm()
			if err != nil {

			}
		})
	}
	r.Run()
}
