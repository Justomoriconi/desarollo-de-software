package dto

import "time"

type CrearCuponRequest struct {
	Codigo           string    `json:"codigo" binding:"required"`
	TipoDescuento    string    `json:"tipo_descuento" binding:"required"`
	ValorDescuento   float64   `json:"valor_descuento" binding:"required"`
	FechaVencimiento time.Time `json:"fecha_vencimiento" binding:"required"`
	LimiteUsos       int       `json:"limite_usos" binding:"required"`
	EventoID         uint      `json:"evento_id"`
}

type ActualizarCuponRequest struct {
	TipoDescuento    string    `json:"tipo_descuento"`
	ValorDescuento   float64   `json:"valor_descuento"`
	FechaVencimiento time.Time `json:"fecha_vencimiento"`
	LimiteUsos       int       `json:"limite_usos"`
	Estado           string    `json:"estado"`
}

type ValidarCuponRequest struct {
	Codigo        string `json:"codigo" binding:"required"`
	TipoEntradaID uint   `json:"tipo_entrada_id" binding:"required"`
}

type ValidarCuponResponse struct {
	CuponID        uint    `json:"cupon_id"`
	Codigo         string  `json:"codigo"`
	TipoDescuento  string  `json:"tipo_descuento"`
	ValorDescuento float64 `json:"valor_descuento"`
	PrecioOriginal float64 `json:"precio_original"`
	PrecioFinal    float64 `json:"precio_final"`
	Descuento      float64 `json:"descuento"`
}

type ComprarConCuponRequest struct {
	TipoEntradaID uint   `json:"tipo_entrada_id" binding:"required"`
	CodigoCupon   string `json:"codigo_cupon"`
}
