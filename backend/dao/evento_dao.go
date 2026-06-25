package dao

import "backend/domain"

func GetEventos(nombre string, estado string) ([]domain.Evento, error) {
	var eventos []domain.Evento

	query := DB

	if nombre != "" {
		query = query.Where("nombre LIKE ?", "%"+nombre+"%")
	}

	if estado != "" {
		query = query.Where("estado = ?", estado)
	}

	result := query.Find(&eventos)
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

func GetTiposEntradaByEvento(eventoID uint) ([]domain.TipoEntrada, error) {
	var tipos []domain.TipoEntrada

	result := DB.Where("evento_id = ?", eventoID).Find(&tipos)
	return tipos, result.Error
}

func GetTiposEntradaByID(id uint) (*domain.TipoEntrada, error) {
	var tipo domain.TipoEntrada
	result := DB.First(&tipo, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tipo, nil
}
