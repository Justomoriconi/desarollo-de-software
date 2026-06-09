package dao

import "backend/domain"

func GetUsuarioByEmail(email string) (*domain.Usuario, error) {
	var usuario domain.Usuario

	result := DB.Where("email = ?", email).First(&usuario)

	if result.Error != nil {
		return nil, result.Error
	}

	return &usuario, nil
}

func CreateUsuario(usuario *domain.Usuario) error {
	result := DB.Create(usuario)

	return result.Error
}
