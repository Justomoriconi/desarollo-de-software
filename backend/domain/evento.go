package domain

import "time"

type Evento struct {
	ID          uint   `gorm:"primaryKey"`
	Nombre      string `gorm:"not null"`
	Descripcion string
	Fecha       time.Time `gorm:"not null"`
	Lugar       string    `gorm:"not null"`
	Estado      string    `gorm:"type:enum('ACTIVO','CANCELADO');not null"`

	TiposEntrada []TipoEntrada `gorm:"foreignKey:EventoID"`
	Cupones      []Cupon       `gorm:"foreignKey:EventoID"`
}
