package services

import (
	"backend/dao"
	"backend/domain"
	"backend/utils"
	"errors"
)

func Register(nombre string, email string, password string) (uint, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return 0, err
	}

	usuario := domain.Usuario{
		Nombre:       nombre,
		Email:        email,
		PasswordHash: hashedPassword,
		Rol:          "CLIENTE",
	}

	err = dao.CreateUsuario(&usuario)
	if err != nil {
		return 0, err
	}

	return usuario.ID, nil

}

func Login(email string, password string) (string, error) {
	usuario, err := dao.GetUsuarioByEmail(email)
	if err != nil {
		return "", errors.New("usuario o contraseña incorrectos")
	}

	passwordCorrecta := utils.CheckPassword(password, usuario.PasswordHash)
	if !passwordCorrecta {
		return "", errors.New("usuario o contraseña incorrectos")
	}

	token, err := utils.GenerateToken(usuario.ID, usuario.Rol)
	if err != nil {
		return "", err
	}

	return token, nil
}
