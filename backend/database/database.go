package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL драйвер
)

const (
	// Параметры подключения
	host     = "localhost"
	//port     = 5433
	port     = 5432
	user     = "postgres"
	password = "password" // Замените на ваш пароль
	dbname   = "web_database"   // Замените на имя вашей базы данных
)

var db *sql.DB

// ConnectToServer - подключение к базе данных
func ConnectToServer() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
		return nil, err
	}

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных:", err)
		return nil, err
	}

	fmt.Println("Успешное подключение к базе данных!")
	return db, nil
}
