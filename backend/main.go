package main

import (
	"backend/controllers"
	"backend/dao"
	"github.com/gin-gonic/gin"
)

func main() {

	dao.InitDB()

	r := gin.Default()

	r.GET("/eventos", controllers.GetEventos)

	r.Run(":8080")
}
