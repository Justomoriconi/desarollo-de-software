package services_test

import (
	"backend/services"
	"testing"
)

func TestRegister_Success(t *testing.T) {
	setupTestDB(t)

	id, err := services.Register("Gaston", "gaston@test.com", "123456")

	if err != nil {
		t.Fatalf("no se esperaba error, se obtuvo: %v", err)
	}

	if id == 0 {
		t.Errorf("se esperaba un id valido, se obtuvo 0")
	}
}

func TestRegister_EmailDuplicado(t *testing.T) {
	setupTestDB(t)

	_, err := services.Register("Gaston", "gaston@test.com", "123456")
	if err != nil {
		t.Fatalf("no se esperaba error en el primer registro: %v", err)
	}

	_, err = services.Register("Otro", "gaston@test.com", "654321")
	if err == nil {
		t.Errorf("se esperaba error por email duplicado")
	}
}

func TestLogin_Success(t *testing.T) {
	setupTestDB(t)

	_, err := services.Register("Gaston", "gaston@test.com", "123456")
	if err != nil {
		t.Fatalf("error en registro previo: %v", err)
	}

	token, err := services.Login("gaston@test.com", "123456")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if token == "" {
		t.Errorf("se esperaba un token no vacio")
	}
}

func TestLogin_PasswordIncorrecta(t *testing.T) {
	setupTestDB(t)

	_, err := services.Register("Gaston", "gaston@test.com", "123456")
	if err != nil {
		t.Fatalf("error en registro previo: %v", err)
	}

	_, err = services.Login("gaston@test.com", "incorrecta")
	if err == nil {
		t.Errorf("se esperaba error por password incorrecta")
	}
}

func TestLogin_UsuarioInexistente(t *testing.T) {
	setupTestDB(t)

	_, err := services.Login("noexiste@test.com", "123456")
	if err == nil {
		t.Errorf("se esperaba error por usuario inexistente")
	}
}
