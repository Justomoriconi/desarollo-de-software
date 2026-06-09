package main

import (
	"backend/controllers"
	"backend/dao"
	"backend/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {

	dao.InitDB()

	r := gin.Default()

	r.GET("/eventos", controllers.GetEventos)
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET(
		"/perfil",
		middlewares.AuthMiddleware(),
		controllers.GetPerfil,
	)

	r.Run(":8080")
}
