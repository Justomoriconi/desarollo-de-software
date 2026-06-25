package services_test

import (
	"backend/dao"
	"testing"

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
