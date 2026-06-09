package services

import (
	"backend/dao"
	"backend/domain"
)

func GetEventos() ([]domain.Evento, error) {
	return dao.GetEventos()
}
