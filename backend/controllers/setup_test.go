package controllers_test

import (
	"backend/controllers"
	"backend/dao"
	"backend/domain"
	"backend/middlewares"
	"backend/utils"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("error abriendo db de test: %v", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS usuarios (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		rol TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS eventos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT NOT NULL,
		descripcion TEXT,
		fecha DATETIME NOT NULL,
		lugar TEXT NOT NULL,
		estado TEXT NOT NULL DEFAULT 'ACTIVO'
	);
	CREATE TABLE IF NOT EXISTS tipo_entradas (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT NOT NULL,
		precio REAL NOT NULL,
		stock_disponible INTEGER NOT NULL,
		evento_id INTEGER NOT NULL
	);
	CREATE TABLE IF NOT EXISTS cupons (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		codigo TEXT NOT NULL UNIQUE,
		tipo_descuento TEXT NOT NULL,
		valor_descuento REAL NOT NULL,
		fecha_vencimiento DATETIME NOT NULL,
		limite_usos INTEGER NOT NULL,
		usos_actuales INTEGER NOT NULL DEFAULT 0,
		estado TEXT NOT NULL DEFAULT 'ACTIVO',
		evento_id INTEGER NOT NULL
	);
	CREATE TABLE IF NOT EXISTS tickets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		fecha_compra DATETIME NOT NULL,
		precio_pagado REAL NOT NULL,
		estado TEXT NOT NULL DEFAULT 'ACTIVO',
		usuario_id INTEGER NOT NULL,
		tipo_entrada_id INTEGER NOT NULL,
		cupon_id INTEGER
	);
	`
	if err := db.Exec(schema).Error; err != nil {
		t.Fatalf("error creando schema de test: %v", err)
	}

	dao.DB = db
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

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
	r.POST("/cupones/validar", middlewares.AuthMiddleware(), controllers.ValidarCupon)

	admin := r.Group("/admin", middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		admin.POST("/eventos", controllers.CrearEvento)
		admin.PUT("/eventos/:id", controllers.ActualizarEvento)
		admin.PUT("/eventos/:id/cancelar", controllers.CancelarEventoAdmin)
		admin.GET("/eventos/:id/reporte", controllers.GetReporteEvento)
		admin.POST("/eventos/:id/tipos-entrada", controllers.CrearTipoEntrada)
		admin.DELETE("/eventos/:id/tipos-entrada/:tipoId", controllers.EliminarTipoEntrada)
		admin.GET("/cupones", controllers.GetCupones)
		admin.POST("/cupones", controllers.CrearCupon)
		admin.PUT("/cupones/:id", controllers.ActualizarCupon)
		admin.PUT("/cupones/:id/desactivar", controllers.DesactivarCupon)
	}

	return r
}

func crearAdminYToken(t *testing.T, email string) (uint, string) {
	usuario := domain.Usuario{
		Nombre: "Admin", Email: email,
		PasswordHash: "hashfalso", Rol: "ADMIN",
	}
	if err := dao.DB.Create(&usuario).Error; err != nil {
		t.Fatalf("error creando admin: %v", err)
	}
	token, err := utils.GenerateToken(usuario.ID, usuario.Rol)
	if err != nil {
		t.Fatalf("error generando token: %v", err)
	}
	return usuario.ID, token
}

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "secret_de_prueba")
	os.Exit(m.Run())
}
