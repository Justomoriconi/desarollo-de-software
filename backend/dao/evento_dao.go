package dao

import "backend/domain"

func GetEventos() ([]domain.Evento, error) {
	var eventos []domain.Evento
	result := DB.Find(&eventos)
	return eventos, result.Error
}
