package services

import (
	"backend/dao"
	"backend/domain"
	"backend/dto"
	"errors"
	"time"
)

func CrearEvento(req dto.CrearEventoRequest) (*domain.Evento, error) {
	if req.Nombre == "" || req.Lugar == "" {
		return nil, errors.New("nombre y lugar son obligatorios")
	}

	if req.Fecha.Before(time.Now()) {
		return nil, errors.New("la fecha del evento debe ser futura")
	}

	estado := req.Estado
	if estado == "" {
		estado = "ACTIVO"
	}

	evento := &domain.Evento{
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		Fecha:       req.Fecha,
		Lugar:       req.Lugar,
		Estado:      estado,
	}

	if err := dao.CrearEvento(evento); err != nil {
		return nil, errors.New("error al crear el evento")
	}

	return evento, nil
}

func ActualizarEvento(id uint, req dto.ActualizarEventoRequest) (*domain.Evento, error) {
	evento, err := dao.GetEventoByID(id)
	if err != nil {
		return nil, errors.New("evento no encontrado")
	}

	if req.Nombre != "" {
		evento.Nombre = req.Nombre
	}
	if req.Descripcion != "" {
		evento.Descripcion = req.Descripcion
	}
	if !req.Fecha.IsZero() {
		evento.Fecha = req.Fecha
	}
	if req.Lugar != "" {
		evento.Lugar = req.Lugar
	}
	if req.Estado != "" {
		evento.Estado = req.Estado
	}

	if err := dao.ActualizarEvento(evento); err != nil {
		return nil, errors.New("error al actualizar el evento")
	}

	return evento, nil
}

func CancelarEventoAdmin(id uint) error {
	_, err := dao.GetEventoByID(id)
	if err != nil {
		return errors.New("evento no encontrado")
	}

	return dao.CancelarEvento(id)
}

func GetReporteEvento(eventoID uint) (*dto.ReporteEventoResponse, error) {
	evento, tickets, err := dao.GetReporteEvento(eventoID)
	if err != nil {
		return nil, errors.New("evento no encontrado")
	}

	reporte := &dto.ReporteEventoResponse{
		EventoID:   evento.ID,
		Nombre:     evento.Nombre,
		Estado:     evento.Estado,
		TotalTipos: len(evento.TiposEntrada),
	}

	// Compradores
	for _, ticket := range tickets {
		reporte.TotalVendidas++
		if ticket.Estado == "CANCELADO" {
			reporte.TotalCanceladas++
		} else {
			reporte.TotalActivas++
		}

		reporte.Compradores = append(reporte.Compradores, dto.CompradorResponse{
			UsuarioID:   ticket.UsuarioID,
			Nombre:      ticket.Usuario.Nombre,
			Email:       ticket.Usuario.Email,
			TipoEntrada: ticket.TipoEntrada.Nombre,
			Estado:      ticket.Estado,
		})
	}

	// Tipos de entrada con métricas
	for _, tipo := range evento.TiposEntrada {
		vendidas := 0
		for _, t := range tickets {
			if t.TipoEntradaID == tipo.ID {
				vendidas++
			}
		}

		reporte.TiposEntrada = append(reporte.TiposEntrada, dto.TipoEntradaReporte{
			ID:              tipo.ID,
			Nombre:          tipo.Nombre,
			Precio:          tipo.Precio,
			StockDisponible: tipo.StockDisponible,
			Vendidas:        vendidas,
		})
	}

	return reporte, nil
}

func CrearTipoEntrada(eventoID uint, req dto.CrearTipoEntradaRequest) (*domain.TipoEntrada, error) {
	_, err := dao.GetEventoByID(eventoID)
	if err != nil {
		return nil, errors.New("evento no encontrado")
	}

	tipo := &domain.TipoEntrada{
		Nombre:          req.Nombre,
		Precio:          req.Precio,
		StockDisponible: req.StockDisponible,
		EventoID:        eventoID,
	}

	if err := dao.CrearTipoEntrada(tipo); err != nil {
		return nil, errors.New("error al crear el tipo de entrada")
	}

	return tipo, nil
}

func EliminarTipoEntrada(id uint) error {
	return dao.EliminarTipoEntrada(id)
}
