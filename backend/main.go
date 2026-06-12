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

	r.Use(middlewares.CORSMiddleware())

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	r.GET("/eventos", controllers.GetEventos)
	r.GET("/eventos/:id", controllers.GetEventoByID)
	r.GET("/eventos/:id/tipos-entrada", controllers.GetTiposEntradaByEvento)

	r.GET("/perfil", middlewares.AuthMiddleware(), controllers.GetPerfil)

	r.POST("/tickets", middlewares.AuthMiddleware(), controllers.ComprarTicket)
	r.GET("/tickets", middlewares.AuthMiddleware(), controllers.GetMisTickets)
	r.PUT("/tickets/:id/cancelar", middlewares.AuthMiddleware(), controllers.CancelarTicket)
	r.PUT("/tickets/:id/transferir", middlewares.AuthMiddleware(), controllers.TransferirTicket)

	r.Run(":8080")
}
