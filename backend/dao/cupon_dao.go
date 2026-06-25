package dao

import (
	"backend/domain"
	"errors"
)

func GetCuponByCodigo(codigo string) (*domain.Cupon, error) {
	var cupon domain.Cupon
	result := DB.Where("codigo = ?", codigo).First(&cupon)
	if result.Error != nil {
		return nil, errors.New("cupón no encontrado")
	}
	return &cupon, nil
}

func GetCupones() ([]domain.Cupon, error) {
	var cupones []domain.Cupon
	result := DB.Find(&cupones)
	return cupones, result.Error
}

func GetCuponByID(id uint) (*domain.Cupon, error) {
	var cupon domain.Cupon
	result := DB.First(&cupon, id)
	if result.Error != nil {
		return nil, errors.New("cupón no encontrado")
	}
	return &cupon, nil
}

func CrearCupon(cupon *domain.Cupon) error {
	return DB.Create(cupon).Error
}

func ActualizarCupon(cupon *domain.Cupon) error {
	return DB.Save(cupon).Error
}

func IncrementarUsosCupon(cuponID uint) error {
	return DB.Model(&domain.Cupon{}).
		Where("id = ?", cuponID).
		Update("usos_actuales", DB.Raw("usos_actuales + 1")).Error
}
