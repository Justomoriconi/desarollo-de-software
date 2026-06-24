package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegister_HTTP_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	body, _ := json.Marshal(map[string]string{
		"nombre": "Gaston", "email": "gaston@test.com", "password": "123456",
	})
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("se esperaba 201, se obtuvo %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestRegister_HTTP_JSONMalFormado(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString("{esto no es json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("se esperaba 400 con JSON mal formado, se obtuvo %d", w.Code)
	}
}

func TestLogin_HTTP_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	registerBody, _ := json.Marshal(map[string]string{
		"nombre": "Gaston", "email": "gaston@test.com", "password": "123456",
	})
	reqReg, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(registerBody))
	reqReg.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), reqReg)

	loginBody, _ := json.Marshal(map[string]string{
		"email": "gaston@test.com", "password": "123456",
	})
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["token"] == "" {
		t.Errorf("se esperaba un token en la respuesta")
	}
}

func TestLogin_HTTP_PasswordIncorrecta(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	registerBody, _ := json.Marshal(map[string]string{
		"nombre": "Gaston", "email": "gaston@test.com", "password": "123456",
	})
	reqReg, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(registerBody))
	reqReg.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), reqReg)

	loginBody, _ := json.Marshal(map[string]string{
		"email": "gaston@test.com", "password": "incorrecta",
	})
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("se esperaba 401, se obtuvo %d", w.Code)
	}
}
