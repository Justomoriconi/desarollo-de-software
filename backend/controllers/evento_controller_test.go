package controllers_test

import (
	"backend/dao"
	"backend/domain"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func crearEvento(t *testing.T, nombre, estado string) domain.Evento {
	evento := domain.Evento{
		Nombre: nombre, Descripcion: "desc", Fecha: time.Now(),
		Lugar: "lugar", Estado: estado,
	}
	if err := dao.DB.Create(&evento).Error; err != nil {
		t.Fatalf("error creando evento: %v", err)
	}
	return evento
}

func TestGetEventos_HTTP_SinFiltros(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	crearEvento(t, "Recital", "ACTIVO")
	crearEvento(t, "Teatro", "CANCELADO")

	req, _ := http.NewRequest("GET", "/eventos", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", w.Code)
	}

	var eventos []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &eventos)
	if len(eventos) != 2 {
		t.Errorf("se esperaban 2 eventos, se obtuvieron %d", len(eventos))
	}
}

func TestGetEventos_HTTP_FiltroEstado(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	crearEvento(t, "Recital", "ACTIVO")
	crearEvento(t, "Teatro", "CANCELADO")

	req, _ := http.NewRequest("GET", "/eventos?estado=ACTIVO", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var eventos []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &eventos)
	if len(eventos) != 1 {
		t.Errorf("se esperaba 1 evento, se obtuvieron %d", len(eventos))
	}
}

func TestGetEventoByID_HTTP_Existe(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	crearEvento(t, "Recital", "ACTIVO")

	req, _ := http.NewRequest("GET", "/eventos/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("se esperaba 200, se obtuvo %d", w.Code)
	}
}

func TestGetEventoByID_HTTP_NoExiste(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/eventos/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("se esperaba 404, se obtuvo %d", w.Code)
	}
}
