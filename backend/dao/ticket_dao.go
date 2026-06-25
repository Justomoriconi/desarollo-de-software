package dao

import (
	"backend/domain"
	"errors"
	"time"

	"gorm.io/gorm"
)

func ComprarTicket(usuarioID uint, tipoEntradaID uint, cuponID *uint, precioFinal float64) (*domain.Ticket, error) {
	var ticket domain.Ticket

	err := DB.Transaction(func(tx *gorm.DB) error {
		var tipoEntrada domain.TipoEntrada

		if err := tx.Preload("Evento").First(&tipoEntrada, tipoEntradaID).Error; err != nil {
			return errors.New("tipo de entrada no encontrado")
		}

		if tipoEntrada.Evento.Estado != "ACTIVO" {
			return errors.New("el evento no está activo")
		}

		if tipoEntrada.StockDisponible <= 0 {
			return errors.New("no hay stock disponible")
		}

		if err := tx.Model(&tipoEntrada).
			Update("stock_disponible", tipoEntrada.StockDisponible-1).Error; err != nil {
			return err
		}

		// Si se usó cupón, incrementar usos
		if cuponID != nil {
			if err := tx.Model(&domain.Cupon{}).
				Where("id = ?", *cuponID).
				Update("usos_actuales", gorm.Expr("usos_actuales + 1")).Error; err != nil {
				return err
			}
		}

		ticket = domain.Ticket{
			FechaCompra:   time.Now(),
			PrecioPagado:  precioFinal,
			Estado:        "ACTIVO",
			UsuarioID:     usuarioID,
			TipoEntradaID: tipoEntradaID,
			CuponID:       cuponID,
		}

		return tx.Create(&ticket).Error
	})

	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func GetTicketsByUsuario(usuarioID uint) ([]domain.Ticket, error) {
	var tickets []domain.Ticket

	result := DB.Preload("TipoEntrada").Preload("TipoEntrada.Evento").
		Where("usuario_id = ?", usuarioID).
		Find(&tickets)

	return tickets, result.Error
}

func GetTicketByID(id uint) (*domain.Ticket, error) {
	var ticket domain.Ticket

	result := DB.Preload("TipoEntrada").Preload("TipoEntrada.Evento").First(&ticket, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &ticket, nil
}

func CancelarTicket(ticket *domain.Ticket) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(ticket).Update("estado", "CANCELADO").Error; err != nil {
			return err
		}

		return tx.Model(&domain.TipoEntrada{}).
			Where("id = ?", ticket.TipoEntradaID).
			Update("stock_disponible", gorm.Expr("stock_disponible + 1")).Error
	})
}

func TransferirTicket(ticket *domain.Ticket, nuevoUsuarioID uint) error {
	return DB.Model(ticket).Update("usuario_id", nuevoUsuarioID).Error
}
