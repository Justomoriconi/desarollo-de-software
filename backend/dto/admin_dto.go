package dto

import "time"

type CrearEventoRequest struct {
	Nombre      string    `json:"nombre" binding:"required"`
	Descripcion string    `json:"descripcion"`
	Fecha       time.Time `json:"fecha" binding:"required"`
	Lugar       string    `json:"lugar" binding:"required"`
	Estado      string    `json:"estado"`
}

type ActualizarEventoRequest struct {
	Nombre      string    `json:"nombre"`
	Descripcion string    `json:"descripcion"`
	Fecha       time.Time `json:"fecha"`
	Lugar       string    `json:"lugar"`
	Estado      string    `json:"estado"`
}

type ReporteEventoResponse struct {
	EventoID        uint                 `json:"evento_id"`
	Nombre          string               `json:"nombre"`
	Estado          string               `json:"estado"`
	TotalTipos      int                  `json:"total_tipos_entrada"`
	TotalVendidas   int                  `json:"total_vendidas"`
	TotalCanceladas int                  `json:"total_canceladas"`
	TotalActivas    int                  `json:"total_activas"`
	Compradores     []CompradorResponse  `json:"compradores"`
	TiposEntrada    []TipoEntradaReporte `json:"tipos_entrada"`
}

type CompradorResponse struct {
	UsuarioID   uint   `json:"usuario_id"`
	Nombre      string `json:"nombre"`
	Email       string `json:"email"`
	TipoEntrada string `json:"tipo_entrada"`
	Estado      string `json:"estado"`
}

type TipoEntradaReporte struct {
	ID              uint    `json:"id"`
	Nombre          string  `json:"nombre"`
	Precio          float64 `json:"precio"`
	StockInicial    int     `json:"stock_inicial"`
	StockDisponible int     `json:"stock_disponible"`
	Vendidas        int     `json:"vendidas"`
}

type CrearTipoEntradaRequest struct {
	Nombre          string  `json:"nombre" binding:"required"`
	Precio          float64 `json:"precio" binding:"required"`
	StockDisponible int     `json:"stock_disponible" binding:"required"`
}
