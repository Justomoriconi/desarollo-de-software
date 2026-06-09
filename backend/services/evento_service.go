package services

import (
	"backend/dao"
	"backend/domain"
)

func GetEventos() ([]domain.Evento, error) {
	return dao.GetEventos()
}

func GetEventoByID(id uint) (*domain.Evento, error) {
	return dao.GetEventoByID(id)
}
