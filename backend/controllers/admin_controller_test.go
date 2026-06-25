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

// ---- helpers admin ----

func crearEventoEnDB(t *testing.T) domain.Evento {
	evento := domain.Evento{
		Nombre: "Evento Test", Descripcion: "desc",
		Fecha: time.Now().Add(24 * time.Hour),
		Lugar: "Córdoba", Estado: "ACTIVO",
	}
	if err := dao.DB.Create(&evento).Error; err != nil {
		t.Fatalf("error creando evento: %v", err)
	}
	return evento
}

// ---- CrearEvento ----

func TestCrearEvento_HTTP_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin@test.com")

	body, _ := json.Marshal(map[string]interface{}{
		"nombre": "Nuevo Evento",
		"lugar":  "Córdoba",
		"fecha":  time.Now().Add(48 * time.Hour).Format(time.RFC3339),
	})
	req, _ := http.NewRequest("POST", "/admin/eventos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("se esperaba 201, se obtuvo %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestCrearEvento_HTTP_SinToken(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	body, _ := json.Marshal(map[string]interface{}{
		"nombre": "Evento",
		"lugar":  "Córdoba",
		"fecha":  time.Now().Add(48 * time.Hour).Format(time.RFC3339),
	})
	req, _ := http.NewRequest("POST", "/admin/eventos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("se esperaba 401, se obtuvo %d", w.Code)
	}
}

func TestCrearEvento_HTTP_RolCliente(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearUsuarioYToken(t, "cliente@test.com")

	body, _ := json.Marshal(map[string]interface{}{
		"nombre": "Evento",
		"lugar":  "Córdoba",
		"fecha":  time.Now().Add(48 * time.Hour).Format(time.RFC3339),
	})
	req, _ := http.NewRequest("POST", "/admin/eventos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("se esperaba 403, se obtuvo %d", w.Code)
	}
}

func TestCrearEvento_HTTP_BodyInvalido(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin2@test.com")

	req, _ := http.NewRequest("POST", "/admin/eventos", bytes.NewBufferString("{roto"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400, se obtuvo %d", w.Code)
	}
}

// ---- ActualizarEvento ----

func TestActualizarEvento_HTTP_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin@test.com")
	evento := crearEventoEnDB(t)

	body, _ := json.Marshal(map[string]string{"nombre": "Actualizado"})
	req, _ := http.NewRequest("PUT", "/admin/eventos/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("se esperaba 200, se obtuvo %d. Body: %s", w.Code, w.Body.String())
	}
	_ = evento
}

func TestActualizarEvento_HTTP_IDInvalido(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin@test.com")

	body, _ := json.Marshal(map[string]string{"nombre": "X"})
	req, _ := http.NewRequest("PUT", "/admin/eventos/abc", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400, se obtuvo %d", w.Code)
	}
}

// ---- CancelarEventoAdmin ----

func TestCancelarEventoAdmin_HTTP_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin@test.com")
	crearEventoEnDB(t)

	req, _ := http.NewRequest("PUT", "/admin/eventos/1/cancelar", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("se esperaba 200, se obtuvo %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestCancelarEventoAdmin_HTTP_IDInvalido(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin@test.com")

	req, _ := http.NewRequest("PUT", "/admin/eventos/abc/cancelar", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400, se obtuvo %d", w.Code)
	}
}

// ---- GetReporteEvento ----

func TestGetReporteEvento_HTTP_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin@test.com")
	crearEventoEnDB(t)

	req, _ := http.NewRequest("GET", "/admin/eventos/1/reporte", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("se esperaba 200, se obtuvo %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestGetReporteEvento_HTTP_NoExiste(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin@test.com")

	req, _ := http.NewRequest("GET", "/admin/eventos/999/reporte", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("se esperaba 404, se obtuvo %d", w.Code)
	}
}

// ---- CrearTipoEntrada ----

func TestCrearTipoEntrada_HTTP_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin@test.com")
	crearEventoEnDB(t)

	body, _ := json.Marshal(map[string]interface{}{
		"nombre": "General", "precio": 5000, "stock_disponible": 100,
	})
	req, _ := http.NewRequest("POST", "/admin/eventos/1/tipos-entrada", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("se esperaba 201, se obtuvo %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestCrearTipoEntrada_HTTP_IDInvalido(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin@test.com")

	body, _ := json.Marshal(map[string]interface{}{
		"nombre": "General", "precio": 5000, "stock_disponible": 100,
	})
	req, _ := http.NewRequest("POST", "/admin/eventos/abc/tipos-entrada", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400, se obtuvo %d", w.Code)
	}
}

// ---- Cupones ----

func TestCrearCupon_HTTP_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin@test.com")
	crearEventoEnDB(t)

	body, _ := json.Marshal(map[string]interface{}{
		"codigo":            "TEST20",
		"tipo_descuento":    "PORCENTAJE",
		"valor_descuento":   20,
		"fecha_vencimiento": time.Now().Add(48 * time.Hour).Format(time.RFC3339),
		"limite_usos":       100,
		"evento_id":         1,
	})
	req, _ := http.NewRequest("POST", "/admin/cupones", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("se esperaba 201, se obtuvo %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestGetCupones_HTTP_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin@test.com")

	req, _ := http.NewRequest("GET", "/admin/cupones", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("se esperaba 200, se obtuvo %d", w.Code)
	}
}

func TestDesactivarCupon_HTTP_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearAdminYToken(t, "admin@test.com")
	crearEventoEnDB(t)

	cupon := &domain.Cupon{
		Codigo: "DESACT", TipoDescuento: "PORCENTAJE",
		ValorDescuento: 10, FechaVencimiento: time.Now().Add(24 * time.Hour),
		LimiteUsos: 100, UsosActuales: 0, Estado: "ACTIVO", EventoID: 1,
	}
	dao.DB.Create(cupon)

	req, _ := http.NewRequest("PUT", "/admin/cupones/1/desactivar", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("se esperaba 200, se obtuvo %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestValidarCupon_HTTP_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	_, token := crearUsuarioYToken(t, "cliente@test.com")
	crearEventoEnDB(t)

	tipo := crearTipoEntrada(t, 10)
	cupon := &domain.Cupon{
		Codigo: "VALID20", TipoDescuento: "PORCENTAJE",
		ValorDescuento: 20, FechaVencimiento: time.Now().Add(24 * time.Hour),
		LimiteUsos: 100, UsosActuales: 0, Estado: "ACTIVO", EventoID: tipo.EventoID,
	}
	dao.DB.Create(cupon)

	body, _ := json.Marshal(map[string]interface{}{
		"codigo":          "VALID20",
		"tipo_entrada_id": tipo.ID,
	})
	req, _ := http.NewRequest("POST", "/cupones/validar", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("se esperaba 200, se obtuvo %d. Body: %s", w.Code, w.Body.String())
	}
}
