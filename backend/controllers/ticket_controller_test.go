package controllers_test

import (
	"backend/dao"
	"backend/domain"
	"backend/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func crearUsuarioYToken(t *testing.T, email string) (uint, string) {
	usuario := domain.Usuario{
		Nombre: "Test", Email: email,
		PasswordHash: "hashfalso", Rol: "CLIENTE",
	}
	if err := dao.DB.Create(&usuario).Error; err != nil {
		t.Fatalf("error creando usuario: %v", err)
	}
	token, err := utils.GenerateToken(usuario.ID, usuario.Rol)
	if err != nil {
		t.Fatalf("error generando token: %v", err)
	}
	return usuario.ID, token
}

func crearTipoEntrada(t *testing.T, stock int) domain.TipoEntrada {
	evento := domain.Evento{
		Nombre: "Evento", Fecha: time.Now(),
		Lugar: "lugar", Estado: "ACTIVO",
	}
	dao.DB.Create(&evento)
	tipo := domain.TipoEntrada{
		Nombre: "General", Precio: 1000,
		StockDisponible: stock, EventoID: evento.ID,
	}
	dao.DB.Create(&tipo)
	return tipo
}

func TestComprarTicket_HTTP_SinToken(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	body, _ := json.Marshal(map[string]uint{"tipo_entrada_id": 1})
	req, _ := http.NewRequest("POST", "/tickets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("se esperaba 401 sin token, se obtuvo %d", w.Code)
	}
}

func TestComprarTicket_HTTP_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	_, token := crearUsuarioYToken(t, "comprador@test.com")
	tipo := crearTipoEntrada(t, 10)

	body, _ := json.Marshal(map[string]uint{"tipo_entrada_id": tipo.ID})
	req, _ := http.NewRequest("POST", "/tickets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestGetMisTickets_HTTP_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	_, token := crearUsuarioYToken(t, "comprador@test.com")

	req, _ := http.NewRequest("GET", "/tickets", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("se esperaba 200, se obtuvo %d", w.Code)
	}
}

func TestCancelarTicket_HTTP_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	usuarioID, token := crearUsuarioYToken(t, "comprador@test.com")
	tipo := crearTipoEntrada(t, 10)

	ticket := domain.Ticket{
		FechaCompra: time.Now(), PrecioPagado: 1000, Estado: "ACTIVO",
		UsuarioID: usuarioID, TipoEntradaID: tipo.ID,
	}
	dao.DB.Create(&ticket)

	req, _ := http.NewRequest("PUT", "/tickets/1/cancelar", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("se esperaba 200, se obtuvo %d. Body: %s", w.Code, w.Body.String())
	}
}

func itoa(n uint) string {
	return fmt.Sprintf("%d", n)
}
