package dao

import (
	"backend/domain"
)

func CrearEvento(evento *domain.Evento) error {
	return DB.Create(evento).Error
}

func ActualizarEvento(evento *domain.Evento) error {
	return DB.Save(evento).Error
}

func CancelarEvento(id uint) error {
	return DB.Model(&domain.Evento{}).
		Where("id = ?", id).
		Update("estado", "CANCELADO").Error
}

func GetReporteEvento(eventoID uint) (*domain.Evento, []domain.Ticket, error) {
	var evento domain.Evento
	if err := DB.Preload("TiposEntrada").First(&evento, eventoID).Error; err != nil {
		return nil, nil, err
	}

	var tickets []domain.Ticket
	err := DB.Preload("Usuario").Preload("TipoEntrada").
		Where("tipo_entrada_id IN (SELECT id FROM tipo_entradas WHERE evento_id = ?)", eventoID).
		Find(&tickets).Error

	return &evento, tickets, err
}
