package services

import (
	"backend/dao"
	"backend/domain"
	"errors"
)

func ComprarTicket(usuarioID uint, tipoEntradaID uint) (*domain.Ticket, error) {
	return dao.ComprarTicket(usuarioID, tipoEntradaID)
}

func GetMisTickets(usuarioID uint) ([]domain.Ticket, error) {
	return dao.GetTicketsByUsuario(usuarioID)
}

func CancelarTicket(usuarioID uint, ticketID uint) error {
	ticket, err := dao.GetTicketByID(ticketID)
	if err != nil {
		return errors.New("ticket no encontrado")
	}

	if ticket.UsuarioID != usuarioID {
		return errors.New("el ticket no pertenece al usuario")
	}

	if ticket.Estado != "ACTIVO" {
		return errors.New("el ticket ya fue cancelado")
	}

	return dao.CancelarTicket(ticket)
}

func TransferirTicket(usuarioID uint, ticketID uint, emailDestino string) error {
	ticket, err := dao.GetTicketByID(ticketID)
	if err != nil {
		return errors.New("ticket no encontrado")
	}

	if ticket.UsuarioID != usuarioID {
		return errors.New("el ticket no pertenece al usuario")
	}

	if ticket.Estado != "ACTIVO" {
		return errors.New("no se puede transferir un ticket cancelado")
	}

	usuarioDestino, err := dao.GetUsuarioByEmail(emailDestino)
	if err != nil {
		return errors.New("usuario destino no encontrado")
	}

	if usuarioDestino.ID == usuarioID {
		return errors.New("no se puede transferir un ticket a uno mismo")
	}

	return dao.TransferirTicket(ticket, usuarioDestino.ID)
}