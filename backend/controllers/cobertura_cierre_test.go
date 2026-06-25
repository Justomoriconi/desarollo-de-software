package controllers_test

import (
	"backend/dao"
	"backend/domain"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestValidarCupon_HTTP_NoExiste(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearUsuarioYToken(t, "cliente@test.com")

	body, _ := json.Marshal(map[string]interface{}{
		"codigo":          "NOEXISTE",
		"tipo_entrada_id": 1,
	})
	req, _ := http.NewRequest("POST", "/cupones/validar", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400, se obtuvo %d", w.Code)
	}
}

func TestValidarCupon_HTTP_BodyInvalido(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearUsuarioYToken(t, "cliente@test.com")

	req, _ := http.NewRequest("POST", "/cupones/validar", bytes.NewBufferString("{roto"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400, se obtuvo %d", w.Code)
	}
}

func TestCrearCupon_HTTP_BodyInvalido(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin@test.com")

	req, _ := http.NewRequest("POST", "/admin/cupones", bytes.NewBufferString("{roto"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400, se obtuvo %d", w.Code)
	}
}

func TestGetCupones_HTTP_SinToken(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/admin/cupones", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("se esperaba 401, se obtuvo %d", w.Code)
	}
}

func TestComprarConCuponInvalido_HTTP(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	crearEventoEnDB(t)
	tipo := crearTipoEntrada(t, 10)
	_, token := crearUsuarioYToken(t, "comprador@test.com")

	body, _ := json.Marshal(map[string]interface{}{
		"tipo_entrada_id": tipo.ID,
		"codigo_cupon":    "NOEXISTE",
	})
	req, _ := http.NewRequest("POST", "/tickets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400 con cupón inválido, se obtuvo %d", w.Code)
	}
}

func TestActualizarCupon_HTTP_BodyInvalido(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin@test.com")
	crearEventoEnDB(t)

	cupon := &domain.Cupon{
		Codigo: "BODYINV", TipoDescuento: "PORCENTAJE",
		ValorDescuento: 10, FechaVencimiento: time.Now().Add(24 * time.Hour),
		LimiteUsos: 100, Estado: "ACTIVO", EventoID: 1,
	}
	dao.DB.Create(cupon)

	req, _ := http.NewRequest("PUT", "/admin/cupones/1", bytes.NewBufferString("{roto"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400, se obtuvo %d", w.Code)
	}
}

func TestCrearTipoEntrada_HTTP_BodyInvalido(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin@test.com")
	crearEventoEnDB(t)

	req, _ := http.NewRequest("POST", "/admin/eventos/1/tipos-entrada", bytes.NewBufferString("{roto"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400, se obtuvo %d", w.Code)
	}
}

func TestActualizarEvento_HTTP_BodyInvalido(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin@test.com")
	crearEventoEnDB(t)

	req, _ := http.NewRequest("PUT", "/admin/eventos/1", bytes.NewBufferString("{roto"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400, se obtuvo %d", w.Code)
	}
}

func TestGetReporteEvento_HTTP_IDInvalido(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin@test.com")

	req, _ := http.NewRequest("GET", "/admin/eventos/abc/reporte", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400, se obtuvo %d", w.Code)
	}
}
