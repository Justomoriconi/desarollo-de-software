package domain

import "time"

type Cupon struct {
	ID               uint      `gorm:"primaryKey"`
	Codigo           string    `gorm:"not null;unique"`
	TipoDescuento    string    `gorm:"type:enum('PORCENTAJE','MONTO_FIJO');not null"`
	ValorDescuento   float64   `gorm:"not null"`
	FechaVencimiento time.Time `gorm:"not null"`
	LimiteUsos       int       `gorm:"not null"`
	UsosActuales     int       `gorm:"not null;default:0"`
	Estado           string    `gorm:"type:enum('ACTIVO','INACTIVO');not null"`

	EventoID uint   `gorm:"not null"`
	Evento   Evento `gorm:"foreignKey:EventoID"`

	Tickets []Ticket `gorm:"foreignKey:CuponID"`
}
