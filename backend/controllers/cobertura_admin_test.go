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

func TestEliminarTipoEntrada_HTTP_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin@test.com")
	crearEventoEnDB(t)
	tipo := crearTipoEntrada(t, 10)

	req, _ := http.NewRequest("DELETE", "/admin/eventos/1/tipos-entrada/"+itoa(tipo.ID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("se esperaba 200, se obtuvo %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestEliminarTipoEntrada_HTTP_IDInvalido(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin@test.com")

	req, _ := http.NewRequest("DELETE", "/admin/eventos/1/tipos-entrada/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400, se obtuvo %d", w.Code)
	}
}

func TestActualizarCupon_HTTP_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin@test.com")
	crearEventoEnDB(t)

	cupon := &domain.Cupon{
		Codigo: "ACT20", TipoDescuento: "PORCENTAJE",
		ValorDescuento: 10, FechaVencimiento: time.Now().Add(24 * time.Hour),
		LimiteUsos: 100, UsosActuales: 0, Estado: "ACTIVO", EventoID: 1,
	}
	dao.DB.Create(cupon)

	body, _ := json.Marshal(map[string]interface{}{
		"valor_descuento": 30,
	})
	req, _ := http.NewRequest("PUT", "/admin/cupones/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("se esperaba 200, se obtuvo %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestActualizarCupon_HTTP_IDInvalido(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin@test.com")

	body, _ := json.Marshal(map[string]interface{}{"valor_descuento": 30})
	req, _ := http.NewRequest("PUT", "/admin/cupones/abc", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400, se obtuvo %d", w.Code)
	}
}

func TestComprarConCupon_HTTP_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	crearEventoEnDB(t)
	tipo := crearTipoEntrada(t, 10)
	_, token := crearUsuarioYToken(t, "comprador@test.com")

	cupon := &domain.Cupon{
		Codigo: "COMP20", TipoDescuento: "PORCENTAJE",
		ValorDescuento: 20, FechaVencimiento: time.Now().Add(24 * time.Hour),
		LimiteUsos: 100, UsosActuales: 0, Estado: "ACTIVO", EventoID: tipo.EventoID,
	}
	dao.DB.Create(cupon)

	body, _ := json.Marshal(map[string]interface{}{
		"tipo_entrada_id": tipo.ID,
		"codigo_cupon":    "COMP20",
	})
	req, _ := http.NewRequest("POST", "/tickets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("se esperaba 201, se obtuvo %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestDesactivarCupon_HTTP_IDInvalido(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin@test.com")

	req, _ := http.NewRequest("PUT", "/admin/cupones/abc/desactivar", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400, se obtuvo %d", w.Code)
	}
}
