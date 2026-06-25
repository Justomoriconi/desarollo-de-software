package services_test

import (
	"backend/dao"
	"backend/domain"
	"backend/dto"
	"backend/services"
	"testing"
	"time"
)

func TestGetCupones_Success(t *testing.T) {
	setupTestDB(t)

	eventoReq := dto.CrearEventoRequest{
		Nombre: "Evento", Lugar: "Córdoba",
		Fecha: time.Now().Add(24 * time.Hour),
	}
	evento, _ := services.CrearEvento(eventoReq)

	dao.DB.Create(&domain.Cupon{
		Codigo: "TEST", TipoDescuento: "PORCENTAJE",
		ValorDescuento: 10, FechaVencimiento: time.Now().Add(24 * time.Hour),
		LimiteUsos: 100, Estado: "ACTIVO", EventoID: evento.ID,
	})

	cupones, err := services.GetCupones()
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(cupones) != 1 {
		t.Errorf("se esperaba 1 cupón, se obtuvieron %d", len(cupones))
	}
}

func TestActualizarCupon_Success(t *testing.T) {
	setupTestDB(t)

	eventoReq := dto.CrearEventoRequest{
		Nombre: "Evento", Lugar: "Córdoba",
		Fecha: time.Now().Add(24 * time.Hour),
	}
	evento, _ := services.CrearEvento(eventoReq)

	cupon, _ := services.CrearCupon(dto.CrearCuponRequest{
		Codigo: "UPD10", TipoDescuento: "PORCENTAJE",
		ValorDescuento: 10, FechaVencimiento: time.Now().Add(24 * time.Hour),
		LimiteUsos: 100, EventoID: evento.ID,
	})

	actualizado, err := services.ActualizarCupon(cupon.ID, dto.ActualizarCuponRequest{
		ValorDescuento: 30,
	})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if actualizado.ValorDescuento != 30 {
		t.Errorf("valor no actualizado: %.0f", actualizado.ValorDescuento)
	}
}

func TestActualizarCupon_NoExiste(t *testing.T) {
	setupTestDB(t)

	_, err := services.ActualizarCupon(999, dto.ActualizarCuponRequest{Estado: "INACTIVO"})
	if err == nil {
		t.Errorf("se esperaba error por cupón inexistente")
	}
}

func TestGetTiposEntradaByEvento_Success(t *testing.T) {
	setupTestDB(t)

	eventoReq := dto.CrearEventoRequest{
		Nombre: "Evento", Lugar: "Córdoba",
		Fecha: time.Now().Add(24 * time.Hour),
	}
	evento, _ := services.CrearEvento(eventoReq)
	services.CrearTipoEntrada(evento.ID, dto.CrearTipoEntradaRequest{
		Nombre: "General", Precio: 5000, StockDisponible: 100,
	})

	tipos, err := services.GetTiposEntradaByEvento(evento.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(tipos) != 1 {
		t.Errorf("se esperaba 1 tipo, se obtuvieron %d", len(tipos))
	}
}

func TestComprarTicketConCupon_Success(t *testing.T) {
	setupTestDB(t)

	eventoReq := dto.CrearEventoRequest{
		Nombre: "Evento", Lugar: "Córdoba",
		Fecha: time.Now().Add(24 * time.Hour),
	}
	evento, _ := services.CrearEvento(eventoReq)

	tipo, _ := services.CrearTipoEntrada(evento.ID, dto.CrearTipoEntradaRequest{
		Nombre: "General", Precio: 10000, StockDisponible: 10,
	})

	cupon, _ := services.CrearCupon(dto.CrearCuponRequest{
		Codigo: "COMP20", TipoDescuento: "PORCENTAJE",
		ValorDescuento: 20, FechaVencimiento: time.Now().Add(24 * time.Hour),
		LimiteUsos: 100, EventoID: evento.ID,
	})

	usuarioID := crearUsuarioDePrueba(t, "comprador@test.com")

	ticket, err := services.ComprarTicketConCupon(usuarioID, tipo.ID, cupon.Codigo)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if ticket.PrecioPagado != 8000 {
		t.Errorf("precio esperado 8000, se obtuvo %.0f", ticket.PrecioPagado)
	}
}

func TestGetMisTickets_Success(t *testing.T) {
	setupTestDB(t)

	eventoReq := dto.CrearEventoRequest{
		Nombre: "Evento", Lugar: "Córdoba",
		Fecha: time.Now().Add(24 * time.Hour),
	}
	evento, _ := services.CrearEvento(eventoReq)
	tipo, _ := services.CrearTipoEntrada(evento.ID, dto.CrearTipoEntradaRequest{
		Nombre: "General", Precio: 5000, StockDisponible: 10,
	})

	usuarioID := crearUsuarioDePrueba(t, "usuario@test.com")
	services.ComprarTicket(usuarioID, tipo.ID)

	tickets, err := services.GetMisTickets(usuarioID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(tickets) != 1 {
		t.Errorf("se esperaba 1 ticket, se obtuvieron %d", len(tickets))
	}
}

func TestGetReporteEvento_ConTickets(t *testing.T) {
	setupTestDB(t)

	eventoReq := dto.CrearEventoRequest{
		Nombre: "Evento con tickets", Lugar: "Córdoba",
		Fecha: time.Now().Add(24 * time.Hour),
	}
	evento, _ := services.CrearEvento(eventoReq)
	tipo, _ := services.CrearTipoEntrada(evento.ID, dto.CrearTipoEntradaRequest{
		Nombre: "General", Precio: 5000, StockDisponible: 10,
	})

	usuarioID := crearUsuarioDePrueba(t, "comprador2@test.com")
	services.ComprarTicket(usuarioID, tipo.ID)

	reporte, err := services.GetReporteEvento(evento.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if reporte.TotalVendidas != 1 {
		t.Errorf("se esperaba 1 vendida, se obtuvo %d", reporte.TotalVendidas)
	}
}

func TestEliminarTipoEntrada_Success(t *testing.T) {
	setupTestDB(t)

	eventoReq := dto.CrearEventoRequest{
		Nombre: "Evento", Lugar: "Córdoba",
		Fecha: time.Now().Add(24 * time.Hour),
	}
	evento, _ := services.CrearEvento(eventoReq)
	tipo, _ := services.CrearTipoEntrada(evento.ID, dto.CrearTipoEntradaRequest{
		Nombre: "General", Precio: 5000, StockDisponible: 10,
	})

	err := services.EliminarTipoEntrada(tipo.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}
