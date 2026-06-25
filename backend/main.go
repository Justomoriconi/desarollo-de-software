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

	// Rutas públicas
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/eventos", controllers.GetEventos)
	r.GET("/eventos/:id", controllers.GetEventoByID)
	r.GET("/eventos/:id/tipos-entrada", controllers.GetTiposEntradaByEvento)

	// Rutas cliente
	r.GET("/perfil", middlewares.AuthMiddleware(), controllers.GetPerfil)
	r.POST("/tickets", middlewares.AuthMiddleware(), controllers.ComprarTicket)
	r.GET("/tickets", middlewares.AuthMiddleware(), controllers.GetMisTickets)
	r.PUT("/tickets/:id/cancelar", middlewares.AuthMiddleware(), controllers.CancelarTicket)
	r.PUT("/tickets/:id/transferir", middlewares.AuthMiddleware(), controllers.TransferirTicket)

	// Rutas admin (token + rol ADMIN)
	admin := r.Group("/admin", middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		admin.POST("/eventos", controllers.CrearEvento)
		admin.PUT("/eventos/:id", controllers.ActualizarEvento)
		admin.PUT("/eventos/:id/cancelar", controllers.CancelarEventoAdmin)
		admin.GET("/eventos/:id/reporte", controllers.GetReporteEvento)
	}

	r.Run(":8080")
}
