package services_test

import (
	"backend/dto"
	"backend/services"
	"testing"
	"time"
)

func TestCrearEvento_Success(t *testing.T) {
	setupTestDB(t)

	req := dto.CrearEventoRequest{
		Nombre: "Evento Test",
		Lugar:  "Córdoba",
		Fecha:  time.Now().Add(24 * time.Hour),
	}

	evento, err := services.CrearEvento(req)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if evento.Nombre != "Evento Test" {
		t.Errorf("nombre incorrecto: %s", evento.Nombre)
	}
}

func TestCrearEvento_FechaPassada(t *testing.T) {
	setupTestDB(t)

	req := dto.CrearEventoRequest{
		Nombre: "Evento Viejo",
		Lugar:  "Córdoba",
		Fecha:  time.Now().Add(-24 * time.Hour),
	}

	_, err := services.CrearEvento(req)
	if err == nil {
		t.Errorf("se esperaba error por fecha pasada")
	}
}

func TestActualizarEvento_Success(t *testing.T) {
	setupTestDB(t)

	req := dto.CrearEventoRequest{
		Nombre: "Original",
		Lugar:  "Córdoba",
		Fecha:  time.Now().Add(24 * time.Hour),
	}
	evento, _ := services.CrearEvento(req)

	upd := dto.ActualizarEventoRequest{Nombre: "Actualizado"}
	actualizado, err := services.ActualizarEvento(evento.ID, upd)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if actualizado.Nombre != "Actualizado" {
		t.Errorf("nombre no actualizado: %s", actualizado.Nombre)
	}
}

func TestActualizarEvento_NoExiste(t *testing.T) {
	setupTestDB(t)

	_, err := services.ActualizarEvento(999, dto.ActualizarEventoRequest{Nombre: "X"})
	if err == nil {
		t.Errorf("se esperaba error por evento inexistente")
	}
}

func TestCancelarEventoAdmin_Success(t *testing.T) {
	setupTestDB(t)

	req := dto.CrearEventoRequest{
		Nombre: "Evento a Cancelar",
		Lugar:  "Córdoba",
		Fecha:  time.Now().Add(24 * time.Hour),
	}
	evento, _ := services.CrearEvento(req)

	err := services.CancelarEventoAdmin(evento.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestCancelarEventoAdmin_NoExiste(t *testing.T) {
	setupTestDB(t)

	err := services.CancelarEventoAdmin(999)
	if err == nil {
		t.Errorf("se esperaba error por evento inexistente")
	}
}

func TestGetReporteEvento_Success(t *testing.T) {
	setupTestDB(t)

	req := dto.CrearEventoRequest{
		Nombre: "Evento Reporte",
		Lugar:  "Córdoba",
		Fecha:  time.Now().Add(24 * time.Hour),
	}
	evento, _ := services.CrearEvento(req)

	reporte, err := services.GetReporteEvento(evento.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if reporte.EventoID != evento.ID {
		t.Errorf("ID de evento incorrecto")
	}
}

func TestGetReporteEvento_NoExiste(t *testing.T) {
	setupTestDB(t)

	_, err := services.GetReporteEvento(999)
	if err == nil {
		t.Errorf("se esperaba error por evento inexistente")
	}
}

func TestCrearTipoEntrada_Success(t *testing.T) {
	setupTestDB(t)

	req := dto.CrearEventoRequest{
		Nombre: "Evento",
		Lugar:  "Córdoba",
		Fecha:  time.Now().Add(24 * time.Hour),
	}
	evento, _ := services.CrearEvento(req)

	tipoReq := dto.CrearTipoEntradaRequest{
		Nombre:          "General",
		Precio:          5000,
		StockDisponible: 100,
	}

	tipo, err := services.CrearTipoEntrada(evento.ID, tipoReq)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if tipo.Nombre != "General" {
		t.Errorf("nombre incorrecto: %s", tipo.Nombre)
	}
}

func TestCrearTipoEntrada_EventoInexistente(t *testing.T) {
	setupTestDB(t)

	tipoReq := dto.CrearTipoEntradaRequest{
		Nombre:          "General",
		Precio:          5000,
		StockDisponible: 100,
	}

	_, err := services.CrearTipoEntrada(999, tipoReq)
	if err == nil {
		t.Errorf("se esperaba error por evento inexistente")
	}
}
