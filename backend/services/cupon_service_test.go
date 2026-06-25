package services_test

import (
	"backend/dao"
	"backend/domain"
	"backend/dto"
	"backend/services"
	"testing"
	"time"
)

func crearCuponDePrueba(t *testing.T, eventoID uint, codigo string, limiteUsos int) *domain.Cupon {
	cupon := &domain.Cupon{
		Codigo:           codigo,
		TipoDescuento:    "PORCENTAJE",
		ValorDescuento:   20,
		FechaVencimiento: time.Now().Add(24 * time.Hour),
		LimiteUsos:       limiteUsos,
		UsosActuales:     0,
		Estado:           "ACTIVO",
		EventoID:         eventoID,
	}
	if err := dao.DB.Create(cupon).Error; err != nil {
		t.Fatalf("error creando cupón de prueba: %v", err)
	}
	return cupon
}

func TestCrearCupon_Success(t *testing.T) {
	setupTestDB(t)

	eventoReq := dto.CrearEventoRequest{
		Nombre: "Evento",
		Lugar:  "Córdoba",
		Fecha:  time.Now().Add(24 * time.Hour),
	}
	evento, _ := services.CrearEvento(eventoReq)

	req := dto.CrearCuponRequest{
		Codigo:           "TEST10",
		TipoDescuento:    "PORCENTAJE",
		ValorDescuento:   10,
		FechaVencimiento: time.Now().Add(24 * time.Hour),
		LimiteUsos:       100,
		EventoID:         evento.ID,
	}

	cupon, err := services.CrearCupon(req)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if cupon.Codigo != "TEST10" {
		t.Errorf("código incorrecto: %s", cupon.Codigo)
	}
}

func TestCrearCupon_TipoInvalido(t *testing.T) {
	setupTestDB(t)

	req := dto.CrearCuponRequest{
		Codigo:           "TEST",
		TipoDescuento:    "INVALIDO",
		ValorDescuento:   10,
		FechaVencimiento: time.Now().Add(24 * time.Hour),
		LimiteUsos:       100,
	}

	_, err := services.CrearCupon(req)
	if err == nil {
		t.Errorf("se esperaba error por tipo de descuento inválido")
	}
}

func TestCrearCupon_FechaVencida(t *testing.T) {
	setupTestDB(t)

	req := dto.CrearCuponRequest{
		Codigo:           "TEST",
		TipoDescuento:    "PORCENTAJE",
		ValorDescuento:   10,
		FechaVencimiento: time.Now().Add(-24 * time.Hour),
		LimiteUsos:       100,
	}

	_, err := services.CrearCupon(req)
	if err == nil {
		t.Errorf("se esperaba error por fecha vencida")
	}
}

func TestValidarCupon_Success(t *testing.T) {
	setupTestDB(t)

	eventoReq := dto.CrearEventoRequest{
		Nombre: "Evento",
		Lugar:  "Córdoba",
		Fecha:  time.Now().Add(24 * time.Hour),
	}
	evento, _ := services.CrearEvento(eventoReq)

	tipoReq := dto.CrearTipoEntradaRequest{
		Nombre: "General", Precio: 10000, StockDisponible: 100,
	}
	tipo, _ := services.CrearTipoEntrada(evento.ID, tipoReq)

	crearCuponDePrueba(t, evento.ID, "VALID20", 100)

	resp, err := services.ValidarCupon("VALID20", tipo.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if resp.PrecioFinal != 8000 {
		t.Errorf("precio final esperado 8000, se obtuvo %.2f", resp.PrecioFinal)
	}
}

func TestValidarCupon_NoExiste(t *testing.T) {
	setupTestDB(t)

	_, err := services.ValidarCupon("NOEXISTE", 1)
	if err == nil {
		t.Errorf("se esperaba error por cupón inexistente")
	}
}

func TestValidarCupon_Inactivo(t *testing.T) {
	setupTestDB(t)

	eventoReq := dto.CrearEventoRequest{
		Nombre: "Evento", Lugar: "Córdoba",
		Fecha: time.Now().Add(24 * time.Hour),
	}
	evento, _ := services.CrearEvento(eventoReq)

	cupon := crearCuponDePrueba(t, evento.ID, "INACTIVO", 100)
	dao.DB.Model(cupon).Update("estado", "INACTIVO")

	_, err := services.ValidarCupon("INACTIVO", 1)
	if err == nil {
		t.Errorf("se esperaba error por cupón inactivo")
	}
}

func TestValidarCupon_LimiteAlcanzado(t *testing.T) {
	setupTestDB(t)

	eventoReq := dto.CrearEventoRequest{
		Nombre: "Evento", Lugar: "Córdoba",
		Fecha: time.Now().Add(24 * time.Hour),
	}
	evento, _ := services.CrearEvento(eventoReq)

	cupon := crearCuponDePrueba(t, evento.ID, "AGOTADO", 5)
	dao.DB.Model(cupon).Update("usos_actuales", 5)

	_, err := services.ValidarCupon("AGOTADO", 1)
	if err == nil {
		t.Errorf("se esperaba error por límite de usos alcanzado")
	}
}

func TestDesactivarCupon_Success(t *testing.T) {
	setupTestDB(t)

	eventoReq := dto.CrearEventoRequest{
		Nombre: "Evento", Lugar: "Córdoba",
		Fecha: time.Now().Add(24 * time.Hour),
	}
	evento, _ := services.CrearEvento(eventoReq)

	cupon := crearCuponDePrueba(t, evento.ID, "DESACT", 100)

	err := services.DesactivarCupon(cupon.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestDesactivarCupon_NoExiste(t *testing.T) {
	setupTestDB(t)

	err := services.DesactivarCupon(999)
	if err == nil {
		t.Errorf("se esperaba error por cupón inexistente")
	}
}
