package services_test

import (
	"backend/dao"
	"backend/domain"
	"backend/services"
	"testing"
	"time"
)

func crearEventoDePrueba(t *testing.T, nombre string, estado string) {
	evento := domain.Evento{
		Nombre:      nombre,
		Descripcion: "Descripcion de prueba",
		Fecha:       time.Now(),
		Lugar:       "Lugar de prueba",
		Estado:      estado,
	}

	if err := dao.DB.Create(&evento).Error; err != nil {
		t.Fatalf("error creando evento de prueba: %v", err)
	}
}

func TestGetEventos_SinFiltros(t *testing.T) {
	setupTestDB(t)

	crearEventoDePrueba(t, "Recital A", "ACTIVO")
	crearEventoDePrueba(t, "Recital B", "CANCELADO")

	eventos, err := services.GetEventos("", "")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if len(eventos) != 2 {
		t.Errorf("se esperaban 2 eventos, se obtuvieron %d", len(eventos))
	}
}

func TestGetEventos_FiltroPorEstado(t *testing.T) {
	setupTestDB(t)

	crearEventoDePrueba(t, "Recital A", "ACTIVO")
	crearEventoDePrueba(t, "Recital B", "CANCELADO")

	eventos, err := services.GetEventos("", "ACTIVO")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if len(eventos) != 1 {
		t.Errorf("se esperaba 1 evento, se obtuvieron %d", len(eventos))
	}

	if eventos[0].Estado != "ACTIVO" {
		t.Errorf("se esperaba estado ACTIVO, se obtuvo %s", eventos[0].Estado)
	}
}

func TestGetEventos_FiltroPorNombre(t *testing.T) {
	setupTestDB(t)

	crearEventoDePrueba(t, "Recital de Rock", "ACTIVO")
	crearEventoDePrueba(t, "Obra de Teatro", "ACTIVO")

	eventos, err := services.GetEventos("recital", "")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if len(eventos) != 1 {
		t.Errorf("se esperaba 1 evento, se obtuvieron %d", len(eventos))
	}
}

func TestGetEventoByID_NoExiste(t *testing.T) {
	setupTestDB(t)

	_, err := services.GetEventoByID(999)
	if err == nil {
		t.Errorf("se esperaba error por evento inexistente")
	}
}
