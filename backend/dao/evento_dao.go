package dao

import "backend/domain"

func GetEventos() ([]domain.Evento, error) {
	var eventos []domain.Evento
	result := DB.Find(&eventos)
	return eventos, result.Error
}

func GetEventoByID(id uint) (*domain.Evento, error) {
	var evento domain.Evento

	result := DB.First(&evento, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &evento, nil
}
