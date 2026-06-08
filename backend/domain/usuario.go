package domain

type Usuario struct {
	ID           uint   `gorm:"primaryKey"`
	Nombre       string `gorm:"not null"`
	Email        string `gorm:"not null;unique"`
	PasswordHash string `gorm:"not null"`
	Rol          string `gorm:"type:enum('CLIENTE','ADMIN');not null"`

	Tickets []Ticket `gorm:"foreignKey:UsuarioID"`
}
