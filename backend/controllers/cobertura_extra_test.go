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

func TestGetPerfil_HTTP_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	_, token := crearUsuarioYToken(t, "perfil@test.com")

	req, _ := http.NewRequest("GET", "/perfil", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("se esperaba 200, se obtuvo %d", w.Code)
	}
}

func TestGetPerfil_HTTP_SinToken(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/perfil", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("se esperaba 401 sin token, se obtuvo %d", w.Code)
	}
}

func TestGetTiposEntrada_HTTP_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	crearTipoEntrada(t, 10)

	req, _ := http.NewRequest("GET", "/eventos/1/tipos-entrada", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", w.Code)
	}
}

func TestGetTiposEntrada_HTTP_IDInvalido(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/eventos/abc/tipos-entrada", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400 con ID inválido, se obtuvo %d", w.Code)
	}
}

func TestTransferirTicket_HTTP_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	usuarioID, token := crearUsuarioYToken(t, "origen@test.com")
	crearUsuarioYToken(t, "destino@test.com")
	tipo := crearTipoEntrada(t, 10)

	ticket := domain.Ticket{
		FechaCompra: time.Now(), PrecioPagado: 1000, Estado: "ACTIVO",
		UsuarioID: usuarioID, TipoEntradaID: tipo.ID,
	}
	dao.DB.Create(&ticket)

	body, _ := json.Marshal(map[string]string{"email_destino": "destino@test.com"})
	req, _ := http.NewRequest("PUT", "/tickets/1/transferir", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("se esperaba 200, se obtuvo %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestTransferirTicket_HTTP_SinToken(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	body, _ := json.Marshal(map[string]string{"email_destino": "x@test.com"})
	req, _ := http.NewRequest("PUT", "/tickets/1/transferir", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("se esperaba 401 sin token, se obtuvo %d", w.Code)
	}
}

func TestTransferirTicket_HTTP_IDInvalido(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	_, token := crearUsuarioYToken(t, "origen@test.com")

	body, _ := json.Marshal(map[string]string{"email_destino": "x@test.com"})
	req, _ := http.NewRequest("PUT", "/tickets/abc/transferir", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400 con ID inválido, se obtuvo %d", w.Code)
	}
}

func TestTransferirTicket_HTTP_DestinoInexistente(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	usuarioID, token := crearUsuarioYToken(t, "origen@test.com")
	tipo := crearTipoEntrada(t, 10)

	ticket := domain.Ticket{
		FechaCompra: time.Now(), PrecioPagado: 1000, Estado: "ACTIVO",
		UsuarioID: usuarioID, TipoEntradaID: tipo.ID,
	}
	dao.DB.Create(&ticket)

	body, _ := json.Marshal(map[string]string{"email_destino": "noexiste@test.com"})
	req, _ := http.NewRequest("PUT", "/tickets/1/transferir", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400 con destino inexistente, se obtuvo %d", w.Code)
	}
}

func TestComprarTicket_HTTP_TipoInexistente(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	_, token := crearUsuarioYToken(t, "comprador@test.com")

	body, _ := json.Marshal(map[string]uint{"tipo_entrada_id": 999})
	req, _ := http.NewRequest("POST", "/tickets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400 con tipo inexistente, se obtuvo %d", w.Code)
	}
}

func TestComprarTicket_HTTP_SinStock(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	_, token := crearUsuarioYToken(t, "comprador@test.com")
	tipo := crearTipoEntrada(t, 0)

	body, _ := json.Marshal(map[string]uint{"tipo_entrada_id": tipo.ID})
	req, _ := http.NewRequest("POST", "/tickets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400 sin stock, se obtuvo %d", w.Code)
	}
}

func TestCancelarTicket_HTTP_IDInvalido(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	_, token := crearUsuarioYToken(t, "comprador@test.com")

	req, _ := http.NewRequest("PUT", "/tickets/abc/cancelar", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400 con ID inválido, se obtuvo %d", w.Code)
	}
}

func TestCancelarTicket_HTTP_Ajeno(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	otroID, _ := crearUsuarioYToken(t, "otro@test.com")
	_, token := crearUsuarioYToken(t, "comprador@test.com")
	tipo := crearTipoEntrada(t, 10)

	ticket := domain.Ticket{
		FechaCompra: time.Now(), PrecioPagado: 1000, Estado: "ACTIVO",
		UsuarioID: otroID, TipoEntradaID: tipo.ID,
	}
	dao.DB.Create(&ticket)

	req, _ := http.NewRequest("PUT", "/tickets/1/cancelar", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400 al cancelar ticket ajeno, se obtuvo %d", w.Code)
	}
}

func TestComprarTicket_HTTP_TokenInvalido(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	body, _ := json.Marshal(map[string]uint{"tipo_entrada_id": 1})
	req, _ := http.NewRequest("POST", "/tickets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token.invalido.aca")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("se esperaba 401 con token inválido, se obtuvo %d", w.Code)
	}
}
func TestGetMisTickets_HTTP_ConTickets(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	usuarioID, token := crearUsuarioYToken(t, "contickets@test.com")
	tipo := crearTipoEntrada(t, 10)

	ticket := domain.Ticket{
		FechaCompra: time.Now(), PrecioPagado: 1000, Estado: "ACTIVO",
		UsuarioID: usuarioID, TipoEntradaID: tipo.ID,
	}
	dao.DB.Create(&ticket)

	req, _ := http.NewRequest("GET", "/tickets", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d", w.Code)
	}

	var tickets []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &tickets)
	if len(tickets) != 1 {
		t.Errorf("se esperaba 1 ticket, se obtuvieron %d", len(tickets))
	}
}

func TestGetMisTickets_HTTP_TokenInvalido(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/tickets", nil)
	req.Header.Set("Authorization", "Bearer token.malo.aca")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("se esperaba 401 con token inválido, se obtuvo %d", w.Code)
	}
}

func TestComprarTicket_HTTP_BodyInvalido(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	_, token := crearUsuarioYToken(t, "bodyinvalido@test.com")

	req, _ := http.NewRequest("POST", "/tickets", bytes.NewBufferString("{no es json"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400 con body inválido, se obtuvo %d", w.Code)
	}
}

func TestTransferirTicket_HTTP_BodyInvalido(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	usuarioID, token := crearUsuarioYToken(t, "transfbody@test.com")
	tipo := crearTipoEntrada(t, 10)

	ticket := domain.Ticket{
		FechaCompra: time.Now(), PrecioPagado: 1000, Estado: "ACTIVO",
		UsuarioID: usuarioID, TipoEntradaID: tipo.ID,
	}
	dao.DB.Create(&ticket)

	req, _ := http.NewRequest("PUT", "/tickets/1/transferir", bytes.NewBufferString("{roto"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400 con body inválido, se obtuvo %d", w.Code)
	}
}
