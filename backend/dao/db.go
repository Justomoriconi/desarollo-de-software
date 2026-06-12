package dao

import (
	"backend/domain"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func InitDB() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No se encontró .env, usando variables de entorno del sistema")
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar a la base de datos :", err)
	}

	err = database.AutoMigrate(
		&domain.Usuario{},
		&domain.Evento{},
		&domain.TipoEntrada{},
		&domain.Cupon{},
		&domain.Ticket{},
	)

	if err != nil {
		log.Fatal("Error migrando la base de datos :", err)
	}

	DB = database
	log.Println("Base de datos conectada y migrada correctamente")

}
