package dto

import "time"

type EventoDetailResponse struct {
	ID          uint      `json:"id"`
	Nombre      string    `json:"nombre"`
	Descripcion string    `json:"descripcion"`
	Fecha       time.Time `json:"fecha"`
	Lugar       string    `json:"lugar"`
	Estado      string    `json:"estado"`
}
