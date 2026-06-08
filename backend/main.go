package main

import (
	"backend/dao"
	"github.com/gin-gonic/gin"
)

func main() {

	dao.InitDB()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API funcionando",
		})
	})
	r.Run(":8080")
}
