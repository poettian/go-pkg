package main

import "github.com/gin-gonic/gin"

type user struct {
	Name string `json:"name" form:"name"`
	Age int `json:"age" form:"age"`
}

func main() {
	router := gin.Default()
	router.GET("/test", func(c *gin.Context) {
		var user user
		if err := c.ShouldBind(&user);err != nil {
			c.JSON(1, gin.H{
				"msg":"error param",
			})
		}
		c.JSON(0, gin.H{
			"name":user.Name,
			"age":user.Age,
		})
	})
	router.Run()
}