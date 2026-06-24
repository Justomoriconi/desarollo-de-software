package controllers_test

import (
	"backend/controllers"
	"backend/dao"
	"backend/middlewares"
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
	CREATE TABLE usuarios (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		rol TEXT NOT NULL
	);
	CREATE TABLE eventos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT NOT NULL,
		descripcion TEXT,
		fecha DATETIME NOT NULL,
		lugar TEXT NOT NULL,
		estado TEXT NOT NULL
	);
	CREATE TABLE tipo_entradas (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT NOT NULL,
		precio REAL NOT NULL,
		stock_disponible INTEGER NOT NULL,
		evento_id INTEGER NOT NULL
	);
	CREATE TABLE cupons (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		codigo TEXT NOT NULL UNIQUE,
		tipo_descuento TEXT NOT NULL,
		valor_descuento REAL NOT NULL,
		fecha_vencimiento DATETIME NOT NULL,
		limite_usos INTEGER NOT NULL,
		usos_actuales INTEGER NOT NULL DEFAULT 0,
		estado TEXT NOT NULL,
		evento_id INTEGER NOT NULL
	);
	CREATE TABLE tickets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		fecha_compra DATETIME NOT NULL,
		precio_pagado REAL NOT NULL,
		estado TEXT NOT NULL,
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

	return r
}

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "secret_de_prueba")
	os.Exit(m.Run())
}
