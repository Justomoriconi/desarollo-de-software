package domain

type TipoEntrada struct {
	ID              uint    `gorm:"primaryKey"`
	Nombre          string  `gorm:"not null"`
	Precio          float64 `gorm:"not null"`
	StockDisponible int     `gorm:"not null"`

	EventoID uint   `gorm:"not null"`
	Evento   Evento `gorm:"foreignKey:EventoID"`

	Tickets []Ticket `gorm:"foreignKey:TipoEntradaID"`
}
