package services_test

import (
	"backend/dao"
	"backend/domain"
	"backend/services"
	"testing"
	"time"
)

func crearEventoConTipoEntrada(t *testing.T, stock int) (uint, uint) {
	evento := domain.Evento{
		Nombre: "Evento Test",
		Fecha:  time.Now(),
		Lugar:  "Lugar Test",
		Estado: "ACTIVO",
	}
	if err := dao.DB.Create(&evento).Error; err != nil {
		t.Fatalf("error creando evento: %v", err)
	}

	tipoEntrada := domain.TipoEntrada{
		Nombre:          "General",
		Precio:          1000,
		StockDisponible: stock,
		EventoID:        evento.ID,
	}
	if err := dao.DB.Create(&tipoEntrada).Error; err != nil {
		t.Fatalf("error creando tipo de entrada: %v", err)
	}

	return evento.ID, tipoEntrada.ID
}

func crearUsuarioDePrueba(t *testing.T, email string) uint {
	id, err := services.Register("Usuario Test", email, "123456")
	if err != nil {
		t.Fatalf("error creando usuario: %v", err)
	}
	return id
}

func TestComprarTicket_Success(t *testing.T) {
	setupTestDB(t)

	usuarioID := crearUsuarioDePrueba(t, "comprador@test.com")
	_, tipoEntradaID := crearEventoConTipoEntrada(t, 10)

	ticket, err := services.ComprarTicket(usuarioID, tipoEntradaID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if ticket.Estado != "ACTIVO" {
		t.Errorf("se esperaba estado ACTIVO, se obtuvo %s", ticket.Estado)
	}

	var tipoEntrada domain.TipoEntrada
	dao.DB.First(&tipoEntrada, tipoEntradaID)

	if tipoEntrada.StockDisponible != 9 {
		t.Errorf("se esperaba stock 9, se obtuvo %d", tipoEntrada.StockDisponible)
	}
}

func TestComprarTicket_SinStock(t *testing.T) {
	setupTestDB(t)

	usuarioID := crearUsuarioDePrueba(t, "comprador@test.com")
	_, tipoEntradaID := crearEventoConTipoEntrada(t, 0)

	_, err := services.ComprarTicket(usuarioID, tipoEntradaID)
	if err == nil {
		t.Errorf("se esperaba error por falta de stock")
	}
}

func TestComprarTicket_TipoEntradaInexistente(t *testing.T) {
	setupTestDB(t)

	usuarioID := crearUsuarioDePrueba(t, "comprador@test.com")

	_, err := services.ComprarTicket(usuarioID, 999)
	if err == nil {
		t.Errorf("se esperaba error por tipo de entrada inexistente")
	}
}

func TestCancelarTicket_Success(t *testing.T) {
	setupTestDB(t)

	usuarioID := crearUsuarioDePrueba(t, "comprador@test.com")
	_, tipoEntradaID := crearEventoConTipoEntrada(t, 10)

	ticket, err := services.ComprarTicket(usuarioID, tipoEntradaID)
	if err != nil {
		t.Fatalf("error en compra previa: %v", err)
	}

	err = services.CancelarTicket(usuarioID, ticket.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	var tipoEntrada domain.TipoEntrada
	dao.DB.First(&tipoEntrada, tipoEntradaID)

	if tipoEntrada.StockDisponible != 10 {
		t.Errorf("se esperaba que el stock vuelva a 10, se obtuvo %d", tipoEntrada.StockDisponible)
	}
}

func TestCancelarTicket_NoPertenece(t *testing.T) {
	setupTestDB(t)

	usuarioID := crearUsuarioDePrueba(t, "comprador@test.com")
	otroUsuarioID := crearUsuarioDePrueba(t, "otro@test.com")
	_, tipoEntradaID := crearEventoConTipoEntrada(t, 10)

	ticket, err := services.ComprarTicket(usuarioID, tipoEntradaID)
	if err != nil {
		t.Fatalf("error en compra previa: %v", err)
	}

	err = services.CancelarTicket(otroUsuarioID, ticket.ID)
	if err == nil {
		t.Errorf("se esperaba error porque el ticket no pertenece al usuario")
	}
}

func TestCancelarTicket_YaCancelado(t *testing.T) {
	setupTestDB(t)

	usuarioID := crearUsuarioDePrueba(t, "comprador@test.com")
	_, tipoEntradaID := crearEventoConTipoEntrada(t, 10)

	ticket, err := services.ComprarTicket(usuarioID, tipoEntradaID)
	if err != nil {
		t.Fatalf("error en compra previa: %v", err)
	}

	if err := services.CancelarTicket(usuarioID, ticket.ID); err != nil {
		t.Fatalf("error en primera cancelacion: %v", err)
	}

	err = services.CancelarTicket(usuarioID, ticket.ID)
	if err == nil {
		t.Errorf("se esperaba error por ticket ya cancelado")
	}
}

func TestTransferirTicket_Success(t *testing.T) {
	setupTestDB(t)

	usuarioID := crearUsuarioDePrueba(t, "comprador@test.com")
	destinoID := crearUsuarioDePrueba(t, "destino@test.com")
	_, tipoEntradaID := crearEventoConTipoEntrada(t, 10)

	ticket, err := services.ComprarTicket(usuarioID, tipoEntradaID)
	if err != nil {
		t.Fatalf("error en compra previa: %v", err)
	}

	err = services.TransferirTicket(usuarioID, ticket.ID, "destino@test.com")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	var ticketActualizado domain.Ticket
	dao.DB.First(&ticketActualizado, ticket.ID)

	if ticketActualizado.UsuarioID != destinoID {
		t.Errorf("se esperaba usuario_id %d, se obtuvo %d", destinoID, ticketActualizado.UsuarioID)
	}
}

func TestTransferirTicket_UsuarioDestinoInexistente(t *testing.T) {
	setupTestDB(t)

	usuarioID := crearUsuarioDePrueba(t, "comprador@test.com")
	_, tipoEntradaID := crearEventoConTipoEntrada(t, 10)

	ticket, err := services.ComprarTicket(usuarioID, tipoEntradaID)
	if err != nil {
		t.Fatalf("error en compra previa: %v", err)
	}

	err = services.TransferirTicket(usuarioID, ticket.ID, "noexiste@test.com")
	if err == nil {
		t.Errorf("se esperaba error por usuario destino inexistente")
	}
}

func TestTransferirTicket_AUnoMismo(t *testing.T) {
	setupTestDB(t)

	usuarioID := crearUsuarioDePrueba(t, "comprador@test.com")
	_, tipoEntradaID := crearEventoConTipoEntrada(t, 10)

	ticket, err := services.ComprarTicket(usuarioID, tipoEntradaID)
	if err != nil {
		t.Fatalf("error en compra previa: %v", err)
	}

	err = services.TransferirTicket(usuarioID, ticket.ID, "comprador@test.com")
	if err == nil {
		t.Errorf("se esperaba error por transferencia a uno mismo")
	}
}
