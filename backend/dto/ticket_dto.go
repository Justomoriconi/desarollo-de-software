package dto

import "time"

type ComprarTicketRequest struct {
	TipoEntradaID uint `json:"tipo_entrada_id"`
}

type TransferirTicketRequest struct {
	EmailDestino string `json:"email_destino"`
}

type TicketResponse struct {
	ID           uint      `json:"id"`
	EventoNombre string    `json:"evento_nombre"`
	TipoEntrada  string    `json:"tipo_entrada"`
	PrecioPagado float64   `json:"precio_pagado"`
	FechaCompra  time.Time `json:"fecha_compra"`
	Estado       string    `json:"estado"`
}