package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	DB = db
	fmt.Println("Conectado ao Banco de Dados")
	return nil

}

func ConnectWithRetry() error {
	var err error

	for i := 0; i < 10; i++ {
		err = Connect()
		if err == nil {
			return nil
		}

		fmt.Println("Tentando conectar ao banco...")
		time.Sleep(2 * time.Second)
	}

	return err
}
