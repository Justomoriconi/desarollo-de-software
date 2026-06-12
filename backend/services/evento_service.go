package services

import (
	"backend/dao"
	"backend/domain"
)

func GetEventos(nombre string, estado string) ([]domain.Evento, error) {
	return dao.GetEventos(nombre, estado)
}

func GetEventoByID(id uint) (*domain.Evento, error) {
	return dao.GetEventoByID(id)
}

func GetTiposEntradaByEvento(eventoID uint) ([]domain.TipoEntrada, error) {
	return dao.GetTiposEntradaByEvento(eventoID)
}
