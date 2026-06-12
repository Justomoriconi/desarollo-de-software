package domain

import "time"

type Ticket struct {
	ID           uint      `gorm:"primaryKey"`
	FechaCompra  time.Time `gorm:"not null"`
	PrecioPagado float64   `gorm:"not null"`
	Estado       string    `gorm:"type:enum('ACTIVO','CANCELADO');not null"`

	UsuarioID uint    `gorm:"not null"`
	Usuario   Usuario `gorm:"foreignKey:UsuarioID"`

	TipoEntradaID uint        `gorm:"not null"`
	TipoEntrada   TipoEntrada `gorm:"foreignKey:TipoEntradaID"`

	CuponID *uint
	Cupon   *Cupon `gorm:"foreignKey:CuponID"` //no es necesario pero simplifica
}
